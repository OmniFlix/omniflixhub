package types

import (
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

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

func NewMsgCreateCampaign(name, description string,
	interaction InteractionType, claimType ClaimType,
	nftDenomId string,
	maxAllowedClaims uint64,
	tokensPerClaim, tokenDeposit Tokens,
	nftMintDetails *NFTDetails,
	distribution *Distribution,
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
		TokensPerClaim:   tokensPerClaim,
		Deposit:          tokenDeposit,
		NftMintDetails:   nftMintDetails,
		Distribution:     distribution,
		StartTime:        startTime,
		Duration:         duration,
		Creator:          creator,
	}
}

func (msg MsgCreateCampaign) Route() string { return MsgRoute }

func (msg MsgCreateCampaign) Type() string { return TypeMsgCreateCampaign }

func (msg MsgCreateCampaign) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if err := ValidateClaimType(msg.ClaimType); err != nil {
		return err
	}
	if err := ValidateInteractionType(msg.Interaction); err != nil {
		return err
	}
	if err := ValidateTokensWithClaimType(msg.ClaimType, msg.Deposit); err != nil {
		return err
	}
	if err := ValidateTokensWithClaimType(msg.ClaimType, msg.TokensPerClaim); err != nil {
		return err
	}
	if msg.ClaimType == CLAIM_TYPE_NFT {
		if err := validateNFTMintDetails(msg.NftMintDetails); err != nil {
			return err
		}
	}
	if msg.ClaimType == CLAIM_TYPE_FT && msg.Distribution.Type == DISTRIBUTION_TYPE_VEST {
		if err := ValidateDistribution(msg.Distribution); err != nil {
			return err
		}
	}
	if msg.MaxAllowedClaims == 0 {
		return sdkerrors.Wrapf(ErrInValidMaxAllowedClaims,
			"max allowed claims must be a positive number (%d)", msg.MaxAllowedClaims)
	}
	if msg.NftMintDetails != nil {
		if err := validateNFTMintDetails(msg.NftMintDetails); err != nil {
			return err
		}
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg *MsgCreateCampaign) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
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

func (msg MsgCancelCampaign) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg *MsgCancelCampaign) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
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

func (msg MsgCampaignDeposit) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid depositor address (%s)", err)
	}
	if !msg.Amount.IsValid() && msg.Amount.IsNegative() && msg.Amount.IsZero() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins,
			"amount must be valid and positive (%s)", msg.Amount.String())
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg *MsgCampaignDeposit) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
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

func (msg MsgClaim) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Claimer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid claimer address (%s)", err)
	}
	return ValidateInteractionType(msg.Interaction)
}

// GetSignBytes Implements Msg.
func (msg *MsgClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgClaim) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Claimer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}
