package ics721nft

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	onftkeeper "github.com/OmniFlix/onft/keeper"
	onfttypes "github.com/OmniFlix/onft/types"
	nfttransfer "github.com/bianjieai/nft-transfer/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/tendermint/tendermint/libs/log"
)

// NewICS721NftKeeper creates a new ics721 Keeper instance
func NewICS721NftKeeper(cdc codec.Codec,
	k onftkeeper.Keeper,
	ak AccountKeeper,
) ICS721NftKeeper {
	return ICS721NftKeeper{
		nk:  k.NFTkeeper(),
		cdc: cdc,
		ak:  ak,
		cb:  onfttypes.NewClassBuilder(cdc, ak.GetModuleAddress),
		tb:  onfttypes.NewTokenBuilder(cdc),
	}
}

// CreateOrUpdateClass implement the method of ICS721Keeper.CreateOrUpdateClass
func (icsnk ICS721NftKeeper) CreateOrUpdateClass(ctx sdk.Context,
	classID,
	classURI,
	classData string,
) error {
	var (
		class nft.Class
		err   error
	)
	if len(classData) != 0 {
		class, err = icsnk.cb.Build(classID, classURI, classData)
		if err != nil {
			return err
		}
	} else {
		var denomMetadata = &onfttypes.DenomMetadata{
			Creator:    icsnk.ak.GetModuleAddress(types.ModuleName).String(),
			PreviewUri: "",
			Schema:     "",
		}

		metadata, err := codectypes.NewAnyWithValue(denomMetadata)
		if err != nil {
			return err
		}
		class = nft.Class{
			Id:   classID,
			Uri:  classURI,
			Data: metadata,
		}
	}

	class = nft.Class{
		Id:  classID,
		Uri: classURI,
	}

	if icsnk.nk.HasClass(ctx, classID) {
		return icsnk.nk.UpdateClass(ctx, class)
	}
	return icsnk.nk.SaveClass(ctx, class)
}

// Mint implement the method of ICS721Keeper.Mint
func (icsnk ICS721NftKeeper) Mint(ctx sdk.Context,
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
	return icsnk.nk.Mint(ctx, token, receiver)
}

// Transfer implement the method of ICS721Keeper.Transfer
func (icsnk ICS721NftKeeper) Transfer(
	ctx sdk.Context,
	classID,
	tokenID,
	tokenData string,
	receiver sdk.AccAddress,
) error {
	if err := icsnk.nk.Transfer(ctx, classID, tokenID, receiver); err != nil {
		return err
	}
	if len(tokenData) == 0 {
		return nil
	}
	nft, _ := icsnk.nk.GetNFT(ctx, classID, tokenID)
	token, err := icsnk.tb.Build(classID, tokenID, nft.Uri, tokenData)
	if err != nil {
		return err
	}

	return icsnk.nk.Update(ctx, token)
}

// GetClass implement the method of ICS721Keeper.GetClass
func (icsnk ICS721NftKeeper) GetClass(ctx sdk.Context, classID string) (nfttransfer.Class, bool) {
	class, exist := icsnk.nk.GetClass(ctx, classID)
	if !exist {
		return nil, false
	}

	return ICS721Class{
		ID:  classID,
		URI: class.Uri,
	}, true
}

// GetNFT implement the method of ICS721Keeper.GetNFT
func (icsnk ICS721NftKeeper) GetNFT(ctx sdk.Context, classID, tokenID string) (nfttransfer.NFT, bool) {
	nft, has := icsnk.nk.GetNFT(ctx, classID, tokenID)
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
func (icsnk ICS721NftKeeper) Burn(ctx sdk.Context, classID string, tokenID string) error {
	return icsnk.nk.Burn(ctx, classID, tokenID)
}

// GetOwner implement the method of ICS721Keeper.GetOwner
func (icsnk ICS721NftKeeper) GetOwner(ctx sdk.Context, classID string, tokenID string) sdk.AccAddress {
	return icsnk.nk.GetOwner(ctx, classID, tokenID)
}

// HasClass implement the method of ICS721Keeper.HasClass
func (icsnk ICS721NftKeeper) HasClass(ctx sdk.Context, classID string) bool {
	return icsnk.nk.HasClass(ctx, classID)
}

// Logger returns a module-specific logger.
func (icsnk ICS721NftKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "ics721/NFTKeeper")
}
