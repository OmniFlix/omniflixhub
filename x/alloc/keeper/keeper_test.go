package keeper_test

import (
	"testing"
	"time"

	"github.com/OmniFlix/omniflixhub/app"
	"github.com/OmniFlix/omniflixhub/testutil/simapp"
	"github.com/OmniFlix/omniflixhub/x/alloc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

type KeeperTestSuite struct {
	suite.Suite
	ctx sdk.Context
	app *app.OmniFlixApp
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = simapp.New(suite.T().TempDir())
	suite.ctx = suite.app.BaseApp.NewContext(
		false,
		tmproto.Header{Height: 1, ChainID: "omniflixhub-1", Time: time.Now().UTC()},
	)
	suite.app.AllocKeeper.SetParams(suite.ctx, types.DefaultParams())
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func FundAccount(bankKeeper bankkeeper.Keeper, ctx sdk.Context, addr sdk.AccAddress, amounts sdk.Coins) error {
	if err := bankKeeper.MintCoins(ctx, minttypes.ModuleName, amounts); err != nil {
		return err
	}
	return bankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, amounts)
}

func FundModuleAccount(bankKeeper bankkeeper.Keeper, ctx sdk.Context, recipientMod string, amounts sdk.Coins) error {
	if err := bankKeeper.MintCoins(ctx, minttypes.ModuleName, amounts); err != nil {
		return err
	}
	return bankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, recipientMod, amounts)
}

func (suite *KeeperTestSuite) TestDistribution() {
	suite.SetupTest()

	denom := suite.app.StakingKeeper.BondDenom(suite.ctx)
	allocKeeper := suite.app.AllocKeeper
	params := suite.app.AllocKeeper.GetParams(suite.ctx)
	nftIncentivesReceiver := sdk.AccAddress([]byte("addr1a--------------"))
	nodeHostIncentivesReceiver := sdk.AccAddress([]byte("addr1b--------------"))
	devRewardsReceiver := sdk.AccAddress([]byte("addr1c--------------"))
	params.DistributionProportions.StakingRewards = sdk.NewDecWithPrec(60, 2)
	params.DistributionProportions.NodeHostsIncentives = sdk.NewDecWithPrec(5, 2)
	params.DistributionProportions.NftIncentives = sdk.NewDecWithPrec(15, 2)
	params.DistributionProportions.CommunityPool = sdk.NewDecWithPrec(5, 2)
	params.DistributionProportions.DeveloperRewards = sdk.NewDecWithPrec(15, 2)
	params.WeightedNftIncentivesReceivers = []types.WeightedAddress{
		{
			Address: nftIncentivesReceiver.String(),
			Weight:  sdk.NewDec(1),
		},
	}
	params.WeightedNodeHostsIncentivesReceivers = []types.WeightedAddress{
		{
			Address: nodeHostIncentivesReceiver.String(),
			Weight:  sdk.NewDec(1),
		},
	}
	params.WeightedDeveloperRewardsReceivers = []types.WeightedAddress{
		{
			Address: devRewardsReceiver.String(),
			Weight:  sdk.NewDec(1),
		},
	}
	suite.app.AllocKeeper.SetParams(suite.ctx, params)

	feePool := suite.app.DistrKeeper.GetFeePool(suite.ctx)
	feeCollector := suite.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	suite.Equal(
		"0",
		suite.app.BankKeeper.GetAllBalances(suite.ctx, feeCollector).AmountOf(denom).String())
	suite.Equal(
		sdk.NewDec(0),
		feePool.CommunityPool.AmountOf(denom))

	mintCoin := sdk.NewCoin(denom, sdk.NewInt(100_000))
	mintCoins := sdk.Coins{mintCoin}
	feeCollectorAccount := suite.app.AccountKeeper.GetModuleAccount(suite.ctx, authtypes.FeeCollectorName)
	suite.Require().NotNil(feeCollectorAccount)

	suite.Require().NoError(FundModuleAccount(suite.app.BankKeeper, suite.ctx, feeCollectorAccount.GetName(), mintCoins))

	feeCollector = suite.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	suite.Equal(
		mintCoin.Amount.String(),
		suite.app.BankKeeper.GetAllBalances(suite.ctx, feeCollector).AmountOf(denom).String())

	suite.Equal(
		sdk.NewDec(0),
		feePool.CommunityPool.AmountOf(denom))

	_ = allocKeeper.DistributeMintedCoins(suite.ctx)

	feeCollector = suite.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	modulePortion := params.DistributionProportions.NftIncentives.
		Add(params.DistributionProportions.DeveloperRewards).
		Add(params.DistributionProportions.NodeHostsIncentives).
		Add(params.DistributionProportions.CommunityPool) // 40%

	// remaining going to next module should be 100% - 40% = 60%
	suite.Equal(
		mintCoin.Amount.ToDec().Mul(sdk.NewDecWithPrec(100, 2).Sub(modulePortion)).RoundInt().String(),
		suite.app.BankKeeper.GetAllBalances(suite.ctx, feeCollector).AmountOf(denom).String())

	suite.Equal(
		mintCoin.Amount.ToDec().Mul(params.DistributionProportions.DeveloperRewards).TruncateInt(),
		suite.app.BankKeeper.GetBalance(suite.ctx, devRewardsReceiver, denom).Amount)

	// since the NFT incentives are not setup yet, funds go into the community pool
	feePool = suite.app.DistrKeeper.GetFeePool(suite.ctx)
	communityPoolPortion := params.DistributionProportions.CommunityPool // 5%

	suite.Equal(
		mintCoin.Amount.ToDec().Mul(communityPoolPortion),
		feePool.CommunityPool.AmountOf(denom))
}
