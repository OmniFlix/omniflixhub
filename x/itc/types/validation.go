package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ValidateCampaign checks campaign is valid or not
func ValidateCampaign(campaign Campaign) error {
	if _, err := sdk.AccAddressFromBech32(campaign.Creator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if err := ValidateClaimType(campaign.ClaimType); err != nil {
		return err
	}
	if err := ValidateInteractionType(campaign.Interaction); err != nil {
		return err
	}
	if err := ValidateTokensWithClaimType(campaign.ClaimType, campaign.TotalTokens); err != nil {
		return err
	}
	if err := ValidateTokensWithClaimType(campaign.ClaimType, campaign.ClaimableTokens); err != nil {
		return err
	}
	if campaign.ClaimType == CLAIM_TYPE_NFT {
		if err := validateNFTMintDetails(campaign.NftMintDetails); err != nil {
			return err
		}
	}
	if campaign.ClaimType == CLAIM_TYPE_FT && campaign.Distribution.Type == DISTRIBUTION_TYPE_VEST {
		if err := ValidateDistribution(campaign.Distribution); err != nil {
			return err
		}
	}
	if campaign.MaxAllowedClaims == 0 {
		return sdkerrors.Wrapf(ErrInValidMaxAllowedClaims,
			"max allowed claims must be a positive number (%d)", campaign.MaxAllowedClaims)
	}
	return nil
}

// ValidateTokens validates tokens
func ValidateTokens(tokens Tokens) error {
	if tokens.Fungible != nil && (tokens.Fungible.IsZero() || tokens.Fungible.IsNegative()) {
		return sdkerrors.Wrapf(
			ErrInvalidTokens,
			"invalid tokens %s, only accepts positive amount",
			tokens.String(),
		)
	}
	return nil
}

func ValidateDuration(t interface{}) error {
	duration, ok := t.(time.Duration)
	if !ok {
		return sdkerrors.Wrapf(ErrInvalidDuration, "invalid value for duration: %T", t)
	}
	if duration.Nanoseconds() <= 0 {
		return sdkerrors.Wrapf(ErrInvalidDuration,
			"invalid duration %s, only accepts positive value", duration.String())
	}
	return nil
}

func ValidateTimestamp(t interface{}) error {
	_, ok := t.(time.Time)
	if !ok {
		return sdkerrors.Wrapf(ErrInvalidTimestamp, "invalid timestamp: %T", t)
	}
	return nil
}

func ValidCampaignStatus(status CampaignStatus) bool {
	if status == CAMPAIGN_STATUS_INACTIVE ||
		status == CAMPAIGN_STATUS_ACTIVE {
		return true
	}
	return false
}

func ValidateClaimType(claimType ClaimType) error {
	if claimType == CLAIM_TYPE_FT || claimType == CLAIM_TYPE_NFT || claimType == CLAIM_TYPE_FT_AND_NFT {
		return nil
	}
	return sdkerrors.Wrapf(ErrInvalidClaimType, "unknown claim type (%s)", claimType)
}

func ValidateClaim(claim Claim) error {
	if _, err := sdk.AccAddressFromBech32(claim.Address); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid claimer address (%s)", err)
	}
	return nil
}

func ValidateTokensWithClaimType(claimType ClaimType, tokens Tokens) error {
	if claimType == CLAIM_TYPE_FT {
		if !(tokens.Fungible.IsValid() && tokens.Fungible.IsPositive()) {
			return sdkerrors.Wrapf(
				ErrInvalidTokens,
				"invalid tokens %s, only accepts positive amount",
				tokens.String(),
			)
		}
	}
	if claimType == CLAIM_TYPE_NFT {
		if !(tokens.NonFungible.IsValid() && tokens.NonFungible.IsPositive()) {
			return sdkerrors.Wrapf(
				ErrInvalidTokens,
				"invalid tokens %s, only accepts positive amount",
				tokens.String(),
			)
		}
	}
	return nil
}

func validateNFTMintDetails(details *NFTDetails) error {
	if details == nil || len(details.Name) == 0 {
		return sdkerrors.Wrapf(
			ErrInvalidTokens,
			"invalid nft mint details, details should not be nil and name can not be empty.")
	}
	return nil
}

func ValidateDistribution(distribution *Distribution) error {
	if err := ValidateDuration(distribution.VestedDistributionEpochDuration); err != nil {
		return err
	}
	if err := ValidateDuration(distribution.VestingDuration); err != nil {
		return err
	}
	return nil
}

func ValidateInteractionType(interaction InteractionType) error {
	if interaction == INTERACTION_TYPE_BURN ||
		interaction == INTERACTION_TYPE_TRANSFER || interaction == INTERACTION_TYPE_HOLD {
		return nil
	}
	return sdkerrors.Wrapf(ErrInvalidClaimType, "unknown interaction type (%s)", interaction)
}
