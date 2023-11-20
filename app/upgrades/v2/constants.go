package v2

import (
	"github.com/OmniFlix/omniflixhub/v2/app/upgrades"
	ibcnfttransfertypes "github.com/bianjieai/nft-transfer/types"
	store "github.com/cosmos/cosmos-sdk/store/types"
)

const UpgradeName = "v2-dev3"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV2UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			ibcnfttransfertypes.ModuleName,
		},
	},
}
