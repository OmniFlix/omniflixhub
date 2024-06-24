package v012

import (
	store "cosmossdk.io/store/types"
	"github.com/OmniFlix/omniflixhub/v5/app/upgrades"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/types"
	icahosttypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/types"
)

const UpgradeName = "v0.12.x"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV012UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{icahosttypes.StoreKey, packetforwardtypes.StoreKey},
	},
}
