package app

import (
	"fmt"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"

	marketplacetypes "github.com/OmniFlix/marketplace/x/marketplace/types"
	store "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// next upgrade name
const upgradeName = "upgrade_1"

func equalTraces(dtA, dtB ibctransfertypes.DenomTrace) bool {
	return dtA.BaseDenom == dtB.BaseDenom && dtA.Path == dtB.Path
}

// RegisterUpgradeHandlers returns upgrade handlers
func (app *App) RegisterUpgradeHandlers(cfg module.Configurator) {
	app.UpgradeKeeper.SetUpgradeHandler(upgradeName, func(ctx sdk.Context, plan upgradetypes.Plan, _ module.VersionMap) (module.VersionMap, error) {
		app.MarketplaceKeeper.SetParams(ctx, marketplacetypes.DefaultParams())
		app.MarketplaceKeeper.SetNextAuctionNumber(ctx, 1)

		fromVM := app.mm.GetVersionMap()

		var newTraces []ibctransfertypes.DenomTrace
		app.TransferKeeper.IterateDenomTraces(ctx,
			func(dt ibctransfertypes.DenomTrace) bool {
				// check if the new way of splitting FullDenom
				// into Trace and BaseDenom passes validation and
				// is the same as the current DenomTrace.
				// If it isn't then store the new DenomTrace in the list of new traces.
				newTrace := ibctransfertypes.ParseDenomTrace(dt.GetFullDenomPath())
				if err := newTrace.Validate(); err == nil && !equalTraces(newTrace, dt) {
					newTraces = append(newTraces, newTrace)
				}
				return false
			})
		// replace the outdated traces with the new trace information
		for _, nt := range newTraces {
			app.TransferKeeper.SetDenomTrace(ctx, nt)
		}

		ctx.Logger().Info("running migrations ...")
		return app.mm.RunMigrations(ctx, cfg, fromVM)
	})

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if upgradeInfo.Name == upgradeName && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := store.StoreUpgrades{
			Added: []string{},
		}
		// configure store loader that checks if hight == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}
