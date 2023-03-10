package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	MsgRoute = "itc"

	TypeMsgCreateCampaign  = "create_campaign"
	TypeMsgCancelCampaign  = "cancel_campaign"
	TypeMsgCampaignDeposit = "deposit_campaign"
	TypeMsgClaim           = "claim"

	// DoNotModify used to indicate that some field should not be updated
	DoNotModify = "[do-not-modify]"
)

var _ sdk.Msg = &MsgCreateCampaign{}

// TODO: update CreateCampaign Msg

func NewMsgCreateCampaign(name, description string,
	interaction InteractionType, claimType ClaimType,
	nftDenomId string,
	maxAllowedClaims uint64,
	claimableTokens, totalTokens Tokens,
	nftMintDetails *NFTDetails,
	distribution Distribution,
	startTime time.Time,
	duration time.Duration,
	creator string,
) *MsgCreateCampaign {
	return &MsgCreateCampaign{
		Name:             name,
		Description:      description,
		Interaction:      interaction,
		ClaimType:        claimType,
		NftDenomId:       nftDenomId,
		MaxAllowedClaims: maxAllowedClaims,
		ClaimableTokens:  claimableTokens,
		TotalTokens:      totalTokens,
		NftMintDetails:   nftMintDetails,
		Distribution:     distribution,
		StartTime:        startTime,
		Duration:         duration,
		Creator:          creator,
	}
}

func (msg MsgCreateCampaign) Route() string { return MsgRoute }

func (msg MsgCreateCampaign) Type() string { return TypeMsgCreateCampaign }

// TODO: validations required

func (msg MsgCreateCampaign) ValidateBasic() error {
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgCreateCampaign) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgCreateCampaign) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

func NewMsgCancelCampaign(id uint64, creator string) *MsgCancelCampaign {
	return &MsgCancelCampaign{
		CampaignId: id,
		Creator:    creator,
	}
}

func (msg MsgCancelCampaign) Route() string { return MsgRoute }

func (msg MsgCancelCampaign) Type() string { return TypeMsgCancelCampaign }

// TODO: validations required

func (msg MsgCancelCampaign) ValidateBasic() error {
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgCancelCampaign) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgCancelCampaign) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

func NewMsgCampaignDeposit(id uint64, amount sdk.Coin, depositor string) *MsgCampaignDeposit {
	return &MsgCampaignDeposit{
		CampaignId: id,
		Amount:     amount,
		Depositor:  depositor,
	}
}

func (msg MsgCampaignDeposit) Route() string { return MsgRoute }

func (msg MsgCampaignDeposit) Type() string { return TypeMsgCampaignDeposit }

// TODO: validations required

func (msg MsgCampaignDeposit) ValidateBasic() error {
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgCampaignDeposit) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgCampaignDeposit) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

func NewMsgClaim(id uint64, nftId string,
	interaction InteractionType, claimer string,
) *MsgClaim {
	return &MsgClaim{
		CampaignId:  id,
		NftId:       nftId,
		Interaction: interaction,
		Claimer:     claimer,
	}
}

func (msg MsgClaim) Route() string { return MsgRoute }

func (msg MsgClaim) Type() string { return TypeMsgClaim }

// TODO: validations required

func (msg MsgClaim) ValidateBasic() error {
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgClaim) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgClaim) GetSigners() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}
