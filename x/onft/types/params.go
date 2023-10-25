package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultDenomCreationFee Default period for closing bids for an auction
var DefaultDenomCreationFee = sdk.NewInt64Coin("uflix", 100_000_000) // 100FLIX

func NewONFTParams(denomCreationFee sdk.Coin) Params {
	return Params{
		DenomCreationFee: denomCreationFee,
	}
}

// DefaultParams returns default onft parameters
func DefaultParams() Params {
	return NewONFTParams(
		DefaultDenomCreationFee,
	)
}

// ValidateBasic performs basic validation on onft parameters.
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
		return errorsmod.Wrapf(ErrInvalidDenomCreationFee, "invalid fee amount %s, only accepts positive amounts", fee.String())
	}
	return nil
}
