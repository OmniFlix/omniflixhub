package types

import (
	"github.com/OmniFlix/omniflix/v2/x/marketplace/exported"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

var (
	_ proto.Message     = &Listing{}
	_ exported.ListingI = &Listing{}
)

func NewListing(id, nftId, denomId string, price sdk.Coin, owner sdk.AccAddress, splitShares []WeightedAddress) Listing {
	return Listing{
		Id:          id,
		NftId:       nftId,
		DenomId:     denomId,
		Price:       price,
		Owner:       owner.String(),
		SplitShares: splitShares,
	}
}

func (l Listing) GetId() string {
	return l.Id
}

func (l Listing) GetDenomId() string {
	return l.DenomId
}

func (l Listing) GetNftId() string {
	return l.NftId
}

func (l Listing) GetPrice() sdk.Coin {
	return l.Price
}

func (l Listing) GetOwner() sdk.AccAddress {
	owner, _ := sdk.AccAddressFromBech32(l.Owner)
	return owner
}

func (l Listing) GetSplitShares() interface{} {
	return l.SplitShares
}
