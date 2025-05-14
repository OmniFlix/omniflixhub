package types

import (
	"context"
	"time"

	sdkmath "cosmossdk.io/math"

	nft "github.com/OmniFlix/omniflixhub/v6/x/onft/exported"
	nfttypes "github.com/OmniFlix/omniflixhub/v6/x/onft/types"
	streampaytypes "github.com/OmniFlix/streampay/v2/x/streampay/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the expected account keeper
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	GetModuleAccount(ctx context.Context, name string) sdk.ModuleAccountI
	GetModuleAddress(module string) sdk.AccAddress
}

// BankKeeper defines the expected bank keeper
type BankKeeper interface {
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SendCoins(ctx context.Context, from sdk.AccAddress, to sdk.AccAddress, amount sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, formModule string, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx context.Context, formModule string, toModule string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, fromAddr sdk.AccAddress, toModule string, amt sdk.Coins) error
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
		royaltyShare sdkmath.LegacyDec,
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
	FundCommunityPool(ctx context.Context, amount sdk.Coins, sender sdk.AccAddress) error
}
