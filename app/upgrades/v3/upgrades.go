package v3

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/OmniFlix/omniflixhub/v3/app/keepers"
	"github.com/OmniFlix/omniflixhub/v3/app/upgrades"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateV3UpgradeHandler(
	mm *module.Manager,
	cfg module.Configurator,
	_ upgrades.BaseAppParamManager,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("running migrations ...")
		// Run migrations before applying any other state changes.
		// NOTE: DO NOT PUT ANY STATE CHANGES BEFORE RunMigrations().
		versionMap, err := mm.RunMigrations(ctx, cfg, fromVM)
		if err != nil {
			return nil, err
		}

		params := wasmtypes.DefaultParams()
		// change this to AllowNoBody for mainnet
		params.CodeUploadAccess = wasmtypes.AllowEverybody
		params.InstantiateDefaultPermission = wasmtypes.AccessTypeEverybody
		if err := keepers.WasmKeeper.SetParams(ctx, params); err != nil {
			return nil, err
		}

		ctx.Logger().Info("Upgrade complete")
		return versionMap, nil
	}
}
