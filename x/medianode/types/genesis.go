package types

import (
	"fmt"
)

func NewGenesisState(medianodes []MediaNode,
	leases []Lease,
	nextMediaNodeNumber uint64,
	params Params,
) *GenesisState {
	return &GenesisState{
		MediaNodes: medianodes,
		Leases:     leases,
		Params:     params,
		LastNodeId: nextMediaNodeNumber,
	}
}

// DefaultGenesis returns default genesis state as raw bytes for the medianode
// module.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:     DefaultParams(),
		MediaNodes: []MediaNode{},
		Leases:     []Lease{},
		LastNodeId: 0,
	}
}

// ValidateGenesis performs basic validation of medianode genesis data
func (gs GenesisState) ValidateGenesis() error {
	if err := gs.Params.Validate(); err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}

	// Validate MediaNodes
	for _, node := range gs.MediaNodes {
		if err := node.Validate(); err != nil {
			return fmt.Errorf("invalid media node: %w", err)
		}
	}

	// Validate Leases
	for _, lease := range gs.Leases {
		if err := lease.Validate(); err != nil {
			return fmt.Errorf("invalid lease: %w", err)
		}
	}

	// Validate LastNodeId is non-negative
	if gs.LastNodeId < 0 {
		return fmt.Errorf("last node ID cannot be negative")
	}

	return nil
}
