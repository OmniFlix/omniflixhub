package rest

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkquery "github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/OmniFlix/omniflixhub/v2/x/marketplace/types"
)

func registerQueryRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/%s/parameters", types.ModuleName),
		queryParams(cliCtx),
	).Methods("GET")

	r.HandleFunc(fmt.Sprintf("/%s/listing/{%s}",
		types.ModuleName, RestParamListingId),
		queryListing(cliCtx),
	).Methods("GET")

	r.HandleFunc(fmt.Sprintf("/%s/listings", types.ModuleName),
		queryAllListings(cliCtx),
	).Methods("GET")

	r.HandleFunc(fmt.Sprintf("/%s/listings/{%s}", types.ModuleName, RestParamOwner),
		queryListingsByOwner(cliCtx),
	).Methods("GET")
}

// queryParams
func queryParams(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParams), nil,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// queryListing
func queryListing(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		listingId := strings.TrimSpace(vars[RestParamListingId])
		params := types.QueryListingParams{
			Id: listingId,
		}

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryListing), bz,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// queryAllListings
func queryAllListings(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		var (
			qc    = types.NewQueryClient(cliCtx)
			query = r.URL.Query()
		)

		_, page, limit, err := rest.ParseHTTPArgs(r)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		pageReq := sdkquery.PageRequest{
			Offset:     uint64((page - 1) * limit),
			Limit:      uint64(limit),
			CountTotal: true,
		}
		owner := query.Get("owner")
		priceDenom := query.Get("priceDenom")

		listings, err := qc.Listings(
			context.Background(),
			&types.QueryListingsRequest{
				Owner:      owner,
				PriceDenom: priceDenom,
				Pagination: &pageReq,
			},
		)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		rest.PostProcessResponse(w, cliCtx, listings)
	}
}

// queryListingsByOwner
func queryListingsByOwner(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ownerStr := r.FormValue(RestParamOwner)

		var err error
		var owner sdk.AccAddress

		if len(ownerStr) > 0 {
			owner, err = sdk.AccAddressFromBech32(ownerStr)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		params := types.QueryListingsByOwnerParams{
			Owner: owner,
		}

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryListingsByOwner), bz,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
