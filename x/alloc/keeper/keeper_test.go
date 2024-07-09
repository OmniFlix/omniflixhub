package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/OmniFlix/omniflixhub/v5/app/apptesting"
	"github.com/OmniFlix/omniflixhub/v5/x/alloc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/stretchr/testify/suite"
	"testing"

	"github.com/OmniFlix/omniflixhub/v5/app"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper
	suite.Suite
	ctx sdk.Context

	app *app.OmniFlixApp
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = app.Setup(suite.T())
	suite.ctx = suite.app.BaseApp.NewContext(false)
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
			expected := suite.app.AllocKeeper.GetParams(suite.ctx)
			err := suite.app.AllocKeeper.SetParams(suite.ctx, tc.input)
			if tc.expectErr {
				suite.Require().Error(err)
			} else {
				expected = tc.input
				suite.Require().NoError(err)
			}

			p := suite.app.AllocKeeper.GetParams(suite.ctx)
			suite.Require().Equal(expected, p)
		})
	}
}

func (suite *KeeperTestSuite) TestDistribution() {
	suite.SetupTest()

	denom, _ := suite.app.StakingKeeper.BondDenom(suite.ctx)
	allocKeeper := suite.app.AllocKeeper
	params := suite.app.AllocKeeper.GetParams(suite.ctx)
	nftIncentivesReceiver := sdk.AccAddress([]byte("addr1a--------------"))
	nodeHostIncentivesReceiver := sdk.AccAddress([]byte("addr1b--------------"))
	devRewardsReceiver := sdk.AccAddress([]byte("addr1c--------------"))
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
	err := suite.app.AllocKeeper.SetParams(suite.ctx, params)
	suite.Require().NoError(err)
	distributionModuleAccount := suite.app.DistrKeeper.GetDistributionAccount(suite.ctx)
	feeCollector := suite.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	suite.Equal(
		"0",
		suite.app.BankKeeper.GetAllBalances(suite.ctx, feeCollector).AmountOf(denom).String())
	suite.Equal(
		sdkmath.LegacyNewDec(0),
		suite.app.BankKeeper.GetAllBalances(suite.ctx, distributionModuleAccount.GetAddress()).AmountOf(denom).String())

	mintCoin := sdk.NewCoin(denom, sdkmath.NewInt(100_000))
	mintCoins := sdk.Coins{mintCoin}
	feeCollectorAccount := suite.app.AccountKeeper.GetModuleAccount(suite.ctx, authtypes.FeeCollectorName)
	suite.Require().NotNil(feeCollectorAccount)

	suite.Require().NoError(FundModuleAccount(suite.app.BankKeeper, suite.ctx, feeCollectorAccount.GetName(), mintCoins))
	feeCollector = suite.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	suite.Equal(
		mintCoin.Amount.String(),
		suite.app.BankKeeper.GetAllBalances(suite.ctx, feeCollector).AmountOf(denom).String())

	suite.Equal(
		sdkmath.LegacyNewDec(0),
		suite.app.BankKeeper.GetAllBalances(suite.ctx, distributionModuleAccount.GetAddress()).AmountOf(denom).String(),
	)

	_ = allocKeeper.DistributeMintedCoins(suite.ctx)

	feeCollector = suite.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	modulePortion := params.DistributionProportions.NftIncentives.
		Add(params.DistributionProportions.DeveloperRewards).
		Add(params.DistributionProportions.NodeHostsIncentives).
		Add(params.DistributionProportions.CommunityPool) // 40%

	// remaining going to next module should be 100% - 40% = 60%
	suite.Equal(
		sdkmath.LegacyNewDecFromInt(mintCoin.Amount).Mul(sdkmath.LegacyNewDecWithPrec(100, 2).Sub(modulePortion)).RoundInt().String(),
		suite.app.BankKeeper.GetAllBalances(suite.ctx, feeCollector).AmountOf(denom).String())

	suite.Equal(
		sdkmath.LegacyNewDecFromInt(mintCoin.Amount).Mul(params.DistributionProportions.DeveloperRewards).TruncateInt(),
		suite.app.BankKeeper.GetBalance(suite.ctx, devRewardsReceiver, denom).Amount)

	// since the NFT incentives are not setup yet, funds go into the community pool
	distributionModuleAccount = suite.app.DistrKeeper.GetDistributionAccount(suite.ctx)
	communityPoolPortion := params.DistributionProportions.CommunityPool // 5%

	suite.Equal(
		sdkmath.LegacyNewDecFromInt(mintCoin.Amount).Mul(communityPoolPortion),
		suite.app.BankKeeper.GetAllBalances(suite.ctx, distributionModuleAccount.GetAddress()).AmountOf(denom).String())
}
