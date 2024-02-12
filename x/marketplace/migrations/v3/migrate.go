package v3

import (
	"github.com/OmniFlix/omniflixhub/v3/x/marketplace/exported"
	"github.com/OmniFlix/omniflixhub/v3/x/marketplace/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "marketplace"
)

var ParamsKey = []byte{0x15}

// Migrate migrates the x/marketplace module state from the consensus version 2 to
// version 3. Specifically, it takes the parameters that are currently stored
// and managed by the x/params modules and stores them directly into the x/marketplace
// module state.
func Migrate(
	ctx sdk.Context,
	store sdk.KVStore,
	legacySubspace exported.Subspace,
	cdc codec.BinaryCodec,
) error {
	var currParams types.Params
	legacySubspace.GetParamSet(ctx, &currParams)

	if err := currParams.ValidateBasic(); err != nil {
		return err
	}
	bz := cdc.MustMarshal(&currParams)
	store.Set(ParamsKey, bz)

	return nil
}
