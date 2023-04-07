package keeper

import (
	"encoding/binary"
	"strings"

	abci "github.com/cometbft/cometbft/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/OmniFlix/omniflixhub/v2/x/onft/types"
)

func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QuerySupply:
			return querySupply(ctx, req, k, legacyQuerierCdc)
		case types.QueryCollection:
			return queryCollection(ctx, req, k, legacyQuerierCdc)
		case types.QueryDenom:
			return queryDenom(ctx, req, k, legacyQuerierCdc)
		case types.QueryDenoms:
			return queryDenoms(ctx, req, k, legacyQuerierCdc)
		case types.QueryONFT:
			return queryONFT(ctx, req, k, legacyQuerierCdc)
		case types.QueryOwner:
			return queryOwnerONFTs(ctx, req, k, legacyQuerierCdc)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func querySupply(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QuerySupplyParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	denom := strings.ToLower(strings.TrimSpace(params.Denom))

	var supply uint64
	switch {
	case params.Owner.Empty() && len(denom) > 0:
		supply = k.GetTotalSupply(ctx, denom)
	default:
		supply = k.GetTotalSupplyOfOwner(ctx, denom, params.Owner)
	}

	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, supply)
	return bz, nil
}

func queryCollection(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryCollectionParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	denom := strings.ToLower(strings.TrimSpace(params.Denom))
	collection, err := k.GetCollection(ctx, denom)
	if err != nil {
		return nil, err
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, collection)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryDenom(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryDenomParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	denom, err := k.GetDenom(ctx, params.ID)

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, denom)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryDenoms(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	denoms := k.GetDenoms(ctx)

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, denoms)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryONFT(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryONFTParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	denom := strings.ToLower(strings.TrimSpace(params.Denom))
	tokenID := strings.ToLower(strings.TrimSpace(params.ONFTID))
	nft, err := k.GetONFT(ctx, denom, tokenID)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrUnknownONFT, "invalid oNFT %s from collection %s", params.ONFTID, params.Denom)
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, nft)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryOwnerONFTs(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryOwnerParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	owner := k.GetOwner(ctx, params.Owner, params.Denom)
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, owner)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
