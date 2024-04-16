package globalfee

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	appparams "github.com/OmniFlix/omniflixhub/v4/app/params"
	globalfeekeeper "github.com/OmniFlix/omniflixhub/v4/x/globalfee/keeper"
	"github.com/OmniFlix/omniflixhub/v4/x/globalfee/types"
)

func TestDefaultGenesis(t *testing.T) {
	encCfg := appparams.MakeEncodingConfig()
	gotJSON := AppModuleBasic{}.DefaultGenesis(encCfg.Marshaler)
	assert.JSONEq(t, `{"params":{"minimum_gas_prices":[],"bypass_min_fee_msg_types":[],"max_total_bypass_min_fee_msg_gas_usage":"2000000"}}`, string(gotJSON), string(gotJSON))
}

func TestValidateGenesis(t *testing.T) {
	encCfg := appparams.MakeEncodingConfig()
	specs := map[string]struct {
		src    string
		expErr bool
	}{
		"all good": {
			src: `{"params":{"minimum_gas_prices":[{"denom":"ALX", "amount":"1"}]}}`,
		},
		"empty minimum": {
			src: `{"params":{"minimum_gas_prices":[]}}`,
		},
		"minimum not set": {
			src: `{"params":{}}`,
		},
		"zero amount allowed": {
			src:    `{"params":{"minimum_gas_prices":[{"denom":"ALX", "amount":"0"}]}}`,
			expErr: false,
		},
		"duplicate denoms not allowed": {
			src:    `{"params":{"minimum_gas_prices":[{"denom":"ALX", "amount":"1"},{"denom":"ALX", "amount":"2"}]}}`,
			expErr: true,
		},
		"negative amounts not allowed": {
			src:    `{"params":{"minimum_gas_prices":[{"denom":"ALX", "amount":"-1"}]}}`,
			expErr: true,
		},
		"denom must be sorted": {
			src:    `{"params":{"minimum_gas_prices":[{"denom":"ZLX", "amount":"1"},{"denom":"ALX", "amount":"2"}]}}`,
			expErr: true,
		},
		"sorted denoms is allowed": {
			src:    `{"params":{"minimum_gas_prices":[{"denom":"ALX", "amount":"1"},{"denom":"ZLX", "amount":"2"}]}}`,
			expErr: false,
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			gotErr := AppModuleBasic{}.ValidateGenesis(encCfg.Marshaler, nil, []byte(spec.src))
			if spec.expErr {
				require.Error(t, gotErr)
				return
			}
			require.NoError(t, gotErr)
		})
	}
}

func TestInitExportGenesis(t *testing.T) {
	specs := map[string]struct {
		src string
		exp types.GenesisState
	}{
		"single fee": {
			src: `{"params":{"minimum_gas_prices":[{"denom":"ALX", "amount":"1"}],"bypass_min_fee_msg_types":[],"max_total_bypass_min_fee_msg_gas_usage":"0"}}`,
			exp: types.GenesisState{
				Params: types.Params{
					MinimumGasPrices:                sdk.NewDecCoins(sdk.NewDecCoin("ALX", sdk.NewInt(1))),
					BypassMinFeeMsgTypes:            []string{},
					MaxTotalBypassMinFeeMsgGasUsage: uint64(0),
				},
			},
		},
		"multiple fee options": {
			src: `{"params":{"minimum_gas_prices":[{"denom":"ALX", "amount":"1"}, {"denom":"BLX", "amount":"0.001"}],"bypass_min_fee_msg_types":[],"max_total_bypass_min_fee_msg_gas_usage":"0"}}`,
			exp: types.GenesisState{
				Params: types.Params{
					MinimumGasPrices: sdk.NewDecCoins(
						sdk.NewDecCoin("ALX", sdk.NewInt(1)),
						sdk.NewDecCoinFromDec("BLX", sdk.NewDecWithPrec(1, 3)),
					),
					BypassMinFeeMsgTypes:            []string{},
					MaxTotalBypassMinFeeMsgGasUsage: uint64(0),
				},
			},
		},
		"no fee set": {
			src: `{"params":{"minimum_gas_prices":[],"bypass_min_fee_msg_types":[],"max_total_bypass_min_fee_msg_gas_usage":"0"}}`,
			exp: types.GenesisState{
				Params: types.Params{
					MinimumGasPrices:                sdk.DecCoins{},
					BypassMinFeeMsgTypes:            []string{},
					MaxTotalBypassMinFeeMsgGasUsage: uint64(0),
				},
			},
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			ctx, encCfg, keeper := setupTestStore(t)
			m := NewAppModule(encCfg.Marshaler, keeper, "uflix")
			m.InitGenesis(ctx, encCfg.Marshaler, []byte(spec.src))
			gotJSON := m.ExportGenesis(ctx, encCfg.Marshaler)
			var got types.GenesisState
			t.Log(got)
			require.NoError(t, encCfg.Marshaler.UnmarshalJSON(gotJSON, &got))
			assert.Equal(t, spec.exp, got, string(gotJSON))
		})
	}
}

func setupTestStore(t *testing.T) (sdk.Context, appparams.EncodingConfig, globalfeekeeper.Keeper) {
	t.Helper()
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	encCfg := appparams.MakeEncodingConfig()
	keyParams := sdk.NewKVStoreKey(types.StoreKey)
	ms.MountStoreWithDB(keyParams, storetypes.StoreTypeIAVL, db)
	require.NoError(t, ms.LoadLatestVersion())

	globalfeeKeeper := globalfeekeeper.NewKeeper(encCfg.Marshaler, keyParams, "omniflix1llyd96levrglxhw6sczgk9wn48t64zkhv4fq0r")

	ctx := sdk.NewContext(ms, tmproto.Header{
		Height:  1234567,
		Time:    time.Date(2020, time.April, 22, 12, 0, 0, 0, time.UTC),
		ChainID: "testing",
	}, false, log.NewNopLogger())

	return ctx, encCfg, globalfeeKeeper
}
