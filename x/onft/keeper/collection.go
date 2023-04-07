package keeper

import (
	"github.com/OmniFlix/omniflixhub/v2/x/onft/exported"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/OmniFlix/omniflixhub/v2/x/onft/types"
)

func (k Keeper) SetCollection(ctx sdk.Context, collection types.Collection) error {
	denom := collection.Denom
	creator, err := sdk.AccAddressFromBech32(denom.Creator)
	if err != nil {
		return err
	}
	if k.HasDenomID(ctx, denom.Id) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denomID %s has already exists", denom.Id)
	}

	if k.HasDenomSymbol(ctx, denom.Symbol) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denomSymbol %s has already exists", denom.Symbol)
	}
	err = k.SetDenom(ctx, types.NewDenom(denom.Id, denom.Symbol, denom.Name, denom.Schema,
		creator, denom.Description, denom.PreviewURI))
	if err != nil {
		return err
	}
	k.setDenomOwner(ctx, denom.Id, creator)

	for _, onft := range collection.ONFTs {
		metadata := types.Metadata{
			Name:        onft.GetName(),
			Description: onft.GetDescription(),
			MediaURI:    onft.GetMediaURI(),
			PreviewURI:  onft.GetPreviewURI(),
		}

		if err := k.MintONFT(ctx,
			collection.Denom.Id,
			onft.GetID(),
			metadata,
			onft.GetData(),
			onft.IsTransferable(),
			onft.IsExtensible(),
			onft.IsNSFW(),
			onft.RoyaltyShare,
			creator,
			onft.GetOwner(),
		); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) GetCollection(ctx sdk.Context, denomID string) (types.Collection, error) {
	denom, err := k.GetDenom(ctx, denomID)
	if err != nil {
		return types.Collection{}, sdkerrors.Wrapf(types.ErrInvalidDenom, "denomID %s not existed ", denomID)
	}

	onfts := k.GetONFTs(ctx, denomID)
	return types.NewCollection(denom, onfts), nil
}

func (k Keeper) GetCollections(ctx sdk.Context) (cs []types.Collection) {
	for _, denom := range k.GetDenoms(ctx) {
		onfts := k.GetONFTs(ctx, denom.Id)
		cs = append(cs, types.NewCollection(denom, onfts))
	}
	return cs
}

func (k Keeper) GetPaginateCollection(ctx sdk.Context,
	request *types.QueryCollectionRequest, denomId string,
) (types.Collection, *query.PageResponse, error) {
	denom, err := k.GetDenom(ctx, denomId)
	if err != nil {
		return types.Collection{}, nil, sdkerrors.Wrapf(types.ErrInvalidDenom, "denomId %s not existed ", denomId)
	}
	var onfts []exported.ONFT
	store := ctx.KVStore(k.storeKey)
	onftStore := prefix.NewStore(store, types.KeyONFT(denomId, ""))
	pagination, err := query.Paginate(onftStore, request.Pagination, func(key []byte, value []byte) error {
		var oNFT types.ONFT
		k.cdc.MustUnmarshal(value, &oNFT)
		onfts = append(onfts, oNFT)
		return nil
	})
	if err != nil {
		return types.Collection{}, nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}
	return types.NewCollection(denom, onfts), pagination, nil
}

func (k Keeper) GetTotalSupply(ctx sdk.Context, denomID string) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyCollection(denomID))
	if len(bz) == 0 {
		return 0
	}
	return types.MustUnMarshalSupply(k.cdc, bz)
}

func (k Keeper) GetTotalSupplyOfOwner(ctx sdk.Context, id string, owner sdk.AccAddress) (supply uint64) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyOwner(owner, id, ""))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		supply++
	}
	return supply
}

func (k Keeper) increaseSupply(ctx sdk.Context, denomID string) {
	supply := k.GetTotalSupply(ctx, denomID)
	supply++

	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalSupply(k.cdc, supply)
	store.Set(types.KeyCollection(denomID), bz)
}

func (k Keeper) decreaseSupply(ctx sdk.Context, denomID string) {
	supply := k.GetTotalSupply(ctx, denomID)
	supply--

	store := ctx.KVStore(k.storeKey)
	if supply == 0 {
		store.Delete(types.KeyCollection(denomID))
		return
	}

	bz := types.MustMarshalSupply(k.cdc, supply)
	store.Set(types.KeyCollection(denomID), bz)
}
