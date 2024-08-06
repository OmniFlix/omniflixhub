package types

import (
	"github.com/OmniFlix/omniflixhub/v5/x/marketplace/exported"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
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
