package types

import (
	nft "github.com/OmniFlix/onft/exported"
	nftypes "github.com/OmniFlix/onft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
)

// AccountKeeper defines the expected account keeper
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	GetModuleAccount(ctx sdk.Context, name string) authtypes.ModuleAccountI
	GetModuleAddress(module string) sdk.AccAddress
}

// BankKeeper defines the expected bank keeper
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SendCoins(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amount sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, formModule string, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, formModule string, toModule string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, fromAddr sdk.AccAddress, toModule string, amt sdk.Coins) error
}

// NftKeeper defines the expected nft keeper
type NftKeeper interface {
	GetONFT(ctx sdk.Context, denomId, onftId string) (nft nft.ONFT, err error)
	GetDenom(ctx sdk.Context, denomId string) (nftypes.Denom, error)
	TransferOwnership(ctx sdk.Context, denomId, nftId string, srcOwner, dstOwner sdk.AccAddress) error
}

type VestingKeeper interface {
	CreateVestingAccount(ctx sdk.Context,
		msg *vestingtypes.MsgCreateVestingAccount,
	) (*vestingtypes.MsgCreateVestingAccountResponse, error)
}
