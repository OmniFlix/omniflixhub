package itc

import (
	"fmt"

	"github.com/OmniFlix/omniflixhub/x/itc/keeper"
	"github.com/OmniFlix/omniflixhub/x/itc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	if err := genState.ValidateGenesis(); err != nil {
		panic(err.Error())
	}
	for _, c := range genState.Campaigns {
		k.SetCampaign(ctx, c)
		k.SetCampaignWithCreator(ctx, c.GetCreator(), c.GetId())
	}
	k.SetNextCampaignNumber(ctx, genState.NextCampaignNumber)
	k.SetParams(ctx, genState.Params)

	for _, cc := range genState.Claims {
		k.SetClaim(ctx, cc)
		k.SetClaimWithNft(ctx, cc.CampaignId, cc.NftId)
	}

	// check if the module account exists
	moduleAcc := k.GetModuleAccountAddress(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}
}

// ExportGenesis exports state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	campaigns := k.GetAllCampaigns(ctx)
	var claims []types.Claim
	for _, campaign := range campaigns {
		claims = append(claims, k.GetClaims(ctx, campaign.GetId())...)
	}
	return types.NewGenesisState(
		campaigns,
		claims,
		k.GetNextCampaignNumber(ctx),
		k.GetParams(ctx),
	)
}

// DefaultGenesisState returns default state
func DefaultGenesisState() *types.GenesisState {
	return types.NewGenesisState(
		[]types.Campaign{},
		[]types.Claim{},
		1,
		types.DefaultParams(),
	)
}
