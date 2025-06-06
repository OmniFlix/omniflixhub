package keeper

import (
	"fmt"

	storetypes "cosmossdk.io/store/types"

	"github.com/OmniFlix/omniflixhub/v6/x/medianode/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/cosmos/gogoproto/types"
)

// SetMediaNode stores the media node in state
func (k Keeper) SetMediaNode(ctx sdk.Context, mediaNode types.MediaNode) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&mediaNode)
	store.Set(types.GetMediaNodeKey(mediaNode.Id), bz)
}

// SetLease stores the lease in state
func (k Keeper) SetLease(ctx sdk.Context, lease types.Lease) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&lease)
	store.Set(types.GetLeaseKey(lease.MediaNodeId), bz)
}

// GetMediaNode returns media node by id
func (k Keeper) GetMediaNode(ctx sdk.Context, id string) (mediaNode types.MediaNode, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetMediaNodeKey(id))
	if bz == nil {
		return mediaNode, false
	}
	k.cdc.MustUnmarshal(bz, &mediaNode)
	return mediaNode, true
}

// GetMediaNode returns media node by id
func (k Keeper) GetMediaNodeLease(ctx sdk.Context, id string) (lease types.Lease, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetLeaseKey(id))
	if bz == nil {
		return lease, false
	}
	k.cdc.MustUnmarshal(bz, &lease)
	return lease, true
}

// SetMediaNodeCount stores the  media node count
func (k Keeper) SetMediaNodeCount(ctx sdk.Context, number uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: number})
	store.Set(types.PrefixMediaNodeCount, bz)
}

// GetMediaNodeCount returns the current media node count
func (k Keeper) GetMediaNodeCount(ctx sdk.Context) (nodesCount uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PrefixMediaNodeCount)
	if bz == nil {
		panic(fmt.Errorf("%s module not initialized", types.ModuleName))
	} else {
		val := gogotypes.UInt64Value{}

		err := k.cdc.Unmarshal(bz, &val)
		if err != nil {
			panic(err)
		}
		nodesCount = val.GetValue()
	}
	return nodesCount
}

// GetAllMediaNodes returns all media nodes from the store
func (k Keeper) GetAllMediaNodes(ctx sdk.Context) (mediaNodes []types.MediaNode) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.PrefixMediaNode)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var mediaNode types.MediaNode
		k.cdc.MustUnmarshal(iterator.Value(), &mediaNode)
		mediaNodes = append(mediaNodes, mediaNode)
	}
	return mediaNodes
}

// GetAllLeases returns all leases from the store
func (k Keeper) GetAllLeases(ctx sdk.Context) (leases []types.Lease) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.PrefixLease)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var lease types.Lease
		k.cdc.MustUnmarshal(iterator.Value(), &lease)
		leases = append(leases, lease)
	}
	return leases
}

// RemoveMediaNode removes the media node from state
func (k Keeper) RemoveMediaNode(ctx sdk.Context, id string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetMediaNodeKey(id))
}

// RemoveLease removes the lease from state
func (k Keeper) RemoveLease(ctx sdk.Context, mediaNodeId string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetLeaseKey(mediaNodeId))
}
