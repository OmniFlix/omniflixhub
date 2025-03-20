package types

import (
	"fmt"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	KeyMinimumLeaseHours        = []byte("MinimumLeaseHours")
	KeyMaximumLeaseHours        = []byte("MaximumLeaseHours")
	KeyMinDeposit               = []byte("MinDeposit")
	KeyInitialDepositPercentage = []byte("InitialDepositPercentage")
	KeyLeaseCommission          = []byte("LeaseCommission")
	KeyCommissionDistribution   = []byte("CommissionDistribution")
	KeyDepositReleasePeriod     = []byte("DepositReleasePeriod")
)

var (
	defaultInitialDepositPercentage = sdkmath.LegacyMustNewDecFromStr("0.1")
	defaultLeaseCommission          = sdkmath.LegacyMustNewDecFromStr("0.01")
	defaultMinDeposit               = types.NewCoin("uflix", sdkmath.NewInt(1000))
	defaultStakingDistribtionPerc   = sdkmath.LegacyMustNewDecFromStr("0.5")
	defaultCPDistributionPerc       = sdkmath.LegacyMustNewDecFromStr("0.5")
	defaultDepositReleasePeriod     = time.Hour * 24 * 7
)

// ParamTable for medianode module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns default medianode parameters
func DefaultParams() Params {
	return Params{
		MinimumLeaseHours:        1,                               // Default minimum lease of 1 hour
		MaximumLeaseHours:        8760,                            // Default maximum lease of 1 year
		MinDeposit:               defaultMinDeposit,               // Default min deposit
		InitialDepositPercentage: defaultInitialDepositPercentage, // Default initial deposit percentage
		LeaseCommission:          defaultLeaseCommission,          // Default lease commission
		CommissionDistribution: Distribution{ // Default commission distribution
			Staking:       defaultStakingDistribtionPerc,
			CommunityPool: defaultCPDistributionPerc,
		},
		DepositReleasePeriod: defaultDepositReleasePeriod,
	}
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMinimumLeaseHours, &p.MinimumLeaseHours, validateMinimumLeaseDays),
		paramtypes.NewParamSetPair(KeyMaximumLeaseHours, &p.MaximumLeaseHours, validateMaximumLeaseDays),
		paramtypes.NewParamSetPair(KeyMinDeposit, &p.MinDeposit, validateMinDeposit),
		paramtypes.NewParamSetPair(KeyInitialDepositPercentage, &p.InitialDepositPercentage, validateInitialDepositPercentage),
		paramtypes.NewParamSetPair(KeyLeaseCommission, &p.LeaseCommission, validateLeaseCommission),
		paramtypes.NewParamSetPair(KeyCommissionDistribution, &p.CommissionDistribution, validateCommissionDistribution),
	}
}

// Validate performs basic validation on medianode parameters
func (p Params) Validate() error {
	if err := validateMinimumLeaseDays(p.MinimumLeaseHours); err != nil {
		return err
	}
	if err := validateMaximumLeaseDays(p.MaximumLeaseHours); err != nil {
		return err
	}
	if p.MinimumLeaseHours >= p.MaximumLeaseHours {
		return fmt.Errorf("minimum lease hours must be less than maximum lease hours: %d >= %d",
			p.MinimumLeaseHours, p.MaximumLeaseHours)
	}
	if err := validateMinDeposit(p.MinDeposit); err != nil {
		return err
	}
	if err := validateInitialDepositPercentage(p.InitialDepositPercentage); err != nil {
		return err
	}
	if err := validateLeaseCommission(p.LeaseCommission); err != nil {
		return err
	}
	if err := validateCommissionDistribution(p.CommissionDistribution); err != nil {
		return err
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
	return nil
}

func validateMinDeposit(i interface{}) error {
	v, ok := i.(types.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if err := v.Validate(); err != nil {
		return err
	}

	if v.IsZero() {
		return fmt.Errorf("min deposit cannot be zero")
	}
	if v.Amount.IsZero() {
		return fmt.Errorf("min deposit amount cannot be zero")
	}
	return nil
}

func validateInitialDepositPercentage(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("initial deposit percentage cannot be negative")
	}
	if v.GT(sdkmath.LegacyMustNewDecFromStr("1.0")) { // should not exceed 100%
		return fmt.Errorf("initial deposit percentage too high: %s", v.String())
	}
	return nil
}

func validateLeaseCommission(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("lease commission cannot be negative")
	}
	if v.GT(sdkmath.LegacyMustNewDecFromStr("1.0")) { // should not exceed 100%
		return fmt.Errorf("lease commission too high: %s", v.String())
	}
	return nil
}

func validateCommissionDistribution(i interface{}) error {
	v, ok := i.(Distribution)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.Staking.IsNegative() || v.CommunityPool.IsNegative() {
		return fmt.Errorf("commission distribution values cannot be negative")
	}
	if v.Staking.Add(v.CommunityPool).GT(sdkmath.LegacyMustNewDecFromStr("1.0")) { // should not exceed 100%
		return fmt.Errorf("total commission distribution too high: Staking: %s, CommunityPool: %s", v.Staking.String(), v.CommunityPool.String())
	}
	return nil
}
