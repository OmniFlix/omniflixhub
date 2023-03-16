package keeper

import (
	"fmt"

	"github.com/OmniFlix/omniflixhub/x/itc/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/gogo/protobuf/types"
)

// GetNextCampaignNumber get the next campaign number
func (k Keeper) GetNextCampaignNumber(ctx sdk.Context) uint64 {
	var nextCampaignNumber uint64
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.PrefixNextCampaignNumber)
	if bz == nil {
		panic(fmt.Errorf("%s module not initialized -- Should have been done in InitGenesis", types.ModuleName))
	} else {
		val := gogotypes.UInt64Value{}

		err := k.cdc.Unmarshal(bz, &val)
		if err != nil {
			panic(err)
		}

		nextCampaignNumber = val.GetValue()
	}
	return nextCampaignNumber
}

// SetNextCampaignNumber set the next campaign number
func (k Keeper) SetNextCampaignNumber(ctx sdk.Context, number uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: number})
	store.Set(types.PrefixNextCampaignNumber, bz)
}

// SetCampaign set a specific campaign in the store
func (k Keeper) SetCampaign(ctx sdk.Context, campaign types.Campaign) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixCampaignId)
	bz := k.cdc.MustMarshal(&campaign)
	store.Set(types.KeyCampaignIdPrefix(campaign.Id), bz)
}

// GetCampaign returns a campaign by its id
func (k Keeper) GetCampaign(ctx sdk.Context, id uint64) (val types.Campaign, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixCampaignId)
	b := store.Get(types.KeyCampaignIdPrefix(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveCampaign removes a campaign from the store
func (k Keeper) RemoveCampaign(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixCampaignId)
	store.Delete(types.KeyCampaignIdPrefix(id))
}

// GetAllCampaigns returns all campaigns
func (k Keeper) GetAllCampaigns(ctx sdk.Context) (list []types.Campaign) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixCampaignId)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Campaign
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetCampaignsByCreator returns all campaigns created by specific address
func (k Keeper) GetCampaignsByCreator(ctx sdk.Context, creator sdk.AccAddress) (campaigns []types.Campaign) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, append(types.PrefixCampaignCreator, creator.Bytes()...))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var id gogotypes.UInt64Value
		k.cdc.MustUnmarshal(iterator.Value(), &id)
		campaign, found := k.GetCampaign(ctx, id.Value)
		if !found {
			continue
		}
		campaigns = append(campaigns, campaign)
	}

	return
}

func (k Keeper) GetClaims(ctx sdk.Context, campaignId uint64) (claims []types.Claim) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyClaimPrefix(campaignId))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Claim
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		claims = append(claims, val)
	}
	return claims
}

func (k Keeper) HasCampaign(ctx sdk.Context, id uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.KeyCampaignIdPrefix(id))
}

func (k Keeper) SetCampaignWithCreator(ctx sdk.Context, owner sdk.AccAddress, id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: id})

	store.Set(types.KeyCampaignCreatorPrefix(owner, id), bz)
}

func (k Keeper) UnsetCampaignWithCreator(ctx sdk.Context, owner sdk.AccAddress, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyCampaignCreatorPrefix(owner, id))
}

func (k Keeper) SetInactiveCampaign(ctx sdk.Context, campaignId uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: campaignId})

	store.Set(types.KeyInActiveCampaignPrefix(campaignId), bz)
}

func (k Keeper) UnsetInactiveCampaign(ctx sdk.Context, campaignId uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyInActiveCampaignPrefix(campaignId))
}

func (k Keeper) SetActiveCampaign(ctx sdk.Context, campaignId uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: campaignId})

	store.Set(types.KeyActiveCampaignPrefix(campaignId), bz)
}

func (k Keeper) UnsetActiveCampaign(ctx sdk.Context, campaignId uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyActiveCampaignPrefix(campaignId))
}

func (k Keeper) IterateInactiveCampaigns(ctx sdk.Context, fn func(index int, item types.Campaign) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixInactiveCampaign)
	iter := sdk.KVStorePrefixIterator(store, []byte{})
	defer iter.Close()

	for i := 0; iter.Valid(); iter.Next() {
		var id gogotypes.UInt64Value
		k.cdc.MustUnmarshal(iter.Value(), &id)
		campaign, _ := k.GetCampaign(ctx, id.Value)

		if stop := fn(i, campaign); stop {
			break
		}
		i++
	}
}

func (k Keeper) IterateActiveCampaigns(ctx sdk.Context, fn func(index int, item types.Campaign) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixActiveCampaign)
	iter := sdk.KVStorePrefixIterator(store, []byte{})
	defer iter.Close()

	for i := 0; iter.Valid(); iter.Next() {
		var id gogotypes.UInt64Value
		k.cdc.MustUnmarshal(iter.Value(), &id)
		campaign, _ := k.GetCampaign(ctx, id.Value)

		if stop := fn(i, campaign); stop {
			break
		}
		i++
	}
}

// UpdateCampaignStatuses update all campaign statuses
func (k Keeper) UpdateCampaignStatuses(ctx sdk.Context) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixCampaignId)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var campaign types.Campaign
		k.cdc.MustUnmarshal(iterator.Value(), &campaign)
		// Active campaigns
		if campaign.StartTime.Before(ctx.BlockTime()) {
			// settle campaign after endtime
			if campaign.EndTime.Before(ctx.BlockTime()) {
				k.endCampaign(ctx, campaign)
			}
		}
	}
	return nil
}

// TODO: implement endCampaign sending remaining amounts or nfts to treasury / creator address
func (k Keeper) endCampaign(ctx sdk.Context, campaign types.Campaign) {
	k.RemoveCampaign(ctx, campaign.GetId())
}
