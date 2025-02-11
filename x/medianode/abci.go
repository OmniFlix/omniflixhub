package medianode

import (
	"github.com/OmniFlix/omniflixhub/v6/x/medianode/keeper"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func EndBlock(ctx sdk.Context, k keeper.Keeper) ([]abcitypes.ValidatorUpdate, error) {
	log := k.Logger(ctx)
	log.Info(" Medianode payments settled.. ")
	if err := k.SettleActiveLeases(ctx); err != nil {
		panic(err)
	}
	return []abcitypes.ValidatorUpdate{}, nil
}
