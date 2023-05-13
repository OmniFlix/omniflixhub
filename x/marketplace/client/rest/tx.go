package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/OmniFlix/omniflixhub/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func registerTxRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc(
		fmt.Sprintf("/%s/list-nft", types.ModuleName),
		ListNft(cliCtx),
	).Methods("POST")

	r.HandleFunc(
		fmt.Sprintf("/%s/listings/{%s}/edit-listing", types.ModuleName, RestParamListingId),
		editListing(cliCtx),
	).Methods("PUT")

	r.HandleFunc(
		fmt.Sprintf("/%s/listings/{%s}/de-list-nft", types.ModuleName, RestParamListingId),
		deListNft(cliCtx),
	).Methods("PUT")

	r.HandleFunc(
		fmt.Sprintf("/%s/listings/{%s}/buy-nft", types.ModuleName, RestParamListingId),
		buyNft(cliCtx),
	).Methods("POST")
}

func ListNft(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req listNftReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}
		owner, err := sdk.AccAddressFromBech32(req.Owner)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest,
				fmt.Sprintf("invalid bech32 account address: %s", req.Owner))
			return
		}

		price, err := sdk.ParseCoinNormalized(req.Price)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to parse min amount: %s", req.Price))
			return
		}

		msg := types.NewMsgListNFT(req.DenomId, req.NftId, price, owner, req.SplitShares)

		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func editListing(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		listingId := strings.TrimSpace(vars[RestParamListingId])
		var req editListingReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}
		owner, err := sdk.AccAddressFromBech32(req.Owner)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest,
				fmt.Sprintf("invalid bech32 account address: %s", req.Owner))
			return
		}

		price, err := sdk.ParseCoinNormalized(req.Price)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to parse min amount: %s", req.Price))
			return
		}

		msg := types.NewMsgEditListing(listingId, price, owner)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func deListNft(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		listingId := strings.TrimSpace(vars[RestParamListingId])
		var req deListNftReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}
		owner, err := sdk.AccAddressFromBech32(req.Owner)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest,
				fmt.Sprintf("invalid bech32 account address: %s", req.Owner))
			return
		}

		msg := types.NewMsgDeListNFT(listingId, owner)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func buyNft(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		listingId := strings.TrimSpace(vars[RestParamListingId])

		var req buyNftReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}
		buyer, err := sdk.AccAddressFromBech32(req.Buyer)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest,
				fmt.Sprintf("invalid bech32 account address: %s", req.Buyer))
			return
		}

		price, err := sdk.ParseCoinNormalized(req.Price)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to parse min amount: %s", req.Price))
			return
		}

		msg := types.NewMsgBuyNFT(listingId, price, buyer)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}
