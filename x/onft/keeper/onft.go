package keeper

import (
	"github.com/OmniFlix/omniflixhub/v2/x/onft/exported"
	"github.com/OmniFlix/omniflixhub/v2/x/onft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) GetONFT(ctx sdk.Context, denomID, onftID string) (nft exported.ONFT, err error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyONFT(denomID, onftID))
	if bz == nil {
		return nil, sdkerrors.Wrapf(types.ErrUnknownCollection, "not found oNFT: %s", denomID)
	}

	var oNFT types.ONFT
	k.cdc.MustUnmarshal(bz, &oNFT)
	return oNFT, nil
}

func (k Keeper) GetONFTs(ctx sdk.Context, denom string) (onfts []exported.ONFT) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.KeyONFT(denom, ""))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var oNFT types.ONFT
		k.cdc.MustUnmarshal(iterator.Value(), &oNFT)
		onfts = append(onfts, oNFT)
	}
	return onfts
}

func (k Keeper) GetOwnerONFTs(ctx sdk.Context, denom string, owner string) (onfts []*types.ONFT) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.KeyONFT(denom, ""))
	defer iterator.Close()
	var onftList []*types.ONFT
	for ; iterator.Valid(); iterator.Next() {
		var oNFT types.ONFT
		k.cdc.MustUnmarshal(iterator.Value(), &oNFT)
		if oNFT.Owner == owner {
			onftList = append(onftList, &oNFT)
		}
	}
	return onftList
}

func (k Keeper) Authorize(ctx sdk.Context, denomID, onftID string, owner sdk.AccAddress) (types.ONFT, error) {
	onft, err := k.GetONFT(ctx, denomID, onftID)
	if err != nil {
		return types.ONFT{}, err
	}

	if !owner.Equals(onft.GetOwner()) {
		return types.ONFT{}, sdkerrors.Wrap(types.ErrUnauthorized, owner.String())
	}
	return onft.(types.ONFT), nil
}

func (k Keeper) HasONFT(ctx sdk.Context, denomID, onftID string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.KeyONFT(denomID, onftID))
}

func (k Keeper) setONFT(ctx sdk.Context, denomID string, onft types.ONFT) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&onft)
	store.Set(types.KeyONFT(denomID, onft.GetID()), bz)
}

func (k Keeper) deleteONFT(ctx sdk.Context, denomID string, onft exported.ONFT) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyONFT(denomID, onft.GetID()))
}
