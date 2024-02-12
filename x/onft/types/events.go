package types

const (
	EventTypeCreateONFTDenom   = "create_onft_denom"
	EventTypeUpdateONFTDenom   = "update_onft_denom"
	EventTypeTransferONFTDenom = "transfer_onft_denom"
	EventTypePurgeONFTDenom    = "purge_onft_denom"

	EventTypeMintONFT     = "mint_onft"
	EventTypeTransferONFT = "transfer_onft"
	EventTypeBurnONFT     = "burn_onft"

	AttributeValueCategory       = ModuleName
	AttributeKeySender           = "sender"
	AttributeKeyCreator          = "creator"
	AttributeKeyOwner            = "owner"
	AttributeKeyRecipient        = "recipient"
	AttributeKeyNFTID            = "nft-id"
	AttributeKeyDenomID          = "denom-id"
	AttributeKeySymbol           = "symbol"
	AttributeKeyName             = "name"
	AttributeKeyDescription      = "description"
	AttributeKeyMediaURI         = "media-uri"
	AttributeKeyPreviewURI       = "preview-uri"
	AttributeKeyRoyaltyReceivers = "royalty-receivers"
)
