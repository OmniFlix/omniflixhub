package v4

import (
	store "cosmossdk.io/store/types"
	"github.com/OmniFlix/omniflixhub/v6/app/upgrades"
	ibchookstypes "github.com/cosmos/ibc-apps/modules/ibc-hooks/v8/types"
)

const UpgradeName = "v4"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV4UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{ibchookstypes.StoreKey},
	},
}
