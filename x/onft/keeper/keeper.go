package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	"github.com/cometbft/cometbft/libs/log"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	nftkeeper "github.com/cosmos/cosmos-sdk/x/nft/keeper"

	"github.com/OmniFlix/omniflixhub/v3/x/onft/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec

	accountKeeper      types.AccountKeeper
	bankKeeper         types.BankKeeper
	distributionKeeper types.DistributionKeeper
	nk                 nftkeeper.Keeper
	authority          string
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	distrKeeper types.DistributionKeeper,
	authority string,
) Keeper {
	// ensure oNFT module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		storeKey:           storeKey,
		cdc:                cdc,
		accountKeeper:      accountKeeper,
		bankKeeper:         bankKeeper,
		distributionKeeper: distrKeeper,
		nk:                 nftkeeper.NewKeeper(storeKey, cdc, accountKeeper, bankKeeper),
		authority:          authority,
	}
}

// GetAuthority returns the onft module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("OmniFlix/%s", types.ModuleName))
}

// NFTkeeper returns a cosmos-sdk nftkeeper.Keeper.
func (k Keeper) NFTkeeper() nftkeeper.Keeper {
	return k.nk
}

func (k Keeper) ValidateRoyaltyReceiverAddresses(splitShares []*types.WeightedAddress) error {
	for _, share := range splitShares {
		addr, err := sdk.AccAddressFromBech32(share.Address)
		if err != nil {
			return err
		}
		if k.bankKeeper.BlockedAddr(addr) {
			return errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "%s is a blocked address and not allowed receive funds", addr)
		}
	}
	return nil
}
