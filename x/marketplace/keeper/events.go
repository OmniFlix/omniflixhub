package keeper

import (
	"fmt"

	"github.com/OmniFlix/omniflixhub/v3/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) createListNftEvent(ctx sdk.Context, owner sdk.AccAddress, listingId, denomId, nftId string, price sdk.Coin) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeListNFT,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyListingId, listingId),
			sdk.NewAttribute(types.AttributeKeyDenomId, denomId),
			sdk.NewAttribute(types.AttributeKeyNftId, nftId),
			sdk.NewAttribute(types.AttributeKeyAmount, price.String()),
		),
	})
}

func (k *Keeper) createDeListNftEvent(ctx sdk.Context, sender sdk.AccAddress, listingId string) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDeListNFT,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
			sdk.NewAttribute(types.AttributeKeyListingId, listingId),
		),
	})
}

func (k *Keeper) createEditListingEvent(ctx sdk.Context, sender sdk.AccAddress, listingId string, price sdk.Coin) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEditListing,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
			sdk.NewAttribute(types.AttributeKeyListingId, listingId),
			sdk.NewAttribute(types.AttributeKeyAmount, price.String()),
		),
	})
}

func (k *Keeper) createBuyNftEvent(ctx sdk.Context, buyer sdk.AccAddress, listId, nftId string, price sdk.Coin) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBuyNFT,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyBuyer, buyer.String()),
			sdk.NewAttribute(types.AttributeKeyListingId, listId),
			sdk.NewAttribute(types.AttributeKeyNftId, nftId),
			sdk.NewAttribute(types.AttributeKeyAmount, price.String()),
		),
	})
}

func (k *Keeper) createRoyaltyShareTransferEvent(ctx sdk.Context, sender, recipient sdk.AccAddress, amount sdk.Coin) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRoyaltyShareTransfer,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, recipient.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})
}

func (k *Keeper) createSplitShareTransferEvent(ctx sdk.Context, sender, recipient sdk.AccAddress, amount sdk.Coin) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSplitShareTransfer,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, recipient.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})
}

func (k *Keeper) createSaleCommissionTransferEvent(ctx sdk.Context, sender, recipient sdk.AccAddress, amount sdk.Coin) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSaleCommissionTransfer,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, recipient.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})
}

func (k *Keeper) createAuctionEvent(ctx sdk.Context, auction types.AuctionListing) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateAuction,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyOwner, auction.GetOwner().String()),
			sdk.NewAttribute(types.AttributeKeyAuctionId, fmt.Sprint(auction.GetId())),
			sdk.NewAttribute(types.AttributeKeyDenomId, auction.GetDenomId()),
			sdk.NewAttribute(types.AttributeKeyNftId, auction.GetNftId()),
			sdk.NewAttribute(types.AttributeKeyStartPrice, auction.GetStartPrice().String()),
		),
	})
}

func (k *Keeper) cancelAuctionEvent(ctx sdk.Context, auction types.AuctionListing) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCancelAuction,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyOwner, auction.GetOwner().String()),
			sdk.NewAttribute(types.AttributeKeyAuctionId, fmt.Sprint(auction.GetId())),
			sdk.NewAttribute(types.AttributeKeyDenomId, auction.GetDenomId()),
			sdk.NewAttribute(types.AttributeKeyNftId, auction.GetNftId()),
		),
	})
}

func (k *Keeper) placeBidEvent(ctx sdk.Context, auction types.AuctionListing, bid types.Bid) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypePlaceBid,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyBidder, bid.GetBidder().String()),
			sdk.NewAttribute(types.AttributeKeyAuctionId, fmt.Sprint(auction.GetId())),
			sdk.NewAttribute(types.AttributeKeyDenomId, auction.GetDenomId()),
			sdk.NewAttribute(types.AttributeKeyNftId, auction.GetNftId()),
			sdk.NewAttribute(types.AttributeKeyAmount, bid.GetAmount().String()),
		),
	})
}

func (k *Keeper) removeAuctionEvent(ctx sdk.Context, auction types.AuctionListing) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRemoveAuction,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyAuctionId, fmt.Sprint(auction.GetId())),
			sdk.NewAttribute(types.AttributeKeyDenomId, auction.GetDenomId()),
			sdk.NewAttribute(types.AttributeKeyNftId, auction.GetNftId()),
		),
	})
}

func (k *Keeper) processBidEvent(ctx sdk.Context, auction types.AuctionListing, bid types.Bid) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeProcessBid,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyAuctionId, fmt.Sprint(auction.GetId())),
			sdk.NewAttribute(types.AttributeKeyDenomId, auction.GetDenomId()),
			sdk.NewAttribute(types.AttributeKeyNftId, auction.GetNftId()),
			sdk.NewAttribute(types.AttributeKeyBidder, bid.GetBidder().String()),
			sdk.NewAttribute(types.AttributeKeyAmount, bid.GetAmount().String()),
		),
	})
}
