package keeper

import (
	"github.com/OmniFlix/omniflixhub/v5/x/onft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetONFTModuleAccount returns oNFT ModuleAccount
func (k Keeper) GetONFTModuleAccount(ctx sdk.Context) sdk.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
}
