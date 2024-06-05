package keeper

import (
	"context"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	errorsmod "cosmossdk.io/errors"

	"github.com/OmniFlix/omniflixhub/v5/x/itc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	if msg.Deposit.Denom != msg.TokensPerClaim.Denom {
		return nil, errorsmod.Wrapf(types.ErrInvalidTokens, "mismatched token denoms")
	}
	// StartTime must be after current time
	if msg.StartTime.Before(ctx.BlockTime()) {
		return nil, errorsmod.Wrapf(types.ErrInvalidDuration, "start time must be in future")
	}
	endTime := msg.StartTime.Add(msg.Duration)
	if endTime.Before(msg.StartTime) || endTime.Equal(msg.StartTime) {
		return nil, errorsmod.Wrapf(types.ErrInvalidDuration, "duration must be positive or nil")
	}
	if msg.Duration > m.Keeper.GetMaxCampaignDuration(ctx) {
		return nil, errorsmod.Wrapf(types.ErrInvalidDuration,
			"duration must be less than max campaign duration (%d)", m.Keeper.GetMaxCampaignDuration(ctx))
	}
	campaignCreationFee := m.Keeper.GetCampaignCreationFee(ctx)
	if !msg.CreationFee.Equal(campaignCreationFee) {
		if msg.CreationFee.Denom != campaignCreationFee.Denom {
			return nil, errorsmod.Wrapf(types.ErrInvalidFeeDenom, "invalid creation fee denom %s",
				msg.CreationFee.Denom)
		}
		if msg.CreationFee.Amount.LT(campaignCreationFee.Amount) {
			return nil, errorsmod.Wrapf(types.ErrNotEnoughFeeAmount,
				"%s fee is not enough, to create %s fee is required",
				msg.CreationFee.String(), campaignCreationFee.String())
		}
		return nil, errorsmod.Wrapf(types.ErrInvalidCreationFee,
			"given fee (%s) not matched with  campaign creation fee. %s required to create itc campaign",
			msg.CreationFee.String(), campaignCreationFee.String())
	}

	if (msg.ClaimType == types.CLAIM_TYPE_FT || msg.ClaimType == types.CLAIM_TYPE_FT_AND_NFT) && msg.Distribution == nil {
		return nil, errorsmod.Wrapf(
			types.ErrInvalidNFTMintDetails,
			"distribution config is required for ft claim type",
		)
	}
	if (msg.ClaimType == types.CLAIM_TYPE_NFT || msg.ClaimType == types.CLAIM_TYPE_FT_AND_NFT) && msg.NftMintDetails == nil {
		return nil, errorsmod.Wrapf(
			types.ErrInvalidNFTMintDetails,
			"nft mint details are required for nft claim type",
		)
	}

	// cut creation fee amount and fund the community pool
	if err := m.distributionKeeper.FundCommunityPool(
		ctx,
		sdk.NewCoins(msg.CreationFee),
		creator,
	); err != nil {
		return nil, err
	}

	availableTokens := msg.Deposit
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
		msg.TokensPerClaim,
		msg.Deposit,
		availableTokens,
		msg.NftMintDetails,
		msg.Distribution,
	)
	err = m.Keeper.CreateCampaign(ctx, creator, campaign)
	if err != nil {
		return nil, err
	}

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
		return nil, errorsmod.Wrapf(types.ErrCampaignDoesNotExists, "campaign id %d not exists", msg.CampaignId)
	}
	if campaign.StartTime.Unix() > ctx.BlockTime().Unix() {
		return nil, errorsmod.Wrapf(types.ErrInactiveCampaign,
			"cannot claim on inactive campaign %d, ", campaign.Id)
	}
	if msg.Interaction != campaign.Interaction {
		return nil, errorsmod.Wrapf(types.ErrInteractionMismatch,
			"required interaction %s, got %s. ", campaign.Interaction, msg.Interaction)
	}

	claim := types.NewClaim(campaign.Id, claimer.String(), msg.NftId, msg.Interaction)

	err = m.Keeper.Claim(ctx, campaign, claimer, claim)
	if err != nil {
		return nil, err
	}

	return &types.MsgClaimResponse{}, nil
}

func (m msgServer) DepositCampaign(goCtx context.Context,
	msg *types.MsgDepositCampaign,
) (*types.MsgDepositCampaignResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	depositor, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return nil, err
	}
	campaign, found := m.Keeper.GetCampaign(ctx, msg.CampaignId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrCampaignDoesNotExists, "campaign id %d not exists", msg.CampaignId)
	}
	if campaign.ClaimType == types.CLAIM_TYPE_NFT {
		return nil, errorsmod.Wrapf(types.ErrDepositNotAllowed, "deposit not allowed for this type of campaign")
	}
	if msg.Amount.Denom != campaign.TotalTokens.Denom {
		return nil, errorsmod.Wrapf(types.ErrTokenDenomMismatch,
			"token denom mismatch, required %s, got %s", campaign.TotalTokens.Denom, msg.Amount.Denom)
	}
	if err := m.Keeper.DepositCampaign(ctx, msg.CampaignId, depositor, msg.Amount); err != nil {
		return nil, err
	}

	return &types.MsgDepositCampaignResponse{}, nil
}

func (m msgServer) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if m.authority != req.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", m.authority, req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.SetParams(ctx, req.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
