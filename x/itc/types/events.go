package types

const (
	EventTypeCreateCampaign  = "create_campaign"
	EventTypeCancelCampaign  = "cancel_campaign"
	EventTypeClaim           = "claim"
	EventTypeEndCampaign     = "end_campaign"
	EventTypeDepositCampaign = "deposit_campaign"

	AttributeValueCategory      = ModuleName
	AttributeKeyCampaignId      = "campaign-id"
	AttributeKeyNftDenomId      = "nft-denom-id"
	AttributeKeyNftId           = "nft-id"
	AttributeKeyClaimType       = "claim-type"
	AttributeKeyInteractionType = "interaction-type"
	AttributeKeyStartTime       = "start-time"
	AttributeKeyEndTime         = "end-time"
	AttributeKeyClaimer         = "claimer"
	AttributeKeyCreator         = "creator"
	AttributeKeySender          = "sender"
	AttributeKeyDepositor       = "depositor"
	AttributeKeyAmount          = "amount"
)
