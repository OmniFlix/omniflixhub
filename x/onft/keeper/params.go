package keeper

import (
	"github.com/OmniFlix/omniflixhub/v2/x/onft/types"
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

// GetDenomCreationFee returns the current adenom creation fee coins list and amounts.
func (k Keeper) GetDenomCreationFee(ctx sdk.Context) (feeCoin sdk.Coin) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyDenomCreationFee, &feeCoin)
	return feeCoin
}
