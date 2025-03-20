package types

import (
	"fmt"
	"net/url"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Validate performs basic validation of the MediaNode fields
func (m MediaNode) Validate() error {
	if err := validateMediaNodeId(m.Id); err != nil {
		return err
	}

	if err := validateMediaNodeURL(m.Url); err != nil {
		return err
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

	// Validate info
	if err := m.Info.Validate(); err != nil {
		return fmt.Errorf("invalid medianode info")
	}

	// Validate lease amount
	if !m.PricePerHour.IsValid() {
		return fmt.Errorf("invalid lease amount per hour")
	}
	if m.PricePerHour.Amount.IsZero() {
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

func validateMediaNodeURL(URL string) error {
	if len(URL) == 0 {
		return fmt.Errorf("medianode URL cannot be empty")
	}
	if len(URL) > MaxURLLength {
		return fmt.Errorf("medianode URL exceeds %d characters", MaxURLLength)
	}
	u, err := url.Parse(URL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("medianode URL must be a valid URL (e.g., https://example.com)")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("medianode URL must use http or https")
	}
	return nil
}

// Validate checks that the Info fields adhere to the defined limits and constraints.
func (i Info) Validate() error {
	// Ensure the name is non empty.
	if len(i.Moniker) == 0 {
		return fmt.Errorf("name cannot be empty")
	}
	if len(i.Moniker) > MaxMediaNodeNameLength {
		return fmt.Errorf("name length %d exceeds limit of %d", len(i.Moniker), MaxMediaNodeNameLength)
	}
	if len(i.Description) > MaxMediaNodeDescriptionLength {
		return fmt.Errorf("description length %d exceeds limit of %d", len(i.Description), MaxMediaNodeDescriptionLength)
	}
	if len(i.Contact) > MaxContactLength {
		return fmt.Errorf("contact length %d exceeds limit of %d", len(i.Contact), MaxContactLength)
	}
	return nil
}
