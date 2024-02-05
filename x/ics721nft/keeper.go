package ics721nft

import (
	"errors"

	errorsmod "cosmossdk.io/errors"

	onftkeeper "github.com/OmniFlix/omniflixhub/v2/x/onft/keeper"
	onfttypes "github.com/OmniFlix/omniflixhub/v2/x/onft/types"
	nfttransfer "github.com/bianjieai/nft-transfer/types"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/nft"
	nftkeeper "github.com/cosmos/cosmos-sdk/x/nft/keeper"
	"github.com/cosmos/gogoproto/proto"
)

// Keeper defines the ICS721 Keeper
type Keeper struct {
	nk  nftkeeper.Keeper
	cdc codec.Codec
	ak  AccountKeeper
	bk  BankKeeper
	cb  onfttypes.ClassBuilder
	nb  onfttypes.NFTBuilder
}

// NewKeeper creates a new ics721 Keeper instance
func NewKeeper(cdc codec.Codec,
	k onftkeeper.Keeper,
	ak AccountKeeper,
	bk BankKeeper,
) Keeper {
	return Keeper{
		nk:  k.NFTkeeper(),
		cdc: cdc,
		ak:  ak,
		bk:  bk,
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
			Creator:          k.ak.GetModuleAddress(onfttypes.ModuleName).String(),
			PreviewUri:       "",
			Schema:           "",
			Description:      "",
			Data:             "",
			UriHash:          "",
			RoyaltyReceivers: nil,
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
	var message proto.Message
	if err := k.cdc.UnpackAny(class.Data, &message); err != nil {
		return err
	}
	denomMetadata, ok := message.(*onfttypes.DenomMetadata)
	if !ok {
		return errors.New("unsupported classMetadata")
	}
	if denomMetadata.RoyaltyReceivers != nil && !k.validRoyaltyReceiverAddresses(denomMetadata.RoyaltyReceivers) {
		denomMetadata.RoyaltyReceivers = nil
		dMeta, err := codectypes.NewAnyWithValue(denomMetadata)
		if err != nil {
			return err
		}
		class.Data = dMeta
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
	_ string,
	receiver sdk.AccAddress,
) error {
	_nft, _ := k.nk.GetNFT(ctx, classID, tokenID)
	nftMetadata, err := onfttypes.UnmarshalNFTMetadata(k.cdc, _nft.Data.GetValue())
	if err != nil {
		return err
	}
	if !nftMetadata.Transferable {
		k.Logger(ctx).Error("non-transferable nft")
		return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "non-transferable nft")
	}
	if err := k.nk.Transfer(ctx, classID, tokenID, receiver); err != nil {
		return err
	}
	return nil
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

func (k Keeper) validRoyaltyReceiverAddresses(addresses []*onfttypes.WeightedAddress) bool {
	weightSum := sdk.NewDec(0)
	for _, addr := range addresses {
		address, err := sdk.AccAddressFromBech32(addr.Address)
		if err != nil {
			return false
		}
		if k.bk.BlockedAddr(address) {
			return false
		}
		if !addr.Weight.IsPositive() {
			return false
		}
		if addr.Weight.GT(sdk.NewDec(1)) {
			return false
		}
		weightSum = weightSum.Add(addr.Weight)
	}
	return weightSum.Equal(sdk.NewDec(1))
}
