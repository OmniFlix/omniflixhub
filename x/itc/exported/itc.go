package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

type CampaignI interface {
	GetId() uint64
	GetName() string
	GetDescription() string
	GetStartTime() time.Time
	GetEndTime() time.Time
	GetNftDenomId() string
	GetMaxAllowedClaims() uint64
	GetClaimType() string
	GetInteraction() string
	GetCreator() sdk.AccAddress
	GetClaimableTokens() interface{}
	GetTotalTokens() interface{}
	GetAvailableTokens() interface{}
	GetReceivedNftIds() []string
	GetNftMintDetails() interface{}
	GetDistribution() interface{}
	GetStatus() string
}

type ClaimI interface {
	GetCampaignId() uint64
	GetAddress() sdk.AccAddress
}
