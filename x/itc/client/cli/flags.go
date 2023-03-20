package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagNftDenomId       = "denom-id"
	FlagNftId            = "nft-id"
	FlagName             = "name"
	FlagDescription      = "description"
	FlagMaxAllowedClaims = "max-allowed-claims"
	FlagClaimableTokens  = "claimable-tokens"
	FlagTotalTokens      = "total-tokens"
	FlagCreator          = "creator"
	FlagStatus           = "status"
	FlagClaimer          = "claimer"
	FlagCoinDenom        = "coin-denom"
	Flag
	FlagStartTime        = "start-time"
	FlagDuration         = "duration"
	FlagAmount           = "amount"
	FlagInteractionType  = "interaction-type"
	FlagClaimType        = "claim-type"
	FlagDistributionType = "distribution-type"
)

var (
	FsCreateCampaign = flag.NewFlagSet("", flag.ContinueOnError)
	FsClaim          = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsCreateCampaign.String(FlagName, "", "name of the campaign")
	FsCreateCampaign.String(FlagDescription, "", "description of the campaign")
	FsCreateCampaign.String(FlagInteractionType, "", "interaction to claim")
	FsCreateCampaign.String(FlagClaimType, "", "type of claim")
	FsCreateCampaign.String(FlagNftDenomId, "", "nft denom id")
	FsCreateCampaign.Uint64(FlagMaxAllowedClaims, 0, "maximum allowed claims for campaign")
	FsCreateCampaign.String(FlagClaimableTokens, "", "tokens per claim")
	FsCreateCampaign.String(FlagTotalTokens, "", "tokens to deposit the campaign")
	FsCreateCampaign.String(FlagStartTime, "", "auction start time")
	FsCreateCampaign.String(FlagDuration, "", "auction duration")

	FsClaim.String(FlagNftId, "", "nft id")
	FsClaim.String(FlagInteractionType, "", "type of the interaction")
}
