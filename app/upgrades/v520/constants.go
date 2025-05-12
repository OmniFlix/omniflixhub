package v520

import (
	store "cosmossdk.io/store/types"
	"github.com/OmniFlix/omniflixhub/v6/app/upgrades"
)

const UpgradeName = "v5.2.0"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV520UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{},
	},
}
