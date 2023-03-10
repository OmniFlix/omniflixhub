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
	duration, ok := t.(*time.Duration)
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
	_, ok := t.(*time.Time)
	if !ok {
		return sdkerrors.Wrapf(ErrInvalidTimestamp, "invalid timestamp: %T", t)
	}
	return nil
}

func ValidateCampaignStatus(status CampaignStatus) bool {
	if status == CAMPAIGN_STATUS_INACTIVE ||
		status == CAMPAIGN_STATUS_ACTIVE {
		return true
	}
	return false
}

func ValidateClaim(claim Claim) error {
	if _, err := sdk.AccAddressFromBech32(claim.Address); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid claimer address (%s)", err)
	}
	return nil
}
