package keeper

import (
	"github.com/OmniFlix/omniflixhub/v5/x/alloc/exported"
	v3 "github.com/OmniFlix/omniflixhub/v5/x/alloc/migrations/v3"
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

// Migrate2to3 migrates the x/alloc module state from the consensus version 2 to
// version 3. Specifically, it takes the parameters that are currently stored
// and managed by the x/params modules and stores them directly into the x/alloc
// module state.
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	return v3.Migrate(ctx, ctx.KVStore(m.keeper.storeKey), m.legacySubspace, m.keeper.cdc)
}
