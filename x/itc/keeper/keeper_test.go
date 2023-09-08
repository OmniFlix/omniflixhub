package keeper_test

import (
	"fmt"
	"testing"
	"time"

	onftkeeper "github.com/OmniFlix/onft/keeper"
	onfttypes "github.com/OmniFlix/onft/types"

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
	defaultNftId               = "onfttest"
	defaultNftMintDenomId      = "onftdenomtest010"
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

	secondaryCampaignName        = "test secondary campaign"
	secondaryCampaignDescription = "test secondary description"
	secondaryNftDenomId          = "onftdenomtest002"
	secondaryNftId               = "onfttest2"
	secondaryCampaignClaimType   = types.CLAIM_TYPE_FT_AND_NFT
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
	_, err := suite.msgServer.CreateCampaign(
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
			&defaultDistribution,
			suite.Ctx.BlockTime(),
			defaultDuration,
			suite.TestAccs[0].String(),
			types.DefaultCampaignCreationFee,
		),
	)

	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) CreateSecondaryCampaign() {
	suite.createSecondaryNftDenom()
	suite.createDefaultMintNftDenom()

	_, _ = suite.msgServer.CreateCampaign(
		sdk.WrapSDKContext(suite.Ctx),
		types.NewMsgCreateCampaign(
			secondaryCampaignName,
			secondaryCampaignDescription,
			types.INTERACTION_TYPE_TRANSFER,
			secondaryCampaignClaimType,
			secondaryNftDenomId,
			defaultMaxClaims,
			defaultTokensPerClaim,
			sdk.NewInt64Coin(
				defaultTokensPerClaim.Denom,
				defaultTokensPerClaim.Amount.MulRaw(int64(defaultMaxClaims)).Int64(),
			),
			&defaultNftMintDetails,
			&defaultDistribution,
			suite.Ctx.BlockTime(),
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

func (suite *KeeperTestSuite) createSecondaryNftDenom() {
	createDenomMsg := onfttypes.NewMsgCreateDenom(
		"test12",
		"test12",
		"{}",
		"test description",
		"ipfs://testpreviewuri",
		suite.TestAccs[0].String(),
		onfttypes.DefaultDenomCreationFee,
	)
	createDenomMsg.Id = secondaryNftDenomId

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

func (suite *KeeperTestSuite) mintNFT(denomId, nftId string) {
	mintNftMsg := onfttypes.NewMsgMintONFT(
		denomId,
		suite.TestAccs[0].String(),
		suite.TestAccs[1].String(),
		onfttypes.Metadata{
			Name:        "test",
			Description: "test",
			MediaURI:    "ipfs://testuri",
			PreviewURI:  "ipfs://testpreviewuri",
		},
		"{}",
		true,
		true,
		false,
		sdk.NewDecWithPrec(1, 2),
	)
	mintNftMsg.Id = nftId
	_, _ = suite.nftMsgServer.MintONFT(
		sdk.WrapSDKContext(suite.Ctx),
		mintNftMsg,
	)
}

func (suite *KeeperTestSuite) mintNFTs() {
	for counter := 1; counter <= 10; counter++ {
		suite.mintNFT(defaultNftDenomId, fmt.Sprintf("%s%d", defaultNftId, counter))
		suite.mintNFT(secondaryNftDenomId, fmt.Sprintf("%s%d", secondaryNftId, counter))
	}
}

func (suite *KeeperTestSuite) TestGetAllCampaigns() {
	suite.SetupTest()

	sdkCtx := suite.Ctx
	campaigns := suite.App.ItcKeeper.GetAllCampaigns(sdkCtx)
	suite.Require().Empty(campaigns)

	suite.CreateDefaultCampaign()

	campaigns = suite.App.ItcKeeper.GetAllCampaigns(sdkCtx)
	suite.Require().Equal(len(campaigns), 1)

	suite.CreateSecondaryCampaign()

	campaigns = suite.App.ItcKeeper.GetAllCampaigns(sdkCtx)
	suite.Require().Equal(len(campaigns), 2)
}
func (suite *KeeperTestSuite) TestGetCampaignByCreator() {
	suite.SetupTest()
	sdkCtx := suite.Ctx
	defaultCreator := suite.TestAccs[0]

	suite.CreateDefaultCampaign()

	campaigns := suite.App.ItcKeeper.GetCampaignsByCreator(sdkCtx, defaultCreator)
	suite.Require().Equal(len(campaigns), 1)

	suite.CreateSecondaryCampaign()

	campaigns = suite.App.ItcKeeper.GetCampaignsByCreator(sdkCtx, defaultCreator)
	suite.Require().Equal(len(campaigns), 2)
}
