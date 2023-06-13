package keeper_test

import (
	onftkeeper "github.com/OmniFlix/onft/keeper"
	onfttypes "github.com/OmniFlix/onft/types"
	"testing"
	"time"

	"github.com/OmniFlix/omniflixhub/app/apptesting"
	"github.com/OmniFlix/omniflixhub/x/itc/keeper"
	"github.com/OmniFlix/omniflixhub/x/itc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper

	queryClient  types.QueryClient
	msgServer    types.MsgServer
	nftMsgServer onfttypes.MsgServer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

var (
	defaultCampaignName        = "test campaign"
	defaultCampaignDescription = "test description"
	defaultInteractionType     = types.INTERACTION_TYPE_BURN
	defaultClaimType           = types.CLAIM_TYPE_FT
	defaultNftDenomId          = "onftdenomtest001"
	defaultNftMintDenomId      = "onftdenomtest002"
	defaultMaxClaims           = uint64(10)
	defaultTokensPerClaim      = sdk.NewInt64Coin(types.DefaultCampaignCreationFee.Denom, 10_000_000)
	defaultDuration            = time.Second * 100
	defaultDistribution        = types.Distribution{Type: types.DISTRIBUTION_TYPE_STREAM, StreamDuration: time.Second * 300}
	defaultNftMintDetails      = types.NFTDetails{
		DenomId:      defaultNftMintDenomId,
		Name:         "test mint nft",
		Description:  "test mint nft description",
		MediaUri:     "ipfs://minttesturi",
		PreviewUri:   "ipfs://minttestpreviewuri",
		Data:         "{}",
		RoyaltyShare: sdk.NewDecWithPrec(1, 2),
		Transferable: true,
		Extensible:   true,
		Nsfw:         false,
	}
)

func (suite *KeeperTestSuite) SetupTest() {
	suite.Setup()

	// Fund every TestAcc with some tokens
	fundAccsAmount := sdk.NewCoins(sdk.NewCoin(types.DefaultCampaignCreationFee.Denom, types.DefaultCampaignCreationFee.Amount.MulRaw(100)))
	for _, acc := range suite.TestAccs {
		suite.FundAcc(acc, fundAccsAmount)
	}

	suite.queryClient = types.NewQueryClient(suite.QueryHelper)
	suite.msgServer = keeper.NewMsgServerImpl(suite.App.ItcKeeper)
	suite.nftMsgServer = onftkeeper.NewMsgServerImpl(suite.App.ONFTKeeper)
}

func (suite *KeeperTestSuite) CreateDefaultCampaign() {
	suite.createDefaultNftDenom()

	_, _ = suite.msgServer.CreateCampaign(
		sdk.WrapSDKContext(suite.Ctx),
		types.NewMsgCreateCampaign(
			defaultCampaignName,
			defaultCampaignDescription,
			defaultInteractionType,
			defaultClaimType,
			defaultNftDenomId,
			defaultMaxClaims,
			defaultTokensPerClaim,
			sdk.NewInt64Coin(
				defaultTokensPerClaim.Denom,
				defaultTokensPerClaim.Amount.MulRaw(int64(defaultMaxClaims)).Int64(),
			),
			nil,
			nil,
			time.Now().Add(time.Second*5),
			defaultDuration,
			suite.TestAccs[0].String(),
			types.DefaultCampaignCreationFee,
		),
	)
}

func (suite *KeeperTestSuite) createDefaultNftDenom() {
	createDenomMsg := onfttypes.NewMsgCreateDenom(
		"test11",
		"test11",
		"{}",
		"test description",
		"ipfs://testpreviewuri",
		suite.TestAccs[0].String(),
		onfttypes.DefaultDenomCreationFee,
	)
	createDenomMsg.Id = defaultNftDenomId

	_, _ = suite.nftMsgServer.CreateDenom(
		sdk.WrapSDKContext(suite.Ctx),
		createDenomMsg,
	)
}

func (suite *KeeperTestSuite) createDefaultMintNftDenom() {
	createDenomMsg := onfttypes.NewMsgCreateDenom(
		"test22",
		"test22",
		"{}",
		"test description",
		"ipfs://testpreviewuri",
		suite.TestAccs[0].String(),
		onfttypes.DefaultDenomCreationFee,
	)
	createDenomMsg.Id = defaultNftMintDenomId

	_, _ = suite.nftMsgServer.CreateDenom(
		sdk.WrapSDKContext(suite.Ctx),
		createDenomMsg,
	)
}
