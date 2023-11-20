package ics721nft

import (
	onftkeeper "github.com/OmniFlix/omniflixhub/v2/x/onft/keeper"
	onfttypes "github.com/OmniFlix/omniflixhub/v2/x/onft/types"
	nfttransfer "github.com/bianjieai/nft-transfer/types"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	nftkeeper "github.com/cosmos/cosmos-sdk/x/nft/keeper"
)

// Keeper defines the ICS721 Keeper
type Keeper struct {
	nk  nftkeeper.Keeper
	cdc codec.Codec
	ak  AccountKeeper
	cb  onfttypes.ClassBuilder
	nb  onfttypes.NFTBuilder
}

// NewKeeper creates a new ics721 Keeper instance
func NewKeeper(cdc codec.Codec,
	k onftkeeper.Keeper,
	ak AccountKeeper,
) Keeper {
	return Keeper{
		nk:  k.NFTkeeper(),
		cdc: cdc,
		ak:  ak,
		cb:  onfttypes.NewClassBuilder(cdc, ak.GetModuleAddress),
		nb:  onfttypes.NewNFTBuilder(cdc),
	}
}

// CreateOrUpdateClass implement the method of ICS721Keeper.CreateOrUpdateClass
func (k Keeper) CreateOrUpdateClass(ctx sdk.Context,
	classID,
	classURI,
	classData string,
) error {
	var (
		class nft.Class
		err   error
	)
	if len(classData) != 0 {
		class, err = k.cb.Build(classID, classURI, classData)
		if err != nil {
			k.Logger(ctx).Error("unable to build class from packet data", "error:", err.Error())
			return err
		}
	} else {
		denomMetadata := &onfttypes.DenomMetadata{
			Creator:    k.ak.GetModuleAddress(onfttypes.ModuleName).String(),
			PreviewUri: "",
			Schema:     "",
		}

		metadata, err := codectypes.NewAnyWithValue(denomMetadata)
		if err != nil {
			k.Logger(ctx).Error("unable to build class metadata from packet data", "error:", err.Error())
			return err
		}
		class = nft.Class{
			Id:   classID,
			Uri:  classURI,
			Data: metadata,
		}
	}
	if k.nk.HasClass(ctx, classID) {
		return k.nk.UpdateClass(ctx, class)
	}
	return k.nk.SaveClass(ctx, class)
}

// Mint implement the method of ICS721Keeper.Mint
func (k Keeper) Mint(ctx sdk.Context,
	classID,
	tokenID,
	tokenURI,
	tokenData string,
	receiver sdk.AccAddress,
) error {
	token, err := k.nb.Build(classID, tokenID, tokenURI, tokenData)
	if err != nil {
		k.Logger(ctx).Error("unable to build token from packet data", "error:", err.Error())
		return err
	}
	return k.nk.Mint(ctx, token, receiver)
}

// Transfer implement the method of ICS721Keeper.Transfer
func (k Keeper) Transfer(
	ctx sdk.Context,
	classID,
	tokenID,
	tokenData string,
	receiver sdk.AccAddress,
) error {
	if err := k.nk.Transfer(ctx, classID, tokenID, receiver); err != nil {
		return err
	}
	if len(tokenData) == 0 {
		return nil
	}
	_nft, _ := k.nk.GetNFT(ctx, classID, tokenID)
	token, err := k.nb.Build(classID, tokenID, _nft.GetUri(), tokenData)
	if err != nil {
		k.Logger(ctx).Error("unable to build token on transfer from packet data", "error:", err.Error())
		return err
	}

	return k.nk.Update(ctx, token)
}

// GetClass implement the method of ICS721Keeper.GetClass
func (k Keeper) GetClass(ctx sdk.Context, classID string) (nfttransfer.Class, bool) {
	class, exist := k.nk.GetClass(ctx, classID)
	if !exist {
		return nil, false
	}
	metadata, err := k.cb.BuildMetadata(class)
	if err != nil {
		k.Logger(ctx).Error("encode class data failed")
		return nil, false
	}

	return ICS721Class{
		ID:   classID,
		URI:  class.Uri,
		Data: metadata,
	}, true
}

// GetNFT implement the method of ICS721Keeper.GetNFT
func (k Keeper) GetNFT(ctx sdk.Context, classID, tokenID string) (nfttransfer.NFT, bool) {
	_nft, has := k.nk.GetNFT(ctx, classID, tokenID)
	if !has {
		return nil, false
	}
	nftMetadata, err := onfttypes.UnmarshalNFTMetadata(k.cdc, _nft.Data.GetValue())
	if err != nil {
		return nil, false
	}
	if !nftMetadata.Transferable {
		k.Logger(ctx).Error("non-transferable nft")
		return nil, false
	}
	metadata, err := k.nb.BuildMetadata(_nft)
	if err != nil {
		k.Logger(ctx).Error("encode nft data failed")
		return nil, false
	}
	return ICS721Token{
		ClassID: classID,
		ID:      tokenID,
		URI:     _nft.Uri,
		Data:    metadata,
	}, true
}

// Burn implement the method of ICS721Keeper.Burn
func (k Keeper) Burn(ctx sdk.Context, classID string, tokenID string) error {
	return k.nk.Burn(ctx, classID, tokenID)
}

// GetOwner implement the method of ICS721Keeper.GetOwner
func (k Keeper) GetOwner(ctx sdk.Context, classID string, tokenID string) sdk.AccAddress {
	return k.nk.GetOwner(ctx, classID, tokenID)
}

// HasClass implement the method of ICS721Keeper.HasClass
func (k Keeper) HasClass(ctx sdk.Context, classID string) bool {
	return k.nk.HasClass(ctx, classID)
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "ics721/NFTKeeper")
}
