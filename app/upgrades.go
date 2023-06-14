package app

import (
	"fmt"

	itctypes "github.com/OmniFlix/omniflixhub/x/itc/types"

	streampaytypes "github.com/OmniFlix/streampay/v2/x/streampay/types"

	marketplacetypes "github.com/OmniFlix/marketplace/x/marketplace/types"
	store "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// next upgrade name
const upgradeName = "v0.11.0"

// RegisterUpgradeHandlers returns upgrade handlers
func (app *App) RegisterUpgradeHandlers(cfg module.Configurator) {
	app.UpgradeKeeper.SetUpgradeHandler(
		upgradeName,
		func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			// set marketplace module params
			app.MarketplaceKeeper.SetParams(ctx, marketplacetypes.DefaultParams())

			// set streampay module params
			streamPayParams := streampaytypes.DefaultParams()
			streamPayParams.StreamPaymentFee = sdk.NewInt64Coin("uflix", 50_000_000) // 50 FLIX
			app.StreamPayKeeper.SetParams(ctx, streamPayParams)

			// set itc module params
			app.ItcKeeper.SetParams(ctx, itctypes.DefaultParams())

			ctx.Logger().Info("running migrations ...")
			return app.mm.RunMigrations(ctx, cfg, fromVM)
		})

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if upgradeInfo.Name == upgradeName && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := store.StoreUpgrades{
			Added: []string{streampaytypes.ModuleName, itctypes.ModuleName},
		}
		// configure store loader that checks if height == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}
