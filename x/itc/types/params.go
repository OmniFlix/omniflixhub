package types

import (
	"fmt"
	"time"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	DefaultMaxCampaignDuration time.Duration = time.Hour * 24 * 90 // 90 days
)

var ParamStoreKeyMaxCampaignDuration = []byte("MaxCampaignDuration")

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns default marketplace parameters
func DefaultParams() Params {
	return Params{
		MaxCampaignDuration: DefaultMaxCampaignDuration,
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyMaxCampaignDuration,
			&p.MaxCampaignDuration, validateMaxCampaignDuration),
	}
}

// ValidateBasic performs basic validation on marketplace parameters.
func (p Params) ValidateBasic() error {
	if err := validateMaxCampaignDuration(p.MaxCampaignDuration); err != nil {
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
