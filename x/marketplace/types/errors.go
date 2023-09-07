package types

import (
	errorsmod "cosmossdk.io/errors"
)

// Listing module errors
var (
	ErrListingNotExists         = errorsmod.Register(ModuleName, 2, "Listing does not exist")
	ErrInvalidOwner             = errorsmod.Register(ModuleName, 3, "invalid Listing owner")
	ErrInvalidPrice             = errorsmod.Register(ModuleName, 4, "invalid amount")
	ErrInvalidListing           = errorsmod.Register(ModuleName, 5, "invalid Listing")
	ErrListingAlreadyExists     = errorsmod.Register(ModuleName, 6, "Listing already exists")
	ErrNotEnoughAmount          = errorsmod.Register(ModuleName, 7, "amount is not enough to buy")
	ErrInvalidPriceDenom        = errorsmod.Register(ModuleName, 8, "invalid price denom")
	ErrInvalidListingId         = errorsmod.Register(ModuleName, 9, "invalid Listing id")
	ErrInvalidNftId             = errorsmod.Register(ModuleName, 10, "invalid nft id")
	ErrNftNotExists             = errorsmod.Register(ModuleName, 11, "nft not exists with given details")
	ErrUnauthorized             = errorsmod.Register(ModuleName, 12, "unauthorized")
	ErrNftNonTransferable       = errorsmod.Register(ModuleName, 13, "non-transferable nft")
	ErrListingDoesNotExists     = errorsmod.Register(ModuleName, 14, "listing doesn't exists")
	ErrInvalidSplits            = errorsmod.Register(ModuleName, 15, "invalid split shares")
	ErrNonPositiveNumber        = errorsmod.Register(ModuleName, 16, "non positive number")
	ErrInvalidAuctionId         = errorsmod.Register(ModuleName, 17, "invalid auction id")
	ErrInvalidWhitelistAccounts = errorsmod.Register(ModuleName, 18, "invalid whitelist accounts")
	ErrAuctionDoesNotExists     = errorsmod.Register(ModuleName, 19, "auction listing doesn't exists")
	ErrBidExists                = errorsmod.Register(ModuleName, 20, "bid exists")
	ErrEndedAuction             = errorsmod.Register(ModuleName, 21, "auction ended")
	ErrInActiveAuction          = errorsmod.Register(ModuleName, 22, "inactive auction")
	ErrBidAmountNotEnough       = errorsmod.Register(ModuleName, 23, "amount is not enough to bid")
	ErrBidDoesNotExists         = errorsmod.Register(ModuleName, 24, "bid does not exists")
	ErrInvalidStartTime         = errorsmod.Register(ModuleName, 25, "invalid start time")
	ErrInvalidPercentage        = errorsmod.Register(ModuleName, 26, "invalid percentage decimal value")
	ErrInvalidTime              = errorsmod.Register(ModuleName, 27, "invalid timestamp value")
	ErrInvalidDuration          = errorsmod.Register(ModuleName, 28, "invalid duration")
)
