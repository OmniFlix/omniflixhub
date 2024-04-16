package keeper

import (
	"github.com/OmniFlix/omniflixhub/v4/x/marketplace/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/cosmos/gogoproto/types"
)

// SetBid set a specific bid for an auction listing in the store
func (k Keeper) SetBid(ctx sdk.Context, bid types.Bid) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixBidByAuctionId)
	bz := k.cdc.MustMarshal(&bid)
	store.Set(types.KeyBidPrefix(bid.AuctionId), bz)
}

// GetBid returns a bid of an auction listing by its id
func (k Keeper) GetBid(ctx sdk.Context, id uint64) (val types.Bid, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixBidByAuctionId)
	b := store.Get(types.KeyBidPrefix(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveBid removes a bid of an auction listing from the store
func (k Keeper) RemoveBid(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixBidByAuctionId)
	store.Delete(types.KeyBidPrefix(id))
}

// GetAllBids returns all bids
func (k Keeper) GetAllBids(ctx sdk.Context) (list []types.Bid) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixBidByAuctionId)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Bid
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetBidsByBidder returns all bids of specific bidder
func (k Keeper) GetBidsByBidder(ctx sdk.Context, bidder sdk.AccAddress) (bids []types.Bid) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, append(types.PrefixBidByBidder, bidder.Bytes()...))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var id gogotypes.UInt64Value
		k.cdc.MustUnmarshal(iterator.Value(), &id)
		bid, found := k.GetBid(ctx, id.Value)
		if !found {
			continue
		}
		bids = append(bids, bid)
	}

	return
}

func (k Keeper) HasBid(ctx sdk.Context, id uint64) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixBidByAuctionId)
	return store.Has(types.KeyBidPrefix(id))
}
