package keeper

import (
	"context"
	"github.com/OmniFlix/omniflixhub/x/itc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"time"
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

// CreateCampaign

func (m msgServer) CreateCampaign(goCtx context.Context,
	msg *types.MsgCreateCampaign) (*types.MsgCreateCampaignResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	var endTime time.Time

	endAt := msg.StartTime.Add(msg.Duration)
	endTime = endAt
	if endTime.Before(msg.StartTime) || endTime.Equal(msg.StartTime) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidDuration, "duration must be positive or nil")
	}
	campaignNumber := m.Keeper.GetNextCampaignNumber(ctx)
	campaign := types.NewCampaign(campaignNumber,
		msg.Name, msg.Description,
		msg.StartTime, endTime,
		creator.String(),
		msg.NftDenomId,
		msg.MaxAllowedClaims,
		msg.Interaction,
		msg.ClaimType,
		msg.ClaimableTokens,
		msg.TotalTokens,
		msg.TotalTokens,
		[]string{},
		msg.NftMintDetails,
		&msg.Distribution,
	)
	err = m.Keeper.CreateCampaign(ctx, campaign)
	if err != nil {
		return nil, err
	}
	// Event

	return &types.MsgCreateCampaignResponse{}, nil
}

// CancelCampaign

func (m msgServer) CancelCampaign(goCtx context.Context,
	msg *types.MsgCancelCampaign) (*types.MsgCancelCampaignResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	campaign, found := m.Keeper.GetCampaign(ctx, msg.CampaignId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrCampaignDoesNotExists, "campaign %d not exists", msg.CampaignId)
	}
	if creator.String() != campaign.Creator {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "unauthorized address %s", creator.String())
	}

	err = m.Keeper.CancelCampaign(ctx, campaign)
	if err != nil {
		return nil, err
	}
	// TODO: event

	return &types.MsgCancelCampaignResponse{}, nil
}

// Claim

func (m msgServer) Claim(goCtx context.Context, msg *types.MsgClaim) (*types.MsgClaimResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	claimer, err := sdk.AccAddressFromBech32(msg.Claimer)
	if err != nil {
		return nil, err
	}

	campaign, found := m.Keeper.GetCampaign(ctx, msg.CampaignId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrCampaignDoesNotExists, "campaign id %d not exists", msg.CampaignId)
	}
	if !campaign.StartTime.Before(ctx.BlockTime()) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidTokens, "cannot do claim on inactive campaign %d, ", campaign.Id)
	}

	claim := types.NewClaim(campaign.Id, claimer.String())

	err = m.Keeper.Claim(ctx, campaign, claim)
	if err != nil {
		return nil, err
	}
	// TODO: event

	return &types.MsgClaimResponse{}, nil
}

func (m msgServer) CampaignDeposit(goCtx context.Context,
	msg *types.MsgCampaignDeposit) (*types.MsgCampaignDepositResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)
	_, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return nil, err
	}

	_, found := m.Keeper.GetCampaign(ctx, msg.CampaignId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrCampaignDoesNotExists, "campaign id %d not exists", msg.CampaignId)
	}
	// TODO: add tokens to total tokens

	// TODO: event

	return &types.MsgCampaignDepositResponse{}, nil
}
