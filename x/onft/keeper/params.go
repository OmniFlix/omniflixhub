package keeper

import (
	"github.com/OmniFlix/omniflixhub/v5/x/onft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams gets the parameters for the onft module.
func (k Keeper) GetParams(ctx context.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}
	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams sets the parameters for the onft module.
func (k Keeper) SetParams(ctx context.Context, params types.Params) error {
	if err := params.ValidateBasic(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)
	return nil
}

// GetDenomCreationFee returns the current denom creation fee coins list and amounts.
func (k Keeper) GetDenomCreationFee(ctx context.Context) sdk.Coin {
	params := k.GetParams(ctx)
	return params.DenomCreationFee
}
