package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/OmniFlix/omniflixhub/x/itc/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type (
	Keeper struct {
		cdc      codec.Codec
		storeKey sdk.StoreKey

		accountKeeper   types.AccountKeeper
		bankKeeper      types.BankKeeper
		nftKeeper       types.NftKeeper
		streampayKeeper types.StreamPayKeeper
		paramSpace      paramtypes.Subspace
	}
)

func NewKeeper(
	cdc codec.Codec,
	storeKey sdk.StoreKey,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	nftKeeper types.NftKeeper,
	streampayKeeper types.StreamPayKeeper,
	ps paramtypes.Subspace,
) Keeper {
	// ensure itc module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:             cdc,
		storeKey:        storeKey,
		accountKeeper:   accountKeeper,
		bankKeeper:      bankKeeper,
		nftKeeper:       nftKeeper,
		streampayKeeper: streampayKeeper,
		paramSpace:      ps,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetModuleAccountAddress(ctx sdk.Context) sdk.AccAddress {
	return k.accountKeeper.GetModuleAddress(types.ModuleName)
}
