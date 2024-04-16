package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/OmniFlix/omniflixhub/v3/x/itc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/pflag"
)

// validate nft details object
func validate(nft *types.NFTDetails) error {
	if nft.DenomId == "" {
		return fmt.Errorf("denom id is required")
	}

	if nft.Name == "" {
		return fmt.Errorf("nft name is required")
	}

	if nft.MediaUri == "" {
		return fmt.Errorf("nft media uri is required")
	}
	if nft.StartIndex <= 0 {
		return fmt.Errorf("start index is required and should be greater than 0")
	}
	if nft.NameDelimiter != "" && len(nft.NameDelimiter) > 1 {
		return fmt.Errorf("name delemeter should be a single character string")
	}
	return nil
}

func parseNftDetails(fs *pflag.FlagSet, claimType types.ClaimType) (nftDetails *types.NFTDetails, err error) {
	if claimType == types.CLAIM_TYPE_FT {
		return nftDetails, nil
	}
	nftDetailsFile, err := fs.GetString(FlagNftDetailsFile)
	if err != nil {
		return nil, err
	}
	if nftDetailsFile == "" {
		return nil, fmt.Errorf("file path not provided")
	}
	contents, err := os.ReadFile(nftDetailsFile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(contents, &nftDetails)
	if err != nil {
		return nil, err
	}

	if err := validate(nftDetails); err != nil {
		return nil, err
	}

	return nftDetails, nil
}

func parseCampaignNameAndDescription(fs *pflag.FlagSet) (name, description string, err error) {
	name, err = fs.GetString(FlagName)
	if err != nil {
		return name, description, err
	}
	description, err = fs.GetString(FlagDescription)
	if err != nil {
		return name, description, err
	}
	return name, description, err
}

func parseClaimType(fs *pflag.FlagSet) (types.ClaimType, error) {
	claimTypeStr, err := fs.GetString(FlagClaimType)
	if err != nil {
		return -1, err
	}
	switch strings.ToLower(claimTypeStr) {
	case "fungible":
		return types.CLAIM_TYPE_FT, nil
	case "non-fungible":
		return types.CLAIM_TYPE_NFT, nil
	case "fungible-and-non-fungible":
		return types.CLAIM_TYPE_FT_AND_NFT, nil
	default:
		return -1, fmt.Errorf("invalid claim type")
	}
}

func parseInteractionType(fs *pflag.FlagSet) (types.InteractionType, error) {
	interactTypeStr, err := fs.GetString(FlagInteractionType)
	if err != nil {
		return -1, err
	}
	switch strings.ToLower(interactTypeStr) {
	case "hold":
		return types.INTERACTION_TYPE_HOLD, nil
	case "transfer":
		return types.INTERACTION_TYPE_TRANSFER, nil
	case "burn":
		return types.INTERACTION_TYPE_BURN, nil
	default:
		return -1, fmt.Errorf("invalid interaction type")
	}
}

func parseStartTimeAndDuration(fs *pflag.FlagSet) (startTime time.Time, duration time.Duration, err error) {
	startTimeStr, err := fs.GetString(FlagStartTime)
	if err != nil {
		return startTime, duration, err
	}
	startTime, err = time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		return startTime, duration, fmt.Errorf("unable to parse start time: %s", startTimeStr)
	}

	durationStr, err := fs.GetString(FlagDuration)
	if err != nil {
		return startTime, duration, err
	}
	duration, err = time.ParseDuration(durationStr)
	if err != nil {
		return startTime, duration, fmt.Errorf("unable to parse duration: %s", durationStr)
	}
	return startTime, duration, err
}

func parseDistribution(fs *pflag.FlagSet, claimType types.ClaimType) (*types.Distribution, error) {
	var distribution *types.Distribution
	if claimType == types.CLAIM_TYPE_NFT {
		return distribution, nil
	}
	distributionTypeStr, err := fs.GetString(FlagDistributionType)
	if err != nil {
		return distribution, err
	}
	switch strings.ToLower(distributionTypeStr) {
	case "stream":
		streamDurationStr, err := fs.GetString(FlagStreamDuration)
		if err != nil {
			return distribution, err
		}
		streamDuration, err := time.ParseDuration(streamDurationStr)
		if err != nil {
			return distribution, err
		}
		distribution = &types.Distribution{
			Type:           types.DISTRIBUTION_TYPE_STREAM,
			StreamDuration: streamDuration,
		}
		return distribution, nil
	case "instant":
		distribution = &types.Distribution{
			Type: types.DISTRIBUTION_TYPE_INSTANT,
		}
		return distribution, nil
	default:
		return distribution, fmt.Errorf("invalid distribution type")
	}
}

func parseCampaignTokens(
	fs *pflag.FlagSet,
	claimType types.ClaimType,
) (tokenPerClaim sdk.Coin, tokensDeposited sdk.Coin, err error) {
	if claimType == types.CLAIM_TYPE_NFT {
		return tokenPerClaim, tokensDeposited, err
	}
	tokensPerClaimStr, err := fs.GetString(FlagTokensPerClaim)
	if err != nil {
		return tokenPerClaim, tokensDeposited, err
	}
	tokenPerClaim, err = sdk.ParseCoinNormalized(tokensPerClaimStr)
	if err != nil {
		return tokenPerClaim, tokensDeposited, err
	}

	tokensDepositStr, err := fs.GetString(FlagDeposit)
	if err != nil {
		return tokenPerClaim, tokensDeposited, err
	}
	tokensDeposited, err = sdk.ParseCoinNormalized(tokensDepositStr)
	if err != nil {
		return tokenPerClaim, tokensDeposited, err
	}

	return tokenPerClaim, tokensDeposited, err
}

func parseCampaignStatus(fs *pflag.FlagSet) (campaignStatus types.CampaignStatus, err error) {
	status, err := fs.GetString(FlagStatus)
	if err != nil {
		return campaignStatus, err
	}
	switch strings.ToLower(status) {
	case "active":
		return types.CAMPAIGN_STATUS_ACTIVE, nil
	case "inactive":
		return types.CAMPAIGN_STATUS_INACTIVE, nil
	default:
		return types.CAMPAIGN_STATUS_UNSPECIFIED, nil
	}
}
