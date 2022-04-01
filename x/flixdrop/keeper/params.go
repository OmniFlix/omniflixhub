package keeper

import (
	"github.com/OmniFlix/omniflixhub/x/flixdrop/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams returns the total set of claim parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// SetParams sets claim parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
