package v6_1

import (
	store "cosmossdk.io/store/types"
	"github.com/OmniFlix/omniflixhub/v6/app/upgrades"
)

const UpgradeName = "v6.1"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV6_1UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added:   []string{},
		Deleted: []string{},
	},
}
