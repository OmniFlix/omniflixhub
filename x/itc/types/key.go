package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName          = "itc"
	StoreKey     string = ModuleName
	QuerierRoute string = ModuleName
	RouterKey    string = ModuleName
)

var (
	PrefixCampaignId         = []byte{0x01}
	PrefixCampaignCreator    = []byte{0x02}
	PrefixNextCampaignNumber = []byte{0x03}
	PrefixClaimByNftId       = []byte{0x11}
)

func KeyCampaignIdPrefix(id uint64) []byte {
	return append(PrefixCampaignId, sdk.Uint64ToBigEndian(id)...)
}

func KeyCampaignCreatorPrefix(creator sdk.AccAddress, id uint64) []byte {
	return append(append(PrefixCampaignCreator, creator.Bytes()...), sdk.Uint64ToBigEndian(id)...)
}

func KeyClaimByNftIdPrefix(id uint64, nftId string) []byte {
	return append(append(PrefixClaimByNftId, sdk.Uint64ToBigEndian(id)...), []byte(nftId)...)
}
