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
