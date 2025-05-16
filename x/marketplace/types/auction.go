package types

import (
	"time"

	sdkmath "cosmossdk.io/math"

	"github.com/OmniFlix/omniflixhub/v6/x/marketplace/exported"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

var (
	_ proto.Message            = &AuctionListing{}
	_ exported.AuctionListingI = &AuctionListing{}
)

func NewAuctionListing(id uint64, nftId, denomId string, startTime time.Time, endTime *time.Time, startPrice sdk.Coin, incrementPercentage sdkmath.LegacyDec,
	owner sdk.AccAddress, whitelistAccounts []string, splitShares []WeightedAddress,
) AuctionListing {
	return AuctionListing{
		Id:                  id,
		NftId:               nftId,
		DenomId:             denomId,
		StartTime:           startTime,
		EndTime:             endTime,
		StartPrice:          startPrice,
		IncrementPercentage: incrementPercentage,
		Owner:               owner.String(),
		WhitelistAccounts:   whitelistAccounts,
		SplitShares:         splitShares,
	}
}

func (al AuctionListing) GetId() uint64 {
	return al.Id
}

func (al AuctionListing) GetDenomId() string {
	return al.DenomId
}

func (al AuctionListing) GetNftId() string {
	return al.NftId
}

func (al AuctionListing) GetStartTime() time.Time {
	return al.StartTime
}

func (al AuctionListing) GetStartPrice() sdk.Coin {
	return al.StartPrice
}

func (al AuctionListing) GetIncrementPercentage() sdkmath.LegacyDec {
	return al.IncrementPercentage
}

func (al AuctionListing) GetOwner() sdk.AccAddress {
	owner, _ := sdk.AccAddressFromBech32(al.Owner)
	return owner
}

func (al AuctionListing) GetSplitShares() interface{} {
	return al.SplitShares
}

func ValidAuctionStatus(status AuctionStatus) bool {
	if status == AUCTION_STATUS_INACTIVE ||
		status == AUCTION_STATUS_ACTIVE {
		return true
	}
	return false
}
