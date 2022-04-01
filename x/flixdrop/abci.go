package flixdrop

import (
	"github.com/OmniFlix/omniflixhub/x/flixdrop/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	params := k.GetParams(ctx)
	if !params.IsFlixDropStarted(ctx.BlockTime()) {
		return
	}
	// End Airdrop
	goneTime := ctx.BlockTime().Sub(params.FlixdropStartTime)
	if goneTime > params.DurationUntilDecay+params.DurationOfDecay {
		// airdrop time passed
		err := k.EndAirdrop(ctx)
		if err != nil {
			panic(err)
		}
	}
}
