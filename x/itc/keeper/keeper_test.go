package keeper_test

import (
	"fmt"
	"testing"
	"time"

	onftkeeper "github.com/OmniFlix/omniflixhub/v2/x/onft/keeper"
	onfttypes "github.com/OmniFlix/omniflixhub/v2/x/onft/types"

	"github.com/OmniFlix/omniflixhub/v2/app/apptesting"
	"github.com/OmniFlix/omniflixhub/v2/x/itc/keeper"
	"github.com/OmniFlix/omniflixhub/v2/x/itc/types"
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
		"ipfs://testuri",
		"ipfs://testUriHash",
		"ipfs://testpreviewuri",
		"",
		suite.TestAccs[0].String(),
		onfttypes.DefaultDenomCreationFee,
		[]*onfttypes.WeightedAddress{
			{
				Address: suite.TestAccs[0].String(),
				Weight:  sdk.OneDec(),
			},
		},
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
		"ipfs://testuri",
		"ipfs://testUriHash",
		"ipfs://testpreviewuri",
		"",
		suite.TestAccs[0].String(),
		onfttypes.DefaultDenomCreationFee,
		nil,
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
		"ipfs://testuri",
		"ipfs://testUriHash",
		"ipfs://testpreviewuri",
		"",
		suite.TestAccs[0].String(),
		onfttypes.DefaultDenomCreationFee,
		[]*onfttypes.WeightedAddress{
			{
				Address: suite.TestAccs[0].String(),
				Weight:  sdk.OneDec(),
			},
		},
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

func (suite *KeeperTestSuite) TestGetClaims() {
	suite.SetupTest()
	sdkCtx := suite.Ctx
	defaultAddress := suite.TestAccs[0]
	keeper := suite.App.ItcKeeper

	claimsToSet := []types.Claim{
		types.NewClaim(
			1,
			defaultAddress.String(),
			defaultNftId,
			defaultInteractionType,
		),
		types.NewClaim(
			2,
			defaultAddress.String(),
			defaultNftId,
			defaultInteractionType,
		),
		types.NewClaim(
			3,
			defaultAddress.String(),
			defaultNftId,
			defaultInteractionType,
		),
	}

	for _, claimToSet := range claimsToSet {
		claimToSet := claimToSet
		keeper.SetClaim(sdkCtx, claimToSet)

		got := keeper.GetClaims(sdkCtx, claimToSet.CampaignId)
		suite.Require().Equal(len(got), 1)
		suite.Require().Equal(claimToSet, got[0])
	}
}

func (suite *KeeperTestSuite) TestFinalizeAndEndCampaigns() {
	suite.SetupTest()
	sdkCtx := suite.Ctx.WithBlockTime(suite.Ctx.BlockTime().Add(defaultDuration * 2))
	keeper := suite.App.ItcKeeper

	suite.CreateDefaultCampaign()
	suite.CreateSecondaryCampaign()

	keeper.FinalizeAndEndCampaigns(sdkCtx)

	campaigns := keeper.GetAllCampaigns(sdkCtx)
	suite.Require().Empty(campaigns)
}

func (suite *KeeperTestSuite) TestHasCampaign() {
	suite.SetupTest()
	sdkCtx := suite.Ctx
	keeper := suite.App.ItcKeeper

	suite.CreateDefaultCampaign()
	suite.Require().True(keeper.HasCampaign(sdkCtx, 1))

	suite.CreateSecondaryCampaign()
	suite.Require().True(keeper.HasCampaign(sdkCtx, 2))

	sdkCtx = sdkCtx.WithBlockTime(sdkCtx.BlockTime().Add(defaultDuration * 2))
	keeper.FinalizeAndEndCampaigns(sdkCtx)

	suite.Require().False(keeper.HasCampaign(sdkCtx, 1))
	suite.Require().False(keeper.HasCampaign(sdkCtx, 2))
}

func (suite *KeeperTestSuite) TestParams() {
	testCases := []struct {
		name      string
		input     types.Params
		expectErr bool
	}{
		{
			name: "set invalid max campaign duration",
			input: types.Params{
				MaxCampaignDuration: -1,
				CreationFee:         types.DefaultCampaignCreationFee,
			},
			expectErr: true,
		},
		{
			name: "set invalid creation fee",
			input: types.Params{
				MaxCampaignDuration: types.DefaultMaxCampaignDuration,
				CreationFee:         sdk.Coin{},
			},
			expectErr: true,
		},
		{
			name: "set full valid params",
			input: types.Params{
				CreationFee:         types.DefaultCampaignCreationFee,
				MaxCampaignDuration: types.DefaultMaxCampaignDuration,
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			expected := suite.App.ItcKeeper.GetParams(suite.Ctx)
			err := suite.App.ItcKeeper.SetParams(suite.Ctx, tc.input)
			if tc.expectErr {
				suite.Require().Error(err)
			} else {
				expected = tc.input
				suite.Require().NoError(err)
			}

			p := suite.App.ItcKeeper.GetParams(suite.Ctx)
			suite.Require().Equal(expected, p)
		})
	}
}
