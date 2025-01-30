package itc

import (
	"github.com/OmniFlix/omniflixhub/v6/x/medianode/keeper"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func EndBlock(ctx sdk.Context, k keeper.Keeper) ([]abcitypes.ValidatorUpdate, error) {
	log := k.Logger(ctx)
	log.Info(" Medianode payments settled.. ")
	return []abcitypes.ValidatorUpdate{}, nil
}
