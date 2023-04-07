package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strings"
	"time"
)

var (
	allowedDenoms = []string{}
)

// ValidateListing checks listing is valid or not
func ValidateListing(listing Listing) error {
	if len(listing.Owner) > 0 {
		if _, err := sdk.AccAddressFromBech32(listing.Owner); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
		}
	}
	if err := ValidateId(listing.Id); err != nil {
		return err
	}
	if err := ValidatePrice(listing.Price); err != nil {
		return err
	}
	if err := ValidateSplitShares(listing.SplitShares); err != nil {
		return err
	}
	return nil
}

// ValidatePrice
func ValidatePrice(price sdk.Coin) error {
	if price.IsZero() || price.IsNegative() {
		return sdkerrors.Wrapf(ErrInvalidPrice, "invalid price %s, only accepts positive amount", price.String())
	}
	/*
		if !StringInSlice(price.Denom, allowedDenoms) {
			return sdkerrors.Wrapf(ErrInvalidPriceDenom, "invalid denom %s", price.Denom)
		}
	*/
	return nil
}

func ValidateDuration(t interface{}) error {
	duration, ok := t.(*time.Duration)
	if !ok {
		return sdkerrors.Wrapf(ErrInvalidDuration, "invalid value for duration: %T", t)
	}
	if duration.Nanoseconds() <= 0 {
		return sdkerrors.Wrapf(ErrInvalidDuration, "invalid duration %s, only accepts positive value", duration.String())
	}
	return nil
}

func ValidateId(id string) error {
	id = strings.TrimSpace(id)
	if len(id) < MinListingIdLength || len(id) > MaxListingIdLength {

		return sdkerrors.Wrapf(
			ErrInvalidListingId,
			"invalid id %s, only accepts value [%d, %d]", id, MinListingIdLength, MaxListingIdLength,
		)
	}
	if !IsBeginWithAlpha(id) || !IsAlphaNumeric(id) {
		return sdkerrors.Wrapf(ErrInvalidListingId, "invalid id %s, only accepts alphanumeric characters,and begin with an english letter", id)
	}
	return nil
}

func ValidateSplitShares(splitShares []WeightedAddress) error {
	if len(splitShares) > MaxSplits {
		return sdkerrors.Wrapf(ErrInvalidSplits, "number of splits are more than the limit, len must be less than or equal to %d ", MaxSplits)
	}
	totalWeight := sdk.NewDec(0)
	for _, share := range splitShares {
		_, err := sdk.AccAddressFromBech32(share.Address)
		if err != nil {
			return err
		}
		totalWeight = totalWeight.Add(share.Weight)
	}
	if !totalWeight.LTE(sdk.OneDec()) {
		return sdkerrors.Wrapf(ErrInvalidSplits, "invalid weights, total sum of weights must be less than %d", 1)
	}
	return nil
}

func ValidateWhiteListAccounts(whitelistAccounts []string) error {
	if len(whitelistAccounts) > MaxWhitelistAccounts {
		return sdkerrors.Wrapf(ErrInvalidWhitelistAccounts,
			"number of whitelist accounts are more than the limit, len must be less than or equal to %d ", MaxWhitelistAccounts)
	}
	for _, address := range whitelistAccounts {
		_, err := sdk.AccAddressFromBech32(address)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateIncrementPercentage(increment sdk.Dec) error {
	if !increment.IsPositive() || !increment.LTE(sdk.NewDec(1)) {
		return sdkerrors.Wrapf(ErrInvalidPercentage, "invalid percentage value (%s)", increment.String())
	}
	return nil
}

func validateAuctionId(id uint64) error {
	if id <= 0 {
		return sdkerrors.Wrapf(ErrInvalidAuctionId, "invalid auction id (%d)", id)
	}
	return nil
}

// ValidateAuctionListing checks auction listing is valid or not
func ValidateAuctionListing(auction AuctionListing) error {
	if len(auction.Owner) > 0 {
		if _, err := sdk.AccAddressFromBech32(auction.Owner); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
		}
	}
	if err := validateAuctionId(auction.Id); err != nil {
		return err
	}
	if err := ValidatePrice(auction.StartPrice); err != nil {
		return err
	}
	if err := validateIncrementPercentage(auction.IncrementPercentage); err != nil {
		return err
	}
	if err := ValidateSplitShares(auction.SplitShares); err != nil {
		return err
	}
	if err := ValidateWhiteListAccounts(auction.WhitelistAccounts); err != nil {
		return err
	}
	return nil
}

// ValidateBid checks bid is valid or not
func ValidateBid(bid Bid) error {
	if len(bid.Bidder) > 0 {
		if _, err := sdk.AccAddressFromBech32(bid.Bidder); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid bidder address (%s)", bid.Bidder)
		}
	}
	if err := ValidatePrice(bid.Amount); err != nil {
		return err
	}
	if bid.Time.IsZero() {
		return sdkerrors.Wrapf(ErrInvalidTime, "invalid time (%s)", bid.Time.String())
	}
	return nil
}
