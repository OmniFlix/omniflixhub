package keeper

import (
	"fmt"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/OmniFlix/omniflixhub/v2/x/itc/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type (
	Keeper struct {
		cdc      codec.Codec
		storeKey storetypes.StoreKey

		accountKeeper      types.AccountKeeper
		bankKeeper         types.BankKeeper
		nftKeeper          types.NftKeeper
		streampayKeeper    types.StreamPayKeeper
		distributionKeeper types.DistributionKeeper

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string
	}
)

func NewKeeper(
	cdc codec.Codec,
	storeKey storetypes.StoreKey,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	nftKeeper types.NftKeeper,
	streampayKeeper types.StreamPayKeeper,
	distributionKeeper types.DistributionKeeper,
	authority string,,
) Keeper {
	// ensure itc module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		cdc:                cdc,
		storeKey:           storeKey,
		accountKeeper:      accountKeeper,
		bankKeeper:         bankKeeper,
		nftKeeper:          nftKeeper,
		streampayKeeper:    streampayKeeper,
		distributionKeeper: distributionKeeper,
		authority:         authority,
	}
}

// GetAuthority returns the x/itc module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetModuleAccountAddress(ctx sdk.Context) sdk.AccAddress {
	return k.accountKeeper.GetModuleAddress(types.ModuleName)
}
