package app

import (
	"fmt"

	alloctypes "github.com/OmniFlix/omniflixhub/x/alloc/types"

	onfttypes "github.com/OmniFlix/onft/types"

	store "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// next upgrade name
const upgradeName = "v0.10.x-part2"

const (
	developerRewardsAccount    = "omniflix1ftvf4euvdvq95jpeyvgf6r6j78j5rct2a3jnkn"
	nftIncentivesAccount       = "omniflix139qa9tklr4trzvugqm5ycvky80px90yn5hs3kc"
	nodeHostsIncentivesAccount = "omniflix1djc90zwkk2vaqryne8c68f2tkp6u9ug9qfrnh8"
)

// RegisterUpgradeHandlers returns upgrade handlers
func (app *App) RegisterUpgradeHandlers(cfg module.Configurator) {
	app.UpgradeKeeper.SetUpgradeHandler(upgradeName, func(ctx sdk.Context, plan upgradetypes.Plan, _ module.VersionMap) (module.VersionMap, error) {
		developerRewardsAcc, err := sdk.AccAddressFromBech32(developerRewardsAccount)
		if err != nil {
			panic(err)
		}
		nftIncentivesAcc, err := sdk.AccAddressFromBech32(nftIncentivesAccount)
		if err != nil {
			panic(err)
		}
		nodeHostsIncentivesAcc, err := sdk.AccAddressFromBech32(nodeHostsIncentivesAccount)
		if err != nil {
			panic(err)
		}

		app.ONFTKeeper.SetParams(ctx, onfttypes.DefaultParams())
		app.AllocKeeper.SetParams(ctx, alloctypes.Params{
			DistributionProportions: alloctypes.DistributionProportions{
				StakingRewards:      sdk.NewDecWithPrec(60, 2), // 60%
				NftIncentives:       sdk.NewDecWithPrec(15, 2), // 15%
				NodeHostsIncentives: sdk.NewDecWithPrec(5, 2),  // 5%
				DeveloperRewards:    sdk.NewDecWithPrec(15, 2), // 15%
				CommunityPool:       sdk.NewDecWithPrec(5, 2),  // 5%
			},
			WeightedDeveloperRewardsReceivers: []alloctypes.WeightedAddress{
				{
					Address: developerRewardsAcc.String(),
					Weight:  sdk.NewDecWithPrec(100, 2),
				},
			},
			WeightedNftIncentivesReceivers: []alloctypes.WeightedAddress{
				{
					Address: nftIncentivesAcc.String(),
					Weight:  sdk.NewDecWithPrec(100, 2),
				},
			},
			WeightedNodeHostsIncentivesReceivers: []alloctypes.WeightedAddress{
				{
					Address: nodeHostsIncentivesAcc.String(),
					Weight:  sdk.NewDecWithPrec(100, 2),
				},
			},
		})

		fromVM := app.mm.GetVersionMap()
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
		// configure store loader that checks if height == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}
