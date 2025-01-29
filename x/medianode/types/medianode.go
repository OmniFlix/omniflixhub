package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

var (
	_ proto.Message = &MediaNode{}
)

func NewMediaNode(
	id uint64,
	url string,
	owner string,
	hardwareSpecs HardwareSpecs,
	leaseAmountPerDay sdk.Coin,
	leased bool,
	createdAt, updatedAt time.Time,
) MediaNode {
	return MediaNode{
		Id:                id,
		Url:               url,
		Owner:             owner,
		HardwareSpecs:     hardwareSpecs,
		LeaseAmountPerDay: leaseAmountPerDay,
		Leased:            leased,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}
}
