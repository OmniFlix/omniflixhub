package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrInvalidCollection       = errorsmod.Register(ModuleName, 3, "invalid ONFT collection")
	ErrUnknownCollection       = errorsmod.Register(ModuleName, 4, "unknown ONFT collection")
	ErrInvalidONFT             = errorsmod.Register(ModuleName, 5, "invalid ONFT")
	ErrONFTAlreadyExists       = errorsmod.Register(ModuleName, 6, "ONFT already exists")
	ErrUnknownONFT             = errorsmod.Register(ModuleName, 7, "unknown ONFT")
	ErrEmptyMetaData           = errorsmod.Register(ModuleName, 8, "ONFT MetaData can't be empty")
	ErrUnauthorized            = errorsmod.Register(ModuleName, 9, "unauthorized address")
	ErrInvalidDenom            = errorsmod.Register(ModuleName, 10, "invalid denom")
	ErrInvalidONFTID           = errorsmod.Register(ModuleName, 11, "invalid ID")
	ErrInvalidONFTMeta         = errorsmod.Register(ModuleName, 12, "invalid metadata")
	ErrInvalidMediaURI         = errorsmod.Register(ModuleName, 13, "invalid media URI")
	ErrInvalidPreviewURI       = errorsmod.Register(ModuleName, 14, "invalid preview URI")
	ErrNotTransferable         = errorsmod.Register(ModuleName, 15, "onft is not transferable")
	ErrNotEditable             = errorsmod.Register(ModuleName, 16, "onft is not editable")
	ErrInvalidOption           = errorsmod.Register(ModuleName, 17, "invalid option")
	ErrInvalidName             = errorsmod.Register(ModuleName, 18, "invalid name")
	ErrInvalidDescription      = errorsmod.Register(ModuleName, 19, "invalid description")
	ErrInvalidURI              = errorsmod.Register(ModuleName, 20, "invalid URI")
	ErrInvalidPercentage       = errorsmod.Register(ModuleName, 21, "invalid percentage")
	ErrInvalidDenomCreationFee = errorsmod.Register(ModuleName, 22, "invalid denom creation fee")
	ErrInvalidFeeDenom         = errorsmod.Register(ModuleName, 23, "invalid creation fee denom")
	ErrNotEnoughFeeAmount      = errorsmod.Register(ModuleName, 24, "invalid creation fee amount")
	ErrInvalidONFTMetadata     = errorsmod.Register(ModuleName, 25, "invalid nft data")
	ErrDenomIdExists           = errorsmod.Register(ModuleName, 26, "denom exists")
)
