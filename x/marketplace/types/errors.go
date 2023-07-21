package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Listing module errors
var (
	ErrListingNotExists         = sdkerrors.Register(ModuleName, 2, "Listing does not exist")
	ErrInvalidOwner             = sdkerrors.Register(ModuleName, 3, "invalid Listing owner")
	ErrInvalidPrice             = sdkerrors.Register(ModuleName, 4, "invalid amount")
	ErrInvalidListing           = sdkerrors.Register(ModuleName, 5, "invalid Listing")
	ErrListingAlreadyExists     = sdkerrors.Register(ModuleName, 6, "Listing already exists")
	ErrNotEnoughAmount          = sdkerrors.Register(ModuleName, 7, "amount is not enough to buy")
	ErrInvalidPriceDenom        = sdkerrors.Register(ModuleName, 8, "invalid price denom")
	ErrInvalidListingId         = sdkerrors.Register(ModuleName, 9, "invalid Listing id")
	ErrInvalidNftId             = sdkerrors.Register(ModuleName, 10, "invalid nft id")
	ErrNftNotExists             = sdkerrors.Register(ModuleName, 11, "nft not exists with given details")
	ErrUnauthorized             = sdkerrors.Register(ModuleName, 12, "unauthorized")
	ErrNftNonTransferable       = sdkerrors.Register(ModuleName, 13, "non-transferable nft")
	ErrListingDoesNotExists     = sdkerrors.Register(ModuleName, 14, "listing doesn't exists")
	ErrInvalidSplits            = sdkerrors.Register(ModuleName, 15, "invalid split shares")
	ErrNonPositiveNumber        = sdkerrors.Register(ModuleName, 16, "non positive number")
	ErrInvalidAuctionId         = sdkerrors.Register(ModuleName, 17, "invalid auction id")
	ErrInvalidWhitelistAccounts = sdkerrors.Register(ModuleName, 18, "invalid whitelist accounts")
	ErrAuctionDoesNotExists     = sdkerrors.Register(ModuleName, 19, "auction listing doesn't exists")
	ErrBidExists                = sdkerrors.Register(ModuleName, 20, "bid exists")
	ErrEndedAuction             = sdkerrors.Register(ModuleName, 21, "auction ended")
	ErrInActiveAuction          = sdkerrors.Register(ModuleName, 22, "inactive auction")
	ErrBidAmountNotEnough       = sdkerrors.Register(ModuleName, 23, "amount is not enough to bid")
	ErrBidDoesNotExists         = sdkerrors.Register(ModuleName, 24, "bid does not exists")
	ErrInvalidStartTime         = sdkerrors.Register(ModuleName, 25, "invalid start time")
	ErrInvalidPercentage        = sdkerrors.Register(ModuleName, 26, "invalid percentage decimal value")
	ErrInvalidTime              = sdkerrors.Register(ModuleName, 27, "invalid timestamp value")
	ErrInvalidDuration          = sdkerrors.Register(ModuleName, 28, "invalid duration")
)
