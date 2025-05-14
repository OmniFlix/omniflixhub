package v2

import (
	store "cosmossdk.io/store/types"
	"github.com/OmniFlix/omniflixhub/v6/app/upgrades"
	globalfeetypes "github.com/OmniFlix/omniflixhub/v6/x/globalfee/types"
	tokenfactorytypes "github.com/OmniFlix/omniflixhub/v6/x/tokenfactory/types"
	ibcnfttransfertypes "github.com/bianjieai/nft-transfer/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/group"
	icqtypes "github.com/cosmos/ibc-apps/modules/async-icq/v8/types"
)

const UpgradeName = "v2"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV2UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			consensustypes.ModuleName,
			crisistypes.ModuleName,
			globalfeetypes.ModuleName,
			group.ModuleName,
			icqtypes.ModuleName,
			tokenfactorytypes.ModuleName,
			ibcnfttransfertypes.ModuleName,
		},
	},
}
