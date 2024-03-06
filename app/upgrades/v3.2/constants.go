package v3_2

import (
	"github.com/OmniFlix/omniflixhub/v3/app/upgrades"
	store "github.com/cosmos/cosmos-sdk/store/types"
	ibchookstypes "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7/types"
)

const UpgradeName = "v3.2"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV3_2UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{ibchookstypes.StoreKey},
	},
}
