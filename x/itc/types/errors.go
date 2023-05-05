package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// itc module errors
var (
	ErrInvalidTokens           = sdkerrors.Register(ModuleName, 2, "invalid tokens")
	ErrInvalidDuration         = sdkerrors.Register(ModuleName, 3, "invalid duration")
	ErrInvalidTimestamp        = sdkerrors.Register(ModuleName, 4, "invalid time")
	ErrNonPositiveNumber       = sdkerrors.Register(ModuleName, 5, "non positive number")
	ErrCampaignDoesNotExists   = sdkerrors.Register(ModuleName, 6, "campaign does not exists")
	ErrInactiveCampaign        = sdkerrors.Register(ModuleName, 7, "campaign is not active")
	ErrInvalidNFTMintDetails   = sdkerrors.Register(ModuleName, 8, "invalid nft mint details")
	ErrInValidMaxAllowedClaims = sdkerrors.Register(ModuleName, 9, "invalid max allowed claims")
	ErrInvalidClaimType        = sdkerrors.Register(ModuleName, 10, "invalid claim type")
	ErrInteractionMismatch     = sdkerrors.Register(ModuleName, 11, "interaction mismatch")
	ErrClaimExists             = sdkerrors.Register(ModuleName, 12, "claim exists")
	ErrClaimNotAllowed         = sdkerrors.Register(ModuleName, 13, "claim not allowed")
	ErrTokenDenomMismatch      = sdkerrors.Register(ModuleName, 14, "invalid token denom")
	ErrCancelNotAllowed        = sdkerrors.Register(ModuleName, 15, "cancel not allowed")
	ErrClaimingNFT             = sdkerrors.Register(ModuleName, 16, "claim nft failed")
	ErrDepositNotAllowed       = sdkerrors.Register(ModuleName, 17, "deposit not allowed")
)
