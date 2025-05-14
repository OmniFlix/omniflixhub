package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

var _ proto.Message = &Lease{}

func NewLease(
	mediaNodeId string,
	lessee string,
	startTime time.Time,
	leasedHours uint64,
	totalAmount sdk.Coin,
	active bool,
	status LeaseStatus,
) Lease {
	return Lease{
		MediaNodeId: mediaNodeId,
		Lessee:      lessee,
		StartTime:   startTime,
		LeasedHours: leasedHours,
	}
}
