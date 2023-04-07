package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagName            = "name"
	FlagDescription     = "description"
	FlagMediaURI        = "media-uri"
	FlagPreviewURI      = "preview-uri"
	FlagData            = "data"
	FlagNonTransferable = "non-transferable"
	FlagInExtensible    = "inextensible"
	FlagRecipient       = "recipient"
	FlagOwner           = "owner"
	FlagDenomID         = "denom-id"
	FlagSchema          = "schema"
	FlagNsfw            = "nsfw"
	FlagRoyaltyShare    = "royalty-share"
	FlagCreationFee     = "creation-fee"
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

	FsUpdateDenom.String(FlagName, "[do-not-modify]", "Name of the denom")
	FsUpdateDenom.String(FlagDescription, "[do-not-modify]", "Description for denom")
	FsUpdateDenom.String(FlagPreviewURI, "[do-not-modify]", "Preview image uri for denom")

	FsTransferDenom.String(FlagRecipient, "", "recipient of the denom")

	FsMintONFT.String(FlagMediaURI, "", "Media uri of onft")
	FsMintONFT.String(FlagRecipient, "", "Receiver of the onft. default value is sender address of transaction")
	FsMintONFT.String(FlagPreviewURI, "", "Preview uri of onft")
	FsMintONFT.String(FlagName, "", "Name of onft")
	FsMintONFT.String(FlagDescription, "", "Description of onft")
	FsMintONFT.String(FlagData, "", "custom data of onft")
	FsMintONFT.Bool(FlagNonTransferable, false, "To mint non-transferable onft")
	FsMintONFT.Bool(FlagInExtensible, false, "To mint non-extensisble onft")
	FsMintONFT.Bool(FlagNsfw, false, "not safe for work flag for onft")
	FsMintONFT.String(FlagRoyaltyShare, "", "Royalty share value decimal value between 0 and 1")

	FsTransferONFT.String(FlagRecipient, "", "Receiver of the onft. default value is sender address of transaction")
	FsQuerySupply.String(FlagOwner, "", "The owner of a nft")
	FsQueryOwner.String(FlagDenomID, "", "id of the denom")
}
