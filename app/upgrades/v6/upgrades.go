package v6

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/OmniFlix/omniflixhub/v6/app/keepers"
	"github.com/OmniFlix/omniflixhub/v6/app/upgrades"
	medianodetypes "github.com/OmniFlix/omniflixhub/v6/x/medianode/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

func CreateV6UpgradeHandler(
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

		// Set Default Params

		err = keepers.MedianodeKeeper.SetParams(ctx, medianodetypes.DefaultParams())
		if err != nil {
			return nil, err
		}

		ctx.Logger().Info("Upgrade complete")
		return versionMap, nil
	}
}
