package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/OmniFlix/omniflixhub/v6/x/medianode/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the token MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// RegisterMediaNode handles the registration of a new media node
func (m msgServer) RegisterMediaNode(goCtx context.Context, msg *types.MsgRegisterMediaNode) (*types.MsgRegisterMediaNodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	minDeposit := m.Keeper.GetMinDeposit(ctx)
	initialDepositPerc := m.Keeper.GetInitialDepositPercentage(ctx)
	minInitialDeposit := sdk.NewCoin(minDeposit.Denom, sdkmath.LegacyNewDecFromInt(minDeposit.Amount).Mul(initialDepositPerc).TruncateInt())
	if !msg.Deposit.IsGTE(minInitialDeposit) {
		return nil, errorsmod.Wrapf(types.ErrInsufficientDeposit, "%s of initial deposit is reqquired", minInitialDeposit.String())
	}

	// Create and store the media node
	mediaNode := types.NewMediaNode(msg.Id, msg.Url, sender.String(), msg.HardwareSpecs, msg.Info, msg.PricePerHour)
	// default status on register
	mediaNode.Status = types.STATUS_PENDING

	if msg.Deposit.Amount.GTE(minDeposit.Amount) {
		mediaNode.Status = types.STATUS_ACTIVE
	}

	mediaNodeData, err := m.Keeper.RegisterMediaNode(ctx, mediaNode, *msg.Deposit, sender)
	if err != nil {
		return nil, err
	}

	return &types.MsgRegisterMediaNodeResponse{
		MediaNodeId: mediaNodeData.Id,
		Status:      mediaNodeData.Status.String(),
	}, nil
}

// UpdateMediaNode handles the update of an existing media node
func (m msgServer) UpdateMediaNode(goCtx context.Context, msg *types.MsgUpdateMediaNode) (*types.MsgUpdateMediaNodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	mediaNode, err := m.Keeper.UpdateMediaNode(ctx, msg.Id, msg.Info, msg.HardwareSpecs, msg.PricePerHour, sender)
	if err != nil {
		return nil, err
	}

	return &types.MsgUpdateMediaNodeResponse{
		MediaNode: &mediaNode,
	}, nil
}

// LeaseMediaNode handles leasing a media node
func (m msgServer) LeaseMediaNode(goCtx context.Context, msg *types.MsgLeaseMediaNode) (*types.MsgLeaseMediaNodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	params := m.Keeper.GetParams(ctx)
	if msg.LeaseHours < params.MinimumLeaseHours {
		return nil, errorsmod.Wrapf(types.ErrInvalidLeaseHours, "minimum of %d lease hours required", params.MinimumLeaseHours)
	}
	if msg.LeaseHours > params.MaximumLeaseHours {
		return nil, errorsmod.Wrapf(types.ErrInvalidLeaseHours, "maximum of %d lease hours allowed", params.MaximumLeaseHours)
	}

	mediaNode, found := m.Keeper.GetMediaNode(ctx, msg.MediaNodeId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrMediaNodeDoesNotExist, "not found")
	}
	if mediaNode.Status != types.STATUS_ACTIVE {
		return nil, errorsmod.Wrapf(types.ErrLeaseNotAllowed, "only active medianodes are allowed to lease")
	}

	if mediaNode.IsLeased() {
		return nil, errorsmod.Wrapf(types.ErrMediaNodeAlreadyLeased, "media node %s is already leased", mediaNode.Id)
	}

	expectedLeaseAmount := sdk.NewCoin(mediaNode.PricePerHour.Denom, mediaNode.PricePerHour.Amount.Mul(sdkmath.NewIntFromUint64(msg.LeaseHours)))
	if !msg.Amount.IsEqual(expectedLeaseAmount) {
		return nil, errorsmod.Wrapf(types.ErrInvalidLeaseAmount, "lease amount must be equal to %s", expectedLeaseAmount.String())
	}

	// Lease the media node
	lease, err := m.Keeper.LeaseMediaNode(ctx, mediaNode, msg.LeaseHours, sender, msg.Amount)
	if err != nil {
		return nil, err
	}

	return &types.MsgLeaseMediaNodeResponse{
		Lease: &lease,
	}, nil
}

// LeaseMediaNode handles leasing a media node
func (m msgServer) ExtendLease(goCtx context.Context, msg *types.MsgExtendLease) (*types.MsgExtendLeaseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	params := m.Keeper.GetParams(ctx)
	if msg.LeaseHours < params.MinimumLeaseHours {
		return nil, errorsmod.Wrapf(types.ErrInvalidLeaseHours, "minimum of %d lease hours required", params.MinimumLeaseHours)
	}
	if msg.LeaseHours > params.MaximumLeaseHours {
		return nil, errorsmod.Wrapf(types.ErrInvalidLeaseHours, "maximum of %d lease hours allowed", params.MaximumLeaseHours)
	}

	mediaNodeLease, found := m.Keeper.GetMediaNodeLease(ctx, msg.MediaNodeId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrLeaseNotFound, "not found")
	}

	expectedLeaseAmount := sdk.NewCoin(mediaNodeLease.PricePerHour.Denom, mediaNodeLease.PricePerHour.Amount.Mul(sdkmath.NewIntFromUint64((msg.LeaseHours))))
	if !msg.Amount.IsEqual(expectedLeaseAmount) {
		return nil, errorsmod.Wrapf(types.ErrInvalidLeaseAmount, "lease amount must be equal to %s", expectedLeaseAmount.String())
	}
	// Lease the media node
	lease, err := m.Keeper.ExtendMediaNodeLease(ctx, mediaNodeLease, msg.LeaseHours, msg.Amount, sender)
	if err != nil {
		return nil, err
	}

	return &types.MsgExtendLeaseResponse{
		Lease: &lease,
	}, nil
}

// CancelLease handles canceling a lease for a media node
func (m msgServer) CancelLease(goCtx context.Context, msg *types.MsgCancelLease) (*types.MsgCancelLeaseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// Cancel the lease
	if err := m.Keeper.CancelLease(ctx, msg.MediaNodeId, sender); err != nil {
		return nil, err
	}

	return &types.MsgCancelLeaseResponse{}, nil
}

// DepositMediaNode handles depositing to a media node
func (m msgServer) DepositMediaNode(goCtx context.Context, msg *types.MsgDepositMediaNode) (*types.MsgDepositMediaNodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// Deposit to the media node
	mediaNode, totalDeposit, err := m.Keeper.DepositMediaNode(ctx, msg.MediaNodeId, msg.Amount, sender)
	if err != nil {
		return nil, err
	}

	return &types.MsgDepositMediaNodeResponse{
		MedianodeId:     mediaNode.Id,
		TotalDeposit:    &totalDeposit,
		MedianodeStatus: mediaNode.Status.String(),
	}, nil
}

// CloseMediaNode handles closing a media node
func (m msgServer) CloseMediaNode(goCtx context.Context, msg *types.MsgCloseMediaNode) (*types.MsgCloseMediaNodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	// Close the media node
	if err := m.Keeper.CloseMediaNode(ctx, msg.MediaNodeId, sender); err != nil {
		return nil, err
	}

	return &types.MsgCloseMediaNodeResponse{}, nil
}

func (m msgServer) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if m.Keeper.GetAuthority() != req.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", m.authority, req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.SetParams(ctx, req.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
