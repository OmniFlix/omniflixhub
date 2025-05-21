package v6dev2

import (
	store "cosmossdk.io/store/types"
	"github.com/OmniFlix/omniflixhub/v6/app/upgrades"
	feemarkettypes "github.com/skip-mev/feemarket/x/feemarket/types"
)

const (
	UpgradeName         = "v6dev2"
	GlobalFeeModuleName = "globalfee"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV6Dev2UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			feemarkettypes.ModuleName,
		},
		Deleted: []string{
			GlobalFeeModuleName,
		},
	},
}
