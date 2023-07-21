package keeper

import (
	"time"

	"github.com/OmniFlix/omniflixhub/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams gets the parameters for the marketplace module.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the parameters for the marketplace module.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// GetSaleCommission returns the current sale commission of marketplace.
func (k Keeper) GetSaleCommission(ctx sdk.Context) (percent sdk.Dec) {
	k.paramSpace.Get(ctx, types.ParamStoreKeySaleCommission, &percent)
	return percent
}

// GetMarketplaceDistributionParams returns the current distribution  of marketplace commission.
func (k Keeper) GetMarketplaceDistributionParams(ctx sdk.Context) (distParams types.Distribution) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyDistribution, &distParams)
	return distParams
}

// GetBidCloseDuration returns the closing duration for bid for auctions.
func (k Keeper) GetBidCloseDuration(ctx sdk.Context) (duration time.Duration) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyBidCloseDuration, &duration)
	return duration
}

// GetMaxAuctionDuration returns the maximum duration for auctions.
func (k Keeper) GetMaxAuctionDuration(ctx sdk.Context) (duration time.Duration) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyMaxAuctionDuration, &duration)
	return duration
}
