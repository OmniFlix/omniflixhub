package keeper_test

import (
	"testing"

	"github.com/OmniFlix/omniflixhub/v5/app/apptesting"
	"github.com/OmniFlix/omniflixhub/v5/x/onft/keeper"
	"github.com/OmniFlix/omniflixhub/v5/x/onft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper

	queryClient types.QueryClient
	msgServer   types.MsgServer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	// Fund every TestAcc with some tokens
	fundAccsAmount := sdk.NewCoins(
		sdk.NewCoin(types.DefaultDenomCreationFee.Denom, types.DefaultDenomCreationFee.Amount.MulRaw(100)),
	)
	for _, acc := range suite.TestAccs {
		suite.FundAcc(acc, fundAccsAmount)
	}
	suite.queryClient = types.NewQueryClient(suite.QueryHelper)
	suite.msgServer = keeper.NewMsgServerImpl(suite.App.ONFTKeeper)
}
