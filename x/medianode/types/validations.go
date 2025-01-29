package types

import (
	"fmt"
	"net/url"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Validate performs basic validation of the MediaNode fields
func (m MediaNode) Validate() error {
	if m.Id == 0 {
		return fmt.Errorf("media node ID cannot be 0")
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
	if !m.LeaseAmountPerDay.IsValid() {
		return fmt.Errorf("invalid lease amount per day")
	}
	if m.LeaseAmountPerDay.IsZero() {
		return fmt.Errorf("lease amount per day cannot be zero")
	}

	return nil
}

// Validate performs basic validation of the HardwareSpecs fields
func (h HardwareSpecs) Validate() error {
	if h.Cpus <= 0 {
		return fmt.Errorf("CPUs must be greater than 0")
	}
	if h.Ram <= 0 {
		return fmt.Errorf("RAM must be greater than 0")
	}
	if h.Storage <= 0 {
		return fmt.Errorf("storage must be greater than 0")
	}
	return nil
}

// Validate performs basic validation of the Lease fields
func (l Lease) Validate() error {
	if l.MediaNodeId == 0 {
		return fmt.Errorf("media node ID cannot be 0")
	}

	if l.LeasedTo == "" {
		return fmt.Errorf("leased to address cannot be empty")
	}

	// Validate leasedTo address format
	if _, err := sdk.AccAddressFromBech32(l.LeasedTo); err != nil {
		return fmt.Errorf("invalid leased to address: %w", err)
	}

	// Validate lease amount
	if !l.LeaseAmount.IsValid() {
		return fmt.Errorf("invalid lease amount")
	}
	if l.LeaseAmount.IsZero() {
		return fmt.Errorf("lease amount cannot be zero")
	}

	if l.LeasedDays == 0 {
		return fmt.Errorf("leased days cannot be 0")
	}

	// Validate lease status
	if err := ValidateLeaseStatus(l.LeaseStatus); err != nil {
		return fmt.Errorf("invalid lease status: %w", err)
	}

	// Validate lease expiry is after leased at time
	if !l.LeaseExpiry.After(l.LeasedAt) {
		return fmt.Errorf("lease expiry must be after leased at time")
	}

	return nil
}

// ValidateLeaseStatus validates the lease status string
func ValidateLeaseStatus(status LeaseStatus) error {
	switch status {
	case LeaseStatus_LEASE_STATUS_ACTIVE,
		LeaseStatus_LEASE_STATUS_CANCELLED,
		LeaseStatus_LEASE_STATUS_EXPIRED:
		return nil
	default:
		return fmt.Errorf("invalid lease status: %s", status)
	}
}
