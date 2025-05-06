package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cosmos/cosmos-sdk/baseapp"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"

	tmcfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/prometheus/client_golang/prometheus"

	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/log"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/OmniFlix/omniflixhub/v6/app"
	"github.com/OmniFlix/omniflixhub/v6/app/params"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/pruning"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/snapshot"
	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtxconfig "github.com/cosmos/cosmos-sdk/x/auth/tx/config"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	rosettaCmd "github.com/cosmos/rosetta/cmd"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewRootCmd creates a new root command for omniflixhubd. It is called once in the
// main function.
func NewRootCmd() (*cobra.Command, params.EncodingConfig) {
	// Set config for prefixes
	encodingConfig := app.GetEncodingConfig()

	initAppOptions := viper.New()
	tempDir := tempDir()
	initAppOptions.Set(flags.FlagHome, tempDir)

	tempApp := app.NewOmniFlixApp(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		true,
		nil,
		tempDir,
		0,
		encodingConfig,
		initAppOptions,
		nil,
		baseapp.SetChainID("omniflixhub-1"),
	)
	defer func() {
		if err := tempApp.Close(); err != nil {
			panic(err)
		}
		if tempDir != app.DefaultNodeHome {
			os.RemoveAll(tempDir)
		}
	}()

	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastSync).
		WithHomeDir(app.DefaultNodeHome).
		WithViper("")

	rootCmd := &cobra.Command{
		Use:   app.Name + "d",
		Short: "OmniFlix Hub App",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			// This needs to go after ReadFromClientConfig, as that function
			// sets the RPC client needed for SIGN_MODE_TEXTUAL. This sign mode
			// is only available if the client is online.
			if !initClientCtx.Offline {
				txConfigOpts := tx.ConfigOptions{
					EnabledSignModes:           append(tx.DefaultSignModes, signing.SignMode_SIGN_MODE_TEXTUAL),
					TextualCoinMetadataQueryFn: authtxconfig.NewGRPCCoinMetadataQueryFn(initClientCtx),
				}
				txConfigWithTextual, err := tx.NewTxConfigWithOptions(
					initClientCtx.Codec,
					txConfigOpts,
				)
				if err != nil {
					return err
				}
				initClientCtx = initClientCtx.WithTxConfig(txConfigWithTextual)
			}

			if err = client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			appTemplate, appConfig := initAppConfig()
			customTMConfig := initCometBftConfig()
			return server.InterceptConfigsPreRunHandler(cmd, appTemplate, appConfig, customTMConfig)
		},
	}

	initRootCmd(rootCmd, encodingConfig, tempApp)

	// must need as part of cosmos-sdk 0.50.x migration
	autocliOptions := autoCliOpts(initClientCtx, tempApp)
	if err := autocliOptions.EnhanceRootCommand(rootCmd); err != nil {
		panic(err)
	}

	return rootCmd, encodingConfig
}

// initCometBftConfig helps to override default CometBft Config values.
// return tmcfg.DefaultConfig if no custom configuration is required for the application.
func initCometBftConfig() *tmcfg.Config {
	cfg := tmcfg.DefaultConfig()

	// these values put a higher strain on node memory
	cfg.P2P.MaxNumInboundPeers = 100
	cfg.P2P.MaxNumOutboundPeers = 40

	return cfg
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig params.EncodingConfig, tempApp *app.OmniFlixApp) {
	ac := appCreator{
		encCfg: encodingConfig,
	}

	rootCmd.AddCommand(
		genutilcli.InitCmd(tempApp.ModuleBasics, app.DefaultNodeHome),
		tmcli.NewCompletionCmd(rootCmd, true),
		genutilcli.CollectGenTxsCmd(
			banktypes.GenesisBalancesIterator{},
			app.DefaultNodeHome,
			genutiltypes.DefaultMessageValidator,
			tempApp.StakingKeeper.ValidatorAddressCodec(),
		),
		genutilcli.MigrateGenesisCmd(genutilcli.MigrationMap),
		genutilcli.GenTxCmd(
			tempApp.ModuleBasics,
			encodingConfig.TxConfig,
			banktypes.GenesisBalancesIterator{},
			app.DefaultNodeHome,
			tempApp.StakingKeeper.ValidatorAddressCodec(),
		),
		genutilcli.ValidateGenesisCmd(tempApp.ModuleBasics),
		genutilcli.AddGenesisAccountCmd(app.DefaultNodeHome, tempApp.AccountKeeper.AddressCodec()),
		tmcli.NewCompletionCmd(rootCmd, true),
		addDebugCommands(debug.Cmd()),
		pruning.Cmd(ac.newApp, app.DefaultNodeHome),
		snapshot.Cmd(ac.newApp),
	)

	server.AddCommands(rootCmd, app.DefaultNodeHome, ac.newApp, ac.appExport, addModuleInitFlags)

	// add keybase, auxiliary RPC, query, and tx child commands
	rootCmd.AddCommand(
		server.StatusCommand(),
		queryCommand(),
		txCommand(tempApp.ModuleBasics),
		keys.Commands(),
	)
	// add rosetta
	rootCmd.AddCommand(rosettaCmd.RosettaCommand(encodingConfig.InterfaceRegistry, encodingConfig.Marshaler))
}

func addModuleInitFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd)
	wasm.AddModuleInitFlags(startCmd)
}

func tempDir() string {
	dir, err := os.MkdirTemp("", "."+app.Name+"-temp")
	if err != nil {
		panic(fmt.Sprintf("failed creating temp directory: %s", err.Error()))
	}
	defer os.RemoveAll(dir)

	return dir
}

func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		rpc.ValidatorCommand(),
		server.QueryBlockCmd(),
		authcmd.QueryTxsByEventsCmd(),
		authcmd.QueryTxCmd(),
	)

	app.ModuleBasics.AddQueryCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func txCommand(moduleBasics module.BasicManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetValidateSignaturesCommand(),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
	)

	moduleBasics.AddTxCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

type appCreator struct {
	encCfg params.EncodingConfig
}

// newApp is an AppCreator
func (ac appCreator) newApp(logger log.Logger, db dbm.DB, traceStore io.Writer, appOpts servertypes.AppOptions) servertypes.Application {
	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}

	var wasmOpts []wasmkeeper.Option
	if cast.ToBool(appOpts.Get("telemetry.enabled")) {
		wasmOpts = append(wasmOpts, wasmkeeper.WithVMCacheMetrics(prometheus.DefaultRegisterer))
	}

	baseappOptions := server.DefaultBaseappOptions(appOpts)

	return app.NewOmniFlixApp(
		logger,
		db,
		traceStore,
		true,
		skipUpgradeHeights,
		cast.ToString(appOpts.Get(flags.FlagHome)),
		cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod)),
		ac.encCfg,
		appOpts,
		wasmOpts,
		baseappOptions...,
	)
}

// appExport creates a new app (optionally at a given height)
func (ac appCreator) appExport(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
	modulesToExport []string,
) (servertypes.ExportedApp, error) {
	var anApp *app.OmniFlixApp

	homePath, ok := appOpts.Get(flags.FlagHome).(string)
	if !ok || homePath == "" {
		return servertypes.ExportedApp{}, errors.New("application home not set")
	}

	var emptyWasmOpts []wasmkeeper.Option
	if height != -1 {
		anApp = app.NewOmniFlixApp(
			logger,
			db,
			traceStore,
			false,
			map[int64]bool{},
			homePath,
			uint(1),
			ac.encCfg,
			appOpts,
			emptyWasmOpts,
		)

		if err := anApp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	} else {
		anApp = app.NewOmniFlixApp(
			logger,
			db,
			traceStore,
			true,
			map[int64]bool{},
			homePath,
			uint(1),
			ac.encCfg,
			appOpts,
			emptyWasmOpts,
		)
	}

	return anApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs, modulesToExport)
}

func autoCliOpts(initClientCtx client.Context, tempApp *app.OmniFlixApp) autocli.AppOptions {
	modules := make(map[string]appmodule.AppModule)
	for _, m := range tempApp.ModuleManager().Modules {
		if moduleWithName, ok := m.(module.HasName); ok {
			moduleName := moduleWithName.Name()
			if appModule, ok := moduleWithName.(appmodule.AppModule); ok {
				modules[moduleName] = appModule
			}
		}
	}

	return autocli.AppOptions{
		Modules:               modules,
		ModuleOptions:         runtimeservices.ExtractAutoCLIOptions(tempApp.ModuleManager().Modules),
		AddressCodec:          authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		ValidatorAddressCodec: authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		ConsensusAddressCodec: authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
		ClientCtx:             initClientCtx,
	}
}
