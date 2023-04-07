package keeper

import (
	"encoding/binary"

	"github.com/OmniFlix/omniflixhub/v2/x/marketplace/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/cosmos/gogoproto/types"
)

// GetListingCount get the total number of listings
func (k Keeper) GetListingCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.PrefixListingsCount
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// SetListingCount set the total number of listings
func (k Keeper) SetListingCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.PrefixListingsCount
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// SetListing set a specific listing in the store
func (k Keeper) SetListing(ctx sdk.Context, listing types.Listing) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixListingId)
	b := k.cdc.MustMarshal(&listing)
	store.Set(types.KeyListingIdPrefix(listing.Id), b)
}

// GetListing returns a listing from its id
func (k Keeper) GetListing(ctx sdk.Context, id string) (val types.Listing, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixListingId)
	b := store.Get(types.KeyListingIdPrefix(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetListing returns a listing from its nft id
func (k Keeper) GetListingIdByNftId(ctx sdk.Context, nftId string) (val string, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixListingNFTID)
	b := store.Get(types.KeyListingNFTIDPrefix(nftId))
	if b == nil {
		return val, false
	}
	var listingId gogotypes.StringValue
	k.cdc.MustUnmarshal(b, &listingId)
	return listingId.Value, true
}

// RemoveListing removes a listing from the store
func (k Keeper) RemoveListing(ctx sdk.Context, id string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixListingId)
	store.Delete(types.KeyListingIdPrefix(id))
}

// GetAllListings returns all listings
func (k Keeper) GetAllListings(ctx sdk.Context) (list []types.Listing) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixListingId)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Listing
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetListingsByOwner returns all listings of specific owner
func (k Keeper) GetListingsByOwner(ctx sdk.Context, owner sdk.AccAddress) (listings []types.Listing) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, append(types.PrefixListingOwner, owner.Bytes()...))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var id gogotypes.StringValue
		k.cdc.MustUnmarshal(iterator.Value(), &id)
		listing, found := k.GetListing(ctx, id.Value)
		if !found {
			continue
		}
		listings = append(listings, listing)
	}

	return
}

func (k Keeper) HasListing(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.KeyListingIdPrefix(id))
}

func (k Keeper) SetWithOwner(ctx sdk.Context, owner sdk.AccAddress, id string) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.StringValue{Value: id})

	store.Set(types.KeyListingOwnerPrefix(owner, id), bz)
}

func (k Keeper) UnsetWithOwner(ctx sdk.Context, owner sdk.AccAddress, id string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyListingOwnerPrefix(owner, id))
}

func (k Keeper) SetWithNFTID(ctx sdk.Context, nftId, id string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixListingNFTID)
	bz := k.cdc.MustMarshal(&gogotypes.StringValue{Value: id})
	store.Set(types.KeyListingNFTIDPrefix(nftId), bz)
}

func (k Keeper) UnsetWithNFTID(ctx sdk.Context, nftId string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixListingNFTID)
	store.Delete(types.KeyListingNFTIDPrefix(nftId))
}

func (k Keeper) SetWithPriceDenom(ctx sdk.Context, priceDenom, id string) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.StringValue{Value: id})

	store.Set(types.KeyListingPriceDenomPrefix(priceDenom, id), bz)
}

func (k Keeper) UnsetWithPriceDenom(ctx sdk.Context, priceDenom, id string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyListingPriceDenomPrefix(priceDenom, id))
}
