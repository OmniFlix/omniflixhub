package v5

import (
	store "cosmossdk.io/store/types"
	"github.com/OmniFlix/omniflixhub/v5/app/upgrades"
)

const UpgradeName = "v5"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV5UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{},
	},
}
