package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	DefaultMaxCampaignDuration = time.Hour * 24 * 90 // 90 days
)

var DefaultCampaignCreationFee = sdk.NewInt64Coin("uflix", 10_000_000)

var (
	ParamStoreKeyMaxCampaignDuration = []byte("MaxCampaignDuration")
	ParamStoreKeyCampaignCreationFee = []byte("CampaignCreationFee")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns default marketplace parameters
func DefaultParams() Params {
	return Params{
		MaxCampaignDuration: DefaultMaxCampaignDuration,
		CreationFee:         DefaultCampaignCreationFee,
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyMaxCampaignDuration, &p.MaxCampaignDuration, validateMaxCampaignDuration),
		paramtypes.NewParamSetPair(ParamStoreKeyCampaignCreationFee, &p.CreationFee, validateCampaignCreationFee),
	}
}

// ValidateBasic performs basic validation on marketplace parameters.
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
		return sdkerrors.Wrapf(ErrInvalidCreationFee,
			"invalid fee amount %s, only accepts positive amounts", fee.String())
	}
	return nil
}
