package keeper

import (
	"fmt"

	"github.com/OmniFlix/omniflixhub/v2/x/marketplace/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/cosmos/gogoproto/types"
)

// GetNextAuctionNumber get the next auction number
func (k Keeper) GetNextAuctionNumber(ctx sdk.Context) uint64 {
	var nextAuctionNumber uint64
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.PrefixNextAuctionNumber)
	if bz == nil {
		panic(fmt.Errorf("%s module not initialized -- Should have been done in InitGenesis", types.ModuleName))
	} else {
		val := gogotypes.UInt64Value{}

		err := k.cdc.Unmarshal(bz, &val)
		if err != nil {
			panic(err)
		}

		nextAuctionNumber = val.GetValue()
	}
	return nextAuctionNumber
}

// SetNextAuctionNumber set the next auction number
func (k Keeper) SetNextAuctionNumber(ctx sdk.Context, number uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: number})
	store.Set(types.PrefixNextAuctionNumber, bz)
}

// SetAuctionListing set a specific auction listing in the store
func (k Keeper) SetAuctionListing(ctx sdk.Context, auctionListing types.AuctionListing) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixAuctionId)
	bz := k.cdc.MustMarshal(&auctionListing)
	store.Set(types.KeyAuctionIdPrefix(auctionListing.Id), bz)
}

// GetAuctionListing returns a auction listing by its id
func (k Keeper) GetAuctionListing(ctx sdk.Context, id uint64) (val types.AuctionListing, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixAuctionId)
	b := store.Get(types.KeyAuctionIdPrefix(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetListing returns a listing from its nft id
func (k Keeper) GetAuctionListingIdByNftId(ctx sdk.Context, nftId string) (val uint64, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixAuctionNFTID)
	bz := store.Get(types.KeyAuctionNFTIDPrefix(nftId))
	if bz == nil {
		return val, false
	}
	var auctionId gogotypes.UInt64Value
	k.cdc.MustUnmarshal(bz, &auctionId)
	return auctionId.Value, true
}

// RemoveAuctionListing removes a auction listing from the store
func (k Keeper) RemoveAuctionListing(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixAuctionId)
	store.Delete(types.KeyAuctionIdPrefix(id))
}

// GetAllAuctionListings returns all auction listings
func (k Keeper) GetAllAuctionListings(ctx sdk.Context) (list []types.AuctionListing) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixAuctionId)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.AuctionListing
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetAuctionListingsByOwner returns all auction listings of specific owner
func (k Keeper) GetAuctionListingsByOwner(ctx sdk.Context, owner sdk.AccAddress) (auctionListings []types.AuctionListing) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, append(types.PrefixAuctionOwner, owner.Bytes()...))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var id gogotypes.UInt64Value
		k.cdc.MustUnmarshal(iterator.Value(), &id)
		listing, found := k.GetAuctionListing(ctx, id.Value)
		if !found {
			continue
		}
		auctionListings = append(auctionListings, listing)
	}

	return
}

func (k Keeper) HasAuctionListing(ctx sdk.Context, id uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.KeyAuctionIdPrefix(id))
}

func (k Keeper) SetAuctionListingWithOwner(ctx sdk.Context, owner sdk.AccAddress, id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: id})

	store.Set(types.KeyAuctionOwnerPrefix(owner, id), bz)
}

func (k Keeper) UnsetAuctionListingWithOwner(ctx sdk.Context, owner sdk.AccAddress, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyAuctionOwnerPrefix(owner, id))
}

func (k Keeper) SetAuctionListingWithNFTID(ctx sdk.Context, nftId string, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixAuctionNFTID)
	bz := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: id})
	store.Set(types.KeyAuctionNFTIDPrefix(nftId), bz)
}

func (k Keeper) UnsetAuctionListingWithNFTID(ctx sdk.Context, nftId string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixAuctionNFTID)
	store.Delete(types.KeyAuctionNFTIDPrefix(nftId))
}

func (k Keeper) SetAuctionListingWithPriceDenom(ctx sdk.Context, priceDenom string, id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: id})

	store.Set(types.KeyAuctionPriceDenomPrefix(priceDenom, id), bz)
}

func (k Keeper) UnsetAuctionListingWithPriceDenom(ctx sdk.Context, priceDenom string, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyAuctionPriceDenomPrefix(priceDenom, id))
}

func (k Keeper) SetInactiveAuction(ctx sdk.Context, auctionId uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: auctionId})

	store.Set(types.KeyInActiveAuctionPrefix(auctionId), bz)
}

func (k Keeper) UnsetInactiveAuction(ctx sdk.Context, auctionId uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyInActiveAuctionPrefix(auctionId))
}

func (k Keeper) SetActiveAuction(ctx sdk.Context, auctionId uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: auctionId})

	store.Set(types.KeyActiveAuctionPrefix(auctionId), bz)
}

func (k Keeper) UnsetActiveAuction(ctx sdk.Context, auctionId uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyActiveAuctionPrefix(auctionId))
}

func (k Keeper) IterateInactiveAuctions(ctx sdk.Context, fn func(index int, item types.AuctionListing) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixInactiveAuction)
	iter := sdk.KVStorePrefixIterator(store, []byte{})
	defer iter.Close()

	for i := 0; iter.Valid(); iter.Next() {
		var id gogotypes.UInt64Value
		k.cdc.MustUnmarshal(iter.Value(), &id)
		auction, _ := k.GetAuctionListing(ctx, id.Value)

		if stop := fn(i, auction); stop {
			break
		}
		i++
	}
}

func (k Keeper) IterateActiveAuctions(ctx sdk.Context, fn func(index int, item types.AuctionListing) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixActiveAuction)
	iter := sdk.KVStorePrefixIterator(store, []byte{})
	defer iter.Close()

	for i := 0; iter.Valid(); iter.Next() {
		var id gogotypes.UInt64Value
		k.cdc.MustUnmarshal(iter.Value(), &id)
		auction, _ := k.GetAuctionListing(ctx, id.Value)

		if stop := fn(i, auction); stop {
			break
		}
		i++
	}
}

// UpdateAuctionStatusesAndProcessBids update all auction listings status
func (k Keeper) UpdateAuctionStatusesAndProcessBids(ctx sdk.Context) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixAuctionId)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var auction types.AuctionListing
		k.cdc.MustUnmarshal(iterator.Value(), &auction)
		if auction.StartTime.Before(ctx.BlockTime()) {
			bid, found := k.GetBid(ctx, auction.GetId())
			if !found && auction.EndTime != nil && auction.EndTime.Before(ctx.BlockTime()) {
				err := k.nftKeeper.TransferOwnership(ctx, auction.GetDenomId(), auction.GetNftId(),
					k.accountKeeper.GetModuleAddress(types.ModuleName), auction.GetOwner())
				if err != nil {
					return err
				}
				k.RemoveAuctionListing(ctx, auction.GetId())
				k.removeAuctionEvent(ctx, auction)
			} else if !found && auction.EndTime == nil &&
				ctx.BlockTime().Sub(auction.StartTime).Seconds() > k.GetBidCloseDuration(ctx).Seconds() {
				err := k.nftKeeper.TransferOwnership(ctx, auction.GetDenomId(), auction.GetNftId(),
					k.accountKeeper.GetModuleAddress(types.ModuleName), auction.GetOwner())
				if err != nil {
					return err
				}
				k.RemoveAuctionListing(ctx, auction.GetId())
				k.removeAuctionEvent(ctx, auction)

			} else if found && ctx.BlockTime().Sub(bid.Time).Seconds() > k.GetBidCloseDuration(ctx).Seconds() {
				err := k.processBid(ctx, auction, bid)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (k Keeper) processBid(ctx sdk.Context, auction types.AuctionListing, bid types.Bid) error {
	owner, err := sdk.AccAddressFromBech32(auction.Owner)
	if err != nil {
		return err
	}
	denom, err := k.nftKeeper.GetDenom(ctx, auction.DenomId)
	if err != nil {
		return err
	}
	nft, err := k.nftKeeper.GetONFT(ctx, auction.DenomId, auction.NftId)
	if err != nil {
		return err
	}
	BidAmountCoin := bid.Amount
	auctionSaleAmountCoin := BidAmountCoin
	err = k.nftKeeper.TransferOwnership(ctx, auction.GetDenomId(), auction.GetNftId(),
		k.accountKeeper.GetModuleAddress(types.ModuleName), bid.GetBidder())
	if err != nil {
		return err
	}
	saleCommission := k.GetSaleCommission(ctx)
	marketplaceCoin := k.GetProportions(bid.Amount, saleCommission)
	if marketplaceCoin.Amount.GTE(sdk.OneInt()) {
		err = k.DistributeCommission(ctx, marketplaceCoin)
		if err != nil {
			return err
		}
		auctionSaleAmountCoin = BidAmountCoin.Sub(marketplaceCoin)
	}
	if nft.GetRoyaltyShare().GT(sdk.ZeroDec()) {
		nftRoyaltyShareCoin := k.GetProportions(auctionSaleAmountCoin, nft.GetRoyaltyShare())
		creator, err := sdk.AccAddressFromBech32(denom.Creator)
		if err != nil {
			return err
		}
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creator, sdk.NewCoins(nftRoyaltyShareCoin))
		if err != nil {
			return err
		}
		k.createRoyaltyShareTransferEvent(ctx, k.accountKeeper.GetModuleAddress(types.ModuleName), creator, nftRoyaltyShareCoin)
		auctionSaleAmountCoin = auctionSaleAmountCoin.Sub(nftRoyaltyShareCoin)
	}
	remaining := auctionSaleAmountCoin

	if len(auction.SplitShares) > 0 {
		for _, share := range auction.SplitShares {
			sharePortionCoin := k.GetProportions(auctionSaleAmountCoin, share.Weight)
			sharePortionCoins := sdk.NewCoins(sharePortionCoin)
			if share.Address == "" {
				err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, owner, sharePortionCoins)
				if err != nil {
					return err
				}
			} else {
				saleSplitAddr, err := sdk.AccAddressFromBech32(share.Address)
				if err != nil {
					return err
				}
				err = k.bankKeeper.SendCoinsFromModuleToAccount(
					ctx, types.ModuleName, saleSplitAddr, sharePortionCoins)
				if err != nil {
					return err
				}
				k.createSplitShareTransferEvent(ctx, k.accountKeeper.GetModuleAddress(types.ModuleName), saleSplitAddr, sharePortionCoin)
			}
			remaining = remaining.Sub(sharePortionCoin)
		}
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, owner, sdk.NewCoins(remaining))
		if err != nil {
			return err
		}
	} else {
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, owner, sdk.NewCoins(remaining))
		if err != nil {
			return err
		}
	}
	k.processBidEvent(ctx, auction, bid)
	k.RemoveAuctionListing(ctx, auction.GetId())
	k.RemoveBid(ctx, auction.GetId())
	return nil
}
