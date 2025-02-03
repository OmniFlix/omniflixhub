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
	for _, mn := range genState.MediaNodes {
		k.SetMediaNode(ctx, mn)
	}
	k.SetNextMediaNodeNumber(ctx, genState.LastNodeId)

	for _, lease := range genState.Leases {
		k.SetLease(ctx, lease)
	}

	// check if the module account exists
	moduleAcc := k.GetModuleAccountAddress(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}
}

// DefaultGenesisState returns default state
func DefaultGenesisState() *types.GenesisState {
	return types.NewGenesisState(
		[]types.MediaNode{},
		[]types.Lease{},
		0,
		types.DefaultParams(),
	)
}
