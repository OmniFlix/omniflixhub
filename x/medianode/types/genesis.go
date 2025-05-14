package types

import (
	"fmt"
)

func NewGenesisState(medianodes []MediaNode,
	leases []Lease,
	mediaNodeCounter uint64,
	params Params,
) *GenesisState {
	return &GenesisState{
		Nodes:       medianodes,
		Leases:      leases,
		Params:      params,
		NodeCounter: mediaNodeCounter,
	}
}

// DefaultGenesis returns default genesis state as raw bytes for the medianode
// module.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:      DefaultParams(),
		Nodes:       []MediaNode{},
		Leases:      []Lease{},
		NodeCounter: 0,
	}
}

// ValidateGenesis performs basic validation of medianode genesis data
func (gs GenesisState) ValidateGenesis() error {
	if err := gs.Params.Validate(); err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}

	// Validate MediaNodes
	for _, node := range gs.Nodes {
		if err := node.Validate(); err != nil {
			return fmt.Errorf("invalid medianode: %w", err)
		}
	}

	// Validate Leases
	for _, lease := range gs.Leases {
		if err := lease.Validate(); err != nil {
			return fmt.Errorf("invalid lease: %w", err)
		}
	}

	return nil
}
