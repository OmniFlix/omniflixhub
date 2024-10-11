package types

import (
	"strings"
	"unicode/utf8"

	sdkmath "cosmossdk.io/math"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateDenom   = "create_denom"
	TypeMsgUpdateDenom   = "update_denom"
	TypeMsgTransferDenom = "transfer_denom"

	TypeMsgMintONFT     = "mint_onft"
	TypeMsgTransferONFT = "transfer_onft"
	TypeMsgBurnONFT     = "burn_onft"
)

var (
	_ sdk.Msg = &MsgUpdateParams{}

	_ sdk.Msg = &MsgCreateDenom{}
	_ sdk.Msg = &MsgUpdateDenom{}
	_ sdk.Msg = &MsgTransferDenom{}

	_ sdk.Msg = &MsgMintONFT{}
	_ sdk.Msg = &MsgTransferONFT{}
	_ sdk.Msg = &MsgBurnONFT{}
)

func NewMsgCreateDenom(
	symbol,
	name, schema, description, uri, uriHash, previewUri, data, sender string,
	creationFee sdk.Coin,
	royaltyReceivers []*WeightedAddress,
	updatableData bool,
) *MsgCreateDenom {
	return &MsgCreateDenom{
		Sender:           sender,
		Id:               GenUniqueID(DenomPrefix),
		Symbol:           symbol,
		Name:             name,
		Schema:           schema,
		Description:      description,
		PreviewURI:       previewUri,
		Uri:              uri,
		UriHash:          uriHash,
		Data:             data,
		CreationFee:      creationFee,
		RoyaltyReceivers: royaltyReceivers,
		UpdatableData:    updatableData,
	}
}

func (msg MsgCreateDenom) Route() string { return RouterKey }

func (msg MsgCreateDenom) Type() string { return TypeMsgCreateDenom }

func (msg MsgCreateDenom) ValidateBasic() error {
	if err := ValidateDenomID(msg.Id); err != nil {
		return err
	}
	if err := ValidateDenomSymbol(msg.Symbol); err != nil {
		return err
	}
	if len(msg.Name) > 0 && !utf8.ValidString(msg.Name) {
		return errorsmod.Wrap(ErrInvalidName, "denom name is invalid")
	}
	if err := ValidateName(msg.Name); err != nil {
		return err
	}
	if len(msg.Description) > 0 && !utf8.ValidString(msg.Description) {
		return errorsmod.Wrap(ErrInvalidDescription, "denom description is invalid")
	}
	if err := ValidateDescription(msg.Description); err != nil {
		return err
	}
	if err := ValidateURI(msg.PreviewURI); err != nil {
		return err
	}
	if err := ValidateCreationFee(msg.CreationFee); err != nil {
		return err
	}
	if msg.RoyaltyReceivers != nil {
		if err := ValidateWeightedAddresses(msg.RoyaltyReceivers); err != nil {
			return errorsmod.Wrap(ErrInvalidRoyaltyReceivers, "royalty receivers value is invalid")
		}
	}

	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}

func (msg MsgCreateDenom) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

func NewMsgUpdateDenom(
	id, name, description, previewUri, sender string,
	royaltyReceivers []*WeightedAddress,
) *MsgUpdateDenom {
	return &MsgUpdateDenom{
		Id:               id,
		Name:             name,
		Description:      description,
		PreviewURI:       previewUri,
		Sender:           sender,
		RoyaltyReceivers: royaltyReceivers,
	}
}

func (msg MsgUpdateDenom) Route() string { return RouterKey }

func (msg MsgUpdateDenom) Type() string { return TypeMsgUpdateDenom }

func (msg MsgUpdateDenom) ValidateBasic() error {
	if err := ValidateDenomID(msg.Id); err != nil {
		return err
	}
	if len(msg.Name) > 0 && !utf8.ValidString(msg.Name) {
		return errorsmod.Wrap(ErrInvalidName, "denom name is invalid")
	}
	if err := ValidateName(msg.Name); err != nil {
		return err
	}
	if len(msg.Description) > 0 && !utf8.ValidString(msg.Description) {
		return errorsmod.Wrap(ErrInvalidDescription, "denom description is invalid")
	}
	if err := ValidateDescription(msg.Description); err != nil {
		return err
	}
	if err := ValidateURI(msg.PreviewURI); err != nil {
		return err
	}
	if msg.RoyaltyReceivers != nil {
		if err := ValidateWeightedAddresses(msg.RoyaltyReceivers); err != nil {
			return errorsmod.Wrap(ErrInvalidRoyaltyReceivers, "royalty receivers value is invalid")
		}
	}

	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}

func (msg MsgUpdateDenom) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

func NewMsgTransferDenom(id, sender, recipient string) *MsgTransferDenom {
	return &MsgTransferDenom{
		Id:        id,
		Sender:    sender,
		Recipient: recipient,
	}
}

func (msg MsgTransferDenom) Route() string { return RouterKey }

func (msg MsgTransferDenom) Type() string { return TypeMsgTransferDenom }

func (msg MsgTransferDenom) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address; %s", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Recipient); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid recipient address; %s", err)
	}
	return nil
}

func (msg MsgTransferDenom) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

func NewMsgPurgeDenom(id, sender string) *MsgPurgeDenom {
	return &MsgPurgeDenom{
		Id:     id,
		Sender: sender,
	}
}

func (msg MsgPurgeDenom) Route() string { return RouterKey }

func (msg MsgPurgeDenom) Type() string { return TypeMsgTransferDenom }

func (msg MsgPurgeDenom) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address; %s", err)
	}
	return nil
}

func (msg MsgPurgeDenom) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

func NewMsgMintONFT(
	denomId, sender, recipient string, metadata Metadata, data string,
	transferable, extensible, nsfw bool, royaltyShare sdkmath.LegacyDec,
) *MsgMintONFT {
	return &MsgMintONFT{
		Id:           GenUniqueID(IDPrefix),
		DenomId:      denomId,
		Metadata:     metadata,
		Data:         data,
		Transferable: transferable,
		Extensible:   extensible,
		Nsfw:         nsfw,
		RoyaltyShare: royaltyShare,
		Sender:       sender,
		Recipient:    recipient,
	}
}

func (msg MsgMintONFT) Route() string { return RouterKey }

func (msg MsgMintONFT) Type() string { return TypeMsgMintONFT }

func (msg MsgMintONFT) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address; %s", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Recipient); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid recipient address; %s", err)
	}
	if err := ValidateName(msg.Metadata.Name); err != nil {
		return err
	}
	if err := ValidateDescription(msg.Metadata.Description); err != nil {
		return err
	}
	if err := ValidateMediaURI(msg.Metadata.MediaURI); err != nil {
		return err
	}
	if err := ValidateURI(msg.Metadata.PreviewURI); err != nil {
		return err
	}
	if msg.RoyaltyShare.IsNegative() || msg.RoyaltyShare.GTE(sdkmath.LegacyNewDec(1)) {
		return errorsmod.Wrapf(ErrInvalidPercentage, "invalid royalty share percentage decimal value; %d, must be positive and less than 1", msg.RoyaltyShare)
	}

	return ValidateONFTID(msg.Id)
}

func (msg MsgMintONFT) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

func NewMsgTransferONFT(id, denomId, sender, recipient string) *MsgTransferONFT {
	return &MsgTransferONFT{
		Id:        strings.ToLower(strings.TrimSpace(id)),
		DenomId:   strings.TrimSpace(denomId),
		Sender:    sender,
		Recipient: recipient,
	}
}

func (msg MsgTransferONFT) Route() string { return RouterKey }

func (msg MsgTransferONFT) Type() string { return TypeMsgTransferONFT }

func (msg MsgTransferONFT) ValidateBasic() error {
	if err := ValidateDenomID(msg.DenomId); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address; %s", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Recipient); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid recipient address; %s", err)
	}
	return ValidateONFTID(msg.Id)
}

func (msg MsgTransferONFT) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

func NewMsgBurnONFT(denomId, id, sender string) *MsgBurnONFT {
	return &MsgBurnONFT{
		DenomId: denomId,
		Id:      id,
		Sender:  sender,
	}
}

func (msg MsgBurnONFT) Route() string { return RouterKey }

func (msg MsgBurnONFT) Type() string { return TypeMsgBurnONFT }

func (msg MsgBurnONFT) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address; %s", err)
	}
	return ValidateONFTID(msg.Id)
}

func (msg MsgBurnONFT) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// MsgUpdateParams

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (m *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (m *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	return m.Params.ValidateBasic()
}
