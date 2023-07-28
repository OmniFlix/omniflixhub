package types

import (
	"github.com/OmniFlix/omniflixhub/v2/x/itc/exported"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateCampaign{}, "OmniFlix/itc/MsgCreateCampaign", nil)
	cdc.RegisterConcrete(&MsgCancelCampaign{}, "OmniFlix/itc/MsgCancelCampaign", nil)
	cdc.RegisterConcrete(&MsgClaim{}, "OmniFlix/itc/MsgClaim", nil)
	cdc.RegisterConcrete(&MsgDepositCampaign{}, "OmniFlix/itc/MsgDepositCampaign", nil)

	cdc.RegisterInterface((*exported.CampaignI)(nil), nil)
	cdc.RegisterConcrete(&Campaign{}, "OmniFlix/itc/Campaign", nil)
	cdc.RegisterInterface((*exported.ClaimI)(nil), nil)
	cdc.RegisterConcrete(&Claim{}, "OmniFlix/itc/Claim", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateCampaign{},
		&MsgCancelCampaign{},
		&MsgClaim{},
		&MsgDepositCampaign{},
	)

	registry.RegisterImplementations((*exported.CampaignI)(nil),
		&Campaign{},
	)
	registry.RegisterImplementations((*exported.ClaimI)(nil),
		&Claim{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
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
