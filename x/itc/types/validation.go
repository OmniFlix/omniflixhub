package types

import (
	"time"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ValidateCampaign checks campaign is valid or not
func ValidateCampaign(campaign Campaign) error {
	if _, err := sdk.AccAddressFromBech32(campaign.Creator); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if err := ValidateClaimType(campaign.ClaimType); err != nil {
		return err
	}
	if err := ValidateInteractionType(campaign.Interaction); err != nil {
		return err
	}
	if campaign.ClaimType != CLAIM_TYPE_NFT {
		if err := ValidateTokens(campaign.TotalTokens, campaign.TokensPerClaim); err != nil {
			return err
		}
		if err := ValidateTokens(campaign.TotalTokens, campaign.AvailableTokens); err != nil {
			return err
		}
	}
	if campaign.ClaimType == CLAIM_TYPE_NFT || campaign.ClaimType == CLAIM_TYPE_FT_AND_NFT {
		if err := validateNFTMintDetails(campaign.NftMintDetails); err != nil {
			return err
		}
	}
	if campaign.ClaimType == CLAIM_TYPE_FT || campaign.ClaimType == CLAIM_TYPE_FT_AND_NFT {
		if err := ValidateDistribution(campaign.Distribution); err != nil {
			return err
		}
	}
	if campaign.MaxAllowedClaims == 0 {
		return errorsmod.Wrapf(ErrInValidMaxAllowedClaims,
			"max allowed claims must be a positive number (%d)", campaign.MaxAllowedClaims)
	}
	return nil
}

func ValidateDuration(t interface{}) error {
	duration, ok := t.(time.Duration)
	if !ok {
		return errorsmod.Wrapf(ErrInvalidDuration, "invalid value for duration: %T", t)
	}
	if duration.Nanoseconds() <= 0 {
		return errorsmod.Wrapf(ErrInvalidDuration,
			"invalid duration %s, only accepts positive value", duration.String())
	}
	return nil
}

func ValidateTimestamp(t interface{}) error {
	_, ok := t.(time.Time)
	if !ok {
		return errorsmod.Wrapf(ErrInvalidTimestamp, "invalid timestamp: %T", t)
	}
	return nil
}

func ValidCampaignStatus(status CampaignStatus) bool {
	if status == CAMPAIGN_STATUS_INACTIVE ||
		status == CAMPAIGN_STATUS_ACTIVE || status == CAMPAIGN_STATUS_UNSPECIFIED {
		return true
	}
	return false
}

func ValidateClaimType(claimType ClaimType) error {
	if claimType == CLAIM_TYPE_FT || claimType == CLAIM_TYPE_NFT || claimType == CLAIM_TYPE_FT_AND_NFT {
		return nil
	}
	return errorsmod.Wrapf(ErrInvalidClaimType, "unknown claim type (%s)", claimType)
}

func ValidateClaim(claim Claim) error {
	if _, err := sdk.AccAddressFromBech32(claim.Address); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid claimer address (%s)", err)
	}
	return nil
}

func ValidateTokens(tokensA, tokensB sdk.Coin) error {
	if tokensA.IsNil() || tokensB.IsNil() {
		return errorsmod.Wrapf(
			ErrInvalidTokens,
			"invalid tokens, only accepts positive amount",
		)
	}
	if !tokensA.IsValid() {
		return errorsmod.Wrapf(
			ErrInvalidTokens,
			"invalid tokens %s, only accepts positive amount",
			tokensA.String(),
		)
	}
	if !tokensB.IsValid() {
		return errorsmod.Wrapf(
			ErrInvalidTokens,
			"invalid tokens %s, only accepts positive amount",
			tokensB.String(),
		)
	}
	if tokensA.Denom != tokensB.Denom {
		return errorsmod.Wrapf(
			ErrInvalidTokens,
			"mismatched token denoms (%s, %s)",
			tokensA.Denom,
			tokensB.Denom,
		)
	}
	return nil
}

func validateNFTMintDetails(details *NFTDetails) error {
	if details == nil || len(details.Name) == 0 || len(details.DenomId) == 0 || len(details.MediaUri) == 0 {
		return errorsmod.Wrapf(
			ErrInvalidNFTMintDetails,
			"invalid nft mint details, details should not be nil and name, media_uri can not be empty.")
	}
	if details.StartIndex <= 0 {
		return errorsmod.Wrapf(
			ErrInvalidNFTMintDetails,
			"invalid nft mint details, Start Index should be a possitive number.")
	}
	if len(details.NameDelimiter) != 1 {
		return errorsmod.Wrapf(
			ErrInvalidNFTMintDetails,
			"invalid nft mint details, Name delemeter should be a symbol string length of one.")
	}

	return nil
}

func ValidateDistribution(distribution *Distribution) error {
	if distribution == nil {
		return errorsmod.Wrapf(ErrInvalidDistribution, "distribution can not be nil")
	}
	if !(distribution.Type == DISTRIBUTION_TYPE_STREAM || distribution.Type == DISTRIBUTION_TYPE_INSTANT) {
		return errorsmod.Wrapf(ErrInvalidClaimType, "invalid distribution type (%s)", distribution.Type)
	}
	if distribution.Type == DISTRIBUTION_TYPE_STREAM {
		if err := ValidateDuration(distribution.StreamDuration); err != nil {
			return err
		}
	}
	return nil
}

func ValidateInteractionType(interaction InteractionType) error {
	if !(interaction == INTERACTION_TYPE_BURN ||
		interaction == INTERACTION_TYPE_TRANSFER || interaction == INTERACTION_TYPE_HOLD) {
		return errorsmod.Wrapf(ErrInteractionMismatch, "unknown interaction type (%s)", interaction)
	}
	return nil
}
