package keeper

import (
	"context"
	"time"

	"cosmossdk.io/store/prefix"
	"github.com/OmniFlix/omniflixhub/v5/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	gogotypes "github.com/cosmos/gogoproto/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

// Params queries params of marketplace module
func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func (k Keeper) Listing(goCtx context.Context, req *types.QueryListingRequest) (*types.QueryListingResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	listing, found := k.GetListing(ctx, req.Id)
	if !found {
		return nil, status.Errorf(codes.NotFound, "listing %s not found", req.Id)
	}

	return &types.QueryListingResponse{Listing: &listing}, nil
}

func (k Keeper) Listings(goCtx context.Context, req *types.QueryListingsRequest) (*types.QueryListingsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	var listings []types.Listing
	var pageRes *query.PageResponse
	store := ctx.KVStore(k.storeKey)

	var owner sdk.AccAddress
	var err error
	if len(req.Owner) > 0 {
		owner, err = sdk.AccAddressFromBech32(req.Owner)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid owner address (%s)", err)
		}
		listingStore := prefix.NewStore(store, append(types.PrefixListingOwner, owner.Bytes()...))
		pageRes, err = query.Paginate(listingStore, req.Pagination, func(key []byte, value []byte) error {
			var listingId gogotypes.StringValue
			k.cdc.MustUnmarshal(value, &listingId)
			listing, found := k.GetListing(ctx, listingId.Value)
			if found {
				listings = append(listings, listing)
			}
			return nil
		})

	} else if len(req.PriceDenom) > 0 {
		listingStore := prefix.NewStore(store, types.KeyListingPriceDenomPrefix(req.PriceDenom, ""))
		pageRes, err = query.Paginate(listingStore, req.Pagination, func(key []byte, value []byte) error {
			var listingId gogotypes.StringValue
			k.cdc.MustUnmarshal(value, &listingId)
			listing, found := k.GetListing(ctx, listingId.Value)
			if found {
				listings = append(listings, listing)
			}
			return nil
		})
	} else {

		listingStore := prefix.NewStore(store, types.PrefixListingId)
		pageRes, err = query.Paginate(listingStore, req.Pagination, func(key []byte, value []byte) error {
			var listing types.Listing
			k.cdc.MustUnmarshal(value, &listing)
			listings = append(listings, listing)
			return nil
		})
	}
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}

	return &types.QueryListingsResponse{Listings: listings, Pagination: pageRes}, nil
}

func (k Keeper) ListingsByOwner(goCtx context.Context, req *types.QueryListingsByOwnerRequest) (*types.QueryListingsByOwnerResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	var owner sdk.AccAddress
	var err error
	if len(req.Owner) > 0 {
		owner, err = sdk.AccAddressFromBech32(req.Owner)
		if err != nil || owner == nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid owner address (%s)", err)
		}
	}

	var listings []types.Listing
	var pageRes *query.PageResponse
	store := ctx.KVStore(k.storeKey)

	listingStore := prefix.NewStore(store, append(types.PrefixListingOwner, owner.Bytes()...))
	pageRes, err = query.Paginate(listingStore, req.Pagination, func(key []byte, value []byte) error {
		var listingId gogotypes.StringValue
		k.cdc.MustUnmarshal(value, &listingId)
		listing, found := k.GetListing(ctx, listingId.Value)
		if found {
			listings = append(listings, listing)
		}
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}

	return &types.QueryListingsByOwnerResponse{Listings: listings, Pagination: pageRes}, nil
}

func (k Keeper) ListingsByPriceDenom(
	goCtx context.Context,
	req *types.QueryListingsByPriceDenomRequest,
) (*types.QueryListingsByPriceDenomResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	var err error

	var listings []types.Listing
	var pageRes *query.PageResponse
	store := ctx.KVStore(k.storeKey)

	listingStore := prefix.NewStore(store, types.KeyListingPriceDenomPrefix(req.PriceDenom, ""))
	pageRes, err = query.Paginate(listingStore, req.Pagination, func(key []byte, value []byte) error {
		var listingId gogotypes.StringValue
		k.cdc.MustUnmarshal(value, &listingId)
		listing, found := k.GetListing(ctx, listingId.Value)
		if found {
			listings = append(listings, listing)
		}
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}

	return &types.QueryListingsByPriceDenomResponse{Listings: listings, Pagination: pageRes}, nil
}

func (k Keeper) ListingByNftId(
	goCtx context.Context,
	req *types.QueryListingByNFTIDRequest,
) (*types.QueryListingResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}
	if req.NftId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "need nft id to request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	listingId, found := k.GetListingIdByNftId(ctx, req.NftId)
	if found {
		listing, err := k.Listing(goCtx, &types.QueryListingRequest{
			Id: listingId,
		})
		if err != nil {
			return nil, err
		}
		return listing, nil
	}
	return nil, status.Errorf(codes.NotFound, "listing not found with given nft id")
}

func (k Keeper) Auctions(goCtx context.Context, req *types.QueryAuctionsRequest) (*types.QueryAuctionsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	var filteredAuctions []types.AuctionListing
	var pageRes *query.PageResponse
	store := ctx.KVStore(k.storeKey)
	auctionStore := prefix.NewStore(store, types.PrefixAuctionId)
	pageRes, err := query.FilteredPaginate(auctionStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var al types.AuctionListing
		k.cdc.MustUnmarshal(value, &al)
		matchOwner, matchPriceDenom, matchStatus := true, true, true
		// match status (if supplied/valid)
		if types.ValidAuctionStatus(req.Status) {
			if req.Status == types.AUCTION_STATUS_ACTIVE {
				matchStatus = al.StartTime.Before(time.Now())
			} else {
				matchStatus = al.StartTime.After(time.Now())
			}
		}

		// match owner address (if supplied)
		if len(req.Owner) > 0 {
			owner, err := sdk.AccAddressFromBech32(req.Owner)
			if err != nil {
				return false, err
			}

			matchOwner = al.Owner == owner.String()
		}

		// match Price Denom (if supplied)
		if len(req.PriceDenom) > 0 {
			matchPriceDenom = al.StartPrice.Denom == req.PriceDenom
		}

		if matchOwner && matchPriceDenom && matchStatus {
			if accumulate {
				filteredAuctions = append(filteredAuctions, al)
			}

			return true, nil
		}

		return false, nil
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}

	return &types.QueryAuctionsResponse{Auctions: filteredAuctions, Pagination: pageRes}, nil
}

func (k Keeper) Auction(goCtx context.Context, req *types.QueryAuctionRequest) (*types.QueryAuctionResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	auction, found := k.GetAuctionListing(ctx, req.Id)
	if !found {
		return nil, status.Errorf(codes.NotFound, "auction %d not found", req.Id)
	}
	return &types.QueryAuctionResponse{Auction: &auction}, nil
}

func (k Keeper) AuctionsByOwner(goCtx context.Context, req *types.QueryAuctionsByOwnerRequest) (*types.QueryAuctionsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	var owner sdk.AccAddress
	var err error
	if len(req.Owner) > 0 {
		owner, err = sdk.AccAddressFromBech32(req.Owner)
		if err != nil || owner == nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid owner address (%s)", err)
		}
	}

	var auctions []types.AuctionListing
	var pageRes *query.PageResponse
	store := ctx.KVStore(k.storeKey)

	auctionStore := prefix.NewStore(store, append(types.PrefixAuctionOwner, owner.Bytes()...))
	pageRes, err = query.Paginate(auctionStore, req.Pagination, func(key []byte, value []byte) error {
		var auctionId gogotypes.UInt64Value
		k.cdc.MustUnmarshal(value, &auctionId)
		auction, found := k.GetAuctionListing(ctx, auctionId.Value)
		if found {
			auctions = append(auctions, auction)
		}
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}

	return &types.QueryAuctionsResponse{Auctions: auctions, Pagination: pageRes}, nil
}

func (k Keeper) AuctionsByPriceDenom(goCtx context.Context, req *types.QueryAuctionsByPriceDenomRequest) (*types.QueryAuctionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	var err error

	var auctions []types.AuctionListing
	var pageRes *query.PageResponse
	store := ctx.KVStore(k.storeKey)

	auctionStore := prefix.NewStore(store, append(types.PrefixAuctionPriceDenom, []byte(req.PriceDenom)...))
	pageRes, err = query.Paginate(auctionStore, req.Pagination, func(key []byte, value []byte) error {
		var auctionId gogotypes.UInt64Value
		k.cdc.MustUnmarshal(value, &auctionId)
		auction, found := k.GetAuctionListing(ctx, auctionId.Value)
		if found {
			auctions = append(auctions, auction)
		}
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}

	return &types.QueryAuctionsResponse{Auctions: auctions, Pagination: pageRes}, nil
}

func (k Keeper) AuctionByNftId(goCtx context.Context, req *types.QueryAuctionByNFTIDRequest) (*types.QueryAuctionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if req.NftId == "" {
		return nil, status.Error(codes.InvalidArgument, "need nft id to request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	auctionId, found := k.GetAuctionListingIdByNftId(ctx, req.NftId)
	if found {
		auction, err := k.Auction(goCtx, &types.QueryAuctionRequest{
			Id: auctionId,
		})
		if err != nil {
			return nil, err
		}
		return auction, nil
	}
	return nil, status.Error(codes.NotFound, "auction not found with given nft id")
}

func (k Keeper) Bid(goCtx context.Context, req *types.QueryBidRequest) (*types.QueryBidResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	bid, found := k.GetBid(ctx, req.Id)
	if !found {
		return nil, status.Errorf(codes.NotFound, "bid not found for auction %d", req.Id)
	}
	return &types.QueryBidResponse{Bid: &bid}, nil
}

func (k Keeper) Bids(goCtx context.Context, req *types.QueryBidsRequest) (*types.QueryBidsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	var bids []types.Bid
	var pageRes *query.PageResponse
	store := ctx.KVStore(k.storeKey)

	if len(req.Bidder) > 0 {
		_, err := sdk.AccAddressFromBech32(req.Bidder)
		if err != nil {
			return nil, err
		}
	}

	bidStore := prefix.NewStore(store, types.PrefixBidByAuctionId)
	pageRes, err := query.Paginate(bidStore, req.Pagination, func(key []byte, value []byte) error {
		var bid types.Bid
		k.cdc.MustUnmarshal(value, &bid)
		if len(req.Bidder) > 0 {
			if bid.Bidder == req.Bidder {
				bids = append(bids, bid)
			}
		} else {
			bids = append(bids, bid)
		}
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}
	return &types.QueryBidsResponse{Bids: bids, Pagination: pageRes}, nil
}
