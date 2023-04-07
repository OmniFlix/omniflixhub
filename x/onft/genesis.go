package onft

import (
	"github.com/OmniFlix/omniflixhub/v2/x/onft/keeper"
	"github.com/OmniFlix/omniflixhub/v2/x/onft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	if err := types.ValidateGenesis(data); err != nil {
		panic(err.Error())
	}

	for _, c := range data.Collections {
		if err := k.SetCollection(ctx, c); err != nil {
			panic(err)
		}
	}
	k.SetParams(ctx, data.Params)
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return types.NewGenesisState(k.GetCollections(ctx), k.GetParams(ctx))
}

func DefaultGenesisState() *types.GenesisState {
	return types.NewGenesisState([]types.Collection{}, types.DefaultParams())
}
