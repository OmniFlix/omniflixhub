package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Listing module errors
var (
	ErrInvalidTokens     = sdkerrors.Register(ModuleName, 2, "Listing does not exist")
	ErrInvalidDuration   = sdkerrors.Register(ModuleName, 4, "invalid duration")
	ErrInvalidTimestamp  = sdkerrors.Register(ModuleName, 6, "invalid time")
	ErrNonPositiveNumber = sdkerrors.Register(ModuleName, 8, "non positive number")
)
