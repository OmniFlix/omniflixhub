package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	QueryParams          = "params"
	QueryListing         = "listing"
	QueryAllListings     = "listings"
	QueryListingsByOwner = "listings-by-owner"
	QueryAuction         = "auction"
	QueryAllAuctions     = "auctions"
	QueryBid             = "bid"
	QueryAllBids         = "bids"
	QueryAuctionsByOwner = "auctions-by-owner"
)

// QueryListingParams is the query parameters for '/marketplace/listings/{id}'
type QueryListingParams struct {
	Id string
}

// NewQueryListingParams
func NewQueryListingParams(id string) QueryListingParams {
	return QueryListingParams{
		Id: id,
	}
}

// QueryAllListingsParams is the query parameters for 'marketplace/listings'
type QueryAllListingsParams struct{}

// NewQueryAllListingsParams
func NewQueryAllListingsParams() QueryAllListingsParams {
	return QueryAllListingsParams{}
}

// QueryListingsByOwnerParams is the query parameters for 'marketplace/listings/{owner}'
type QueryListingsByOwnerParams struct {
	Owner sdk.AccAddress
}

// NewQueryListingsByOwnerParams
func NewQueryListingsByOwnerParams(owner sdk.AccAddress) QueryListingsByOwnerParams {
	return QueryListingsByOwnerParams{
		Owner: owner,
	}
}

// QueryAuctionParams is the query parameters for '/marketplace/auctions/{id}'
type QueryAuctionParams struct {
	Id uint64
}

// NewQueryAuctionParams
func NewQueryAuctionParams(id uint64) QueryAuctionParams {
	return QueryAuctionParams{
		Id: id,
	}
}

// QueryAllListingsParams is the query parameters for 'marketplace/auctions'
type QueryAllAuctionsParams struct{}

// NewQueryAllListingsParams
func NewQueryAllAuctionsParams() QueryAllAuctionsParams {
	return QueryAllAuctionsParams{}
}

// QueryListingsByOwnerParams is the query parameters for 'marketplace/auctions/{owner}'
type QueryAuctionsByOwnerParams struct {
	Owner sdk.AccAddress
}

// NewQueryAuctionsByOwnerParams
func NewQueryAuctionsByOwnerParams(owner sdk.AccAddress) QueryAuctionsByOwnerParams {
	return QueryAuctionsByOwnerParams{
		Owner: owner,
	}
}

// QueryBidParams is the query parameters for '/marketplace/bids/{id}'
type QueryBidParams struct {
	Id uint64
}

// NewQueryAuctionParams
func NewQueryBidParams(id uint64) QueryBidParams {
	return QueryBidParams{
		Id: id,
	}
}

// QueryAllBidsParams is the query parameters for 'marketplace/bids'
type QueryAllBidsParams struct{}

// NewQueryAllListingsParams
func NewQueryAllBidsParams() QueryAllBidsParams {
	return QueryAllBidsParams{}
}
