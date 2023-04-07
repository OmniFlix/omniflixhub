package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/cosmos/gogoproto/types"

	"github.com/OmniFlix/omniflixhub/v2/x/marketplace/exported"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgListNFT{}, "OmniFlix/marketplace/MsgListNFT", nil)
	cdc.RegisterConcrete(&MsgEditListing{}, "OmniFlix/marketplace/MsgEditListing", nil)
	cdc.RegisterConcrete(&MsgDeListNFT{}, "OmniFlix/marketplace/MsgDeListNFT", nil)
	cdc.RegisterConcrete(&MsgBuyNFT{}, "OmniFlix/marketplace/MsgBuyNFT", nil)
	cdc.RegisterConcrete(&MsgCreateAuction{}, "OmniFlix/marketplace/MsgCreateAuction", nil)
	cdc.RegisterConcrete(&MsgCancelAuction{}, "OmniFlix/marketplace/MsgCancelAuction", nil)
	cdc.RegisterConcrete(&MsgPlaceBid{}, "OmniFlix/marketplace/MsgPlaceBid", nil)

	cdc.RegisterInterface((*exported.ListingI)(nil), nil)
	cdc.RegisterConcrete(&Listing{}, "OmniFlix/marketplace/Listing", nil)
	cdc.RegisterInterface((*exported.AuctionListingI)(nil), nil)
	cdc.RegisterConcrete(&AuctionListing{}, "OmniFlix/marketplace/AuctionListing", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgListNFT{},
		&MsgEditListing{},
		&MsgDeListNFT{},
		&MsgBuyNFT{},
		&MsgCreateAuction{},
		&MsgCancelAuction{},
		&MsgPlaceBid{},
	)

	registry.RegisterImplementations((*exported.ListingI)(nil),
		&Listing{},
	)
	registry.RegisterImplementations((*exported.AuctionListingI)(nil),
		&AuctionListing{},
	)
}

var (
	amino = codec.NewLegacyAmino()

	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}

func MustMarshalListingID(cdc codec.BinaryCodec, listingId string) []byte {
	listingIdWrap := gogotypes.StringValue{Value: listingId}
	return cdc.MustMarshal(&listingIdWrap)
}

func MustUnMarshalListingID(cdc codec.BinaryCodec, value []byte) string {
	var listingIdWrap gogotypes.StringValue
	cdc.MustUnmarshal(value, &listingIdWrap)
	return listingIdWrap.Value
}
