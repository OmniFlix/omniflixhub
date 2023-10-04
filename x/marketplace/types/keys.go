package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName               = "marketplace"
	StoreKey          string = ModuleName
	QuerierRoute      string = ModuleName
	RouterKey         string = ModuleName
	DefaultParamspace        = ModuleName
)

var (
	PrefixListingId         = []byte{0x01}
	PrefixListingOwner      = []byte{0x02}
	PrefixListingsCount     = []byte{0x03}
	PrefixListingNFTID      = []byte{0x04}
	PrefixListingPriceDenom = []byte{0x05}
	PrefixAuctionId         = []byte{0x06}
	PrefixAuctionOwner      = []byte{0x07}
	PrefixAuctionNFTID      = []byte{0x08}
	PrefixAuctionPriceDenom = []byte{0x09}
	PrefixNextAuctionNumber = []byte{0x10}
	PrefixBidByAuctionId    = []byte{0x11}
	PrefixBidByBidder       = []byte{0x12}
	PrefixInactiveAuction   = []byte{0x13}
	PrefixActiveAuction     = []byte{0x14}

	ParamsKey = []byte{0x15}
)

func KeyListingIdPrefix(id string) []byte {
	return append(PrefixListingId, []byte(id)...)
}

func KeyListingOwnerPrefix(owner sdk.AccAddress, id string) []byte {
	return append(append(PrefixListingOwner, owner.Bytes()...), []byte(id)...)
}

func KeyListingNFTIDPrefix(nftId string) []byte {
	return append(PrefixListingNFTID, []byte(nftId)...)
}

func KeyListingPriceDenomPrefix(priceDenom, id string) []byte {
	return append(append(PrefixListingPriceDenom, []byte(priceDenom)...), []byte(id)...)
}

func KeyAuctionIdPrefix(id uint64) []byte {
	return append(PrefixAuctionId, sdk.Uint64ToBigEndian(id)...)
}

func KeyAuctionOwnerPrefix(owner sdk.AccAddress, id uint64) []byte {
	return append(append(PrefixAuctionOwner, owner.Bytes()...), sdk.Uint64ToBigEndian(id)...)
}

func KeyAuctionNFTIDPrefix(nftId string) []byte {
	return append(PrefixAuctionNFTID, []byte(nftId)...)
}

func KeyAuctionPriceDenomPrefix(priceDenom string, id uint64) []byte {
	return append(append(PrefixAuctionPriceDenom, []byte(priceDenom)...), sdk.Uint64ToBigEndian(id)...)
}

func KeyBidPrefix(id uint64) []byte {
	return append(PrefixBidByAuctionId, sdk.Uint64ToBigEndian(id)...)
}

func KeyInActiveAuctionPrefix(id uint64) []byte {
	return append(PrefixInactiveAuction, sdk.Uint64ToBigEndian(id)...)
}

func KeyActiveAuctionPrefix(id uint64) []byte {
	return append(PrefixActiveAuction, sdk.Uint64ToBigEndian(id)...)
}
