package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagName             = "name"
	FlagDescription      = "description"
	FlagMediaURI         = "media-uri"
	FlagPreviewURI       = "preview-uri"
	FlagData             = "data"
	FlagNonTransferable  = "non-transferable"
	FlagInExtensible     = "inextensible"
	FlagRecipient        = "recipient"
	FlagOwner            = "owner"
	FlagDenomID          = "denom-id"
	FlagSchema           = "schema"
	FlagNsfw             = "nsfw"
	FlagRoyaltyShare     = "royalty-share"
	FlagCreationFee      = "creation-fee"
	FlagURI              = "uri"
	FlagURIHash          = "uri-hash"
	FlagRoyaltyReceivers = "royalty-receivers"
)

var (
	FsCreateDenom   = flag.NewFlagSet("", flag.ContinueOnError)
	FsUpdateDenom   = flag.NewFlagSet("", flag.ContinueOnError)
	FsTransferDenom = flag.NewFlagSet("", flag.ContinueOnError)
	FsMintONFT      = flag.NewFlagSet("", flag.ContinueOnError)
	FsTransferONFT  = flag.NewFlagSet("", flag.ContinueOnError)
	FsQuerySupply   = flag.NewFlagSet("", flag.ContinueOnError)
	FsQueryOwner    = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsCreateDenom.String(FlagSchema, "", "Denom schema")
	FsCreateDenom.String(FlagName, "", "Name of the denom")
	FsCreateDenom.String(FlagDescription, "", "Description for denom")
	FsCreateDenom.String(FlagPreviewURI, "", "Preview image uri for denom")
	FsCreateDenom.String(FlagCreationFee, "", "fee amount for creating denom")
	FsCreateDenom.String(FlagRoyaltyReceivers, "", "royalty receivers with comma separated ex: \"address:percentage,address:percentage\"")

	FsUpdateDenom.String(FlagName, "[do-not-modify]", "Name of the denom")
	FsUpdateDenom.String(FlagDescription, "[do-not-modify]", "Description for denom")
	FsUpdateDenom.String(FlagPreviewURI, "[do-not-modify]", "Preview image uri for denom")
	FsUpdateDenom.String(FlagRoyaltyReceivers, "", "royalty receivers with comma separated ex: \"address:percentage,address:percentage\"")
	FsCreateDenom.String(FlagURI, "", "uri for denom")
	FsCreateDenom.String(FlagURIHash, "", "uri hash for denom")
	FsCreateDenom.String(FlagData, "", "json data of the denom")

	FsTransferDenom.String(FlagRecipient, "", "recipient of the denom")

	FsMintONFT.String(FlagMediaURI, "", "Media uri of onft")
	FsMintONFT.String(FlagRecipient, "", "Receiver of the onft. default value is sender address of transaction")
	FsMintONFT.String(FlagPreviewURI, "", "Preview uri of onft")
	FsMintONFT.String(FlagName, "", "Name of onft")
	FsMintONFT.String(FlagDescription, "", "Description of onft")
	FsMintONFT.String(FlagData, "", "custom data of onft")
	FsMintONFT.Bool(FlagNonTransferable, false, "To mint non-transferable onft")
	FsMintONFT.Bool(FlagInExtensible, false, "To mint non-extensible onft")
	FsMintONFT.Bool(FlagNsfw, false, "not safe for work flag for onft")
	FsMintONFT.String(FlagRoyaltyShare, "", "Royalty share value decimal value between 0 and 1")
	FsMintONFT.String(FlagURIHash, "", "uri hash for the nft")

	FsTransferONFT.String(FlagRecipient, "", "Receiver of the onft. default value is sender address of transaction")
	FsQuerySupply.String(FlagOwner, "", "The owner of a nft")
	FsQueryOwner.String(FlagDenomID, "", "id of the denom")
}
