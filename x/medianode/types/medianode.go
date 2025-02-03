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
		Url:           url,
		HardwareSpecs: hardwareSpecs,
		PricePerDay:   leaseAmountPerDay,
		Leased:        false,
		Owner:         owner,
		Deposits:      []*Deposit{},
	}
}

// IsLeased returns true if the media node is currently leased
func (m MediaNode) IsLeased() bool {
	return m.Leased
}
