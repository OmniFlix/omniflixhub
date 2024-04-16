package v2

import (
	"github.com/OmniFlix/omniflixhub/v4/x/onft/exported"
	"github.com/OmniFlix/omniflixhub/v4/x/onft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "onft"
)

var ParamsKey = []byte{0x07}

func Migrate(
	ctx sdk.Context,
	storeKey storetypes.StoreKey,
	legacySubspace exported.Subspace,
	cdc codec.BinaryCodec,
	nftKeeper NFTKeeper,
) error {
	k := keeper{
		storeKey:  storeKey,
		cdc:       cdc,
		nftKeeper: nftKeeper,
	}

	if err := MigrateParams(ctx, legacySubspace, k); err != nil {
		return err
	}
	return MigrateCollections(ctx, storeKey, ctx.Logger(), k)
}

// MigrateParams migrates the onft module state from the consensus version 1 to
// version 2. Specifically, it takes the parameters that are currently stored
// and managed by the x/params modules and stores them directly into the onft
// module state.
func MigrateParams(
	ctx sdk.Context,
	legacySubspace exported.Subspace,
	k keeper,
) error {
	var currParams types.Params
	legacySubspace.GetParamSet(ctx, &currParams)

	if err := currParams.ValidateBasic(); err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&currParams)
	store.Set(ParamsKey, bz)

	return nil
}
