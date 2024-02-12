package keeper

import (
	"github.com/OmniFlix/omniflixhub/v3/x/onft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// GetONFTModuleAccount returns oNFT ModuleAccount
func (k Keeper) GetONFTModuleAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
}
