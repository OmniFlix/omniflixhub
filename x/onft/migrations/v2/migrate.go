package v2

import (
	"github.com/OmniFlix/omniflixhub/v2/x/onft/exported"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "onft"
)

var ParamsKey = []byte{0x07}

// Migrate migrates the onft module state from the consensus version 1 to
// version 2. Specifically, it takes the parameters that are currently stored
// and managed by the x/params modules and stores them directly into the onft
// module state.
func Migrate(
	ctx sdk.Context,
	storeKey storetypes.StoreKey,
	legacySubspace exported.Subspace,
	cdc codec.BinaryCodec,
	nftKeeper NFTKeeper,
) error {
	/*
		var currParams types.Params
		legacySubspace.GetParamSet(ctx, &currParams)
	*/

	k := keeper{
		storeKey:  storeKey,
		cdc:       cdc,
		nftKeeper: nftKeeper,
	}
	/*
		if err := currParams.ValidateBasic(); err != nil {
			return err
		}
		store := ctx.KVStore(k.storeKey)

		bz := cdc.MustMarshal(&currParams)
		store.Set(ParamsKey, bz)

	*/

	return MigrateCollections(ctx, storeKey, cdc, ctx.Logger(), k)
}
