package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgRegisterMediaNode{}, "OmniFlix/medianode/MsgRegisterMediaNode")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateMediaNode{}, "OmniFlix/medianode/MsgUpdateMediaNode")
	legacy.RegisterAminoMsg(cdc, &MsgLeaseMediaNode{}, "OmniFlix/medianode/MsgLeaseMediaNode")
	legacy.RegisterAminoMsg(cdc, &MsgCancelLease{}, "OmniFlix/medianode/MsgCancelLease")

	cdc.RegisterConcrete(&MediaNode{}, "OmniFlix/medianode/MediaNode", nil)
	cdc.RegisterConcrete(&Lease{}, "OmniFlix/medianode/Lease", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRegisterMediaNode{},
		&MsgUpdateMediaNode{},
		&MsgLeaseMediaNode{},
		&MsgCancelLease{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
