package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/nft module sentinel errors
var (
	ErrInvalidCollection       = sdkerrors.Register(ModuleName, 3, "invalid ONFT collection")
	ErrUnknownCollection       = sdkerrors.Register(ModuleName, 4, "unknown ONFT collection")
	ErrInvalidONFT             = sdkerrors.Register(ModuleName, 5, "invalid ONFT")
	ErrONFTAlreadyExists       = sdkerrors.Register(ModuleName, 6, "ONFT already exists")
	ErrUnknownONFT             = sdkerrors.Register(ModuleName, 7, "unknown ONFT")
	ErrEmptyMetaData           = sdkerrors.Register(ModuleName, 8, "ONFT MetaData can't be empty")
	ErrUnauthorized            = sdkerrors.Register(ModuleName, 9, "unauthorized address")
	ErrInvalidDenom            = sdkerrors.Register(ModuleName, 10, "invalid denom")
	ErrInvalidONFTID           = sdkerrors.Register(ModuleName, 11, "invalid ID")
	ErrInvalidONFTMeta         = sdkerrors.Register(ModuleName, 12, "invalid metadata")
	ErrInvalidMediaURI         = sdkerrors.Register(ModuleName, 13, "invalid media URI")
	ErrInvalidPreviewURI       = sdkerrors.Register(ModuleName, 14, "invalid preview URI")
	ErrNotTransferable         = sdkerrors.Register(ModuleName, 15, "onft is not transferable")
	ErrNotEditable             = sdkerrors.Register(ModuleName, 16, "onft is not editable")
	ErrInvalidOption           = sdkerrors.Register(ModuleName, 17, "invalid option")
	ErrInvalidName             = sdkerrors.Register(ModuleName, 18, "invalid name")
	ErrInvalidDescription      = sdkerrors.Register(ModuleName, 19, "invalid description")
	ErrInvalidURI              = sdkerrors.Register(ModuleName, 20, "invalid URI")
	ErrInvalidPercentage       = sdkerrors.Register(ModuleName, 21, "invalid percentage")
	ErrInvalidDenomCreationFee = sdkerrors.Register(ModuleName, 22, "invalid denom creation fee")
	ErrInvalidFeeDenom         = sdkerrors.Register(ModuleName, 23, "invalid creation fee denom")
	ErrNotEnoughFeeAmount      = sdkerrors.Register(ModuleName, 24, "invalid creation fee amount")

	// this line is used by starport scaffolding # ibc/errors
)
