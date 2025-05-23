package v6

import (
	store "cosmossdk.io/store/types"
	"github.com/OmniFlix/omniflixhub/v6/app/upgrades"
	medianodetypes "github.com/OmniFlix/omniflixhub/v6/x/medianode/types"
	feemarkettypes "github.com/skip-mev/feemarket/x/feemarket/types"
)

const (
	UpgradeName         = "v6"
	GlobalFeeModuleName = "globalfee"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV6UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			medianodetypes.StoreKey,
			feemarkettypes.StoreKey,
		},
		Deleted: []string{
			GlobalFeeModuleName,
		},
	},
}
