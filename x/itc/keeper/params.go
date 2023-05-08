package keeper

import (
	"time"

	"github.com/OmniFlix/omniflixhub/x/itc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams gets the parameters for the itc module.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the parameters for the itc module.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// GetMaxCampaignDuration returns the maximum duration allowed to create a campaign.
func (k Keeper) GetMaxCampaignDuration(ctx sdk.Context) (duration time.Duration) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyMaxCampaignDuration, &duration)
	return duration
}

// GetCampaignCreationFee returns the creation fee required to create a campaign.
func (k Keeper) GetCampaignCreationFee(ctx sdk.Context) (creationFee sdk.Coin) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyCampaignCreationFee, &creationFee)
	return creationFee
}
