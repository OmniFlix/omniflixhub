package v6

import (
	"context"
	"time"

	sdkmath "cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/OmniFlix/omniflixhub/v6/app/keepers"
	"github.com/OmniFlix/omniflixhub/v6/app/upgrades"
	medianodetypes "github.com/OmniFlix/omniflixhub/v6/x/medianode/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	feemarkettypes "github.com/skip-mev/feemarket/x/feemarket/types"
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

		// Configure Medianode Module
		err = ConfigureMedianodeModule(ctx, keepers)
		if err != nil {
			return versionMap, err
		}

		err = ConfigureFeeMarketModule(ctx, keepers)
		if err != nil {
			return versionMap, err
		}

		ctx.Logger().Info("Upgrade complete")
		return versionMap, nil
	}
}

func ConfigureMedianodeModule(ctx sdk.Context, keepers *keepers.AppKeepers) error {
	params := medianodetypes.DefaultParams()

	params.MinDeposit = sdk.NewCoin("uflix", sdkmath.NewInt(5000_000_000)) // 5000 FLIX
	params.InitialDepositPercentage = sdkmath.LegacyNewDecWithPrec(10, 2)
	params.MinimumLeaseHours = 3                                // 3 hours
	params.MaximumLeaseHours = 8760                             // 1 year
	params.DepositReleasePeriod = time.Hour * 24 * 7            // 1 week
	params.LeaseCommission = sdkmath.LegacyNewDecWithPrec(1, 2) // 1%
	params.CommissionDistribution = medianodetypes.Distribution{
		Staking:       sdkmath.LegacyNewDecWithPrec(50, 2), // 50%
		CommunityPool: sdkmath.LegacyNewDecWithPrec(50, 2), // 50%
	}

	if err := keepers.MedianodeKeeper.SetParams(ctx, params); err != nil {
		return err
	}

	return nil
}

func ConfigureFeeMarketModule(ctx sdk.Context, keepers *keepers.AppKeepers) error {
	params, err := keepers.FeeMarketKeeper.GetParams(ctx)
	if err != nil {
		return err
	}

	params.Enabled = true
	params.FeeDenom = "uflix"
	params.DistributeFees = true
	params.MinBaseGasPrice = sdkmath.LegacyMustNewDecFromStr("0.005")
	params.MaxBlockUtilization = feemarkettypes.DefaultMaxBlockUtilization
	if err := keepers.FeeMarketKeeper.SetParams(ctx, params); err != nil {
		return err
	}

	state, err := keepers.FeeMarketKeeper.GetState(ctx)
	if err != nil {
		return err
	}

	state.BaseGasPrice = sdkmath.LegacyMustNewDecFromStr("0.005")

	return keepers.FeeMarketKeeper.SetState(ctx, state)
}
