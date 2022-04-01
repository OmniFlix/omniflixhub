package keeper

import (
	"github.com/OmniFlix/omniflixhub/x/flixdrop/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
)

// SetActions set actions details
func (k Keeper) SetActions(ctx sdk.Context, actions []types.WeightedAction) error {
	for _, action := range actions {
		err := k.SetAction(ctx, action)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetAction sets a action  in store
func (k Keeper) SetAction(ctx sdk.Context, action types.WeightedAction) error {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, []byte{})

	bz, err := proto.Marshal(&action)
	if err != nil {
		return err
	}

	prefixStore.Set(types.KeyActionPrefix(action.Index), bz)
	return nil
}

// GetActions get actions for genesis export
func (k Keeper) GetActions(ctx sdk.Context) []types.WeightedAction {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.ActionsPrefix)

	iterator := prefixStore.Iterator(nil, nil)
	defer iterator.Close()

	actions := []types.WeightedAction{}
	for ; iterator.Valid(); iterator.Next() {

		action := types.WeightedAction{}

		err := proto.Unmarshal(iterator.Value(), &action)
		if err != nil {
			panic(err)
		}

		actions = append(actions, action)
	}
	return actions
}

// GetAction returns the action details
func (k Keeper) GetAction(ctx sdk.Context, index uint64) (types.WeightedAction, error) {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, []byte{})
	if !prefixStore.Has(types.KeyActionPrefix(index)) {
		return types.WeightedAction{}, nil
	}
	bz := prefixStore.Get(types.KeyActionPrefix(index))

	action := types.WeightedAction{}
	err := proto.Unmarshal(bz, &action)
	if err != nil {
		return types.WeightedAction{}, err
	}

	return action, nil
}

func (k Keeper) IsActionClaimStarted(ctx sdk.Context, index uint64) bool {
	action, err := k.GetAction(ctx, index)
	if err != nil {
		return false
	}
	if action.StartTime.IsZero() {
		return false
	}
	if ctx.BlockTime().Before(action.StartTime) {
		return false
	}
	return true
}
