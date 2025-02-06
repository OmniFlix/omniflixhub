package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
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

	// Validate the media node registration details
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Create and store the media node
	mediaNode := types.NewMediaNode(msg.Url, sender.String(), msg.HardwareSpecs, msg.PricePerDay)
	if err := m.Keeper.RegisterMediaNode(ctx, mediaNode); err != nil {
		return nil, err
	}

	return &types.MsgRegisterMediaNodeResponse{}, nil
}

// UpdateMediaNode handles the update of an existing media node
func (m msgServer) UpdateMediaNode(goCtx context.Context, msg *types.MsgUpdateMediaNode) (*types.MsgUpdateMediaNodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// Validate the media node update details
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Retrieve the existing media node
	existingMediaNode, found := m.Keeper.GetMediaNode(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrMediaNodeDoesNotExist, "not found")
	}
	existingMediaNode.HardwareSpecs = msg.HardwareSpecs
	existingMediaNode.PricePerDay = msg.PricePerDay
	m.Keeper.UpdateMediaNode(ctx, existingMediaNode, sender)

	return &types.MsgUpdateMediaNodeResponse{}, nil
}

// LeaseMediaNode handles leasing a media node
func (m msgServer) LeaseMediaNode(goCtx context.Context, msg *types.MsgLeaseMediaNode) (*types.MsgLeaseMediaNodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// Validate the lease details
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Lease the media node
	if err := m.Keeper.LeaseMediaNode(ctx, msg.MediaNodeId, msg.LeaseDays, sender); err != nil {
		return nil, err
	}

	return &types.MsgLeaseMediaNodeResponse{}, nil
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

	// Validate the deposit details
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Deposit to the media node
	if err := m.Keeper.DepositMediaNode(ctx, msg.MediaNodeId, msg.Amount, sender); err != nil {
		return nil, err
	}

	return &types.MsgDepositMediaNodeResponse{}, nil
}

// CloseMediaNode handles closing a media node
func (m msgServer) CloseMediaNode(goCtx context.Context, msg *types.MsgCloseMediaNode) (*types.MsgCloseMediaNodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// Validate the close details
	if err := msg.ValidateBasic(); err != nil {
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
