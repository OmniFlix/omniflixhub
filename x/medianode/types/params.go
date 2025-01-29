package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	KeyMinimumLeaseDays = []byte("MinimumLeaseDays")
	KeyMaximumLeaseDays = []byte("MaximumLeaseDays")
)

// ParamTable for medianode module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns default medianode parameters
func DefaultParams() Params {
	return Params{
		MinimumLeaseDays: 1,   // Default minimum lease of 1 day
		MaximumLeaseDays: 365, // Default maximum lease of 1 year
	}
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMinimumLeaseDays, &p.MinimumLeaseDays, validateMinimumLeaseDays),
		paramtypes.NewParamSetPair(KeyMaximumLeaseDays, &p.MaximumLeaseDays, validateMaximumLeaseDays),
	}
}

// Validate performs basic validation on medianode parameters
func (p Params) Validate() error {
	if err := validateMinimumLeaseDays(p.MinimumLeaseDays); err != nil {
		return err
	}
	if err := validateMaximumLeaseDays(p.MaximumLeaseDays); err != nil {
		return err
	}
	if p.MinimumLeaseDays >= p.MaximumLeaseDays {
		return fmt.Errorf("minimum lease days must be less than maximum lease days: %d >= %d",
			p.MinimumLeaseDays, p.MaximumLeaseDays)
	}
	return nil
}

func validateMinimumLeaseDays(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("minimum lease days cannot be 0")
	}

	if v > 365 { // arbitrary upper limit for minimum lease
		return fmt.Errorf("minimum lease days too high: %d", v)
	}

	return nil
}

func validateMaximumLeaseDays(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("maximum lease days cannot be 0")
	}

	if v > 3650 { // arbitrary upper limit of 10 years
		return fmt.Errorf("maximum lease days too high: %d", v)
	}

	return nil
}
