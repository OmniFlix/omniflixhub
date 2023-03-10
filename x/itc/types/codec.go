package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/gogo/protobuf/types"

	"github.com/OmniFlix/omniflixhub/x/itc/exported"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateCampaign{}, "OmniFlix/itc/MsgCreateCampaign", nil)
	cdc.RegisterConcrete(&MsgCancelCampaign{}, "OmniFlix/itc/MsgCancelCampaign", nil)
	cdc.RegisterConcrete(&MsgClaim{}, "OmniFlix/itc/MsgClaim", nil)
	cdc.RegisterConcrete(&MsgCampaignDeposit{}, "OmniFlix/itc/MsgCampaignDeposit", nil)

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
		&MsgCampaignDeposit{},
	)

	registry.RegisterImplementations((*exported.CampaignI)(nil),
		&Campaign{},
	)
	registry.RegisterImplementations((*exported.ClaimI)(nil),
		&Claim{},
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

func MustMarshalListingID(cdc codec.BinaryCodec, listingId string) []byte {
	listingIdWrap := gogotypes.StringValue{Value: listingId}
	return cdc.MustMarshal(&listingIdWrap)
}

func MustUnMarshalListingID(cdc codec.BinaryCodec, value []byte) string {
	var listingIdWrap gogotypes.StringValue
	cdc.MustUnmarshal(value, &listingIdWrap)
	return listingIdWrap.Value
}
