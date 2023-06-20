package app

import (
	"fmt"

	appparams "github.com/OmniFlix/omniflixhub/app/params"

	itctypes "github.com/OmniFlix/omniflixhub/x/itc/types"

	maketplacetypes "github.com/OmniFlix/marketplace/x/marketplace/types"
	streampaytypes "github.com/OmniFlix/streampay/v2/x/streampay/types"

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
			ctx.Logger().Info("running migrations ...")
			// Run migrations before applying any other state changes.
			versionMap, err := app.mm.RunMigrations(ctx, cfg, fromVM)
			if err != nil {
				return nil, err
			}
			// update marketplace module new parameters
			marketplaceParamSubspace, ok := app.ParamsKeeper.GetSubspace(maketplacetypes.ModuleName)
			if !ok {
				panic("marketplace params subspace not found")
			}
			marketplaceParamSubspace.Set(
				ctx,
				maketplacetypes.ParamStoreKeyMaxAuctionDuration,
				maketplacetypes.DefaultMaxAuctionDuration,
			)
			// set streampay params
			streampayParams := streampaytypes.DefaultParams()
			streampayParams.StreamPaymentFee = sdk.NewInt64Coin(appparams.BondDenom, 50_000_000)
			app.StreamPayKeeper.SetParams(ctx, streampayParams)

			// set itc module params
			app.ItcKeeper.SetParams(ctx, itctypes.DefaultParams())

			return versionMap, nil
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
