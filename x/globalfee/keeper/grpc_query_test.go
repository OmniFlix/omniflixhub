package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	globalfeekeeper "github.com/OmniFlix/omniflixhub/v5/x/globalfee/keeper"
	"github.com/OmniFlix/omniflixhub/v5/x/globalfee/types"
)

func TestQueryGlobalFeeParamMinGasPrices(t *testing.T) {
	specs := map[string]struct {
		setupStore func(ctx sdk.Context, k globalfeekeeper.Keeper)
		expMin     sdk.DecCoins
	}{
		"single fee coin": {
			setupStore: func(ctx sdk.Context, k globalfeekeeper.Keeper) {
				err := k.SetParams(ctx, types.Params{
					MinimumGasPrices: sdk.NewDecCoins(sdk.NewDecCoin("uflix", sdkmath.OneInt())),
				})
				require.NoError(t, err)
			},
			expMin: sdk.NewDecCoins(sdk.NewDecCoin("uflix", sdkmath.OneInt())),
		},
		"multiple fee coins": {
			setupStore: func(ctx sdk.Context, k globalfeekeeper.Keeper) {
				err := k.SetParams(ctx, types.Params{
					MinimumGasPrices: sdk.NewDecCoins(sdk.NewDecCoin("uflix", sdkmath.OneInt()), sdk.NewDecCoin("test", sdkmath.NewInt(2))),
				})
				require.NoError(t, err)
			},
			expMin: sdk.NewDecCoins(sdk.NewDecCoin("uflix", sdkmath.OneInt()), sdk.NewDecCoin("test", sdkmath.NewInt(2))),
		},
		"no coins": {
			setupStore: func(ctx sdk.Context, k globalfeekeeper.Keeper) {
				err := k.SetParams(ctx, types.Params{})
				require.NoError(t, err)
			},
		},
		"no param set": {
			setupStore: func(ctx sdk.Context, k globalfeekeeper.Keeper) {
			},
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			ctx, _, keeper := setupTestStore(t)
			spec.setupStore(ctx, keeper)
			q := globalfeekeeper.NewGrpcQuerier(keeper)
			gotResp, gotErr := q.Params(sdk.WrapSDKContext(ctx), nil)
			require.NoError(t, gotErr)
			require.NotNil(t, gotResp)
			assert.Equal(t, spec.expMin, gotResp.Params.MinimumGasPrices)
		})
	}
}
