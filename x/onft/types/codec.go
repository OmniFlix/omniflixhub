package types

import (
	"github.com/OmniFlix/omniflixhub/v2/x/onft/exported"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	authzcodec "github.com/cosmos/cosmos-sdk/x/authz/codec"
	"github.com/cosmos/gogoproto/proto"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateDenom{}, "OmniFlix/onft/MsgCreateDenom", nil)
	cdc.RegisterConcrete(&MsgUpdateDenom{}, "OmniFlix/onft/MsgUpdateDenom", nil)
	cdc.RegisterConcrete(&MsgTransferDenom{}, "OmniFlix/onft/MsgTransferDenom", nil)
	cdc.RegisterConcrete(&MsgTransferONFT{}, "OmniFlix/onft/MsgTransferONFT", nil)
	cdc.RegisterConcrete(&MsgMintONFT{}, "OmniFlix/onft/MsgMintONFT", nil)
	cdc.RegisterConcrete(&MsgBurnONFT{}, "OmniFlix/onft/MsgBurnONFT", nil)
	cdc.RegisterConcrete(&MsgUpdateParams{}, "OmniFlix/onft/MsgUpdateParams", nil)

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

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
)

func init() {
	RegisterLegacyAminoCodec(amino)

	sdk.RegisterLegacyAminoCodec(amino)
	// Register all Amino interfaces and concrete types on the authz Amino codec so that this can later be
	// used to properly serialize MsgGrant and MsgExec instances
	RegisterLegacyAminoCodec(authzcodec.Amino)
	amino.Seal()
}
