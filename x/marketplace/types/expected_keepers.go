package types

import (
	"context"
	nft "github.com/OmniFlix/omniflixhub/v5/x/onft/exported"
	nftypes "github.com/OmniFlix/omniflixhub/v5/x/onft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AccountKeeper interface {
	// Methods imported from account should be defined here
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	GetModuleAccount(ctx context.Context, name string) sdk.ModuleAccountI
	GetModuleAddress(module string) sdk.AccAddress
}

type BankKeeper interface {
	// Methods imported from bank should be defined here
	BlockedAddr(recipient sdk.AccAddress) bool
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SendCoins(ctx context.Context, from sdk.AccAddress, to sdk.AccAddress, amount sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, formModule string, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx context.Context, formModule string, toModule string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, fromAddr sdk.AccAddress, toModule string, amt sdk.Coins) error
}

// NftKeeper methods imported from nft should be defined here
type NftKeeper interface {
	GetONFT(ctx sdk.Context, denomId, onftId string) (nft nft.ONFTI, err error)
	GetDenomInfo(ctx sdk.Context, denomId string) (*nftypes.Denom, error)
	TransferOwnership(ctx sdk.Context, denomId, nftId string, srcOwner, dstOwner sdk.AccAddress) error
}

// DistributionKeeper defines the expected distribution keeper
type DistributionKeeper interface {
	FundCommunityPool(ctx context.Context, amount sdk.Coins, sender sdk.AccAddress) error
}
