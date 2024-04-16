package keeper

import (
	"context"

	"github.com/OmniFlix/omniflixhub/v4/x/globalfee/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.QueryServer = &GrpcQuerier{}

type GrpcQuerier struct {
	keeper Keeper
}

func NewGrpcQuerier(k Keeper) GrpcQuerier {
	return GrpcQuerier{
		keeper: k,
	}
}

// Params return global-fee params
func (g GrpcQuerier) Params(stdCtx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(stdCtx)

	p := g.keeper.GetParams(ctx)

	return &types.QueryParamsResponse{
		Params: p,
	}, nil
}
