package types

const (
	EventTypeRegisterMediaNode       = "register_medianode"
	EventTypeUpdateMediaNode         = "update_medianode"
	EventTypeLeaseMediaNode          = "lease_medianode"
	EventTypeExtendLease             = "extend_medianode_lease"
	EventTypeCalcelLease             = "cancel_medianode_lease"
	EventTypeDepositMediaNode        = "deposit_medianode"
	EventTypeCloseMediaNode          = "close_medianode"
	EventTypeExpireLease             = "expire_medianode_lease"
	EventTypeRefundDeposit           = "refund_medianode_deposit"
	EventTypeLeaseCommissionTransfer = "lease_commission_transfer"
	EventTypeLeasePaymentTransfer    = "lease_payment_transfer"
	EventTypeSettleLeasePayment      = "settle_lease_payment"

	AttributeValueCategory    = ModuleName
	AttributeKeyMediaNodeId   = "medianode-id"
	AttributeKeyMediaNodeURL  = "medianode-url"
	AttributeKeyPricePerDay   = "price-per-day"
	AttributeKeyStatus        = "status"
	AttributeKeyOwner         = "owner"
	AttributeKeyDepositor     = "depositor"
	AttributeKeyRecipient     = "recipient"
	AttributeKeyAmount        = "amount"
	AttributeKeyStartTime     = "start-time"
	AttributeKeyLessee        = "lessee"
	AttributeKeySettledAmount = "settled-amount"
	AttributeKeyRefundAmount  = "refund-amount"
)
