package types

import (
	"time"

	nft "github.com/OmniFlix/omniflixhub/v4/x/onft/exported"
	nfttypes "github.com/OmniFlix/omniflixhub/v4/x/onft/types"
	streampaytypes "github.com/OmniFlix/streampay/v2/x/streampay/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
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
	GetONFT(ctx sdk.Context, denomId, onftId string) (nft nft.ONFTI, err error)
	GetDenomInfo(ctx sdk.Context, denomId string) (*nfttypes.Denom, error)
	MintONFT(
		ctx sdk.Context,
		denomID,
		nftID,
		name,
		description,
		mediaURI,
		uriHash,
		previewURI,
		nftData string,
		createdAt time.Time,
		transferable,
		extensible,
		nsfw bool,
		royaltyShare sdk.Dec,
		receiver sdk.AccAddress,
	) error
	TransferOwnership(ctx sdk.Context, denomId, nftId string, srcOwner, dstOwner sdk.AccAddress) error
	BurnONFT(ctx sdk.Context, denomId, nftId string, owner sdk.AccAddress) error
}

type StreamPayKeeper interface {
	CreateStreamPayment(
		ctx sdk.Context,
		sender, recipient sdk.AccAddress,
		amount sdk.Coin,
		paymentType streampaytypes.StreamType,
		duration time.Duration,
		periods []*streampaytypes.Period,
		cancellable bool,
	) (string, error)
}

// DistributionKeeper defines the expected distribution keeper
type DistributionKeeper interface {
	FundCommunityPool(ctx sdk.Context, amount sdk.Coins, sender sdk.AccAddress) error
}
