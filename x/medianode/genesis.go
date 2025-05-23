package medianode

import (
	"fmt"

	"github.com/OmniFlix/omniflixhub/v6/x/medianode/keeper"
	"github.com/OmniFlix/omniflixhub/v6/x/medianode/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	if err := genState.ValidateGenesis(); err != nil {
		panic(err.Error())
	}

	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}

	for _, mn := range genState.Nodes {
		k.SetMediaNode(ctx, mn)
	}
	for _, lease := range genState.Leases {
		k.SetLease(ctx, lease)
	}

	k.SetMediaNodeCount(ctx, genState.NodeCounter)

	// check if the module account exists
	moduleAcc := k.GetModuleAccountAddress(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}
}

// ExportGenesis exports state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	medianodes := k.GetAllMediaNodes(ctx)
	leases := k.GetAllLeases(ctx)
	return types.NewGenesisState(
		medianodes,
		leases,
		k.GetMediaNodeCount(ctx),
		k.GetParams(ctx),
	)
}

// DefaultGenesisState returns default state
func DefaultGenesisState() *types.GenesisState {
	return types.DefaultGenesis()
}
