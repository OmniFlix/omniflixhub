package keeper_test

import (
	"fmt"
	"time"

	"github.com/OmniFlix/omniflixhub/x/itc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestCreateCampaign() {
	// create default nft denom
	suite.createDefaultNftDenom()
	suite.createDefaultMintNftDenom()

	for _, tc := range []struct {
		desc                  string
		campaignName          string
		campaignDescription   string
		nftDenomId            string
		interactionType       types.InteractionType
		claimType             types.ClaimType
		startTime             time.Time
		duration              time.Duration
		maxAllowedClaims      uint64
		tokensPerClaim        sdk.Coin
		deposit               sdk.Coin
		nftMintDetails        *types.NFTDetails
		distribution          *types.Distribution
		creator               string
		creationFee           sdk.Coin
		valid                 bool
		expectedMessageEvents int
	}{
		{
			desc:                  "valid default create campaign",
			nftDenomId:            defaultNftDenomId,
			campaignName:          defaultCampaignName,
			campaignDescription:   defaultCampaignDescription,
			interactionType:       types.INTERACTION_TYPE_BURN,
			claimType:             types.CLAIM_TYPE_FT,
			startTime:             time.Now(),
			duration:              time.Hour,
			maxAllowedClaims:      10,
			tokensPerClaim:        sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10),
			deposit:               sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10),
			nftMintDetails:        nil,
			distribution:          &defaultDistribution,
			creator:               suite.TestAccs[0].String(),
			creationFee:           types.DefaultCampaignCreationFee,
			valid:                 true,
			expectedMessageEvents: 1,
		},
		{
			desc:                  "valid transfer interaction create campaign",
			nftDenomId:            defaultNftDenomId,
			campaignName:          defaultCampaignName,
			campaignDescription:   defaultCampaignDescription,
			interactionType:       types.INTERACTION_TYPE_TRANSFER,
			claimType:             types.CLAIM_TYPE_FT,
			startTime:             time.Now(),
			duration:              time.Hour,
			maxAllowedClaims:      10,
			tokensPerClaim:        sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10),
			deposit:               sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10),
			nftMintDetails:        nil,
			distribution:          &defaultDistribution,
			creator:               suite.TestAccs[0].String(),
			creationFee:           types.DefaultCampaignCreationFee,
			valid:                 true,
			expectedMessageEvents: 1,
		},
		{
			desc:                  "valid ft & nft claim type  campaign",
			nftDenomId:            defaultNftDenomId,
			campaignName:          defaultCampaignName,
			campaignDescription:   defaultCampaignDescription,
			interactionType:       types.INTERACTION_TYPE_TRANSFER,
			claimType:             types.CLAIM_TYPE_FT_AND_NFT,
			startTime:             time.Now(),
			duration:              time.Hour,
			maxAllowedClaims:      10,
			tokensPerClaim:        sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10),
			deposit:               sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 100),
			nftMintDetails:        &defaultNftMintDetails,
			distribution:          &defaultDistribution,
			creator:               suite.TestAccs[0].String(),
			creationFee:           types.DefaultCampaignCreationFee,
			valid:                 true,
			expectedMessageEvents: 1,
		},
		{
			desc:                  "nft claim type without nft mint details",
			nftDenomId:            defaultNftDenomId,
			campaignName:          defaultCampaignName,
			campaignDescription:   defaultCampaignDescription,
			interactionType:       types.INTERACTION_TYPE_BURN,
			claimType:             types.CLAIM_TYPE_NFT,
			startTime:             time.Now(),
			duration:              time.Hour,
			maxAllowedClaims:      10,
			tokensPerClaim:        sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10),
			deposit:               sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10),
			nftMintDetails:        nil,
			distribution:          &defaultDistribution,
			creator:               suite.TestAccs[0].String(),
			creationFee:           types.DefaultCampaignCreationFee,
			valid:                 false,
			expectedMessageEvents: 0,
		},
		{
			desc:                  "invalid case nil distribution",
			nftDenomId:            defaultNftDenomId,
			campaignName:          defaultCampaignName,
			campaignDescription:   defaultCampaignDescription,
			interactionType:       types.INTERACTION_TYPE_BURN,
			claimType:             types.CLAIM_TYPE_FT,
			startTime:             time.Now(),
			duration:              time.Hour,
			maxAllowedClaims:      10,
			tokensPerClaim:        sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 0),
			deposit:               sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 0),
			nftMintDetails:        nil,
			distribution:          nil,
			creator:               suite.TestAccs[0].String(),
			creationFee:           types.DefaultCampaignCreationFee,
			valid:                 false,
			expectedMessageEvents: 0,
		},
		{
			desc:                  "invalid case - past start time",
			nftDenomId:            defaultNftDenomId,
			campaignName:          defaultCampaignName,
			campaignDescription:   defaultCampaignDescription,
			interactionType:       types.INTERACTION_TYPE_BURN,
			claimType:             types.CLAIM_TYPE_FT,
			startTime:             time.Now().Add(-time.Hour),
			duration:              time.Hour,
			maxAllowedClaims:      10,
			tokensPerClaim:        sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10),
			deposit:               sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10),
			nftMintDetails:        nil,
			distribution:          &defaultDistribution,
			creator:               suite.TestAccs[0].String(),
			creationFee:           types.DefaultCampaignCreationFee,
			valid:                 false,
			expectedMessageEvents: 0,
		},
		{
			desc:                  "invalid case - not enough creation fee",
			nftDenomId:            defaultNftDenomId,
			campaignName:          defaultCampaignName,
			campaignDescription:   defaultCampaignDescription,
			interactionType:       types.INTERACTION_TYPE_BURN,
			claimType:             types.CLAIM_TYPE_FT,
			startTime:             time.Now(),
			duration:              time.Hour,
			maxAllowedClaims:      10,
			tokensPerClaim:        sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10),
			deposit:               sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10),
			nftMintDetails:        nil,
			distribution:          &defaultDistribution,
			creator:               suite.TestAccs[0].String(),
			creationFee:           types.DefaultCampaignCreationFee.SubAmount(sdk.NewInt(1_000_000)),
			valid:                 false,
			expectedMessageEvents: 0,
		},
		{
			desc:                  "invalid case - denoms mismatch",
			nftDenomId:            defaultNftDenomId,
			campaignName:          defaultCampaignName,
			campaignDescription:   defaultCampaignDescription,
			interactionType:       types.INTERACTION_TYPE_BURN,
			claimType:             types.CLAIM_TYPE_FT,
			startTime:             time.Now(),
			duration:              time.Hour,
			maxAllowedClaims:      10,
			tokensPerClaim:        sdk.NewInt64Coin("test", 10),
			deposit:               sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10),
			nftMintDetails:        nil,
			distribution:          &defaultDistribution,
			creator:               suite.TestAccs[0].String(),
			creationFee:           types.DefaultCampaignCreationFee,
			valid:                 false,
			expectedMessageEvents: 0,
		},
		{
			desc:                  "invalid case - invalid duration",
			nftDenomId:            defaultNftDenomId,
			campaignName:          defaultCampaignName,
			campaignDescription:   defaultCampaignDescription,
			interactionType:       types.INTERACTION_TYPE_BURN,
			claimType:             types.CLAIM_TYPE_FT,
			startTime:             time.Now(),
			duration:              -time.Second,
			maxAllowedClaims:      10,
			tokensPerClaim:        sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10),
			deposit:               sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10),
			nftMintDetails:        nil,
			distribution:          &defaultDistribution,
			creator:               suite.TestAccs[0].String(),
			creationFee:           types.DefaultCampaignCreationFee,
			valid:                 false,
			expectedMessageEvents: 0,
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
			suite.AssertEventEmitted(ctx, types.EventTypeCreateCampaign, tc.expectedMessageEvents)
		})
	}
}
