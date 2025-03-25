package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgRegisterMediaNode{}, "OmniFlix/register-medianode")
	legacy.RegisterAminoMsg(cdc, &MsgDepositMediaNode{}, "OmniFlix/deposit-medianode")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateMediaNode{}, "OmniFlix/update-medianode")
	legacy.RegisterAminoMsg(cdc, &MsgLeaseMediaNode{}, "OmniFlix/lease-medianode")
	legacy.RegisterAminoMsg(cdc, &MsgExtendLease{}, "OmniFlix/extend-medianode-lease")
	legacy.RegisterAminoMsg(cdc, &MsgCancelLease{}, "OmniFlix/cancel-medianode-lease")
	legacy.RegisterAminoMsg(cdc, &MsgCloseMediaNode{}, "OmniFlix/close-medianode")

	cdc.RegisterConcrete(&MediaNode{}, "OmniFlix/medianode", nil)
	cdc.RegisterConcrete(&Lease{}, "OmniFlix/medianode-lease", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRegisterMediaNode{},
		&MsgDepositMediaNode{},
		&MsgUpdateMediaNode{},
		&MsgLeaseMediaNode{},
		&MsgExtendLease{},
		&MsgCancelLease{},
		&MsgCloseMediaNode{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
