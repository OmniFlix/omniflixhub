package keeper

import (
	"fmt"
	"time"

	"github.com/OmniFlix/omniflixhub/v5/x/itc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) emitCreateCampaignEvent(ctx context.Context, campaign types.Campaign) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateCampaign,
			sdk.NewAttribute(types.AttributeKeyCampaignId, fmt.Sprint(campaign.GetId())),
			sdk.NewAttribute(types.AttributeKeyCreator, campaign.Creator),
			sdk.NewAttribute(types.AttributeKeyNftDenomId, campaign.NftDenomId),
			sdk.NewAttribute(types.AttributeKeyClaimType, campaign.ClaimType.String()),
			sdk.NewAttribute(types.AttributeKeyInteractionType, campaign.Interaction.String()),
			sdk.NewAttribute(types.AttributeKeyStartTime, campaign.StartTime.Format(time.RFC3339)),
			sdk.NewAttribute(types.AttributeKeyEndTime, campaign.EndTime.Format(time.RFC3339)),
		),
	)
}

func (k Keeper) emitDepositCampaignEvent(ctx context.Context, campaignId uint64, depositor string, amount sdk.Coin) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDepositCampaign,
			sdk.NewAttribute(types.AttributeKeyCampaignId, fmt.Sprint(campaignId)),
			sdk.NewAttribute(types.AttributeKeyDepositor, depositor),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	)
}

func (k Keeper) emitCancelCampaignEvent(ctx context.Context, campaignId uint64, sender string) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCancelCampaign,
			sdk.NewAttribute(types.AttributeKeyCampaignId, fmt.Sprint(campaignId)),
			sdk.NewAttribute(types.AttributeKeySender, sender),
		),
	)
}

func (k Keeper) emitClaimEvent(ctx context.Context, campaign uint64, claimer string, nftId string) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeClaim,
			sdk.NewAttribute(types.AttributeKeyCampaignId, fmt.Sprint(campaign)),
			sdk.NewAttribute(types.AttributeKeyClaimer, claimer),
			sdk.NewAttribute(types.AttributeKeyNftId, nftId),
		),
	)
}

func (k Keeper) emitEndCampaignEvent(ctx context.Context, campaignId uint64) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEndCampaign,
			sdk.NewAttribute(types.AttributeKeyCampaignId, fmt.Sprint(campaignId)),
		),
	)
}
