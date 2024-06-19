package v2

import (
	"github.com/OmniFlix/omniflixhub/v5/app/upgrades"
	globalfeetypes "github.com/OmniFlix/omniflixhub/v5/x/globalfee/types"
	tokenfactorytypes "github.com/OmniFlix/omniflixhub/v5/x/tokenfactory/types"
	ibcnfttransfertypes "github.com/bianjieai/nft-transfer/types"
	store "github.com/cosmos/cosmos-sdk/store/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/group"
	icqtypes "github.com/cosmos/ibc-apps/modules/async-icq/v7/types"
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
