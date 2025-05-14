package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

var _ proto.Message = &MediaNode{}

func NewMediaNode(
	id string,
	url string,
	owner string,
	hardwareSpecs HardwareSpecs,
	info Info,
	leaseAmountPerHour sdk.Coin,
) MediaNode {
	return MediaNode{
		Id:            id,
		Url:           url,
		HardwareSpecs: hardwareSpecs,
		PricePerHour:  leaseAmountPerHour,
		Leased:        false,
		Owner:         owner,
		Deposits:      []*Deposit{},
		Info:          info,
	}
}

// IsLeased returns true if the media node is currently leased
func (m MediaNode) IsLeased() bool {
	return m.Leased
}
