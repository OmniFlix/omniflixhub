package keeper

import (
	errorsmod "cosmossdk.io/errors"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/nft"

	"github.com/OmniFlix/omniflixhub/v4/x/onft/types"
)

// SaveDenom saves a denom
func (k Keeper) SaveDenom(
	ctx sdk.Context,
	id,
	symbol,
	name,
	schema string,
	creator sdk.AccAddress,
	description,
	previewUri string,
	uri,
	uriHash,
	data string,
	royaltyReceivers []*types.WeightedAddress,
) error {
	denomMetadata := &types.DenomMetadata{
		Creator:          creator.String(),
		Schema:           schema,
		PreviewUri:       previewUri,
		Data:             data,
		RoyaltyReceivers: royaltyReceivers,
	}
	metadata, err := codectypes.NewAnyWithValue(denomMetadata)
	if err != nil {
		return err
	}
	err = k.nk.SaveClass(ctx, nft.Class{
		Id:          id,
		Name:        name,
		Symbol:      symbol,
		Description: description,
		Uri:         uri,
		UriHash:     uriHash,
		Data:        metadata,
	})
	if err != nil {
		return err
	}
	// emit events
	k.emitCreateONFTDenomEvent(ctx, id, symbol, name, creator.String())

	return nil
}

// TransferDenomOwner transfers the ownership to new address
func (k Keeper) TransferDenomOwner(
	ctx sdk.Context,
	denomID string,
	srcOwner,
	dstOwner sdk.AccAddress,
) error {
	denom, err := k.GetDenomInfo(ctx, denomID)
	if err != nil {
		return err
	}
	sender := srcOwner.String()
	recipient := dstOwner.String()

	// authorize
	if sender != denom.Creator {
		return errorsmod.Wrapf(
			sdkerrors.ErrUnauthorized,
			"%s is not allowed to transfer denom %s", sender,
			denomID,
		)
	}

	denomMetadata := &types.DenomMetadata{
		Creator:          recipient,
		Schema:           denom.Schema,
		PreviewUri:       denom.PreviewURI,
		Data:             denom.Data,
		RoyaltyReceivers: denom.RoyaltyReceivers,
	}
	data, err := codectypes.NewAnyWithValue(denomMetadata)
	if err != nil {
		return err
	}
	class := nft.Class{
		Id:          denom.Id,
		Name:        denom.Name,
		Symbol:      denom.Symbol,
		Description: denom.Description,
		Uri:         denom.Uri,
		UriHash:     denom.UriHash,
		Data:        data,
	}

	err = k.nk.UpdateClass(ctx, class)
	if err != nil {
		return err
	}
	k.emitTransferONFTDenomEvent(ctx, denomID, denom.Symbol, sender, recipient)
	return nil
}

func (k Keeper) HasDenom(ctx sdk.Context, id string) bool {
	return k.nk.HasClass(ctx, id)
}

func (k Keeper) UpdateDenom(ctx sdk.Context, msg *types.MsgUpdateDenom) error {
	denom, err := k.GetDenomInfo(ctx, msg.Id)
	if err != nil {
		return err
	}

	// authorize
	if msg.Sender != denom.Creator {
		return errorsmod.Wrapf(
			sdkerrors.ErrUnauthorized,
			"%s is not allowed to transfer denom %s", msg.Sender,
			denom.Id,
		)
	}

	denomMetadata := &types.DenomMetadata{
		Creator:          denom.Creator,
		Schema:           denom.Schema,
		PreviewUri:       denom.PreviewURI,
		Data:             denom.Data,
		RoyaltyReceivers: denom.RoyaltyReceivers,
	}
	if msg.PreviewURI != types.DoNotModify {
		denomMetadata.PreviewUri = msg.PreviewURI
	}
	if msg.RoyaltyReceivers != nil {
		denomMetadata.RoyaltyReceivers = msg.RoyaltyReceivers
	}
	data, err := codectypes.NewAnyWithValue(denomMetadata)
	if err != nil {
		return err
	}
	class := nft.Class{
		Id:          denom.Id,
		Name:        denom.Name,
		Symbol:      denom.Symbol,
		Description: denom.Description,
		Uri:         denom.Uri,
		UriHash:     denom.UriHash,
		Data:        data,
	}
	if msg.Name != types.DoNotModify {
		class.Name = msg.Name
	}
	if msg.Description != types.DoNotModify {
		class.Description = msg.Description
	}
	k.emitUpdateONFTDenomEvent(ctx, class.Id, class.Name, class.Description, denomMetadata.PreviewUri, msg.Sender)
	return k.nk.UpdateClass(ctx, class)
}

func (k Keeper) GetDenoms(ctx sdk.Context) (denoms []types.Denom, err error) {
	classes := k.nk.GetClasses(ctx)
	for _, class := range classes {
		var denomMetadata types.DenomMetadata
		if err := k.cdc.Unmarshal(class.Data.GetValue(), &denomMetadata); err != nil {
			return nil, err
		}
		denoms = append(denoms, types.Denom{
			Id:               class.Id,
			Name:             class.Name,
			Schema:           denomMetadata.Schema,
			Creator:          denomMetadata.Creator,
			Symbol:           class.Symbol,
			Description:      class.Description,
			PreviewURI:       denomMetadata.PreviewUri,
			Uri:              class.Uri,
			UriHash:          class.UriHash,
			RoyaltyReceivers: denomMetadata.RoyaltyReceivers,
		})
	}
	return denoms, nil
}

func (k Keeper) AuthorizeDenomCreator(ctx sdk.Context, id string, creator sdk.AccAddress) error {
	denom, err := k.GetDenomInfo(ctx, id)
	if err != nil {
		return err
	}

	if creator.String() != denom.Creator {
		return errorsmod.Wrap(types.ErrUnauthorized, creator.String())
	}
	return nil
}

func (k Keeper) HasPermissionToMint(ctx sdk.Context, denomID string, sender sdk.AccAddress) bool {
	denom, err := k.GetDenomInfo(ctx, denomID)
	if err != nil {
		return false
	}

	if sender.String() == denom.Creator {
		return true
	}
	return false
}

func (k Keeper) GetDenomInfo(ctx sdk.Context, denomID string) (*types.Denom, error) {
	class, ok := k.nk.GetClass(ctx, denomID)
	if !ok {
		return nil, errorsmod.Wrapf(types.ErrInvalidDenom, "denom ID %s not exists", denomID)
	}

	var denomMetadata types.DenomMetadata
	if err := k.cdc.Unmarshal(class.Data.GetValue(), &denomMetadata); err != nil {
		return nil, err
	}
	return &types.Denom{
		Id:               class.Id,
		Name:             class.Name,
		Schema:           denomMetadata.Schema,
		Creator:          denomMetadata.Creator,
		Symbol:           class.Symbol,
		Description:      class.Description,
		PreviewURI:       denomMetadata.PreviewUri,
		Uri:              class.Uri,
		UriHash:          class.UriHash,
		Data:             denomMetadata.Data,
		RoyaltyReceivers: denomMetadata.RoyaltyReceivers,
	}, nil
}

// PurgeDenom deletes the denom if no nfts in it
func (k Keeper) PurgeDenom(
	ctx sdk.Context,
	denomID string,
	sender sdk.AccAddress,
) error {
	denom, err := k.GetDenomInfo(ctx, denomID)
	if err != nil {
		return err
	}

	// authorize
	if sender.String() != denom.Creator {
		return errorsmod.Wrapf(
			sdkerrors.ErrUnauthorized,
			"%s is not allowed to purge denom %s", sender,
			denomID,
		)
	}
	if k.nk.GetTotalSupply(ctx, denomID) != 0 {
		return errorsmod.Wrapf(
			types.ErrNotAllowed,
			"can not purge denom (%s) with nfts available in it",
			denomID,
		)
	}
	// delete the denom
	k.DeleteDenomFromStore(ctx, denomID)

	k.emitPurgeONFTDenomEvent(ctx, denomID)
	return nil
}

func (k Keeper) DeleteDenomFromStore(ctx sdk.Context, denomId string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(classStoreKey(denomId))
}
