package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagDenomId     = "denom-id"
	FlagNftId       = "nft-id"
	FlagPrice       = "price"
	FlagOwner       = "owner"
	FlagBidder      = "bidder"
	FlagPriceDenom  = "price-denom"
	FlagSplitShares = "split-shares"
	FlagWhiteListAccounts = "whitelist-accounts"
	FlagStartTime   = "start-time"
	FlagStartPrice  = "start-price"
	FlagIncrementPercentage = "increment-percentage"
	FlagDuration = "duration"
	FlagAmount = "amount"
)

var (
	FsListNft     = flag.NewFlagSet("", flag.ContinueOnError)
	FsEditListing = flag.NewFlagSet("", flag.ContinueOnError)
	FsBuyNFT      = flag.NewFlagSet("", flag.ContinueOnError)

	FsCreateAuction = flag.NewFlagSet("", flag.ContinueOnError)
	FsPlaceBid      = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsListNft.String(FlagDenomId, "", "nft denom id")
	FsListNft.String(FlagNftId, "", "nft id")
	FsListNft.String(FlagPrice, "", "listing price of nft")
	FsListNft.String(FlagSplitShares, "", "split shares for listing")

	FsEditListing.String(FlagPrice, "", "listing price of nft")

	FsBuyNFT.String(FlagPrice, "", "buying price of nft")

	FsCreateAuction.String(FlagDenomId, "", "nft denom id")
	FsCreateAuction.String(FlagNftId, "", "nft id")
	FsCreateAuction.String(FlagStartPrice, "", "auction bid start price of nft")
	FsCreateAuction.String(FlagStartTime, "", "auction start time")
	FsCreateAuction.String(FlagDuration, "", "auction duration")
	FsCreateAuction.String(FlagIncrementPercentage, "0.01", "bid increment percentage")
	FsCreateAuction.String(FlagWhiteListAccounts, "", "whitelist accounts for private auction")
	FsCreateAuction.String(FlagSplitShares, "", "split shares for listing")

	FsPlaceBid.String(FlagAmount, "", "auction bid amount")
}
