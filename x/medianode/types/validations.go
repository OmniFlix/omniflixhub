package types

import (
	"fmt"
	"net/url"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Validate performs basic validation of the MediaNode fields
func (m MediaNode) Validate() error {
	if m.Id == "" {
		return fmt.Errorf("media node ID cannot be empty")
	}

	if m.Url == "" {
		return fmt.Errorf("media node URL cannot be empty")
	}

	// Validate URL format
	if _, err := url.Parse(m.Url); err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	if m.Owner == "" {
		return fmt.Errorf("owner address cannot be empty")
	}

	// Validate owner address format
	if _, err := sdk.AccAddressFromBech32(m.Owner); err != nil {
		return fmt.Errorf("invalid owner address: %w", err)
	}

	// Validate hardware specs
	if err := m.HardwareSpecs.Validate(); err != nil {
		return fmt.Errorf("invalid hardware specs: %w", err)
	}

	// Validate lease amount
	if !m.PricePerHour.IsValid() {
		return fmt.Errorf("invalid lease amount per hour")
	}
	if m.PricePerHour.IsZero() {
		return fmt.Errorf("lease amount per hour cannot be zero")
	}

	return nil
}

// Validate performs basic validation of the HardwareSpecs fields
func (h HardwareSpecs) Validate() error {
	if h.Cpus <= 0 {
		return fmt.Errorf("CPUs must be greater than 0")
	}
	if h.RamInGb <= 0 {
		return fmt.Errorf("RAM must be greater than 0")
	}
	if h.StorageInGb <= 0 {
		return fmt.Errorf("storage must be greater than 0")
	}
	return nil
}

// Validate performs basic validation of the Lease fields
func (l Lease) Validate() error {
	if l.MediaNodeId == "" {
		return fmt.Errorf("media node ID cannot be empty")
	}

	if l.Lessee == "" {
		return fmt.Errorf("leased to address cannot be empty")
	}

	// Validate leasedTo address format
	if _, err := sdk.AccAddressFromBech32(l.Lessee); err != nil {
		return fmt.Errorf("invalid leased to address: %w", err)
	}

	// Validate lease amount
	if !l.TotalLeaseAmount.IsValid() {
		return fmt.Errorf("invalid lease amount")
	}
	if l.TotalLeaseAmount.IsZero() {
		return fmt.Errorf("lease amount cannot be zero")
	}

	if l.LeasedHours == 0 {
		return fmt.Errorf("leased days cannot be 0")
	}

	return nil
}

// ValidateLeaseStatus validates the lease status string
func ValidateLeaseStatus(status LeaseStatus) error {
	switch status {
	case LEASE_STATUS_ACTIVE,
		LEASE_STATUS_CANCELLED,
		LEASE_STATUS_EXPIRED:
		return nil
	default:
		return fmt.Errorf("invalid lease status: %s", status)
	}
}

func validateMediaNodeId(id string) error {
	// Validate media node ID format
	if len(id) != MediaNodeIdLength || id[:2] != MediaNodeIdPrefix {
		return fmt.Errorf("media node ID must start with '%s' and be %d characters long", MediaNodeIdPrefix, MediaNodeIdLength)
	}
	return nil
}
