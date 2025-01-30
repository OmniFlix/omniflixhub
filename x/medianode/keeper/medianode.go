package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/OmniFlix/omniflixhub/v6/x/medianode/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gogotypes "github.com/cosmos/gogoproto/types"
)

// CancelLease cancels an existing lease for a media node
func (k Keeper) CancelLease(ctx sdk.Context, mediaNodeId uint64, sender sdk.AccAddress) error {
	mediaNode, found := k.GetMediaNode(ctx, mediaNodeId)
	if !found {
		return errorsmod.Wrapf(types.ErrMediaNodeDoesNotExist, "media node %d does not exist", mediaNodeId)
	}
	lease, found := k.GetMediaNodeLease(ctx, mediaNodeId)
	if !found {
		return errorsmod.Wrapf(types.ErrLeaseNotFound, "lease for media node %d does not exist", mediaNodeId)
	}

	if sender.String() != lease.LeasedTo {
		return errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "unauthorized address %s", sender.String())
	}

	// Calculate remaining lease days and refund amount
	remainingDays := uint64(ctx.BlockTime().Sub(lease.LeasedAt).Hours() / 24)
	if remainingDays > 0 {
		refundAmount := sdk.NewCoin(
			mediaNode.LeaseAmountPerDay.Denom,
			mediaNode.LeaseAmountPerDay.Amount.MulRaw(int64(remainingDays)),
		)

		// Return remaining funds to lessee
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(
			ctx,
			types.ModuleName,
			sender,
			sdk.NewCoins(refundAmount),
		); err != nil {
			return err
		}
	}

	// Clear lease
	mediaNode.Leased = false
	lease.LeaseStatus = types.LeaseStatus_LEASE_STATUS_CANCELLED
	k.SetMediaNode(ctx, mediaNode)
	k.SetLease(ctx, lease)
	return nil
}

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
func (k Keeper) GetMediaNode(ctx sdk.Context, id uint64) (mediaNode types.MediaNode, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetMediaNodeKey(id))
	if bz == nil {
		return mediaNode, false
	}
	k.cdc.MustUnmarshal(bz, &mediaNode)
	return mediaNode, true
}

// GetMediaNode returns media node by id
func (k Keeper) GetMediaNodeLease(ctx sdk.Context, id uint64) (lease types.Lease, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetLeaseKey(id))
	if bz == nil {
		return lease, false
	}
	k.cdc.MustUnmarshal(bz, &lease)
	return lease, true
}

// SetNextMediaNodeNumber stores the next media node ID to be assigned
func (k Keeper) SetNextMediaNodeNumber(ctx sdk.Context, number uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: number})
	store.Set(types.PrefixNextNodeId, bz)
}

// GetNextMediaNodeNumber returns the next media node ID to be assigned
func (k Keeper) GetNextMediaNodeNumber(ctx sdk.Context) (nextNodeId uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PrefixNextNodeId)
	if bz == nil {
		panic(fmt.Errorf("%s module not initialized -- Should have been done in InitGenesis", types.ModuleName))
	} else {
		val := gogotypes.UInt64Value{}

		err := k.cdc.Unmarshal(bz, &val)
		if err != nil {
			panic(err)
		}
		nextNodeId = val.GetValue()
	}
	return nextNodeId
}
