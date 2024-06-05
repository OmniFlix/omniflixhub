package alloc

import (
	"fmt"
	"time"

	"github.com/OmniFlix/omniflixhub/v5/x/alloc/keeper"
	"github.com/OmniFlix/omniflixhub/v5/x/alloc/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker to distribute block rewards on every begin block
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	if err := k.DistributeMintedCoins(ctx); err != nil {
		panic(fmt.Sprintf("Error distribute minted coins: %s", err.Error()))
	}
}
