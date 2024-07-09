package itc

import (
	"github.com/OmniFlix/omniflixhub/v5/x/itc/keeper"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func EndBlock(ctx context.Context, k keeper.Keeper) []abcitypes.ValidatorUpdate {
	log := k.Logger(ctx)
	err := k.FinalizeAndEndCampaigns(ctx)
	if err != nil {
		panic(err)
	}
	log.Info("Updated and processed campaigns .. ")
	return []abcitypes.ValidatorUpdate{}
}
