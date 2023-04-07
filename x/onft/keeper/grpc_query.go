package keeper

import (
	"context"
	"strings"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/OmniFlix/omniflixhub/v2/x/onft/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Supply(c context.Context, request *types.QuerySupplyRequest) (*types.QuerySupplyResponse, error) {
	denom := strings.ToLower(strings.TrimSpace(request.DenomId))
	ctx := sdk.UnwrapSDKContext(c)

	var supply uint64
	switch {
	case len(request.Owner) == 0 && len(denom) > 0:
		supply = k.GetTotalSupply(ctx, denom)
	default:
		owner, err := sdk.AccAddressFromBech32(request.Owner)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid owner address %s", request.Owner)
		}
		supply = k.GetTotalSupplyOfOwner(ctx, denom, owner)
	}
	return &types.QuerySupplyResponse{
		Amount: supply,
	}, nil
}

func (k Keeper) Collection(c context.Context, request *types.QueryCollectionRequest) (*types.QueryCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	collection, pagination, err := k.GetPaginateCollection(ctx, request, request.DenomId)
	if err != nil {
		return nil, err
	}
	return &types.QueryCollectionResponse{
		Collection: &collection,
		Pagination: pagination,
	}, nil
}

func (k Keeper) Denom(c context.Context, request *types.QueryDenomRequest) (*types.QueryDenomResponse, error) {
	denom := strings.ToLower(strings.TrimSpace(request.DenomId))
	ctx := sdk.UnwrapSDKContext(c)

	denomObject, err := k.GetDenom(ctx, denom)
	if err != nil {
		return nil, err
	}

	return &types.QueryDenomResponse{
		Denom: &denomObject,
	}, nil
}

func (k Keeper) Denoms(c context.Context, request *types.QueryDenomsRequest) (*types.QueryDenomsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	var (
		denoms     []types.Denom
		pagination *query.PageResponse
		err        error
	)
	store := ctx.KVStore(k.storeKey)

	if request.Owner != "" {
		owner, err := sdk.AccAddressFromBech32(request.Owner)
		if err != nil {
			return nil, err
		}
		denomStore := prefix.NewStore(store, types.KeyDenomCreator(owner, ""))
		pagination, err = query.Paginate(denomStore, request.Pagination, func(key []byte, value []byte) error {
			denomId := types.MustUnMarshalDenomID(k.cdc, value)
			denom, _ := k.GetDenom(ctx, denomId)
			denoms = append(denoms, denom)
			return nil
		})
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
		}

	} else {
		denomStore := prefix.NewStore(store, types.KeyDenomID(""))
		pagination, err = query.Paginate(denomStore, request.Pagination, func(key []byte, value []byte) error {
			var denom types.Denom
			k.cdc.MustUnmarshal(value, &denom)
			denoms = append(denoms, denom)
			return nil
		})
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
		}
	}
	return &types.QueryDenomsResponse{
		Denoms:     denoms,
		Pagination: pagination,
	}, nil
}

func (k Keeper) ONFT(c context.Context, request *types.QueryONFTRequest) (*types.QueryONFTResponse, error) {
	denom := strings.ToLower(strings.TrimSpace(request.DenomId))
	onftID := strings.ToLower(strings.TrimSpace(request.Id))
	ctx := sdk.UnwrapSDKContext(c)

	nft, err := k.GetONFT(ctx, denom, onftID)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrUnknownONFT, "invalid ONFT %s from collection %s", request.Id, request.DenomId)
	}

	oNFT, ok := nft.(types.ONFT)
	if !ok {
		return nil, sdkerrors.Wrapf(types.ErrUnknownONFT, "invalid type ONFT %s from collection %s", request.Id, request.DenomId)
	}

	return &types.QueryONFTResponse{
		ONFT: &oNFT,
	}, nil
}

func (k Keeper) OwnerONFTs(c context.Context, request *types.QueryOwnerONFTsRequest) (*types.QueryOwnerONFTsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	address, err := sdk.AccAddressFromBech32(request.Owner)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid owner address %s", request.Owner)
	}

	owner := types.Owner{
		Address:       address.String(),
		IDCollections: types.IDCollections{},
	}
	var ownerCollections []types.OwnerONFTCollection
	idsMap := make(map[string][]string)
	store := ctx.KVStore(k.storeKey)
	onftStore := prefix.NewStore(store, types.KeyOwner(address, request.DenomId, ""))
	pagination, err := query.Paginate(onftStore, request.Pagination, func(key []byte, value []byte) error {
		denomId := request.DenomId
		onftId := string(key)
		if len(denomId) == 0 {
			denomId, onftId, _ = types.SplitKeyDenom(key)
		}
		if ids, ok := idsMap[denomId]; ok {
			idsMap[denomId] = append(ids, onftId)
		} else {
			idsMap[denomId] = []string{onftId}
			owner.IDCollections = append(
				owner.IDCollections,
				types.IDCollection{DenomId: denomId},
			)
		}
		return nil
	})
	for i := 0; i < len(owner.IDCollections); i++ {
		owner.IDCollections[i].OnftIds = idsMap[owner.IDCollections[i].DenomId]
		denom, _ := k.GetDenom(ctx, owner.IDCollections[i].DenomId)
		var onfts []types.ONFT
		for _, onftid := range owner.IDCollections[i].OnftIds {
			onft, _ := k.GetONFT(ctx, denom.Id, onftid)
			onfts = append(onfts, onft.(types.ONFT))
		}
		ownerCollection := types.OwnerONFTCollection{
			Denom: denom,
			Onfts: onfts,
		}
		ownerCollections = append(ownerCollections, ownerCollection)
	}
	return &types.QueryOwnerONFTsResponse{
		Owner:       address.String(),
		Collections: ownerCollections,
		Pagination:  pagination,
	}, nil
}

// Params queries params of oNFT module
func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	var params types.Params
	k.paramSpace.GetParamSet(ctx, &params)

	return &types.QueryParamsResponse{Params: params}, nil
}
