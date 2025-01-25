package types

import (
	"github.com/OmniFlix/omniflixhub/v6/x/onft/exported"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	"github.com/cosmos/gogoproto/proto"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgCreateDenom{}, "OmniFlix/onft/MsgCreateDenom")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateDenom{}, "OmniFlix/onft/MsgUpdateDenom")
	legacy.RegisterAminoMsg(cdc, &MsgTransferDenom{}, "OmniFlix/onft/MsgTransferDenom")
	legacy.RegisterAminoMsg(cdc, &MsgTransferONFT{}, "OmniFlix/onft/MsgTransferONFT")
	legacy.RegisterAminoMsg(cdc, &MsgMintONFT{}, "OmniFlix/onft/MsgMintONFT")
	legacy.RegisterAminoMsg(cdc, &MsgBurnONFT{}, "OmniFlix/onft/MsgBurnONFT")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "OmniFlix/onft/MsgUpdateParams")

	cdc.RegisterConcrete(&Params{}, "OmniFlix/onft/Params", nil)

	cdc.RegisterInterface((*exported.ONFTI)(nil), nil)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateDenom{},
		&MsgUpdateDenom{},
		&MsgTransferDenom{},
		&MsgTransferONFT{},
		&MsgMintONFT{},
		&MsgBurnONFT{},
		&MsgUpdateParams{},
	)

	registry.RegisterInterface(
		"OmniFlix.onft.v1beta1.ONFTI",
		(*exported.ONFTI)(nil),
		&ONFT{},
	)
	registry.RegisterImplementations(
		(*proto.Message)(nil),
		&DenomMetadata{},
		&ONFTMetadata{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
