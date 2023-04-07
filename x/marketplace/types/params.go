package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default period for closing bids for an auction
const (
	DefaultBidClosePeriod time.Duration = time.Hour * 12 // 12 Hours
)

// Parameter keys
var (
	ParamStoreKeySaleCommission   = []byte("SaleCommission")
	ParamStoreKeyDistribution     = []byte("MarketplaceDistribution")
	ParamStoreKeyBidCloseDuration = []byte("BidCloseDuration")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewMarketplaceParams(saleCommission sdk.Dec, distribution Distribution, bidCloseDuration time.Duration) Params {
	return Params{
		SaleCommission:   saleCommission,
		Distribution:     distribution,
		BidCloseDuration: bidCloseDuration,
	}
}

// DefaultParams returns default marketplace parameters
func DefaultParams() Params {
	return NewMarketplaceParams(
		sdk.NewDecWithPrec(1, 2), // 1%
		Distribution{
			Staking:       sdk.NewDecWithPrec(50, 2), // 50%
			CommunityPool: sdk.NewDecWithPrec(50, 2), // 50%
		},
		DefaultBidClosePeriod,
	)
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeySaleCommission, &p.SaleCommission, validateSaleCommission),
		paramtypes.NewParamSetPair(ParamStoreKeyDistribution, &p.Distribution, validateMarketplaceDistributionParams),
		paramtypes.NewParamSetPair(ParamStoreKeyBidCloseDuration, &p.BidCloseDuration, validateBidCloseDuration),
	}
}

// ValidateBasic performs basic validation on marketplace parameters.
func (p Params) ValidateBasic() error {
	if err := validateSaleCommission(p.SaleCommission); err != nil {
		return err
	}
	if err := validateMarketplaceDistributionParams(p.Distribution); err != nil {
		return err
	}
	if err := validateBidCloseDuration(p.BidCloseDuration); err != nil {
		return err
	}
	return nil
}

func validateSaleCommission(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("sale commission must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("sale commission must be positive: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("sale commission too large: %s", v)
	}

	return nil
}

func validateStakingDistribution(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("staking distribution value must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("staking distribution value must be positive: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("staking distribution value too large: %s", v)
	}

	return nil
}

func validateCommunityPoolDistribution(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("community pool distribution value must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("community pool distribution value must be positive: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("community pool distribution value too large: %s", v)
	}

	return nil
}

func validateMarketplaceDistributionParams(i interface{}) error {
	v, ok := i.(Distribution)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	err := validateStakingDistribution(v.Staking)
	if err != nil {
		return err
	}
	err = validateCommunityPoolDistribution(v.CommunityPool)
	if err != nil {
		return err
	}
	if !v.Staking.Add(v.CommunityPool).Equal(sdk.OneDec()) {
		return fmt.Errorf("marketplace distribtution commission params sum must be equal to : %d", 1)
	}
	return nil
}

func validateBidCloseDuration(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.Seconds() <= 0 {
		return fmt.Errorf("bid close duration must be positive: %f", v.Seconds())
	}

	return nil
}
