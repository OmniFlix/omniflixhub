package types

import (
	"fmt"
	"time"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultMaxCampaignDuration = time.Hour * 24 * 90 // 90 days
)

var DefaultCampaignCreationFee = sdk.NewInt64Coin("uflix", 10_000_000)

func NewParams(
	creationFee sdk.Coin,
	maxCampaignDuration time.Duration,
) Params {
	return Params{
		CreationFee:         creationFee,
		MaxCampaignDuration: maxCampaignDuration,
	}
}

// DefaultParams returns default itc parameters
func DefaultParams() Params {
	return Params{
		MaxCampaignDuration: DefaultMaxCampaignDuration,
		CreationFee:         DefaultCampaignCreationFee,
	}
}

// ValidateBasic performs basic validation on itc parameters.
func (p Params) ValidateBasic() error {
	if err := validateMaxCampaignDuration(p.MaxCampaignDuration); err != nil {
		return err
	}
	if err := validateCampaignCreationFee(p.CreationFee); err != nil {
		return err
	}
	return nil
}

func validateMaxCampaignDuration(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.Seconds() <= 0 {
		return fmt.Errorf("max campaign duration must be positive: %f", v.Seconds())
	}

	return nil
}

func validateCampaignCreationFee(i interface{}) error {
	fee, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !fee.IsValid() || fee.IsZero() {
		return errorsmod.Wrapf(ErrInvalidCreationFee,
			"invalid fee amount %s, only accepts positive amounts", fee.String())
	}
	return nil
}
