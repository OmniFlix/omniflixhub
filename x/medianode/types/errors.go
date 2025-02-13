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
	ErrInvalidLeaseDays       = errorsmod.Register(ModuleName, 6, "invalid lease days")
	ErrMediaNodeNotFound      = errorsmod.Register(ModuleName, 7, "media node not found")
	ErrInvalidLeaseDuration   = errorsmod.Register(ModuleName, 8, "invalid lease duration")
	ErrMediaNodeAlreadyLeased = errorsmod.Register(ModuleName, 9, "media node already leased")
	ErrLeaseNotFound          = errorsmod.Register(ModuleName, 10, "lease not found")
	ErrInvalidLeaseStatus     = errorsmod.Register(ModuleName, 11, "invalid lease status")
	ErrUnauthorized           = errorsmod.Register(ModuleName, 12, "unauthorized operation")
	ErrLeaseExpired           = errorsmod.Register(ModuleName, 13, "lease has expired")

	// Additional errors needed for keeper functions
	ErrMediaNodeExists        = errorsmod.Register(ModuleName, 14, "media node already exists")
	ErrMediaNodeDoesNotExist  = errorsmod.Register(ModuleName, 15, "media node does not exist")
	ErrMediaNodeNotLeased     = errorsmod.Register(ModuleName, 16, "media node is not leased")
	ErrInvalidMediaNodeStatus = errorsmod.Register(ModuleName, 17, "media node status not allows deposit")
	ErrLeaseNotActive         = errorsmod.Register(ModuleName, 18, "lease not in active state")
	ErrInsufficientDeposit    = errorsmod.Register(ModuleName, 19, "insufficient deposit")
	ErrLeaseNotAllowed        = errorsmod.Register(ModuleName, 20, "lease not allowed")
)
