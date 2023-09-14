package exported

import (
	"time"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
	GetTokensPerClaim() sdk.Coin
	GetTotalTokens() sdk.Coin
	GetAvailableTokens() sdk.Coin
	GetReceivedNftIds() []string
	GetNftMintDetails() interface{}
	GetDistribution() interface{}
}

type ClaimI interface {
	GetCampaignId() uint64
	GetAddress() sdk.AccAddress
	GetNftId() string
	GetInteractionType() string
}

type (
	ParamSet = paramtypes.ParamSet

	// Subspace defines an interface that implements the legacy x/params Subspace
	// type.
	//
	// NOTE: This is used solely for migration of x/params managed parameters.
	Subspace interface {
		GetParamSet(ctx sdk.Context, ps ParamSet)
	}
)
