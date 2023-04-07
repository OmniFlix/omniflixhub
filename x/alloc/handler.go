package alloc

import (
	"fmt"

	"github.com/OmniFlix/omniflixhub/x/alloc/keeper"
	"github.com/OmniFlix/omniflixhub/x/alloc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler ...
func NewHandler(k keeper.Keeper) sdk.Handler {
	_ = keeper.NewMsgServerImpl(k)
	// this line is used by starport scaffolding # handler/msgServer

	return func(_ sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) { //nolint:gocritic
		// this line is used by starport scaffolding # 1
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}
