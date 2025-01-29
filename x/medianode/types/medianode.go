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

func (m MediaNode) GetId() uint64 {
	return m.Id
}

func (m MediaNode) GetUrl() string {
	return m.Url
}

func (m MediaNode) GetOwner() sdk.AccAddress {
	owner, _ := sdk.AccAddressFromBech32(m.Owner)
	return owner
}

func (m MediaNode) GetHardwareSpecs() HardwareSpecs {
	return m.HardwareSpecs
}

func (m MediaNode) GetLeaseAmountPerDay() sdk.Coin {
	return m.LeaseAmountPerDay
}

func (m MediaNode) IsLeased() bool {
	return m.Leased
}

func (m MediaNode) GetCreatedAt() time.Time {
	return m.CreatedAt
}

func (m MediaNode) GetUpdatedAt() time.Time {
	return m.UpdatedAt
}
