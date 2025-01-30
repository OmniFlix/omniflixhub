package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

var (
	_ proto.Message = &MediaNode{}
)

func NewMediaNode(
	url string,
	owner string,
	hardwareSpecs HardwareSpecs,
	leaseAmountPerDay sdk.Coin,
) MediaNode {
	return MediaNode{
		Url:               url,
		Owner:             owner,
		HardwareSpecs:     hardwareSpecs,
		LeaseAmountPerDay: leaseAmountPerDay,
		Leased:            false,
	}
}

// IsLeased returns true if the media node is currently leased
func (m MediaNode) IsLeased() bool {
	return m.Leased
}
