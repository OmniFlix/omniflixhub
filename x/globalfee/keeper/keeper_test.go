package keeper_test

import (
	"testing"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/store/metrics"
	rootmulti "cosmossdk.io/store/rootmulti"
	storetypes "cosmossdk.io/store/types"
	appparams "github.com/OmniFlix/omniflixhub/v5/app/params"
	globalfeekeeper "github.com/OmniFlix/omniflixhub/v5/x/globalfee/keeper"
	"github.com/OmniFlix/omniflixhub/v5/x/globalfee/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func setupTestStore(t *testing.T) (sdk.Context, appparams.EncodingConfig, globalfeekeeper.Keeper) {
	t.Helper()
	db := dbm.NewMemDB()
	logger := log.NewNopLogger()
	ms := rootmulti.NewStore(db, logger, metrics.NewNoOpMetrics())
	encCfg := appparams.MakeEncodingConfig()
	keyParams := storetypes.NewKVStoreKey(types.StoreKey)
	ms.MountStoreWithDB(keyParams, storetypes.StoreTypeIAVL, db)
	require.NoError(t, ms.LoadLatestVersion())

	globalFeeKeeper := globalfeekeeper.NewKeeper(encCfg.Marshaler, keyParams, "omniflix1llyd96levrglxhw6sczgk9wn48t64zkhv4fq0r")

	ctx := sdk.NewContext(ms, tmproto.Header{
		Height:  1234567,
		Time:    time.Date(2020, time.April, 22, 12, 0, 0, 0, time.UTC),
		ChainID: "omniflixhub-test",
	}, false, log.NewNopLogger())

	return ctx, encCfg, globalFeeKeeper
}
