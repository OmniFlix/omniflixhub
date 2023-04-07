package exported

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ONFT interface {
	GetID() string
	GetOwner() sdk.AccAddress
	GetName() string
	GetDescription() string
	GetMediaURI() string
	GetPreviewURI() string
	GetData() string
	IsTransferable() bool
	IsExtensible() bool
	IsNSFW() bool
	GetCreatedTime() time.Time
	GetRoyaltyShare() sdk.Dec
}
