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

// GetMediaNodeKey returns the store key to retrieve a MediaNode from the index fields
func GetMediaNodeKey(id uint64) []byte {
	return append(PrefixMediaNode, sdk.Uint64ToBigEndian(id)...)
}

// GetLeaseKey returns the store key to retrieve a Lease from the index fields
func GetLeaseKey(id uint64) []byte {
	return append(PrefixLease, sdk.Uint64ToBigEndian(id)...)
}
