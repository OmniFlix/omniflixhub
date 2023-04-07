package keeper

import (
	"github.com/OmniFlix/omniflixhub/v2/x/onft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetOwner gets all the ID collections owned by an address and denom ID
func (k Keeper) GetOwner(ctx sdk.Context, address sdk.AccAddress, denomId string) types.Owner {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyOwner(address, denomId, ""))
	defer iterator.Close()

	owner := types.Owner{
		Address:       address.String(),
		IDCollections: types.IDCollections{},
	}
	idsMap := make(map[string][]string)

	for ; iterator.Valid(); iterator.Next() {
		_, denomID, onftId, _ := types.SplitKeyOwner(iterator.Key())
		if ids, ok := idsMap[denomID]; ok {
			idsMap[denomID] = append(ids, onftId)
		} else {
			idsMap[denomID] = []string{onftId}
			owner.IDCollections = append(
				owner.IDCollections,
				types.IDCollection{DenomId: denomID},
			)
		}
	}

	for i := 0; i < len(owner.IDCollections); i++ {
		owner.IDCollections[i].OnftIds = idsMap[owner.IDCollections[i].DenomId]
	}

	return owner
}

// GetOwners gets all the ID collections
func (k Keeper) GetOwners(ctx sdk.Context) (owners types.Owners) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStoreReversePrefixIterator(store, types.KeyOwner(nil, "", ""))
	defer iterator.Close()

	idcsMap := make(map[string]types.IDCollections)
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		address, denom, id, _ := types.SplitKeyOwner(key)
		if _, ok := idcsMap[address.String()]; !ok {
			idcsMap[address.String()] = types.IDCollections{}
			owners = append(
				owners,
				types.Owner{Address: address.String()},
			)
		}
		idcs := idcsMap[address.String()]
		idcs = idcs.Add(denom, id)
		idcsMap[address.String()] = idcs
	}
	for i, owner := range owners {
		owners[i].IDCollections = idcsMap[owner.Address]
	}

	return owners
}

func (k Keeper) deleteOwner(ctx sdk.Context, denomID, onftId string, owner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyOwner(owner, denomID, onftId))
}

func (k Keeper) setOwner(ctx sdk.Context,
	denomId, onftId string,
	owner sdk.AccAddress,
) {
	store := ctx.KVStore(k.storeKey)

	bz := types.MustMarshalONFTID(k.cdc, onftId)
	store.Set(types.KeyOwner(owner, denomId, onftId), bz)
}

func (k Keeper) swapOwner(ctx sdk.Context, denomID, tokenID string, srcOwner, dstOwner sdk.AccAddress) {
	k.deleteOwner(ctx, denomID, tokenID, srcOwner)
	k.setOwner(ctx, denomID, tokenID, dstOwner)
}
