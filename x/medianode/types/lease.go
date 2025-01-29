package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

var (
	_ proto.Message = &Lease{}
)

func NewLease(
	id uint64,
	mediaNodeId uint64,
	lessee string,
	startTime time.Time,
	leasedDays uint64,
	totalAmount sdk.Coin,
	active bool,
	status LeaseStatus,
) Lease {
	return Lease{
		MediaNodeId: mediaNodeId,
		LeasedTo:    lessee,
		LeasedAt:    startTime,
		LeasedDays:  leasedDays,
		LeaseStatus: status,
	}
}

func (l Lease) GetMediaNodeId() uint64 {
	return l.MediaNodeId
}

func (l Lease) GetLessee() sdk.AccAddress {
	lessee, _ := sdk.AccAddressFromBech32(l.GetLeasedTo())
	return lessee
}

func (l Lease) GetLeasedAt() time.Time {
	return l.LeasedAt
}

func (l Lease) GetLeaseExpiry() time.Time {
	return l.LeaseExpiry
}
