package keeper_test

import (
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/OmniFlix/omniflixhub/v5/app/apptesting"
	"github.com/OmniFlix/omniflixhub/v5/x/alloc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper

	queryClient types.QueryClient
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.Reset()
	suite.queryClient = types.NewQueryClient(suite.QueryHelper)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func FundModuleAccount(bankKeeper bankkeeper.Keeper, ctx sdk.Context, recipientMod string, amounts sdk.Coins) error {
	if err := bankKeeper.MintCoins(ctx, types.ModuleName, amounts); err != nil {
		return err
	}
	return bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, recipientMod, amounts)
}

func (suite *KeeperTestSuite) TestParams() {
	testCases := []struct {
		name      string
		input     types.Params
		expectErr bool
	}{
		{
			name: "set invalid params",
			input: types.Params{
				DistributionProportions: types.DistributionProportions{
					StakingRewards:      sdkmath.LegacyNewDecWithPrec(-1, 2),
					NftIncentives:       sdkmath.LegacyNewDecWithPrec(100, 2),
					NodeHostsIncentives: sdkmath.LegacyNewDecWithPrec(5, 2),
					DeveloperRewards:    sdkmath.LegacyNewDecWithPrec(0, 2),
					CommunityPool:       sdkmath.LegacyNewDecWithPrec(5, 2),
				},
			},
			expectErr: true,
		},
		{
			name: "set full valid params",
			input: types.Params{
				DistributionProportions: types.DistributionProportions{
					StakingRewards:      sdkmath.LegacyNewDecWithPrec(60, 2), // 60%
					NftIncentives:       sdkmath.LegacyNewDecWithPrec(15, 2), // 15%
					NodeHostsIncentives: sdkmath.LegacyNewDecWithPrec(5, 2),  // 5%
					DeveloperRewards:    sdkmath.LegacyNewDecWithPrec(15, 2), // 15%
					CommunityPool:       sdkmath.LegacyNewDecWithPrec(5, 2),  // 5%
				},
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			expected := suite.App.AllocKeeper.GetParams(suite.Ctx)
			err := suite.App.AllocKeeper.SetParams(suite.Ctx, tc.input)
			if tc.expectErr {
				suite.Require().Error(err)
			} else {
				expected = tc.input
				suite.Require().NoError(err)
			}

			p := suite.App.AllocKeeper.GetParams(suite.Ctx)
			suite.Require().Equal(expected, p)
		})
	}
}

func (suite *KeeperTestSuite) TestDistribution() {
	suite.SetupTest()
	denom, _ := suite.App.StakingKeeper.BondDenom(suite.Ctx)
	allocKeeper := suite.App.AllocKeeper
	params := suite.App.AllocKeeper.GetParams(suite.Ctx)
	nftIncentivesReceiver, _ := sdk.AccAddressFromBech32("omniflix139qa9tklr4trzvugqm5ycvky80px90yn5hs3kc")
	nodeHostIncentivesReceiver, _ := sdk.AccAddressFromBech32("omniflix1djc90zwkk2vaqryne8c68f2tkp6u9ug9qfrnh8")
	devRewardsReceiver, _ := sdk.AccAddressFromBech32("omniflix1ftvf4euvdvq95jpeyvgf6r6j78j5rct2a3jnkn")
	params.DistributionProportions.StakingRewards = sdkmath.LegacyNewDecWithPrec(60, 2)
	params.DistributionProportions.NodeHostsIncentives = sdkmath.LegacyNewDecWithPrec(5, 2)
	params.DistributionProportions.NftIncentives = sdkmath.LegacyNewDecWithPrec(15, 2)
	params.DistributionProportions.CommunityPool = sdkmath.LegacyNewDecWithPrec(5, 2)
	params.DistributionProportions.DeveloperRewards = sdkmath.LegacyNewDecWithPrec(15, 2)
	params.WeightedNftIncentivesReceivers = []types.WeightedAddress{
		{
			Address: nftIncentivesReceiver.String(),
			Weight:  sdkmath.LegacyNewDec(1),
		},
	}
	params.WeightedNodeHostsIncentivesReceivers = []types.WeightedAddress{
		{
			Address: nodeHostIncentivesReceiver.String(),
			Weight:  sdkmath.LegacyNewDec(1),
		},
	}
	params.WeightedDeveloperRewardsReceivers = []types.WeightedAddress{
		{
			Address: devRewardsReceiver.String(),
			Weight:  sdkmath.LegacyNewDec(1),
		},
	}
	err := suite.App.AllocKeeper.SetParams(suite.Ctx, params)
	fmt.Println(params, suite.App.AppCodec())
	suite.Require().NoError(err)
	distributionModuleAccount := suite.App.DistrKeeper.GetDistributionAccount(suite.Ctx)
	feeCollector := suite.App.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	suite.Equal(
		"0",
		suite.App.BankKeeper.GetAllBalances(suite.Ctx, feeCollector).AmountOf(denom).String())
	suite.Equal(
		sdkmath.LegacyNewDec(0).TruncateInt().String(),
		suite.App.BankKeeper.GetAllBalances(suite.Ctx, distributionModuleAccount.GetAddress()).AmountOf(denom).String())

	mintCoin := sdk.NewCoin(denom, sdkmath.NewInt(100_000))
	mintCoins := sdk.Coins{mintCoin}
	feeCollectorAccount := suite.App.AccountKeeper.GetModuleAccount(suite.Ctx, authtypes.FeeCollectorName)
	suite.Require().NotNil(feeCollectorAccount)

	suite.Require().NoError(FundModuleAccount(suite.App.BankKeeper, suite.Ctx, feeCollectorAccount.GetName(), mintCoins))
	feeCollector = suite.App.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	suite.Equal(
		mintCoin.Amount.String(),
		suite.App.BankKeeper.GetAllBalances(suite.Ctx, feeCollector).AmountOf(denom).String())

	suite.Equal(
		sdkmath.LegacyNewDec(0).TruncateInt().String(),
		suite.App.BankKeeper.GetAllBalances(suite.Ctx, distributionModuleAccount.GetAddress()).AmountOf(denom).String(),
	)

	_ = allocKeeper.DistributeMintedCoins(suite.Ctx)

	feeCollector = suite.App.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	modulePortion := params.DistributionProportions.NftIncentives.
		Add(params.DistributionProportions.DeveloperRewards).
		Add(params.DistributionProportions.NodeHostsIncentives).
		Add(params.DistributionProportions.CommunityPool) // 40%

	// remaining going to next module should be 100% - 40% = 60%
	suite.Equal(
		sdkmath.LegacyNewDecFromInt(mintCoin.Amount).Mul(sdkmath.LegacyNewDecWithPrec(100, 2).Sub(modulePortion)).TruncateInt().String(),
		suite.App.BankKeeper.GetAllBalances(suite.Ctx, feeCollector).AmountOf(denom).String())

	suite.Equal(
		sdkmath.LegacyNewDecFromInt(mintCoin.Amount).Mul(params.DistributionProportions.DeveloperRewards).TruncateInt().String(),
		suite.App.BankKeeper.GetBalance(suite.Ctx, devRewardsReceiver, denom).Amount.String())

	// since the NFT incentives are not setup yet, funds go into the community pool
	distributionModuleAccount = suite.App.DistrKeeper.GetDistributionAccount(suite.Ctx)
	communityPoolPortion := params.DistributionProportions.CommunityPool // 5%

	suite.Equal(
		sdkmath.LegacyNewDecFromInt(mintCoin.Amount).Mul(communityPoolPortion).TruncateInt().String(),
		suite.App.BankKeeper.GetAllBalances(suite.Ctx, distributionModuleAccount.GetAddress()).AmountOf(denom).String())
}
