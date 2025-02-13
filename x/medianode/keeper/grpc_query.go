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

var _ types.QueryServer = Keeper{}

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

// AvailableNodes returns not leased active media nodes
func (k Keeper) AvailableNodes(goCtx context.Context,
	req *types.QueryAvailableNodesRequest,
) (*types.QueryAvailableNodesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	var availableNodes []types.MediaNode
	var pageRes *query.PageResponse
	store := ctx.KVStore(k.storeKey)
	mediaNodeStore := prefix.NewStore(store, types.PrefixMediaNode)
	pageRes, err := query.FilteredPaginate(mediaNodeStore,
		req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
			var mn types.MediaNode
			k.cdc.MustUnmarshal(value, &mn)

			// Check if the node is active and not leased
			if mn.Status == types.STATUS_ACTIVE && !mn.Leased {
				if accumulate {
					availableNodes = append(availableNodes, mn)
				}
				return true, nil
			}

			return false, nil
		})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}

	return &types.QueryAvailableNodesResponse{MediaNodes: availableNodes, Pagination: pageRes}, nil
}

// Lease returns the lease of a given media node if it exists
func (k Keeper) Lease(goCtx context.Context,
	req *types.QueryLeaseRequest,
) (*types.QueryLeaseResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	lease, found := k.GetMediaNodeLease(ctx, req.MediaNodeId) // Assuming GetLease is a method that retrieves the lease by media node ID
	if !found {
		return nil, status.Errorf(codes.NotFound, "lease for media node %d not found", req.MediaNodeId)
	}

	return &types.QueryLeaseResponse{Lease: lease}, nil
}

// LeasesByLeasee returns all active leases for a specific leasee
func (k Keeper) LeasesByLessee(goCtx context.Context,
	req *types.QueryLeasesByLesseeRequest,
) (*types.QueryLeasesByLesseeResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	var leases []types.Lease // Assuming Lease is the type for lease
	var pageRes *query.PageResponse
	store := ctx.KVStore(k.storeKey)
	leaseStore := prefix.NewStore(store, types.PrefixLease) // Assuming PrefixLease is defined for leases
	pageRes, err := query.FilteredPaginate(leaseStore,
		req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
			var lease types.Lease
			k.cdc.MustUnmarshal(value, &lease)

			// Check if the lease is active and belongs to the leasee
			if lease.Leasee == req.Lessee && lease.Status == types.LEASE_STATUS_ACTIVE { // Assuming STATUS_ACTIVE is defined
				if accumulate {
					leases = append(leases, lease)
				}
				return true, nil
			}

			return false, nil
		})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}

	return &types.QueryLeasesByLesseeResponse{Leases: leases, Pagination: pageRes}, nil
}
