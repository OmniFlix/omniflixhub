package types

import (
	"github.com/OmniFlix/omniflixhub/v5/x/itc/exported"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	authzcodec "github.com/cosmos/cosmos-sdk/x/authz/codec"
	govcodec "github.com/cosmos/cosmos-sdk/x/gov/codec"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgCreateCampaign{}, "OmniFlix/itc/MsgCreateCampaign")
	legacy.RegisterAminoMsg(cdc, &MsgCancelCampaign{}, "OmniFlix/itc/MsgCancelCampaign")
	legacy.RegisterAminoMsg(cdc, &MsgClaim{}, "OmniFlix/itc/MsgClaim")
	legacy.RegisterAminoMsg(cdc, &MsgDepositCampaign{}, "OmniFlix/itc/MsgDepositCampaign")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "OmniFlix/itc/MsgUpdateParams")

	cdc.RegisterInterface((*exported.CampaignI)(nil), nil)
	cdc.RegisterConcrete(&Campaign{}, "OmniFlix/itc/Campaign", nil)
	cdc.RegisterInterface((*exported.ClaimI)(nil), nil)
	cdc.RegisterConcrete(&Claim{}, "OmniFlix/itc/Claim", nil)
	cdc.RegisterConcrete(&Params{}, "OmniFlix/itc/Params", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateCampaign{},
		&MsgCancelCampaign{},
		&MsgClaim{},
		&MsgDepositCampaign{},
		&MsgUpdateParams{},
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
