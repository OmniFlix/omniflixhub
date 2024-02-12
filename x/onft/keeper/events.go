package keeper

import (
	onfttypes "github.com/OmniFlix/omniflixhub/v3/x/onft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) emitCreateONFTDenomEvent(ctx sdk.Context, denomId, symbol, name, creator string) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			onfttypes.EventTypeCreateONFTDenom,
			sdk.NewAttribute(onfttypes.AttributeKeyDenomID, denomId),
			sdk.NewAttribute(onfttypes.AttributeKeySymbol, symbol),
			sdk.NewAttribute(onfttypes.AttributeKeyName, name),
			sdk.NewAttribute(onfttypes.AttributeKeyCreator, creator),
		),
	)
}

func (k Keeper) emitUpdateONFTDenomEvent(ctx sdk.Context, denomId, name, description, previewURI, creator string) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			onfttypes.EventTypeUpdateONFTDenom,
			sdk.NewAttribute(onfttypes.AttributeKeyDenomID, denomId),
			sdk.NewAttribute(onfttypes.AttributeKeyName, name),
			sdk.NewAttribute(onfttypes.AttributeKeyDescription, description),
			sdk.NewAttribute(onfttypes.AttributeKeyPreviewURI, previewURI),
			sdk.NewAttribute(onfttypes.AttributeKeyCreator, creator),
		),
	)
}

func (k Keeper) emitTransferONFTDenomEvent(ctx sdk.Context, denomId, symbol, sender, recipient string) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			onfttypes.EventTypeTransferONFTDenom,
			sdk.NewAttribute(onfttypes.AttributeKeyDenomID, denomId),
			sdk.NewAttribute(onfttypes.AttributeKeySymbol, symbol),
			sdk.NewAttribute(onfttypes.AttributeKeySender, sender),
			sdk.NewAttribute(onfttypes.AttributeKeyRecipient, recipient),
		),
	)
}

func (k Keeper) emitPurgeONFTDenomEvent(ctx sdk.Context, denomId string) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			onfttypes.EventTypePurgeONFTDenom,
			sdk.NewAttribute(onfttypes.AttributeKeyDenomID, denomId),
		),
	)
}

func (k Keeper) emitMintONFTEvent(ctx sdk.Context, nftId, denomId, uri, owner string) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			onfttypes.EventTypeMintONFT,
			sdk.NewAttribute(onfttypes.AttributeKeyNFTID, nftId),
			sdk.NewAttribute(onfttypes.AttributeKeyDenomID, denomId),
			sdk.NewAttribute(onfttypes.AttributeKeyMediaURI, uri),
			sdk.NewAttribute(onfttypes.AttributeKeyOwner, owner),
		),
	)
}

func (k Keeper) emitTransferONFTEvent(ctx sdk.Context, nftId, denomId, sender, recipient string) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			onfttypes.EventTypeTransferONFT,
			sdk.NewAttribute(onfttypes.AttributeKeyNFTID, nftId),
			sdk.NewAttribute(onfttypes.AttributeKeyDenomID, denomId),
			sdk.NewAttribute(onfttypes.AttributeKeySender, sender),
			sdk.NewAttribute(onfttypes.AttributeKeyRecipient, recipient),
		),
	)
}

func (k Keeper) emitBurnONFTEvent(ctx sdk.Context, nftId, denomId, owner string) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			onfttypes.EventTypeBurnONFT,
			sdk.NewAttribute(onfttypes.AttributeKeyNFTID, nftId),
			sdk.NewAttribute(onfttypes.AttributeKeyDenomID, denomId),
			sdk.NewAttribute(onfttypes.AttributeKeyOwner, owner),
		),
	)
}
