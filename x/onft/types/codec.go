package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/cosmos/gogoproto/types"

	"github.com/OmniFlix/omniflixhub/v2/x/onft/exported"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateDenom{}, "OmniFlix/onft/MsgCreateDenom", nil)
	cdc.RegisterConcrete(&MsgUpdateDenom{}, "OmniFlix/onft/MsgUpdateDenom", nil)
	cdc.RegisterConcrete(&MsgTransferDenom{}, "OmniFlix/onft/MsgTransferDenom", nil)
	cdc.RegisterConcrete(&MsgTransferONFT{}, "OmniFlix/onft/MsgTransferONFT", nil)
	cdc.RegisterConcrete(&MsgMintONFT{}, "OmniFlix/onft/MsgMintONFT", nil)
	cdc.RegisterConcrete(&MsgBurnONFT{}, "OmniFlix/onft/MsgBurnONFT", nil)

	cdc.RegisterInterface((*exported.ONFT)(nil), nil)
	cdc.RegisterConcrete(&ONFT{}, "OmniFlix/onft/ONFT", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateDenom{},
		&MsgUpdateDenom{},
		&MsgTransferDenom{},
		&MsgTransferONFT{},
		&MsgMintONFT{},
		&MsgBurnONFT{},
	)

	registry.RegisterImplementations((*exported.ONFT)(nil),
		&ONFT{},
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

func MustMarshalSupply(cdc codec.BinaryCodec, supply uint64) []byte {
	supplyWrap := gogotypes.UInt64Value{Value: supply}
	return cdc.MustMarshal(&supplyWrap)
}

func MustUnMarshalSupply(cdc codec.BinaryCodec, value []byte) uint64 {
	var supplyWrap gogotypes.UInt64Value
	cdc.MustUnmarshal(value, &supplyWrap)
	return supplyWrap.Value
}

func MustMarshalONFTID(cdc codec.BinaryCodec, onftID string) []byte {
	onftIDWrap := gogotypes.StringValue{Value: onftID}
	return cdc.MustMarshal(&onftIDWrap)
}

func MustUnMarshalONFTID(cdc codec.BinaryCodec, value []byte) string {
	var onftIDWrap gogotypes.StringValue
	cdc.MustUnmarshal(value, &onftIDWrap)
	return onftIDWrap.Value
}

func MustMarshalDenomID(cdc codec.BinaryCodec, denomID string) []byte {
	denomIDWrap := gogotypes.StringValue{Value: denomID}
	return cdc.MustMarshal(&denomIDWrap)
}

func MustUnMarshalDenomID(cdc codec.BinaryCodec, value []byte) string {
	var denomIDWrap gogotypes.StringValue
	cdc.MustUnmarshal(value, &denomIDWrap)
	return denomIDWrap.Value
}
