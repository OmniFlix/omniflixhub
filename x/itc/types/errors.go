package types

import (
	errorsmod "cosmossdk.io/errors"
)

// itc module errors
var (
	ErrInvalidTokens           = errorsmod.Register(ModuleName, 2, "invalid tokens")
	ErrInvalidDuration         = errorsmod.Register(ModuleName, 3, "invalid duration")
	ErrInvalidTimestamp        = errorsmod.Register(ModuleName, 4, "invalid time")
	ErrNonPositiveNumber       = errorsmod.Register(ModuleName, 5, "non positive number")
	ErrCampaignDoesNotExists   = errorsmod.Register(ModuleName, 6, "campaign does not exists")
	ErrInactiveCampaign        = errorsmod.Register(ModuleName, 7, "campaign is not active")
	ErrInvalidNFTMintDetails   = errorsmod.Register(ModuleName, 8, "invalid nft mint details")
	ErrInValidMaxAllowedClaims = errorsmod.Register(ModuleName, 9, "invalid max allowed claims")
	ErrInvalidClaimType        = errorsmod.Register(ModuleName, 10, "invalid claim type")
	ErrInteractionMismatch     = errorsmod.Register(ModuleName, 11, "interaction mismatch")
	ErrClaimExists             = errorsmod.Register(ModuleName, 12, "claim exists")
	ErrClaimNotAllowed         = errorsmod.Register(ModuleName, 13, "claim not allowed")
	ErrTokenDenomMismatch      = errorsmod.Register(ModuleName, 14, "invalid token denom")
	ErrClaimingNFT             = errorsmod.Register(ModuleName, 15, "claim nft failed")
	ErrDepositNotAllowed       = errorsmod.Register(ModuleName, 16, "deposit not allowed")
	ErrInvalidCreationFee      = errorsmod.Register(ModuleName, 17, "invalid fee")
	ErrInvalidFeeDenom         = errorsmod.Register(ModuleName, 18, "invalid fee denom")
	ErrNotEnoughFeeAmount      = errorsmod.Register(ModuleName, 19, "not enough fee")
	ErrInvalidDistribution     = errorsmod.Register(ModuleName, 20, "invalid distribution")
	ErrInvalidNftDenomId       = errorsmod.Register(ModuleName, 21, "invalid nft denom id")
)
