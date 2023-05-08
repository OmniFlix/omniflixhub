package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagNftDenomId       = "nft-denom-id"
	FlagNftId            = "nft-id"
	FlagName             = "name"
	FlagDescription      = "description"
	FlagMaxAllowedClaims = "max-allowed-claims"
	FlagTokensPerClaim   = "tokens-per-claim"
	FlagDeposit          = "deposit"
	FlagCreator          = "creator"
	FlagStatus           = "status"
	FlagClaimer          = "claimer"
	FlagStartTime        = "start-time"
	FlagDuration         = "duration"
	FlagAmount           = "amount"
	FlagInteractionType  = "interaction-type"
	FlagClaimType        = "claim-type"
	FlagDistributionType = "distribution-type"
	FlagStreamDuration   = "stream-duration"
	FlagNftDetailsFile   = "nft-details-file"
	FlagCreationFee      = "creation-fee"
)

var (
	FsCreateCampaign  = flag.NewFlagSet("", flag.ContinueOnError)
	FsDepositCampaign = flag.NewFlagSet("", flag.ContinueOnError)
	FsClaim           = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsCreateCampaign.String(FlagName, "", "name of the campaign")
	FsCreateCampaign.String(FlagDescription, "", "description of the campaign")
	FsCreateCampaign.String(FlagInteractionType, "", "interaction to claim")
	FsCreateCampaign.String(FlagClaimType, "", "type of claim")
	FsCreateCampaign.String(FlagNftDenomId, "", "nft denom id")
	FsCreateCampaign.Uint64(FlagMaxAllowedClaims, 0, "maximum allowed claims for campaign")
	FsCreateCampaign.String(FlagTokensPerClaim, "", "tokens per claim")
	FsCreateCampaign.String(FlagDeposit, "", "tokens to deposit the campaign")
	FsCreateCampaign.String(FlagStartTime, "", "campaign start time")
	FsCreateCampaign.String(FlagDuration, "", "campaign duration")
	FsCreateCampaign.String(FlagDistributionType, "", "type of distribution")
	FsCreateCampaign.String(FlagStreamDuration, "", "claimed amount distribution duration")
	FsCreateCampaign.String(FlagNftDetailsFile, "", "nft details file")
	FsCreateCampaign.String(FlagCreationFee, "", "creation fee")

	FsDepositCampaign.String(FlagAmount, "", "deposit amount")

	FsClaim.String(FlagNftId, "", "nft id")
	FsClaim.String(FlagInteractionType, "", "type of the interaction")
	FsClaim.String(FlagClaimer, "", "claimer address")
}
