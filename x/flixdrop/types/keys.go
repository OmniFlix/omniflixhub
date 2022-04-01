package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName defines the module name
	ModuleName = "flixdrop"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_claim"

	// ParamsKey defines the store key for claim module parameters
	ParamsKey = "params"

	// ActionKey defines the store key to store user accomplished actions
	ActionKey = "action"
)

// KVStore keys
var (
	// ClaimRecordsStorePrefix defines the store prefix for the claim records
	ClaimRecordsStorePrefix = []byte{0x01}
	ActionsPrefix           = []byte{0x02}
)

func KeyActionPrefix(index uint64) []byte {
	return append(ActionsPrefix, sdk.Uint64ToBigEndian(index)...)
}
