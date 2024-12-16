package keeper

import (
	"time"

	sdkmath "cosmossdk.io/math"

	"github.com/OmniFlix/omniflixhub/v5/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams gets the parameters for the marketplace module.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}
	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams sets the parameters for the marketplace module.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.ValidateBasic(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)
	return nil
}

// GetSaleCommission returns the current sale commission of marketplace.
func (k Keeper) GetSaleCommission(ctx sdk.Context) (percent sdkmath.LegacyDec) {
	params := k.GetParams(ctx)
	return params.SaleCommission
}

// GetMarketplaceDistributionParams returns the current distribution  of marketplace commission.
func (k Keeper) GetMarketplaceDistributionParams(ctx sdk.Context) (distParams types.Distribution) {
	params := k.GetParams(ctx)
	return params.Distribution
}

// GetBidCloseDuration returns the closing duration for bid for auctions.
func (k Keeper) GetBidCloseDuration(ctx sdk.Context) (duration time.Duration) {
	params := k.GetParams(ctx)
	return params.BidCloseDuration
}

// GetMaxAuctionDuration returns the maximum duration for auctions.
func (k Keeper) GetMaxAuctionDuration(ctx sdk.Context) (duration time.Duration) {
	params := k.GetParams(ctx)
	return params.MaxAuctionDuration
}

// GetBidExtensionWindow returns the bid extension window duration for auctions.
func (k Keeper) GetBidExtensionWindow(ctx sdk.Context) (duration time.Duration) {
	params := k.GetParams(ctx)
	return params.BidExtensionWindow
}

// GetBidExtensionDuration returns the bid extension duration for the auctions.
func (k Keeper) GetBidExtensionDuration(ctx sdk.Context) (duration time.Duration) {
	params := k.GetParams(ctx)
	return params.BidExtensionDuration
}
