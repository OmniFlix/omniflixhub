package ics721nft

import (
	nfttransfer "github.com/bianjieai/nft-transfer/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	nftkeeper "github.com/cosmos/cosmos-sdk/x/nft/keeper"
	"github.com/tendermint/tendermint/libs/log"
)

// NewICS721NftKeeper creates a new ics721 Keeper instance
func NewICS721NftKeeper(cdc codec.Codec,
	nk nftkeeper.Keeper,
	ak AccountKeeper,
) ICS721NftKeeper {
	return ICS721NftKeeper{
		nk:  nk,
		cdc: cdc,
		ak:  ak,
	}
}

// CreateOrUpdateClass implement the method of ICS721Keeper.CreateOrUpdateClass
func (ik ICS721NftKeeper) CreateOrUpdateClass(ctx sdk.Context,
	classID,
	classURI,
	classData string,
) error {
	var (
		class nft.Class
	)

	class = nft.Class{
		Id:  classID,
		Uri: classURI,
	}

	if ik.nk.HasClass(ctx, classID) {
		return ik.nk.UpdateClass(ctx, class)
	}
	return ik.nk.SaveClass(ctx, class)
}

// Mint implement the method of ICS721Keeper.Mint
func (ik ICS721NftKeeper) Mint(ctx sdk.Context,
	classID,
	tokenID,
	tokenURI,
	tokenData string,
	receiver sdk.AccAddress,
) error {

	token := nft.NFT{
		ClassId: classID,
		Id:      tokenID,
		Uri:     tokenURI,
	}
	return ik.nk.Mint(ctx, token, receiver)
}

// Transfer implement the method of ICS721Keeper.Transfer
func (ik ICS721NftKeeper) Transfer(
	ctx sdk.Context,
	classID,
	tokenID,
	tokenData string,
	receiver sdk.AccAddress,
) error {
	if err := ik.nk.Transfer(ctx, classID, tokenID, receiver); err != nil {
		return err
	}
	return nil
}

// GetClass implement the method of ICS721Keeper.GetClass
func (ik ICS721NftKeeper) GetClass(ctx sdk.Context, classID string) (nfttransfer.Class, bool) {
	class, exist := ik.nk.GetClass(ctx, classID)
	if !exist {
		return nil, false
	}

	return ICS721Class{
		ID:  classID,
		URI: class.Uri,
	}, true
}

// GetNFT implement the method of ICS721Keeper.GetNFT
func (ik ICS721NftKeeper) GetNFT(ctx sdk.Context, classID, tokenID string) (nfttransfer.NFT, bool) {
	nft, has := ik.nk.GetNFT(ctx, classID, tokenID)
	if !has {
		return nil, false
	}
	return ICS721Token{
		ClassID: classID,
		ID:      tokenID,
		URI:     nft.Uri,
	}, true
}

// Burn implement the method of ICS721Keeper.Burn
func (ik ICS721NftKeeper) Burn(ctx sdk.Context, classID string, tokenID string) error {
	return ik.nk.Burn(ctx, classID, tokenID)
}

// GetOwner implement the method of ICS721Keeper.GetOwner
func (ik ICS721NftKeeper) GetOwner(ctx sdk.Context, classID string, tokenID string) sdk.AccAddress {
	return ik.nk.GetOwner(ctx, classID, tokenID)
}

// HasClass implement the method of ICS721Keeper.HasClass
func (ik ICS721NftKeeper) HasClass(ctx sdk.Context, classID string) bool {
	return ik.nk.HasClass(ctx, classID)
}

// Logger returns a module-specific logger.
func (ik ICS721NftKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "ics721/NFTKeeper")
}
