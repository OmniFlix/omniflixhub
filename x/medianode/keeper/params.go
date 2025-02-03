package keeper

import (
	sdkmath "cosmossdk.io/math"
	"github.com/OmniFlix/omniflixhub/v6/x/medianode/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams gets the parameters for the medianode module.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}
	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams sets the parameters for the medianode module.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)
	return nil
}

// GetMinimumLeaseDays returns the minimum lease days parameter.
func (k Keeper) GetMinimumLeaseDays(ctx sdk.Context) (minimumLeaseDays uint64) {
	params := k.GetParams(ctx)
	return params.MinimumLeaseDays
}

// GetMaximumLeaseDays returns the maximum lease days parameter.
func (k Keeper) GetMaximumLeaseDays(ctx sdk.Context) (maximumLeaseDays uint64) {
	params := k.GetParams(ctx)
	return params.MaximumLeaseDays
}

// GetMinDeposit returns the minimum deposit parameter.
func (k Keeper) GetMinDeposit(ctx sdk.Context) (minDeposit sdk.Coin) {
	params := k.GetParams(ctx)
	return params.MinDeposit
}

// GetInitialDepositPercentage returns the initial deposit percentage parameter.
func (k Keeper) GetInitialDepositPercentage(ctx sdk.Context) (initialDepositPercentage sdkmath.LegacyDec) {
	params := k.GetParams(ctx)
	return params.InitialDepositPercentage
}

// GetLeaseCommission returns the lease commission parameter.
func (k Keeper) GetLeaseCommission(ctx sdk.Context) (leaseCommission sdkmath.LegacyDec) {
	params := k.GetParams(ctx)
	return params.LeaseCommission
}

// GetCommissionDistribution returns the commission distribution parameter.
func (k Keeper) GetCommissionDistribution(ctx sdk.Context) (commissionDistribution types.Distribution) {
	params := k.GetParams(ctx)
	return params.CommissionDistribution
}
