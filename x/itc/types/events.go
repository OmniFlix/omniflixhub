package types

const (
	EventTypeCreateCampaign  = "create_campaign"
	EventTypeCancelCampaign  = "cancel_campaign"
	EventTypeClaim           = "claim"
	EventTypeEndCampaign     = "end_campaign"
	EventTypeCampaignDeposit = "campaign_deposit"

	AttributeValueCategory      = ModuleName
	AttributeKeyCampaignId      = "campaign-id"
	AttributeKeyDenomId         = "denom-id"
	AttributeKeyNftDenomId      = "nft-denom-id"
	AttributeKeyNftId           = "nft-id"
	AttributeKeyClaimer         = "claimer"
	AttributeKeyCreator         = "creator"
	AttributeKeyDepositor       = "depositor"
	AttributeKeyAvailableTokens = "available-tokens"
	AttributeKeyTotalTokens     = "total-tokens"
	AttributeKeyAmount          = "amount"
)
