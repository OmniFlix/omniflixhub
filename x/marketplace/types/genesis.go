package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

func NewGenesisState(listings []Listing, listingCount uint64, params Params,
	auctions []AuctionListing, bids []Bid, nextAuctionNumber uint64,
) *GenesisState {
	return &GenesisState{
		Listings:          listings,
		ListingCount:      listingCount,
		Params:            params,
		Auctions:          auctions,
		Bids:              bids,
		NextAuctionNumber: nextAuctionNumber,
	}
}

func (m *GenesisState) ValidateGenesis() error {
	for _, l := range m.Listings {
		if l.GetOwner().Empty() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing nft owner")
		}
		if err := ValidateListing(l); err != nil {
			return err
		}
	}
	if m.ListingCount < 0 {
		return sdkerrors.Wrap(ErrNonPositiveNumber, "must be a positive number")
	}
	if err := m.Params.ValidateBasic(); err != nil {
		return err
	}
	for _, auction := range m.Auctions {
		if err := ValidateAuctionListing(auction); err != nil {
			return err
		}
	}
	for _, bid := range m.Bids {
		if err := ValidateBid(bid); err != nil {
			return err
		}
	}
	if m.NextAuctionNumber <= 0 {
		return sdkerrors.Wrap(ErrNonPositiveNumber, "must be a number and greater than 0.")
	}
	return nil
}
