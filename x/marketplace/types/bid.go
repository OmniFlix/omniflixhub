package types

import (
	"github.com/OmniFlix/marketplace/x/marketplace/exported"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"time"
)

var (
	_ proto.Message = &Bid{}
	_ exported.BidI = &Bid{}
)

func NewBid(auctionId uint64, amount sdk.Coin, bidTime time.Time, bidder sdk.AccAddress) Bid {
	return Bid{
		AuctionId: auctionId,
		Amount:    amount,
		Time:      bidTime,
		Bidder:    bidder.String(),
	}
}

func (b Bid) GetAuctionId() uint64 {
	return b.AuctionId
}

func (b Bid) GetAmount() sdk.Coin {
	return b.Amount
}

func (b Bid) GetBidder() sdk.AccAddress {
	bidder, _ := sdk.AccAddressFromBech32(b.Bidder)
	return bidder
}
