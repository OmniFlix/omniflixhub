package v3

import (
	store "cosmossdk.io/store/types"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/OmniFlix/omniflixhub/v6/app/upgrades"
)

const UpgradeName = "v3"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV3UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{wasmtypes.ModuleName},
	},
}
