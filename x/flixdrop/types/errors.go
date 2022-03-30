package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/claim module sentinel errors
var (
	ErrIncorrectModuleAccountBalance = sdkerrors.Register(
		ModuleName,
		1,
		"flixdrop module account balance != sum of all claim record InitialClaimableAmounts",
	)
)