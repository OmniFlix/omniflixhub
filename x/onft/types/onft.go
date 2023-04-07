package types

import (
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/OmniFlix/omniflixhub/v2/x/onft/exported"
)

var _ exported.ONFT = ONFT{}

func NewONFT(
	id string, metadata Metadata, data string, transferable, extensible bool, owner sdk.AccAddress,
	createdTime time.Time, nsfw bool, royaltyShare sdk.Dec,
) ONFT {
	return ONFT{
		Id:           id,
		Metadata:     metadata,
		Data:         data,
		Owner:        owner.String(),
		Transferable: transferable,
		Extensible:   extensible,
		CreatedAt:    createdTime,
		Nsfw:         nsfw,
		RoyaltyShare: royaltyShare,
	}
}

func (onft ONFT) GetID() string {
	return onft.Id
}

func (onft ONFT) GetName() string {
	return onft.Metadata.Name
}

func (onft ONFT) GetDescription() string {
	return onft.Metadata.Description
}

func (onft ONFT) GetMediaURI() string {
	return onft.Metadata.MediaURI
}

func (onft ONFT) GetPreviewURI() string {
	return onft.Metadata.PreviewURI
}

func (onft ONFT) GetOwner() sdk.AccAddress {
	owner, _ := sdk.AccAddressFromBech32(onft.Owner)
	return owner
}

func (onft ONFT) GetMetadata() string {
	return onft.Metadata.String()
}

func (onft ONFT) GetData() string {
	return onft.Data
}

func (onft ONFT) IsTransferable() bool {
	return onft.Transferable
}

func (onft ONFT) IsExtensible() bool {
	return onft.Extensible
}

func (onft ONFT) GetCreatedTime() time.Time {
	return onft.CreatedAt
}

func (onft ONFT) IsNSFW() bool {
	return onft.Nsfw
}

func (onft ONFT) GetRoyaltyShare() sdk.Dec {
	return onft.RoyaltyShare
}

// ONFT

type ONFTs []exported.ONFT

func NewONFTs(onfts ...exported.ONFT) ONFTs {
	if len(onfts) == 0 {
		return ONFTs{}
	}
	return onfts
}

func ValidateONFTID(onftId string) error {
	onftId = strings.TrimSpace(onftId)
	if len(onftId) < MinIDLen || len(onftId) > MaxIDLen {
		return sdkerrors.Wrapf(
			ErrInvalidONFTID,
			"invalid onftId %s, length must be between [%d, %d]", onftId, MinIDLen, MaxIDLen)
	}
	if !IsBeginWithAlpha(onftId) || !IsAlphaNumeric(onftId) {
		return sdkerrors.Wrapf(
			ErrInvalidONFTID,
			"invalid onftId %s, only accepts alphanumeric characters and begin with an english letter", onftId)
	}
	return nil
}
