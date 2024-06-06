package v5

import (
	"github.com/OmniFlix/omniflixhub/v5/app/upgrades"
	store "github.com/cosmos/cosmos-sdk/store/types"
	feeabstypes "github.com/osmosis-labs/fee-abstraction/v7/x/feeabs/types"
)

const UpgradeName = "v5"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV4UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{feeabstypes.StoreKey},
	},
}
