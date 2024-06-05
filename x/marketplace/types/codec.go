package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	authzcodec "github.com/cosmos/cosmos-sdk/x/authz/codec"
	govcodec "github.com/cosmos/cosmos-sdk/x/gov/codec"
	gogotypes "github.com/cosmos/gogoproto/types"

	"github.com/OmniFlix/omniflixhub/v5/x/marketplace/exported"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgListNFT{}, "OmniFlix/marketplace/MsgListNFT")
	legacy.RegisterAminoMsg(cdc, &MsgEditListing{}, "OmniFlix/marketplace/MsgEditListing")
	legacy.RegisterAminoMsg(cdc, &MsgDeListNFT{}, "OmniFlix/marketplace/MsgDeListNFT")
	legacy.RegisterAminoMsg(cdc, &MsgBuyNFT{}, "OmniFlix/marketplace/MsgBuyNFT")
	legacy.RegisterAminoMsg(cdc, &MsgCreateAuction{}, "OmniFlix/marketplace/MsgCreateAuction")
	legacy.RegisterAminoMsg(cdc, &MsgCancelAuction{}, "OmniFlix/marketplace/MsgCancelAuction")
	legacy.RegisterAminoMsg(cdc, &MsgPlaceBid{}, "OmniFlix/marketplace/MsgPlaceBid")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "OmniFlix/marketplace/MsgUpdateParams")

	cdc.RegisterInterface((*exported.ListingI)(nil), nil)
	cdc.RegisterConcrete(&Listing{}, "OmniFlix/marketplace/Listing", nil)
	cdc.RegisterInterface((*exported.AuctionListingI)(nil), nil)
	cdc.RegisterConcrete(&AuctionListing{}, "OmniFlix/marketplace/AuctionListing", nil)
	cdc.RegisterConcrete(&Params{}, "OmniFlix/marketplace/Params", nil)
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
		&MsgUpdateParams{},
	)

	registry.RegisterImplementations((*exported.ListingI)(nil),
		&Listing{},
	)
	registry.RegisterImplementations((*exported.AuctionListingI)(nil),
		&AuctionListing{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	sdk.RegisterLegacyAminoCodec(amino)
	// Register all Amino interfaces and concrete types on the authz Amino codec
	// so that this can later be used to properly serialize MsgGrant and MsgExec
	// instances.
	RegisterLegacyAminoCodec(authzcodec.Amino)
	RegisterLegacyAminoCodec(govcodec.Amino)
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
