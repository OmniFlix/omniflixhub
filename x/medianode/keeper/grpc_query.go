package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/OmniFlix/omniflixhub/v6/x/medianode/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Params returns params of the medianode module.
func (k Keeper) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

// MediaNode returns details of the media node
func (k Keeper) MediaNode(goCtx context.Context,
	req *types.QueryMediaNodeRequest,
) (*types.QueryMediaNodeResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	mediaNode, found := k.GetMediaNode(ctx, req.Id)
	if !found {
		return nil, status.Errorf(codes.NotFound, "media node %d not found", req.Id)
	}

	return &types.QueryMediaNodeResponse{MediaNode: mediaNode}, nil
}

// MediaNodes returns all media nodes with optional pagination
func (k Keeper) MediaNodes(goCtx context.Context,
	req *types.QueryMediaNodesRequest,
) (*types.QueryMediaNodesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	var mediaNodes []types.MediaNode
	var pageRes *query.PageResponse
	store := ctx.KVStore(k.storeKey)
	mediaNodeStore := prefix.NewStore(store, types.PrefixMediaNode)
	pageRes, err := query.FilteredPaginate(mediaNodeStore,
		req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
			var mn types.MediaNode
			k.cdc.MustUnmarshal(value, &mn)

			if accumulate {
				mediaNodes = append(mediaNodes, mn)
			}

			return true, nil
		})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}

	return &types.QueryMediaNodesResponse{MediaNodes: mediaNodes, Pagination: pageRes}, nil
}

// MediaNodesByOwner returns media nodes owned by a specific address
func (k Keeper) MediaNodesByOwner(goCtx context.Context,
	req *types.QueryMediaNodesByOwnerRequest,
) (*types.QueryMediaNodesByOwnerResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	var mediaNodes []types.MediaNode
	var pageRes *query.PageResponse
	store := ctx.KVStore(k.storeKey)
	mediaNodeStore := prefix.NewStore(store, types.PrefixMediaNode)
	pageRes, err := query.FilteredPaginate(mediaNodeStore,
		req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
			var mn types.MediaNode
			k.cdc.MustUnmarshal(value, &mn)

			if mn.Owner == req.Owner {
				if accumulate {
					mediaNodes = append(mediaNodes, mn)
				}
				return true, nil
			}

			return false, nil
		})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}

	return &types.QueryMediaNodesByOwnerResponse{MediaNodes: mediaNodes, Pagination: pageRes}, nil
}
