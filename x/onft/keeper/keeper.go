package keeper

import (
	"fmt"

	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	errorsmod "cosmossdk.io/errors"
	"github.com/cometbft/cometbft/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/OmniFlix/omniflixhub/v2/x/onft/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec

	accountKeeper      types.AccountKeeper
	bankKeeper         types.BankKeeper
	distributionKeeper types.DistributionKeeper
	paramSpace         paramstypes.Subspace
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	distrKeeper types.DistributionKeeper,
	paramSpace paramstypes.Subspace,
) Keeper {
	// ensure oNFT module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:           storeKey,
		cdc:                cdc,
		accountKeeper:      accountKeeper,
		bankKeeper:         bankKeeper,
		distributionKeeper: distrKeeper,
		paramSpace:         paramSpace,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("OmniFlix/%s", types.ModuleName))
}

func (k Keeper) CreateDenom(
	ctx sdk.Context, id, symbol, name, schema string,
	creator sdk.AccAddress, description, previewUri string, fee sdk.Coin,
) error {
	if k.HasDenomID(ctx, id) {
		return errorsmod.Wrapf(types.ErrInvalidDenom, "denomID %s has already exists", id)
	}

	if k.HasDenomSymbol(ctx, symbol) {
		return errorsmod.Wrapf(types.ErrInvalidDenom, "denomSymbol %s has already exists", symbol)
	}

	err := k.distributionKeeper.FundCommunityPool(
		ctx,
		sdk.NewCoins(fee),
		creator,
	)
	if err != nil {
		return err
	}
	err = k.SetDenom(ctx, types.NewDenom(id, symbol, name, schema, creator, description, previewUri))
	if err != nil {
		return err
	}
	k.setDenomOwner(ctx, id, creator)
	return nil
}

func (k Keeper) UpdateDenom(ctx sdk.Context, id, name, description, previewURI string, sender sdk.AccAddress) error {
	if !k.HasDenomID(ctx, id) {
		return errorsmod.Wrapf(types.ErrInvalidDenom, "denom id %s not exists", id)
	}
	denom, err := k.AuthorizeDenomCreator(ctx, id, sender)
	if err != nil {
		return err
	}
	if len(name) > 0 && name != types.DoNotModify {
		denom.Name = name
	}
	if len(description) > 0 && description != types.DoNotModify {
		denom.Description = description
	}
	if len(previewURI) > 0 && previewURI != types.DoNotModify {
		denom.PreviewURI = previewURI
	}
	return k.SetDenom(ctx, denom)
}

func (k Keeper) TransferDenomOwner(ctx sdk.Context, id string, curOwner, newOwner sdk.AccAddress) error {
	denom, err := k.GetDenom(ctx, id)
	if err != nil {
		return errorsmod.Wrapf(types.ErrInvalidDenom, "denom ID %s not exists", id)
	}

	if curOwner.String() != denom.Creator {
		return errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "unauthorized address %s", curOwner.String())
	}

	denom.Creator = newOwner.String()

	err = k.SetDenom(ctx, denom)
	if err != nil {
		return err
	}
	k.swapDenomOwner(ctx, id, curOwner, newOwner)
	return nil
}

func (k Keeper) MintONFT(
	ctx sdk.Context, denomID, onftID string,
	metadata types.Metadata, data string, transferable, extensible, nsfw bool,
	royaltyShare sdk.Dec, sender, recipient sdk.AccAddress,
) error {
	if !k.HasPermissionToMint(ctx, denomID, sender) {
		return errorsmod.Wrapf(types.ErrUnauthorized, "only creator of denom has permission to mint")
	}
	if !k.HasDenomID(ctx, denomID) {
		return errorsmod.Wrapf(types.ErrInvalidDenom, "denomID %s not exists", denomID)
	}

	if k.HasONFT(ctx, denomID, onftID) {
		return errorsmod.Wrapf(types.ErrONFTAlreadyExists, "ONFT %s already exists in collection %s", onftID, denomID)
	}

	k.setONFT(ctx, denomID, types.NewONFT(
		onftID,
		metadata,
		data,
		transferable,
		extensible,
		recipient,
		ctx.BlockHeader().Time,
		nsfw,
		royaltyShare,
	))
	k.setOwner(ctx, denomID, onftID, recipient)
	k.increaseSupply(ctx, denomID)
	return nil
}

func (k Keeper) EditONFT(ctx sdk.Context, denomID, onftID string, owner sdk.AccAddress) error {
	if !k.HasDenomID(ctx, denomID) {
		return errorsmod.Wrapf(types.ErrInvalidDenom, "denomID %s not exists", denomID)
	}
	_, err := k.GetDenom(ctx, denomID)
	if err != nil {
		return err
	}

	onft, err := k.Authorize(ctx, denomID, onftID, owner)
	if err != nil {
		return err
	}

	k.setONFT(ctx, denomID, onft)
	return nil
}

func (k Keeper) TransferOwnership(ctx sdk.Context, denomID, onftID string, srcOwner, dstOwner sdk.AccAddress) error {
	if !k.HasDenomID(ctx, denomID) {
		return errorsmod.Wrapf(types.ErrInvalidDenom, "denomID %s not exists", denomID)
	}

	onft, err := k.Authorize(ctx, denomID, onftID, srcOwner)
	if err != nil {
		return err
	}
	if !onft.IsTransferable() {
		return errorsmod.Wrap(types.ErrNotTransferable, onft.GetID())
	}

	onft.Owner = dstOwner.String()

	k.setONFT(ctx, denomID, onft)
	k.swapOwner(ctx, denomID, onftID, srcOwner, dstOwner)
	return nil
}

func (k Keeper) BurnONFT(ctx sdk.Context,
	denomID, onftID string,
	owner sdk.AccAddress,
) error {
	if !k.HasDenomID(ctx, denomID) {
		return errorsmod.Wrapf(types.ErrInvalidDenom, "denomID %s not exists", denomID)
	}

	onft, err := k.Authorize(ctx, denomID, onftID, owner)
	if err != nil {
		return err
	}

	k.deleteONFT(ctx, denomID, onft)
	k.deleteOwner(ctx, denomID, onftID, owner)
	k.decreaseSupply(ctx, denomID)
	return nil
}
