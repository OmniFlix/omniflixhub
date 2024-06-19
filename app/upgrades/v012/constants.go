package v012

import (
	"github.com/OmniFlix/omniflixhub/v5/app/upgrades"
	store "github.com/cosmos/cosmos-sdk/store/types"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward/types"
	icahosttypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
)

const UpgradeName = "v0.12.x"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV012UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{icahosttypes.StoreKey, packetforwardtypes.StoreKey},
	},
}
