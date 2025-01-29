package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName          = "medianode"
	StoreKey     string = ModuleName
	QuerierRoute string = ModuleName
	RouterKey    string = ModuleName
)

var (
	PrefixMediaNode  = []byte{0x01}
	PrefixLease      = []byte{0x02}
	PrefixNextNodeId = []byte{0x03}

	ParamsKey = []byte{0x12}
)

func KeyMediaNodePrefix(id uint64) []byte {
	return append(PrefixMediaNode, sdk.Uint64ToBigEndian(id)...)
}

func KeyLeasePrefix(nodeId uint64) []byte {
	return append(PrefixLease, sdk.Uint64ToBigEndian(nodeId)...)
}
