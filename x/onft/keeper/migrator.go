package keeper

import (
	"github.com/OmniFlix/omniflixhub/v6/x/onft/exported"
	v2 "github.com/OmniFlix/omniflixhub/v6/x/onft/migrations/v2"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Migrator is a struct for handling in-place state migrations.
type Migrator struct {
	keeper         Keeper
	legacySubspace exported.Subspace
}

func NewMigrator(k Keeper, ss exported.Subspace) Migrator {
	return Migrator{
		keeper:         k,
		legacySubspace: ss,
	}
}

// Migrate1to2 migrates the onft module state from the consensus version 1 to
// version 2. Specifically, it takes the parameters that are currently stored
// and managed by the x/params modules and stores them directly into the onft
// module state.
// and
//
//	migrates onft store to x/nft
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return v2.Migrate(ctx, m.keeper.storeKey, m.legacySubspace, m.keeper.cdc, m.keeper)
}
