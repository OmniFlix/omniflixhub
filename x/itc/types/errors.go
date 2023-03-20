package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// itc module errors
var (
	ErrInvalidTokens           = sdkerrors.Register(ModuleName, 2, "Listing does not exist")
	ErrInvalidDuration         = sdkerrors.Register(ModuleName, 4, "invalid duration")
	ErrInvalidTimestamp        = sdkerrors.Register(ModuleName, 6, "invalid time")
	ErrNonPositiveNumber       = sdkerrors.Register(ModuleName, 8, "non positive number")
	ErrCampaignDoesNotExists   = sdkerrors.Register(ModuleName, 10, "campaign does not exists")
	ErrInactiveCampaign        = sdkerrors.Register(ModuleName, 12, "campaign is not active")
	ErrNFTMintDetails          = sdkerrors.Register(ModuleName, 14, "invalid nft mint details")
	ErrInValidMaxAllowedClaims = sdkerrors.Register(ModuleName, 16, "invalid max allowed claims")
	ErrInvalidClaimType        = sdkerrors.Register(ModuleName, 18, "invalid claim type")
	ErrInteractionMismatched   = sdkerrors.Register(ModuleName, 20, "interaction mismatch")
	ErrClaimExists             = sdkerrors.Register(ModuleName, 22, "claim exists")
	ErrClaimNotAllowed         = sdkerrors.Register(ModuleName, 24, "claim not allowed")
)
