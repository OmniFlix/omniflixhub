package rest

import (
	"github.com/OmniFlix/marketplace/x/marketplace/types"
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

// Rest variable names
// nolint
const (
	RestParamListingId = "listing_id"
	RestParamOwner     = "owner"
	RestParamBuyer     = "buyer"
)

// RegisterHandlers
func RegisterHandlers(cliCtx client.Context, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

type listNftReq struct {
	BaseReq     rest.BaseReq            `json:"base_req"`
	DenomId     string                  `json:"denom_id"`
	NftId       string                  `json:"nft_id"`
	Price       string                  `json:"price"`
	Owner       string                  `json:"owner"`
	SplitShares []types.WeightedAddress `json:"split_shares"`
}
type editListingReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Price   string       `json:"price"`
	Owner   string       `json:"owner"`
}

type deListNftReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Owner   string       `json:"owner"`
}

type buyNftReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Price   string       `json:"price"`
	Buyer   string       `json:"buyer"`
}
