package v4

import (
	"github.com/OmniFlix/omniflixhub/v5/app/upgrades"
	store "github.com/cosmos/cosmos-sdk/store/types"
	ibchookstypes "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7/types"
)

const UpgradeName = "v4"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV4UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{ibchookstypes.StoreKey},
	},
}
