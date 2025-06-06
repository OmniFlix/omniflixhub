package v3

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/OmniFlix/omniflixhub/v6/app/keepers"
	"github.com/OmniFlix/omniflixhub/v6/app/upgrades"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

func CreateV3UpgradeHandler(
	mm *module.Manager,
	cfg module.Configurator,
	_ upgrades.BaseAppParamManager,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(context context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(context)
		// Run migrations before applying any other state changes.
		// NOTE: DO NOT PUT ANY STATE CHANGES BEFORE RunMigrations().
		versionMap, err := mm.RunMigrations(ctx, cfg, fromVM)
		if err != nil {
			return nil, err
		}

		params := wasmtypes.DefaultParams()
		// Permissioned cosmwasm
		params.CodeUploadAccess = wasmtypes.AllowNobody
		params.InstantiateDefaultPermission = wasmtypes.AccessTypeEverybody
		if err := keepers.WasmKeeper.SetParams(ctx, params); err != nil {
			return nil, err
		}

		ctx.Logger().Info("Upgrade complete")
		return versionMap, nil
	}
}
