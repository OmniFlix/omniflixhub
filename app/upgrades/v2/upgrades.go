package v2

import (
	"github.com/OmniFlix/omniflixhub/v2/app/keepers"
	"github.com/OmniFlix/omniflixhub/v2/app/upgrades"
	streampaytypes "github.com/OmniFlix/streampay/v2/x/streampay/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateV2UpgradeHandler(
	mm *module.Manager,
	cfg module.Configurator,
	bpm upgrades.BaseAppParamManager,
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
		if err := keepers.StreamPayKeeper.SetParams(ctx, streampaytypes.DefaultParams()); err != nil {
			return nil, err
		}

		ctx.Logger().Info("Upgrade complete")
		return versionMap, nil
	}
}
