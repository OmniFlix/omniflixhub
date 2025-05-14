package keeper

import (
	"github.com/OmniFlix/omniflixhub/v6/x/marketplace/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetMarketplaceAccount returns marketplace ModuleAccount
func (k Keeper) GetMarketplaceAccount(ctx sdk.Context) sdk.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
}
