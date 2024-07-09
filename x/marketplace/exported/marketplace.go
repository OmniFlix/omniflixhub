package exported

import (
	sdkmath "cosmossdk.io/math"
	"time"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ListingI interface {
	GetId() string
	GetDenomId() string
	GetNftId() string
	GetPrice() sdk.Coin
	GetOwner() sdk.AccAddress
	GetSplitShares() interface{}
}

type AuctionListingI interface {
	GetId() uint64
	GetDenomId() string
	GetNftId() string
	GetStartPrice() sdk.Coin
	GetStartTime() time.Time
	GetIncrementPercentage() sdkmath.LegacyDec
	GetOwner() sdk.AccAddress
	GetSplitShares() interface{}
	GetStatus() string
}

type BidI interface {
	GetAuctionId() uint64
	GetAmount() sdk.Coin
	GetBidder() sdk.AccAddress
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
