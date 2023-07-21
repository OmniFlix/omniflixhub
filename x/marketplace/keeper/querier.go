package keeper

import (
	"fmt"
	"strings"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/OmniFlix/omniflixhub/v2/x/marketplace/types"
)

func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryParams:
			return queryParams(ctx, path[1:], req, k, legacyQuerierCdc)
		case types.QueryListing:
			return queryListing(ctx, req, k, legacyQuerierCdc)
		case types.QueryAllListings:
			return queryAllListings(ctx, req, k, legacyQuerierCdc)
		case types.QueryListingsByOwner:
			return queryListingsByOwner(ctx, req, k, legacyQuerierCdc)
		case types.QueryAuction:
			return queryAuction(ctx, req, k, legacyQuerierCdc)
		case types.QueryAllAuctions:
			return queryAllAuctions(ctx, req, k, legacyQuerierCdc)
		case types.QueryAuctionsByOwner:
			return queryAuctionsByOwner(ctx, req, k, legacyQuerierCdc)
		case types.QueryBid:
			return queryBid(ctx, req, k, legacyQuerierCdc)
		case types.QueryAllBids:
			return queryAllBids(ctx, req, k, legacyQuerierCdc)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryParams(ctx sdk.Context, _ []string, _ abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryListing(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryListingParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	id := strings.ToLower(strings.TrimSpace(params.Id))

	listing, found := k.GetListing(ctx, id)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrListingDoesNotExists, fmt.Sprintf("listing %s does not exist", id))
	}
	return codec.MarshalJSONIndent(legacyQuerierCdc, listing)
}

func queryAllListings(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryAllListingsParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	listings := k.GetAllListings(ctx)

	return codec.MarshalJSONIndent(legacyQuerierCdc, listings)
}

func queryListingsByOwner(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryListingsByOwnerParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	listings := k.GetListingsByOwner(ctx, params.Owner)
	return codec.MarshalJSONIndent(legacyQuerierCdc, listings)
}

func queryAuction(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryAuctionParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	auction, found := k.GetAuctionListing(ctx, params.Id)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrAuctionDoesNotExists, fmt.Sprintf("auction %d does not exist", params.Id))
	}
	return codec.MarshalJSONIndent(legacyQuerierCdc, auction)
}

func queryAllAuctions(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryAllAuctionsParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	auctions := k.GetAllAuctionListings(ctx)

	return codec.MarshalJSONIndent(legacyQuerierCdc, auctions)
}

func queryAuctionsByOwner(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryAuctionsByOwnerParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	auctions := k.GetAuctionListingsByOwner(ctx, params.Owner)
	return codec.MarshalJSONIndent(legacyQuerierCdc, auctions)
}

func queryBid(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryBidParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	bid, found := k.GetBid(ctx, params.Id)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrBidDoesNotExists, fmt.Sprintf("auction %d does not have any bid", params.Id))
	}
	return codec.MarshalJSONIndent(legacyQuerierCdc, bid)
}

func queryAllBids(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryAllBidsParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	bids := k.GetAllBids(ctx)

	return codec.MarshalJSONIndent(legacyQuerierCdc, bids)
}
