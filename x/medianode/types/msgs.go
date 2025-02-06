package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	MsgRoute = "medianode"

	TypeMsgRegisterMediaNode = "register_media_node"
	TypeMsgUpdateMediaNode   = "update_media_node"
	TypeMsgLeaseMediaNode    = "lease_media_node"
	TypeMsgCancelLease       = "cancel_lease"
	TypeMsgDepositMediaNode  = "deposit_media_node"
	TypeMsgCloseMediaNode    = "close_media_node"
)

var (
	_ sdk.Msg = &MsgRegisterMediaNode{}
	_ sdk.Msg = &MsgUpdateMediaNode{}
	_ sdk.Msg = &MsgLeaseMediaNode{}
	_ sdk.Msg = &MsgCancelLease{}
)

// Register Media Node
func NewMsgRegisterMediaNode(url string, hardwareSpecs HardwareSpecs, pricePerDay sdk.Coin, sender string) *MsgRegisterMediaNode {
	return &MsgRegisterMediaNode{
		Url:           url,
		HardwareSpecs: hardwareSpecs,
		PricePerDay:   pricePerDay,
		Sender:        sender,
	}
}

func (msg MsgRegisterMediaNode) Route() string { return MsgRoute }

func (msg MsgRegisterMediaNode) Type() string { return TypeMsgRegisterMediaNode }

func (msg MsgRegisterMediaNode) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if msg.Url == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "url cannot be empty")
	}
	if err := msg.PricePerDay.Validate(); err != nil {
		return errorsmod.Wrap(err, "invalid lease amount per day")
	}
	return nil
}

func (msg MsgRegisterMediaNode) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// Update Media Node
func NewMsgUpdateMediaNode(id uint64, hardwareSpecs HardwareSpecs, leaseAmountPerDay sdk.Coin, sender string) *MsgUpdateMediaNode {
	return &MsgUpdateMediaNode{
		Id:            id,
		HardwareSpecs: hardwareSpecs,
		PricePerDay:   leaseAmountPerDay,
		Sender:        sender,
	}
}

func (msg MsgUpdateMediaNode) Route() string { return MsgRoute }

func (msg MsgUpdateMediaNode) Type() string { return TypeMsgUpdateMediaNode }

func (msg MsgUpdateMediaNode) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if err := msg.PricePerDay.Validate(); err != nil {
		return errorsmod.Wrap(err, "invalid lease amount per day")
	}
	return nil
}

func (msg MsgUpdateMediaNode) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// Lease Media Node
func NewMsgLeaseMediaNode(mediaNodeId uint64, leaseDays uint64, sender string) *MsgLeaseMediaNode {
	return &MsgLeaseMediaNode{
		MediaNodeId: mediaNodeId,
		LeaseDays:   leaseDays,
		Sender:      sender,
	}
}

func (msg MsgLeaseMediaNode) Route() string { return MsgRoute }

func (msg MsgLeaseMediaNode) Type() string { return TypeMsgLeaseMediaNode }

func (msg MsgLeaseMediaNode) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if msg.LeaseDays == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "lease days must be greater than 0")
	}
	return nil
}

func (msg MsgLeaseMediaNode) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// Cancel Lease
func NewMsgCancelLease(mediaNodeId uint64, sender string) *MsgCancelLease {
	return &MsgCancelLease{
		MediaNodeId: mediaNodeId,
		Sender:      sender,
	}
}

func (msg MsgCancelLease) Route() string { return MsgRoute }

func (msg MsgCancelLease) Type() string { return TypeMsgCancelLease }

func (msg MsgCancelLease) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}

func (msg MsgCancelLease) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// Deposit Media Node
func NewMsgDepositMediaNode(mediaNodeId uint64, amount sdk.Coin, sender string) *MsgDepositMediaNode {
	return &MsgDepositMediaNode{
		MediaNodeId: mediaNodeId,
		Amount:      amount,
		Sender:      sender,
	}
}

func (msg MsgDepositMediaNode) Route() string { return MsgRoute }

func (msg MsgDepositMediaNode) Type() string { return TypeMsgDepositMediaNode }

func (msg MsgDepositMediaNode) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if err := msg.Amount.Validate(); err != nil {
		return errorsmod.Wrap(err, "invalid deposit amount")
	}
	return nil
}

func (msg MsgDepositMediaNode) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

func NewMsgCloseMediaNode(mediaNodeId uint64, sender string) *MsgCloseMediaNode {
	return &MsgCloseMediaNode{
		MediaNodeId: mediaNodeId,
		Sender:      sender,
	}
}

func (msg MsgCloseMediaNode) Route() string { return MsgRoute }

func (msg MsgCloseMediaNode) Type() string { return TypeMsgCloseMediaNode }

func (msg MsgCloseMediaNode) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}

func (msg MsgCloseMediaNode) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}
