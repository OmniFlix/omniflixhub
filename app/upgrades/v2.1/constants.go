package v2_1

import (
	store "cosmossdk.io/store/types"
	"github.com/OmniFlix/omniflixhub/v5/app/upgrades"
)

const UpgradeName = "v2.1"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV2UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{},
	},
}
