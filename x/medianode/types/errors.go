package types

import (
	errorsmod "cosmossdk.io/errors"
)

// medianode module errors
var (
	ErrInvalidMediaNodeID     = errorsmod.Register(ModuleName, 2, "invalid media node ID")
	ErrInvalidURL             = errorsmod.Register(ModuleName, 3, "invalid URL")
	ErrInvalidHardwareSpecs   = errorsmod.Register(ModuleName, 4, "invalid hardware specifications")
	ErrInvalidLeaseAmount     = errorsmod.Register(ModuleName, 5, "invalid lease amount")
	ErrMediaNodeNotFound      = errorsmod.Register(ModuleName, 6, "media node not found")
	ErrInvalidLeaseDuration   = errorsmod.Register(ModuleName, 7, "invalid lease duration")
	ErrMediaNodeAlreadyLeased = errorsmod.Register(ModuleName, 8, "media node already leased")
	ErrLeaseNotFound          = errorsmod.Register(ModuleName, 9, "lease not found")
	ErrInvalidLeaseStatus     = errorsmod.Register(ModuleName, 10, "invalid lease status")
	ErrUnauthorized           = errorsmod.Register(ModuleName, 11, "unauthorized operation")
	ErrLeaseExpired           = errorsmod.Register(ModuleName, 12, "lease has expired")
)
