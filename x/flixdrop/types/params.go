package types

import (
	"fmt"
	"strings"
	"time"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	DefaultClaimDenom         = "uflix"
	DefaultDurationUntilDecay = time.Hour
	DefaultDurationOfDecay    = time.Hour * 5
)

// Parameter store keys
var (
	KeyFlixDropStartTime  = []byte("FlixDropStartTime")
	KeyClaimDenom         = []byte("ClaimDenom")
	KeyDurationUntilDecay = []byte("DurationUntilDecay")
	KeyDurationOfDecay    = []byte("DurationOfDecay")
)

func NewParams(claimDenom string, startTime time.Time, durationUntilDecay, durationOfDecay time.Duration) Params {
	return Params{
		ClaimDenom:         claimDenom,
		FlixdropStartTime:  startTime,
		DurationUntilDecay: durationUntilDecay,
		DurationOfDecay:    durationOfDecay,
	}
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyClaimDenom, &p.ClaimDenom, validateDenom),
		paramtypes.NewParamSetPair(KeyFlixDropStartTime, &p.FlixdropStartTime, validateTime),
		paramtypes.NewParamSetPair(KeyDurationUntilDecay, &p.DurationUntilDecay, validateDuration),
		paramtypes.NewParamSetPair(KeyDurationOfDecay, &p.DurationOfDecay, validateDuration),
	}
}

// Validate validates all params
func (p Params) Validate() error {
	err := validateDenom(p.ClaimDenom)
	return err
}

func (p Params) IsFlixDropStarted(t time.Time) bool {
	if p.FlixdropStartTime.IsZero() {
		return false
	}
	if t.Before(p.FlixdropStartTime) {
		return false
	}
	return true
}

// ParamKeyTable for staking module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func validateDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return fmt.Errorf("invalid denom: %s", v)
	}

	return nil
}

func validateTime(i interface{}) error {
	_, ok := i.(time.Time)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateDuration(i interface{}) error {
	d, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if d < 1 {
		return fmt.Errorf("duration must be greater than or equal to 1: %d", d)
	}
	return nil
}
