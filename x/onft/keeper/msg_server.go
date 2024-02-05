package keeper

import (
	"context"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	errorsmod "cosmossdk.io/errors"
	"github.com/OmniFlix/omniflixhub/v2/x/onft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the NFT MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if m.authority != req.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", m.authority, req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.SetParams(ctx, req.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}

func (m msgServer) CreateDenom(goCtx context.Context, msg *types.MsgCreateDenom) (*types.MsgCreateDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	if m.Keeper.HasDenom(ctx, msg.Id) {
		return nil, errorsmod.Wrapf(types.ErrDenomIdExists, "denom id already exists %s", msg.Id)
	}
	denomCreationFee := m.Keeper.GetDenomCreationFee(ctx)
	if !msg.CreationFee.Equal(denomCreationFee) {
		if msg.CreationFee.Denom != denomCreationFee.Denom {
			return nil, errorsmod.Wrapf(types.ErrInvalidFeeDenom, "invalid creation fee denom %s",
				msg.CreationFee.Denom)
		}
		if msg.CreationFee.Amount.LT(denomCreationFee.Amount) {
			return nil, errorsmod.Wrapf(
				types.ErrNotEnoughFeeAmount,
				"%s fee is not enough, to create %s fee is required",
				msg.CreationFee.String(),
				denomCreationFee.String(),
			)
		}
		return nil, errorsmod.Wrapf(
			types.ErrInvalidDenomCreationFee,
			"given fee (%s) not matched with  denom creation fee. %s required to create onft denom",
			msg.CreationFee.String(),
			denomCreationFee.String(),
		)
	}
	err = m.Keeper.distributionKeeper.FundCommunityPool(
		ctx,
		sdk.NewCoins(denomCreationFee),
		sender,
	)
	if err != nil {
		return nil, err
	}
	if err := m.Keeper.SaveDenom(ctx,
		msg.Id,
		msg.Symbol,
		msg.Name,
		msg.Schema,
		sender,
		msg.Description,
		msg.PreviewURI,
		msg.Uri,
		msg.UriHash,
		msg.Data,
		msg.RoyaltyReceivers,
	); err != nil {
		return nil, err
	}

	return &types.MsgCreateDenomResponse{}, nil
}

func (m msgServer) UpdateDenom(goCtx context.Context, msg *types.MsgUpdateDenom) (*types.MsgUpdateDenomResponse, error) {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	err = m.Keeper.UpdateDenom(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgUpdateDenomResponse{}, nil
}

func (m msgServer) TransferDenom(goCtx context.Context, msg *types.MsgTransferDenom) (*types.MsgTransferDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return nil, err
	}

	err = m.Keeper.TransferDenomOwner(ctx, msg.Id, sender, recipient)
	if err != nil {
		return nil, err
	}

	return &types.MsgTransferDenomResponse{}, nil
}

func (m msgServer) MintONFT(goCtx context.Context, msg *types.MsgMintONFT) (*types.MsgMintONFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return nil, err
	}
	if !m.Keeper.HasPermissionToMint(ctx, msg.DenomId, sender) {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrUnauthorized,
			"%s is not allowed to mint nft under denom %s",
			sender.String(),
			msg.DenomId,
		)
	}

	if m.Keeper.HasONFT(ctx, msg.DenomId, msg.Id) {
		return nil, errorsmod.Wrapf(
			types.ErrONFTAlreadyExists,
			"ONFT with id %s already exists in collection %s", msg.Id, msg.DenomId)
	}
	if err := m.Keeper.MintONFT(ctx,
		msg.DenomId,
		msg.Id,
		msg.Metadata.Name,
		msg.Metadata.Description,
		msg.Metadata.MediaURI,
		msg.Metadata.UriHash,
		msg.Metadata.PreviewURI,
		msg.Data,
		ctx.BlockTime(),
		msg.Transferable,
		msg.Extensible,
		msg.Nsfw,
		msg.RoyaltyShare,
		recipient,
	); err != nil {
		return nil, err
	}

	return &types.MsgMintONFTResponse{}, nil
}

func (m msgServer) TransferONFT(goCtx context.Context,
	msg *types.MsgTransferONFT,
) (*types.MsgTransferONFTResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.TransferOwnership(ctx, msg.DenomId, msg.Id,
		sender,
		recipient,
	); err != nil {
		return nil, err
	}

	return &types.MsgTransferONFTResponse{}, nil
}

func (m msgServer) BurnONFT(goCtx context.Context,
	msg *types.MsgBurnONFT,
) (*types.MsgBurnONFTResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.BurnONFT(ctx, msg.DenomId, msg.Id, sender); err != nil {
		return nil, err
	}

	return &types.MsgBurnONFTResponse{}, nil
}

func (m msgServer) PurgeDenom(goCtx context.Context,
	msg *types.MsgPurgeDenom,
) (*types.MsgPurgeDenomResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.PurgeDenom(ctx, msg.Id, sender); err != nil {
		return nil, err
	}

	return &types.MsgPurgeDenomResponse{}, nil
}
