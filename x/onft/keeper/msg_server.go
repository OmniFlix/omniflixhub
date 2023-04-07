package keeper

import (
	"context"

	"github.com/OmniFlix/omniflixhub/v2/x/onft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

func (m msgServer) CreateDenom(goCtx context.Context,
	msg *types.MsgCreateDenom,
) (*types.MsgCreateDenomResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	denomCreationFee := m.Keeper.GetDenomCreationFee(ctx)
	if !msg.CreationFee.Equal(denomCreationFee) {
		if msg.CreationFee.Denom != denomCreationFee.Denom {
			return nil, sdkerrors.Wrapf(types.ErrInvalidFeeDenom, "invalid creation fee denom %s",
				msg.CreationFee.Denom)
		}
		if msg.CreationFee.Amount.LT(denomCreationFee.Amount) {
			return nil, sdkerrors.Wrapf(types.ErrNotEnoughFeeAmount,
				"%s fee is not enough, to create %s fee is required", msg.CreationFee.String(), denomCreationFee.String())
		}
		return nil, sdkerrors.Wrapf(types.ErrInvalidDenomCreationFee,
			"given fee (%s) not matched with  denom creation fee. %s required to create onft denom",
			msg.CreationFee.String(), denomCreationFee.String())
	}
	if err := m.Keeper.CreateDenom(ctx,
		msg.Id,
		msg.Symbol,
		msg.Name,
		msg.Schema,
		sender,
		msg.Description,
		msg.PreviewURI,
		msg.CreationFee,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventCreateDenom{
			Id:      msg.Id,
			Symbol:  msg.Symbol,
			Name:    msg.Name,
			Creator: msg.Sender,
		},
	)

	return &types.MsgCreateDenomResponse{}, nil
}

func (m msgServer) UpdateDenom(goCtx context.Context, msg *types.MsgUpdateDenom) (*types.MsgUpdateDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	err = m.Keeper.UpdateDenom(ctx, msg.Id, msg.Name, msg.Description, msg.PreviewURI, sender)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventUpdateDenom{
			Id:      msg.Id,
			Name:    msg.Name,
			Creator: msg.Sender,
		},
	)

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
	ctx.EventManager().EmitTypedEvent(
		&types.EventTransferDenom{
			Id:        msg.Id,
			Sender:    msg.Sender,
			Recipient: msg.Recipient,
		},
	)

	return &types.MsgTransferDenomResponse{}, nil
}

func (m msgServer) MintONFT(goCtx context.Context,
	msg *types.MsgMintONFT,
) (*types.MsgMintONFTResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.MintONFT(ctx,
		msg.DenomId,
		msg.Id,
		msg.Metadata,
		msg.Data,
		msg.Transferable,
		msg.Extensible,
		msg.Nsfw,
		msg.RoyaltyShare,
		sender,
		recipient,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventMintONFT{
			Id:      msg.Id,
			DenomId: msg.DenomId,
			URI:     msg.Metadata.MediaURI,
			Owner:   msg.Recipient,
		},
	)

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

	ctx.EventManager().EmitTypedEvent(
		&types.EventTransferONFT{
			Id:        msg.Id,
			DenomId:   msg.DenomId,
			Sender:    msg.Sender,
			Recipient: msg.Recipient,
		},
	)

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

	ctx.EventManager().EmitTypedEvent(
		&types.EventBurnONFT{
			Id:      msg.Id,
			DenomId: msg.DenomId,
			Owner:   msg.Sender,
		},
	)

	return &types.MsgBurnONFTResponse{}, nil
}
