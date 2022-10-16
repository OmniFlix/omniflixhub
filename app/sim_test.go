package app_test

/*


import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/OmniFlix/omniflixhub/app/helpers"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/libs/rand"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/store"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	simulation2 "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

func init() {
	simapp.GetSimulatorFlags()
}

// interBlockCacheOpt returns a BaseApp option function that sets the persistent
// inter-block write-through cache.
func interBlockCacheOpt() func(*baseapp.BaseApp) {
	return baseapp.SetInterBlockCache(store.NewCommitKVStoreCacheManager())
}

func TestFullAppSimulation(t *testing.T) {
	config, db, dir, logger, skip, err := simapp.SetupSimulation("leveldb-app-sim", "Simulation")
	if skip {
		t.Skip("skipping application simulation")
	}
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		db.Close()
		require.NoError(t, os.RemoveAll(dir))
	}()

	app := NewOmniFlixApp(logger, db, nil, true, map[int64]bool{}, DefaultNodeHome, simapp.FlagPeriodValue, MakeEncodingConfig(), helpers.EmptyAppOptions{}, fauxMerkleModeOpt)
	require.Equal(t, "onft", app.Name())

	// run randomized simulation
	_, simParams, simErr := simulation.SimulateFromSeed(
		t,
		os.Stdout,
		app.BaseApp,
		AppStateFn(app.AppCodec(), app.SimulationManager()),
		simtypes.RandomAccounts, // Replace with own random account function if using keys other than secp256k1
		simapp.SimulationOperations(app, app.AppCodec(), config),
		app.ModuleAccountAddrs(),
		config,
		app.AppCodec(),
	)

	// export state and simParams before the simulation error is checked
	err = simapp.CheckExportSimulation(app, config, simParams)
	require.NoError(t, err)
	require.NoError(t, simErr)

	if config.Commit {
		simapp.PrintStats(db)
	}
}

// fauxMerkleModeOpt returns a BaseApp option to use a dbStoreAdapter instead of
// an IAVLStore for faster simulation speed.
func fauxMerkleModeOpt(bapp *baseapp.BaseApp) {
	bapp.SetFauxMerkleMode()
}

//// TODO: Make another test for the fuzzer itself, which just has noOp txs
//// and doesn't depend on the application.
func TestAppStateDeterminism(t *testing.T) {
	if !simapp.FlagEnabledValue {
		t.Skip("skipping application simulation")
	}

	config := simapp.NewConfigFromFlags()
	config.InitialBlockHeight = 1
	config.ExportParamsPath = ""
	config.OnOperation = false
	config.AllInvariants = false
	config.ChainID = helpers.SimAppChainID

	numSeeds := 3
	numTimesToRunPerSeed := 5
	appHashList := make([]json.RawMessage, numTimesToRunPerSeed)

	for i := 0; i < numSeeds; i++ {
		config.Seed = rand.Int63()

		for j := 0; j < numTimesToRunPerSeed; j++ {
			var logger log.Logger
			if simapp.FlagVerboseValue {
				logger = log.TestingLogger()
			} else {
				logger = log.NewNopLogger()
			}

			db := dbm.NewMemDB()
			app := NewOmniFlixApp(
				logger, db, nil, true, map[int64]bool{}, DefaultNodeHome,
				simapp.FlagPeriodValue, MakeEncodingConfig(),
				simapp.EmptyAppOptions{}, interBlockCacheOpt())

			fmt.Printf(
				"running non-determinism simulation; seed %d: %d/%d, attempt: %d/%d\n",
				config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
			)

			_, _, err := simulation.SimulateFromSeed(
				t,
				os.Stdout,
				app.BaseApp,
				AppStateFn(app.AppCodec(), app.SimulationManager()),
				simulation2.RandomAccounts, // Replace with own random account function if using keys other than secp256k1
				simapp.SimulationOperations(app, app.AppCodec(), config),
				app.ModuleAccountAddrs(),
				config,
				app.AppCodec(),
			)
			require.NoError(t, err)

			if config.Commit {
				simapp.PrintStats(db)
			}

			appHash := app.LastCommitID().Hash
			appHashList[j] = appHash

			if j != 0 {
				require.Equal(
					t, hex.EncodeToString(appHashList[0]), hex.EncodeToString(appHashList[j]),
					"non-determinism in seed %d: %d/%d, attempt: %d/%d\n", config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
				)
			}
		}
	}
}
*/
