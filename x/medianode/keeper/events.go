package keeper

import (
	"fmt"

	"github.com/OmniFlix/omniflixhub/v6/x/medianode/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) registerMediaNodeEvent(ctx sdk.Context, owner string, mediaNodeId uint64, URL string, pricePerDay sdk.Coin, status types.Status) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRegisterMediaNode,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyOwner, owner),
			sdk.NewAttribute(types.AttributeKeyMediaNodeId, fmt.Sprint(mediaNodeId)),
			sdk.NewAttribute(types.AttributeKeyMediaNodeURL, URL),
			sdk.NewAttribute(types.AttributeKeyStatus, status.String()),
			sdk.NewAttribute(types.AttributeKeyPricePerDay, pricePerDay.String()),
		),
	})
}

func (k *Keeper) updateMediaNodeEvent(ctx sdk.Context, sender string, mediaNodeId uint64) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateMediaNode,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, sender),
			sdk.NewAttribute(types.AttributeKeyMediaNodeId, fmt.Sprint(mediaNodeId)),
		),
	})
}

func (k *Keeper) leaseMediaNodeEvent(ctx sdk.Context, lessee string, mediaNodeId uint64, amount sdk.Coin) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeLeaseMediaNode,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyLessee, lessee),
			sdk.NewAttribute(types.AttributeKeyMediaNodeId, fmt.Sprint(mediaNodeId)),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})
}

func (k *Keeper) cancelLeaseMediaNodeEvent(ctx sdk.Context, lessee string, mediaNodeId uint64) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCalcelLease,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyLessee, lessee),
			sdk.NewAttribute(types.AttributeKeyMediaNodeId, fmt.Sprint(mediaNodeId)),
		),
	})
}

func (k *Keeper) depositMediaNodeEvent(ctx sdk.Context, depositor string, mediaNodeId uint64, amount sdk.Coin) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDepositMediaNode,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyDepositor, depositor),
			sdk.NewAttribute(types.AttributeKeyMediaNodeId, fmt.Sprint(mediaNodeId)),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})
}

func (k *Keeper) closeMediaNodeEvent(ctx sdk.Context, owner string, mediaNodeId uint64) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCloseMediaNode,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyOwner, owner),
			sdk.NewAttribute(types.AttributeKeyMediaNodeId, fmt.Sprint(mediaNodeId)),
		),
	})
}

func (k *Keeper) createLeasePaymentTransferEvent(ctx sdk.Context, sender, recipient sdk.AccAddress, amount sdk.Coin) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeLeasePaymentTransfer,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, recipient.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})
}

func (k *Keeper) createLeaseCommissionTransferEvent(ctx sdk.Context, sender, recipient sdk.AccAddress, amount sdk.Coin) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeLeaseCommissionTransfer,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, recipient.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})
}

func (k *Keeper) createMediaNodeDepositRefundEvent(ctx sdk.Context, sender, recipient sdk.AccAddress, amount sdk.Coin) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRefundDeposit,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, recipient.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})
}
