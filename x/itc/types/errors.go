package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// itc module errors
var (
	ErrInvalidTokens         = sdkerrors.Register(ModuleName, 2, "Listing does not exist")
	ErrInvalidDuration       = sdkerrors.Register(ModuleName, 4, "invalid duration")
	ErrInvalidTimestamp      = sdkerrors.Register(ModuleName, 6, "invalid time")
	ErrNonPositiveNumber     = sdkerrors.Register(ModuleName, 8, "non positive number")
	ErrCampaignDoesNotExists = sdkerrors.Register(ModuleName, 10, "campaign does not exists")
	ErrInactiveCampaign      = sdkerrors.Register(ModuleName, 12, "campaign is not active")
)
