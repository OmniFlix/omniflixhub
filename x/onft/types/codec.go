package types

import (
	"github.com/OmniFlix/omniflixhub/v4/x/onft/exported"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	authzcodec "github.com/cosmos/cosmos-sdk/x/authz/codec"
	govcodec "github.com/cosmos/cosmos-sdk/x/gov/codec"
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

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	sdk.RegisterLegacyAminoCodec(amino)
	// Register all Amino interfaces and concrete types on the authz Amino codec so that this can later be
	// used to properly serialize MsgGrant and MsgExec instances
	RegisterLegacyAminoCodec(authzcodec.Amino)
	RegisterLegacyAminoCodec(govcodec.Amino)
}
