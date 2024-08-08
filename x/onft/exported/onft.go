package exported

import (
	"time"

	sdkmath "cosmossdk.io/math"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ONFTI interface {
	GetID() string
	GetOwner() sdk.AccAddress
	GetName() string
	GetDescription() string
	GetMediaURI() string
	GetURIHash() string
	GetPreviewURI() string
	GetData() string
	IsTransferable() bool
	IsExtensible() bool
	IsNSFW() bool
	GetCreatedTime() time.Time
	GetRoyaltyShare() sdkmath.LegacyDec
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
