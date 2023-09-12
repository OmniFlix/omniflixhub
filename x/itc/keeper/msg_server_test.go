package keeper_test

import (
	"fmt"
	"time"

	onfttypes "github.com/OmniFlix/onft/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/OmniFlix/omniflixhub/x/itc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestCreateCampaign() {
	// create default nft denom
	suite.createDefaultNftDenom()
	suite.createDefaultMintNftDenom()

	for _, tc := range []struct {
		desc                   string
		campaignName           string
		campaignDescription    string
		nftDenomId             string
		interactionType        types.InteractionType
		claimType              types.ClaimType
		startTime              time.Time
		duration               time.Duration
		maxAllowedClaims       uint64
		tokensPerClaim         sdk.Coin
		deposit                sdk.Coin
		nftMintDetails         *types.NFTDetails
		distribution           *types.Distribution
		creator                string
		creationFee            sdk.Coin
		valid                  bool
		expectedMessageEvents  int
		expectedTransferEvents int
	}{
		{
			desc:                   "valid default create campaign",
			nftDenomId:             defaultNftDenomId,
			campaignName:           defaultCampaignName,
			campaignDescription:    defaultCampaignDescription,
			interactionType:        types.INTERACTION_TYPE_BURN,
			claimType:              types.CLAIM_TYPE_FT,
			startTime:              time.Now(),
			duration:               time.Hour,
			maxAllowedClaims:       10,
			tokensPerClaim:         sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10),
			deposit:                sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10),
			nftMintDetails:         nil,
			distribution:           &defaultDistribution,
			creator:                suite.TestAccs[0].String(),
			creationFee:            types.DefaultCampaignCreationFee,
			valid:                  true,
			expectedMessageEvents:  1,
			expectedTransferEvents: 2,
		},
		{
			desc:                   "valid transfer interaction create campaign",
			nftDenomId:             defaultNftDenomId,
			campaignName:           defaultCampaignName,
			campaignDescription:    defaultCampaignDescription,
			interactionType:        types.INTERACTION_TYPE_TRANSFER,
			claimType:              types.CLAIM_TYPE_FT,
			startTime:              time.Now(),
			duration:               time.Hour,
			maxAllowedClaims:       10,
			tokensPerClaim:         sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10),
			deposit:                sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10),
			nftMintDetails:         nil,
			distribution:           &defaultDistribution,
			creator:                suite.TestAccs[0].String(),
			creationFee:            types.DefaultCampaignCreationFee,
			valid:                  true,
			expectedMessageEvents:  1,
			expectedTransferEvents: 2,
		},
		{
			desc:                   "valid ft & nft claim type  campaign",
			nftDenomId:             defaultNftDenomId,
			campaignName:           defaultCampaignName,
			campaignDescription:    defaultCampaignDescription,
			interactionType:        types.INTERACTION_TYPE_TRANSFER,
			claimType:              types.CLAIM_TYPE_FT_AND_NFT,
			startTime:              time.Now(),
			duration:               time.Hour,
			maxAllowedClaims:       10,
			tokensPerClaim:         sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10),
			deposit:                sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 100),
			nftMintDetails:         &defaultNftMintDetails,
			distribution:           &defaultDistribution,
			creator:                suite.TestAccs[0].String(),
			creationFee:            types.DefaultCampaignCreationFee,
			valid:                  true,
			expectedMessageEvents:  1,
			expectedTransferEvents: 2,
		},
		{
			desc:                   "nft claim type without nft mint details",
			nftDenomId:             defaultNftDenomId,
			campaignName:           defaultCampaignName,
			campaignDescription:    defaultCampaignDescription,
			interactionType:        types.INTERACTION_TYPE_BURN,
			claimType:              types.CLAIM_TYPE_NFT,
			startTime:              time.Now(),
			duration:               time.Hour,
			maxAllowedClaims:       10,
			tokensPerClaim:         sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10),
			deposit:                sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10),
			nftMintDetails:         nil,
			distribution:           &defaultDistribution,
			creator:                suite.TestAccs[0].String(),
			creationFee:            types.DefaultCampaignCreationFee,
			valid:                  false,
			expectedMessageEvents:  0,
			expectedTransferEvents: 0,
		},
		{
			desc:                   "invalid case nil distribution",
			nftDenomId:             defaultNftDenomId,
			campaignName:           defaultCampaignName,
			campaignDescription:    defaultCampaignDescription,
			interactionType:        types.INTERACTION_TYPE_BURN,
			claimType:              types.CLAIM_TYPE_FT,
			startTime:              time.Now(),
			duration:               time.Hour,
			maxAllowedClaims:       10,
			tokensPerClaim:         sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 0),
			deposit:                sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 0),
			nftMintDetails:         nil,
			distribution:           nil,
			creator:                suite.TestAccs[0].String(),
			creationFee:            types.DefaultCampaignCreationFee,
			valid:                  false,
			expectedMessageEvents:  0,
			expectedTransferEvents: 0,
		},
		{
			desc:                   "invalid case - past start time",
			nftDenomId:             defaultNftDenomId,
			campaignName:           defaultCampaignName,
			campaignDescription:    defaultCampaignDescription,
			interactionType:        types.INTERACTION_TYPE_BURN,
			claimType:              types.CLAIM_TYPE_FT,
			startTime:              time.Now().Add(-time.Hour),
			duration:               time.Hour,
			maxAllowedClaims:       10,
			tokensPerClaim:         sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10),
			deposit:                sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10),
			nftMintDetails:         nil,
			distribution:           &defaultDistribution,
			creator:                suite.TestAccs[0].String(),
			creationFee:            types.DefaultCampaignCreationFee,
			valid:                  false,
			expectedMessageEvents:  0,
			expectedTransferEvents: 0,
		},
		{
			desc:                   "invalid case - not enough creation fee",
			nftDenomId:             defaultNftDenomId,
			campaignName:           defaultCampaignName,
			campaignDescription:    defaultCampaignDescription,
			interactionType:        types.INTERACTION_TYPE_BURN,
			claimType:              types.CLAIM_TYPE_FT,
			startTime:              time.Now(),
			duration:               time.Hour,
			maxAllowedClaims:       10,
			tokensPerClaim:         sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10),
			deposit:                sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10),
			nftMintDetails:         nil,
			distribution:           &defaultDistribution,
			creator:                suite.TestAccs[0].String(),
			creationFee:            types.DefaultCampaignCreationFee.SubAmount(sdk.NewInt(1_000_000)),
			valid:                  false,
			expectedMessageEvents:  0,
			expectedTransferEvents: 0,
		},
		{
			desc:                   "invalid case - denoms mismatch",
			nftDenomId:             defaultNftDenomId,
			campaignName:           defaultCampaignName,
			campaignDescription:    defaultCampaignDescription,
			interactionType:        types.INTERACTION_TYPE_BURN,
			claimType:              types.CLAIM_TYPE_FT,
			startTime:              time.Now(),
			duration:               time.Hour,
			maxAllowedClaims:       10,
			tokensPerClaim:         sdk.NewInt64Coin("test", 10),
			deposit:                sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10),
			nftMintDetails:         nil,
			distribution:           &defaultDistribution,
			creator:                suite.TestAccs[0].String(),
			creationFee:            types.DefaultCampaignCreationFee,
			valid:                  false,
			expectedMessageEvents:  0,
			expectedTransferEvents: 0,
		},
		{
			desc:                   "invalid case - invalid duration",
			nftDenomId:             defaultNftDenomId,
			campaignName:           defaultCampaignName,
			campaignDescription:    defaultCampaignDescription,
			interactionType:        types.INTERACTION_TYPE_BURN,
			claimType:              types.CLAIM_TYPE_FT,
			startTime:              time.Now(),
			duration:               -time.Second,
			maxAllowedClaims:       10,
			tokensPerClaim:         sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10),
			deposit:                sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10),
			nftMintDetails:         nil,
			distribution:           &defaultDistribution,
			creator:                suite.TestAccs[0].String(),
			creationFee:            types.DefaultCampaignCreationFee,
			valid:                  false,
			expectedMessageEvents:  0,
			expectedTransferEvents: 0,
		},
	} {
		suite.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			ctx := suite.Ctx.WithEventManager(sdk.NewEventManager())
			suite.Require().Equal(0, len(ctx.EventManager().Events()))
			createCampaignMsg := types.NewMsgCreateCampaign(
				tc.campaignName,
				tc.campaignDescription,
				tc.interactionType,
				tc.claimType,
				tc.nftDenomId,
				tc.maxAllowedClaims,
				tc.tokensPerClaim,
				tc.deposit,
				tc.nftMintDetails,
				tc.distribution,
				tc.startTime,
				tc.duration,
				tc.creator,
				tc.creationFee,
			)
			// Test create campaign message
			_, err := suite.msgServer.CreateCampaign(
				sdk.WrapSDKContext(ctx),
				createCampaignMsg,
			)
			if tc.valid {
				suite.Require().NoError(err)
			}
			// Ensure current number and type of event is emitted
			suite.AssertEventEmitted(ctx, types.TypeMsgCreateCampaign, tc.expectedMessageEvents)
			suite.AssertEventEmitted(ctx, banktypes.EventTypeTransfer, tc.expectedTransferEvents)
		})
	}
}

func (suite *KeeperTestSuite) TestCancelCampaign() {
	// create default campaign
	suite.CreateDefaultCampaign()

	// create secondary campaign
	suite.CreateSecondaryCampaign()

	for _, tc := range []struct {
		desc                   string
		campaignId             uint64
		creator                string
		valid                  bool
		expectedMessageEvents  int
		expectedTransferEvents int
	}{
		{
			desc:                   "valid case - cancel default campaign",
			campaignId:             1,
			creator:                suite.TestAccs[0].String(),
			valid:                  true,
			expectedMessageEvents:  1,
			expectedTransferEvents: 1,
		},
		{
			desc:                   "invalid case - cancel campaign from different creator",
			campaignId:             2,
			creator:                suite.TestAccs[1].String(),
			valid:                  false,
			expectedMessageEvents:  0,
			expectedTransferEvents: 0,
		},
		{
			desc:                   "invalid case - cancel not existing campaign",
			campaignId:             20,
			creator:                suite.TestAccs[0].String(),
			valid:                  false,
			expectedMessageEvents:  0,
			expectedTransferEvents: 0,
		},
	} {
		suite.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			ctx := suite.Ctx.WithEventManager(sdk.NewEventManager())
			suite.Require().Equal(0, len(ctx.EventManager().Events()))
			cancelCampaignMsg := types.NewMsgCancelCampaign(
				tc.campaignId,
				tc.creator,
			)

			_, err := suite.msgServer.CancelCampaign(
				sdk.WrapSDKContext(ctx),
				cancelCampaignMsg,
			)

			if tc.valid {
				suite.Require().NoError(err)
			}
			suite.AssertEventEmitted(ctx, types.TypeMsgCancelCampaign, tc.expectedMessageEvents)
			suite.AssertEventEmitted(ctx, banktypes.EventTypeTransfer, tc.expectedTransferEvents)
		})
	}
}

func (suite *KeeperTestSuite) TestClaimCampaign() {
	suite.CreateDefaultCampaign()
	suite.CreateSecondaryCampaign()
	suite.mintNFTs()

	for _, tc := range []struct {
		desc                      string
		campaignId                uint64
		nftId                     string
		interactionType           types.InteractionType
		claimer                   string
		valid                     bool
		expectedMessageEvents     int
		expectedTransferEvents    int
		expectedNftTransferEvents int
		expectedNftMintEvents     int
		expectedNftBurnEvents     int
	}{
		{
			desc:                      "valid case - claim default campaign",
			campaignId:                1,
			nftId:                     "onfttest1",
			interactionType:           types.INTERACTION_TYPE_BURN,
			claimer:                   suite.TestAccs[1].String(),
			valid:                     true,
			expectedMessageEvents:     1,
			expectedTransferEvents:    1,
			expectedNftTransferEvents: 0,
			expectedNftMintEvents:     0,
			expectedNftBurnEvents:     1,
		},
		{
			desc:                      "valid case - claim ft & nft campaign",
			campaignId:                2,
			nftId:                     "onfttest21",
			interactionType:           types.INTERACTION_TYPE_TRANSFER,
			claimer:                   suite.TestAccs[1].String(),
			valid:                     true,
			expectedMessageEvents:     1,
			expectedTransferEvents:    1,
			expectedNftTransferEvents: 1,
			expectedNftMintEvents:     1,
			expectedNftBurnEvents:     0,
		},
		{
			desc:                      "invalid case - claim campaign with invalid nft id",
			campaignId:                1,
			nftId:                     "invalidnftid",
			interactionType:           types.INTERACTION_TYPE_BURN,
			claimer:                   suite.TestAccs[1].String(),
			valid:                     false,
			expectedMessageEvents:     0,
			expectedTransferEvents:    0,
			expectedNftTransferEvents: 0,
			expectedNftMintEvents:     0,
			expectedNftBurnEvents:     0,
		},
		{
			desc:                      "invalid case - claim campaign with different interaction type",
			campaignId:                1,
			nftId:                     "onfttest2",
			interactionType:           types.INTERACTION_TYPE_HOLD,
			claimer:                   suite.TestAccs[1].String(),
			valid:                     false,
			expectedMessageEvents:     0,
			expectedTransferEvents:    0,
			expectedNftTransferEvents: 0,
			expectedNftMintEvents:     0,
			expectedNftBurnEvents:     0,
		},
		{
			desc:                      "invalid case - claim campaign from different address",
			campaignId:                1,
			nftId:                     "onfttest2",
			interactionType:           types.INTERACTION_TYPE_BURN,
			claimer:                   suite.TestAccs[0].String(),
			valid:                     false,
			expectedMessageEvents:     0,
			expectedTransferEvents:    0,
			expectedNftTransferEvents: 0,
			expectedNftMintEvents:     0,
			expectedNftBurnEvents:     0,
		},
		{
			desc:                      "invalid case - non-existent campaignId",
			campaignId:                10,
			nftId:                     "onfttest1",
			interactionType:           types.INTERACTION_TYPE_BURN,
			claimer:                   suite.TestAccs[1].String(),
			valid:                     false,
			expectedMessageEvents:     0,
			expectedTransferEvents:    0,
			expectedNftTransferEvents: 0,
			expectedNftMintEvents:     0,
			expectedNftBurnEvents:     0,
		},
	} {
		suite.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			ctx := suite.Ctx.WithEventManager(sdk.NewEventManager())
			suite.Require().Equal(0, len(ctx.EventManager().Events()))
			claimCampaignMsg := types.NewMsgClaim(
				tc.campaignId,
				tc.nftId,
				tc.interactionType,
				tc.claimer,
			)

			_, err := suite.msgServer.Claim(
				sdk.WrapSDKContext(ctx),
				claimCampaignMsg,
			)

			if tc.valid {
				suite.Require().NoError(err)
			}

			suite.AssertEventEmitted(ctx, types.TypeMsgClaim, tc.expectedMessageEvents)
			suite.AssertEventEmitted(ctx, banktypes.EventTypeTransfer, tc.expectedTransferEvents)
			suite.AssertEventEmitted(ctx, onfttypes.EventTypeBurnONFT, tc.expectedNftBurnEvents)
			suite.AssertEventEmitted(ctx, onfttypes.EventTypeTransferONFT, tc.expectedNftTransferEvents)
		})
	}
}

func (suite *KeeperTestSuite) TestDepositCampaign() {
	suite.CreateDefaultCampaign()

	for _, tc := range []struct {
		desc                   string
		campaignId             uint64
		deposit                sdk.Coin
		depositor              string
		valid                  bool
		expectedMessageEvents  int
		expectedTransferEvents int
	}{
		{
			desc:                   "valid case - deposit default campaign",
			campaignId:             1,
			deposit:                sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10_000_000),
			depositor:              suite.TestAccs[0].String(),
			valid:                  true,
			expectedMessageEvents:  1,
			expectedTransferEvents: 1,
		},
		{
			desc:                   "invalid case - deposit with wrong token denom",
			campaignId:             1,
			deposit:                sdk.NewInt64Coin("test", 10_000_000),
			depositor:              suite.TestAccs[0].String(),
			valid:                  false,
			expectedMessageEvents:  0,
			expectedTransferEvents: 0,
		},
		{
			desc:                   "invalid case - deposit from different address",
			campaignId:             1,
			deposit:                sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10_000_000),
			depositor:              suite.TestAccs[1].String(),
			valid:                  false,
			expectedMessageEvents:  0,
			expectedTransferEvents: 0,
		},
		{
			desc:                   "invalid case - non-existent campaignId",
			campaignId:             10,
			deposit:                sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10_000_000),
			depositor:              suite.TestAccs[0].String(),
			valid:                  false,
			expectedMessageEvents:  0,
			expectedTransferEvents: 0,
		},
	} {
		suite.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			ctx := suite.Ctx.WithEventManager(sdk.NewEventManager())
			suite.Require().Equal(0, len(ctx.EventManager().Events()))
			depositCampaignMsg := types.NewMsgDepositCampaign(
				tc.campaignId,
				tc.deposit,
				tc.depositor,
			)

			_, err := suite.msgServer.DepositCampaign(
				sdk.WrapSDKContext(ctx),
				depositCampaignMsg,
			)

			if tc.valid {
				suite.Require().NoError(err)
			}

			suite.AssertEventEmitted(ctx, types.TypeMsgDepositCampaign, tc.expectedMessageEvents)
		})
	}
}
