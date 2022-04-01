package flixdrop

import (
	"github.com/OmniFlix/omniflixhub/x/flixdrop/keeper"
	"github.com/OmniFlix/omniflixhub/x/flixdrop/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

// InitGenesis initializes the flixdrop module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// If its the chain genesis, set the airdrop start time to be now, and setup the needed module accounts.
	if genState.Params.FlixdropStartTime.Equal(time.Time{}) {
		genState.Params.FlixdropStartTime = ctx.BlockTime()
		k.CreateModuleAccount(ctx, genState.ModuleAccountBalance)
	}
	if err := k.SetClaimRecords(ctx, genState.ClaimRecords); err != nil {
		panic(err)
	}
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the flixdrop module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	params := k.GetParams(ctx)

	genesis := types.DefaultGenesis()
	genesis.ModuleAccountBalance = k.GetModuleAccountBalance(ctx)
	genesis.Params = params
	genesis.Actions = k.GetActions(ctx)
	genesis.ClaimRecords = k.ClaimRecords(ctx)
	return genesis
}
