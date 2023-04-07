package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default period for closing bids for an auction
var (
	DefaultDenomCreationFee = sdk.NewInt64Coin("uflix", 100_000_000) // 100FLIX
)

// Parameter keys
var (
	ParamStoreKeyDenomCreationFee = []byte("DenomCreationFee")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewONFTParams(denomCreationFee sdk.Coin) Params {
	return Params{
		DenomCreationFee: denomCreationFee,
	}
}

// DefaultParams returns default marketplace parameters
func DefaultParams() Params {
	return NewONFTParams(
		DefaultDenomCreationFee,
	)
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyDenomCreationFee, &p.DenomCreationFee, validateDenomCreationFee),
	}
}

// ValidateBasic performs basic validation on marketplace parameters.
func (p Params) ValidateBasic() error {
	if err := validateDenomCreationFee(p.DenomCreationFee); err != nil {
		return err
	}
	return nil
}

// ValidateDenomCreationFee performs validation of denom creation fee

func validateDenomCreationFee(i interface{}) error {
	fee, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !fee.IsValid() || fee.IsZero() {
		return sdkerrors.Wrapf(ErrInvalidDenomCreationFee, "invalid fee amount %s, only accepts positive amounts", fee.String())
	}
	return nil
}
