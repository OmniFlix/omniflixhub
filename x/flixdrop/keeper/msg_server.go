package keeper

import (
	"context"

	"github.com/OmniFlix/omniflixhub/x/flixdrop/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) Claim(goCtx context.Context, msg *types.MsgClaim) (*types.MsgClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	address, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, err
	}
	params := k.GetParams(ctx)
	if !params.IsFlixDropStarted(ctx.BlockTime()) {
		return nil, types.ErrFlixDropNotStarted
	}
	if !k.IsActionClaimStarted(ctx, msg.Action) {
		return nil, types.ErrActionClaimNotStarted
	}
	coins, err := k.Keeper.ClaimCoinsForAction(ctx, address, msg.GetAction())
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})
	return &types.MsgClaimResponse{
		Address:       msg.Address,
		ClaimedAmount: coins,
	}, nil
}
