package keeper

import (
	"context"
	"strings"

	"github.com/cosmos/cosmos-sdk/x/nft"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

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
		supply = k.GetBalance(ctx, denom, owner)
	}
	return &types.QuerySupplyResponse{
		Amount: supply,
	}, nil
}

func (k Keeper) Collection(c context.Context, request *types.QueryCollectionRequest) (*types.QueryCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	denom, err := k.GetDenomInfo(ctx, request.DenomId)
	if err != nil {
		return nil, err
	}

	r := &nft.QueryNFTsRequest{
		ClassId:    request.DenomId,
		Pagination: shapePageRequest(request.Pagination),
	}

	result, err := k.nk.NFTs(c, r)
	if err != nil {
		return nil, err
	}

	var nfts []types.ONFT
	for _, _nft := range result.Nfts {
		owner := k.nk.GetOwner(ctx, request.DenomId, _nft.Id)

		nftMetadata, err := types.UnmarshalNFTMetadata(k.cdc, _nft.Data.GetValue())
		if err != nil {
			return nil, err
		}

		nfts = append(nfts, types.ONFT{
			Id: _nft.Id,
			Metadata: types.Metadata{
				Name:        nftMetadata.Name,
				Description: nftMetadata.Description,
				MediaURI:    _nft.Uri,
				UriHash:     _nft.UriHash,
				PreviewURI:  nftMetadata.PreviewURI,
			},
			Owner:        owner.String(),
			Data:         nftMetadata.Data,
			Transferable: nftMetadata.Transferable,
			Extensible:   nftMetadata.Extensible,
			CreatedAt:    nftMetadata.CreatedAt,
			Nsfw:         nftMetadata.Nsfw,
			RoyaltyShare: nftMetadata.RoyaltyShare,
		})
	}

	collection := &types.Collection{
		Denom: *denom,
		ONFTs: nfts,
	}

	response := &types.QueryCollectionResponse{
		Collection: collection,
		Pagination: result.Pagination,
	}

	return response, nil
}

func (k Keeper) IBCCollection(c context.Context, request *types.QueryIBCCollectionRequest) (*types.QueryCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	denom, err := k.GetDenomInfo(ctx, "ibc/"+request.Hash)
	if err != nil {
		return nil, err
	}

	r := &nft.QueryNFTsRequest{
		ClassId:    denom.Id,
		Pagination: shapePageRequest(request.Pagination),
	}

	result, err := k.nk.NFTs(c, r)
	if err != nil {
		return nil, err
	}
	var onfts []types.ONFT
	for _, _nft := range result.Nfts {
		owner := k.nk.GetOwner(ctx, denom.Id, _nft.Id)

		nftMetadata, err := types.UnmarshalNFTMetadata(k.cdc, _nft.Data.GetValue())
		if err != nil {
			return nil, err
		}

		onfts = append(onfts, types.ONFT{
			Id: _nft.Id,
			Metadata: types.Metadata{
				Name:        nftMetadata.Name,
				Description: nftMetadata.Description,
				MediaURI:    _nft.Uri,
				UriHash:     _nft.UriHash,
				PreviewURI:  nftMetadata.PreviewURI,
			},
			Owner:        owner.String(),
			Data:         nftMetadata.Data,
			Transferable: nftMetadata.Transferable,
			Extensible:   nftMetadata.Extensible,
			CreatedAt:    nftMetadata.CreatedAt,
			Nsfw:         nftMetadata.Nsfw,
			RoyaltyShare: nftMetadata.RoyaltyShare,
		})
	}

	collection := &types.Collection{
		Denom: *denom,
		ONFTs: onfts,
	}
	return &types.QueryCollectionResponse{
		Collection: collection,
		Pagination: result.Pagination,
	}, nil
}

func (k Keeper) Denom(c context.Context, request *types.QueryDenomRequest) (*types.QueryDenomResponse, error) {
	denom := strings.ToLower(strings.TrimSpace(request.DenomId))
	ctx := sdk.UnwrapSDKContext(c)

	denomObject, err := k.GetDenomInfo(ctx, denom)
	if err != nil {
		return nil, err
	}

	return &types.QueryDenomResponse{
		Denom: denomObject,
	}, nil
}

func (k Keeper) IBCDenom(c context.Context, request *types.QueryIBCDenomRequest) (*types.QueryDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	denomObject, err := k.GetDenomInfo(ctx, "ibc/"+request.Hash)
	if err != nil {
		return nil, err
	}

	return &types.QueryDenomResponse{
		Denom: denomObject,
	}, nil
}

func (k Keeper) Denoms(c context.Context, request *types.QueryDenomsRequest) (*types.QueryDenomsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	result, err := k.nk.Classes(c, &nft.QueryClassesRequest{
		Pagination: shapePageRequest(request.Pagination),
	})
	if err != nil {
		return nil, err
	}
	var denoms []types.Denom
	if request.Owner != "" {
		owner, err := sdk.AccAddressFromBech32(request.Owner)
		if err != nil {
			return nil, err
		}
		for _, class := range result.Classes {
			denom, err := k.GetDenomInfo(ctx, class.Id)
			if err != nil {
				return nil, err
			}
			if denom.Creator == owner.String() {
				denoms = append(denoms, *denom)
			}
		}
	} else {
		for _, class := range result.Classes {
			denom, err := k.GetDenomInfo(ctx, class.Id)
			if err != nil {
				return nil, err
			}
			denoms = append(denoms, *denom)
		}
	}
	return &types.QueryDenomsResponse{
		Denoms:     denoms,
		Pagination: result.Pagination,
	}, nil
}

func (k Keeper) ONFT(c context.Context, request *types.QueryONFTRequest) (*types.QueryONFTResponse, error) {
	denom := strings.ToLower(strings.TrimSpace(request.DenomId))
	onftID := strings.ToLower(strings.TrimSpace(request.Id))
	ctx := sdk.UnwrapSDKContext(c)

	onft, err := k.GetONFT(ctx, denom, onftID)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrUnknownONFT, "invalid ONFT %s from collection %s", request.Id, request.DenomId)
	}

	oNFT, ok := onft.(types.ONFT)
	if !ok {
		return nil, errorsmod.Wrapf(types.ErrUnknownONFT, "invalid type ONFT %s from collection %s", request.Id, request.DenomId)
	}

	return &types.QueryONFTResponse{
		ONFT: &oNFT,
	}, nil
}

func (k Keeper) OwnerONFTs(c context.Context, request *types.QueryOwnerONFTsRequest) (*types.QueryOwnerONFTsResponse, error) {
	r := &nft.QueryNFTsRequest{
		ClassId:    request.DenomId,
		Owner:      request.Owner,
		Pagination: shapePageRequest(request.Pagination),
	}

	result, err := k.nk.NFTs(c, r)
	if err != nil {
		return nil, err
	}

	denomMap := make(map[string][]string)
	var denoms []string
	for _, _nft := range result.Nfts {
		if denomMap[_nft.ClassId] == nil {
			denomMap[_nft.ClassId] = []string{}
			denoms = append(denoms, _nft.ClassId)
		}
		denomMap[_nft.ClassId] = append(denomMap[_nft.ClassId], _nft.Id)
	}

	var idc []types.IDCollection
	for _, denomID := range denoms {
		idc = append(idc, types.IDCollection{
			DenomId: denomID,
			OnftIds: denomMap[denomID],
		})
	}

	response := &types.QueryOwnerONFTsResponse{
		Owner: &types.Owner{
			Address:       request.Owner,
			IDCollections: idc,
		},
		Pagination: result.Pagination,
	}

	return response, nil
}

// Params queries params of oNFT module
func (k Keeper) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)
	return &types.QueryParamsResponse{Params: params}, nil
}
