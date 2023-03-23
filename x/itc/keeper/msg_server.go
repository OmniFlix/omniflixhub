package keeper

import (
	"context"
	"github.com/OmniFlix/omniflixhub/x/itc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
	msg *types.MsgCreateCampaign,
) (*types.MsgCreateCampaignResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}
	// StartTime must be after current time
	if msg.StartTime.Before(ctx.BlockTime()) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidDuration, "start time must be in future")
	}
	endTime := msg.StartTime.Add(msg.Duration)
	if endTime.Before(msg.StartTime) || endTime.Equal(msg.StartTime) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidDuration, "duration must be positive or nil")
	}
	if msg.Duration > m.Keeper.GetMaxCampaignDuration(ctx) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidDuration,
			"duration must be less than max campaign duration (%d)", m.Keeper.GetMaxCampaignDuration(ctx))
	}

	availableTokens := msg.TotalTokens
	campaignNumber := m.Keeper.GetNextCampaignNumber(ctx)
	campaign := types.NewCampaign(campaignNumber,
		msg.Name,
		msg.Description,
		msg.StartTime,
		endTime,
		creator.String(),
		msg.NftDenomId,
		msg.MaxAllowedClaims,
		msg.Interaction,
		msg.ClaimType,
		msg.ClaimableTokens,
		msg.TotalTokens,
		availableTokens,
		msg.NftMintDetails,
		msg.Distribution,
	)
	err = m.Keeper.CreateCampaign(ctx, creator, campaign)
	if err != nil {
		return nil, err
	}
	// Event

	return &types.MsgCreateCampaignResponse{}, nil
}

// CancelCampaign

func (m msgServer) CancelCampaign(goCtx context.Context,
	msg *types.MsgCancelCampaign,
) (*types.MsgCancelCampaignResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	err = m.Keeper.CancelCampaign(ctx, msg.CampaignId, creator)
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
		return nil, sdkerrors.Wrapf(types.ErrInactiveCampaign,
			"cannot claim on inactive campaign %d, ", campaign.Id)
	}
	if msg.Interaction != campaign.Interaction {
		return nil, sdkerrors.Wrapf(types.ErrInteractionMismatch,
			"required interaction %s, got %s. ", campaign.Interaction, msg.Interaction)
	}

	claim := types.NewClaim(campaign.Id, claimer.String(), msg.NftId, msg.Interaction)

	err = m.Keeper.Claim(ctx, campaign, claimer, claim)
	if err != nil {
		return nil, err
	}
	// TODO: event

	return &types.MsgClaimResponse{}, nil
}

func (m msgServer) CampaignDeposit(goCtx context.Context,
	msg *types.MsgCampaignDeposit,
) (*types.MsgCampaignDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	depositor, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return nil, err
	}
	campaign, found := m.Keeper.GetCampaign(ctx, msg.CampaignId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrCampaignDoesNotExists, "campaign id %d not exists", msg.CampaignId)
	}
	if msg.Amount.Denom != campaign.TotalTokens.Fungible.Denom {
		return nil, sdkerrors.Wrapf(types.ErrTokenDenomMismatch,
			"token denom mismatch, required %s, got %s", campaign.TotalTokens.Fungible.Denom, msg.Amount.Denom)
	}
	if err := m.Keeper.CampaignDeposit(ctx, msg.CampaignId, depositor, msg.Amount); err != nil {
		return nil, err
	}
	// TODO: event

	return &types.MsgCampaignDepositResponse{}, nil
}
