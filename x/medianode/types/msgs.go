package types

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/OmniFlix/omniflixhub/v6/app/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	MsgRoute = "medianode"

	TypeMsgRegisterMediaNode = "register_media_node"
	TypeMsgUpdateMediaNode   = "update_media_node"
	TypeMsgLeaseMediaNode    = "lease_media_node"
	TypeMsgExtendLease       = "extend_lease"
	TypeMsgCancelLease       = "cancel_lease"
	TypeMsgDepositMediaNode  = "deposit_media_node"
	TypeMsgCloseMediaNode    = "close_media_node"
)

var (
	_ sdk.Msg = &MsgRegisterMediaNode{}
	_ sdk.Msg = &MsgUpdateMediaNode{}
	_ sdk.Msg = &MsgLeaseMediaNode{}
	_ sdk.Msg = &MsgCancelLease{}
	_ sdk.Msg = &MsgExtendLease{}
	_ sdk.Msg = &MsgDepositMediaNode{}
	_ sdk.Msg = &MsgCloseMediaNode{}
)

// Register Media Node
func NewMsgRegisterMediaNode(url string, info Info, hardwareSpecs HardwareSpecs, pricePerHour, deposit sdk.Coin, sender string) (*MsgRegisterMediaNode, error) {
	mediaNodeId, err := GenUniqueID(MediaNodeIdPrefix)
	if err != nil {
		return nil, err
	}
	return &MsgRegisterMediaNode{
		Id:            mediaNodeId,
		Url:           url,
		HardwareSpecs: hardwareSpecs,
		PricePerHour:  pricePerHour,
		Sender:        sender,
		Deposit:       &deposit,
		Info:          info,
	}, nil
}

func (msg MsgRegisterMediaNode) Route() string { return MsgRoute }

func (msg MsgRegisterMediaNode) Type() string { return TypeMsgRegisterMediaNode }

func (msg MsgRegisterMediaNode) ValidateBasic() error {
	if err := validateMediaNodeId(msg.Id); err != nil {
		return err
	}

	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if msg.Url == "" {
		return errorsmod.Wrap(ErrInvalidURL, "url cannot be empty or invalied")
	}
	if err := msg.Info.Validate(); err != nil {
		return errorsmod.Wrap(ErrInvalidInfo, "invalid medianode info")
	}
	if err := msg.PricePerHour.Validate(); err != nil {
		return errorsmod.Wrap(ErrInvalidLeaseAmount, "invalid lease amount per hour")
	}
	if msg.PricePerHour.Denom != params.BondDenom {
		return errorsmod.Wrap(ErrInvalidLeaseAmount, "invalid denom for price per hour only uflix is supported")
	}
	if err := msg.Deposit.Validate(); err != nil {
		return errorsmod.Wrap(ErrInvalidDeposit, "invalid deposit amount")
	}
	if msg.PricePerHour.Amount.IsZero() {
		return errorsmod.Wrap(ErrInvalidLeaseAmount, "price per hour must be positive")
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
func NewMsgUpdateMediaNode(id string, info *Info, hardwareSpecs *HardwareSpecs, leaseAmountPerHour *sdk.Coin, sender string) *MsgUpdateMediaNode {
	return &MsgUpdateMediaNode{
		Id:            id,
		HardwareSpecs: hardwareSpecs,
		PricePerHour:  leaseAmountPerHour,
		Info:          info,
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
	if msg.HardwareSpecs != nil {
		if err := msg.HardwareSpecs.Validate(); err != nil {
			return errorsmod.Wrap(ErrInvalidHardwareSpecs, "invalid medianode info")
		}
	}
	if msg.Info != nil {
		if err := msg.Info.Validate(); err != nil {
			return errorsmod.Wrap(ErrInvalidInfo, "invalid medianode info")
		}
	}
	if msg.PricePerHour != nil {
		if err := msg.PricePerHour.Validate(); err != nil {
			return errorsmod.Wrap(ErrInvalidLeaseAmount, "invalid lease amount per day")
		}
		if msg.PricePerHour.Amount.IsZero() {
			return errorsmod.Wrap(ErrInvalidLeaseAmount, "price per hour amount must be positive")
		}
		if msg.PricePerHour.Denom != params.BondDenom {
			return errorsmod.Wrap(ErrInvalidLeaseAmount, "invalid denom for price per hour, only uflix is supported")
		}
	}
	if msg.HardwareSpecs == nil && msg.Info == nil && msg.PricePerHour == nil {
		return errorsmod.Wrapf(ErrInvalidUpdateMsg, "no updates are present in message")
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
func NewMsgLeaseMediaNode(mediaNodeId string, leaseHours uint64, leaseAmount sdk.Coin, sender string) *MsgLeaseMediaNode {
	return &MsgLeaseMediaNode{
		MediaNodeId: mediaNodeId,
		LeaseHours:  leaseHours,
		Sender:      sender,
		Amount:      leaseAmount,
	}
}

func (msg MsgLeaseMediaNode) Route() string { return MsgRoute }

func (msg MsgLeaseMediaNode) Type() string { return TypeMsgLeaseMediaNode }

func (msg MsgLeaseMediaNode) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if msg.LeaseHours == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "lease hours must be greater than 0")
	}
	if msg.Amount.Amount.IsNegative() || msg.Amount.Amount.IsZero() {
		return errorsmod.Wrap(ErrInvalidLeaseAmount, "lease amount must be greater than 0")
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

// Extend Lease
func NewMsgExtendLease(mediaNodeId string, leaseHours uint64, amount sdk.Coin, sender string) *MsgExtendLease {
	return &MsgExtendLease{
		MediaNodeId: mediaNodeId,
		LeaseHours:  leaseHours,
		Amount:      amount,
		Sender:      sender,
	}
}

func (msg MsgExtendLease) Route() string { return MsgRoute }

func (msg MsgExtendLease) Type() string { return TypeMsgExtendLease }

func (msg MsgExtendLease) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if msg.LeaseHours == 0 {
		return errorsmod.Wrap(ErrInvalidLeaseHours, "lease hours must be greater than 0")
	}
	if err := msg.Amount.Validate(); err != nil {
		return errorsmod.Wrap(ErrInvalidLeaseAmount, "invalid lease amount")
	}
	if msg.Amount.Amount.IsZero() {
		return errorsmod.Wrap(ErrInvalidLeaseAmount, "amount must be possitive")
	}
	return nil
}

func (msg MsgExtendLease) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// Cancel Lease
func NewMsgCancelLease(mediaNodeId string, sender string) *MsgCancelLease {
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
func NewMsgDepositMediaNode(mediaNodeId string, amount sdk.Coin, sender string) *MsgDepositMediaNode {
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
		return errorsmod.Wrap(ErrInvalidDeposit, "invalid deposit amount")
	}
	if msg.Amount.Amount.IsZero() {
		return errorsmod.Wrap(ErrInvalidDeposit, "deposit amount must be greater than zero")
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

func NewMsgCloseMediaNode(mediaNodeId string, sender string) *MsgCloseMediaNode {
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
