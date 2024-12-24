package v510

import (
	store "cosmossdk.io/store/types"
	"github.com/OmniFlix/omniflixhub/v5/app/upgrades"
)

const UpgradeName = "v5.1.0"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV510UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{},
	},
}
