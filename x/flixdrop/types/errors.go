package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/claim module sentinel errors
var (
	ErrFlixDropNotStarted            = sdkerrors.Register(ModuleName, 2, "flixdrop not started")
	ErrIncorrectModuleAccountBalance = sdkerrors.Register(ModuleName, 3, "claim module account balance != sum of all claim record InitialClaimableAmounts")
	ErrActionClaimNotStarted         = sdkerrors.Register(ModuleName, 4, "claim not started for this action")
)
