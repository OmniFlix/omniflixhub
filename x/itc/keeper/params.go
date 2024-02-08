package keeper

import (
	"time"

	"github.com/OmniFlix/omniflixhub/v3/x/itc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams gets the parameters for the itc module.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}
	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams sets the parameters for the itc module.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.ValidateBasic(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)
	return nil
}

// GetMaxCampaignDuration returns the maximum duration allowed to create a campaign.
func (k Keeper) GetMaxCampaignDuration(ctx sdk.Context) (duration time.Duration) {
	params := k.GetParams(ctx)
	return params.MaxCampaignDuration
}

// GetCampaignCreationFee returns the creation fee required to create a campaign.
func (k Keeper) GetCampaignCreationFee(ctx sdk.Context) (creationFee sdk.Coin) {
	params := k.GetParams(ctx)
	return params.CreationFee
}
