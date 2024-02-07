package types

import (
	"github.com/cosmos/cosmos-sdk/codec"

	gogotypes "github.com/cosmos/gogoproto/types"
)

const (
	ModuleName = "onft"
	StoreKey   = ModuleName
	RouterKey  = ModuleName
)

var (
	PrefixONFT       = []byte{0x01}
	PrefixOwners     = []byte{0x02}
	PrefixCollection = []byte{0x03}
	PrefixDenom      = []byte{0x04}

	ParamsKey = []byte{0x07}
)

func MustUnMarshalSupply(cdc codec.BinaryCodec, value []byte) uint64 {
	var supplyWrap gogotypes.UInt64Value
	cdc.MustUnmarshal(value, &supplyWrap)
	return supplyWrap.Value
}

func MustUnMarshalONFTID(cdc codec.BinaryCodec, value []byte) string {
	var onftIDWrap gogotypes.StringValue
	cdc.MustUnmarshal(value, &onftIDWrap)
	return onftIDWrap.Value
}

func MustUnMarshalDenomID(cdc codec.BinaryCodec, value []byte) string {
	var denomIDWrap gogotypes.StringValue
	cdc.MustUnmarshal(value, &denomIDWrap)
	return denomIDWrap.Value
}
