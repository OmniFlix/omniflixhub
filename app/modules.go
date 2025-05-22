package app

import (
	circuittypes "cosmossdk.io/x/circuit/types"
	"cosmossdk.io/x/nft"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	appparams "github.com/OmniFlix/omniflixhub/v6/app/params"
	nfttransfer "github.com/bianjieai/nft-transfer"
	ibcnfttransfertypes "github.com/bianjieai/nft-transfer/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/group"

	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"

	"github.com/cosmos/cosmos-sdk/x/authz"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"

	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/cosmos/ibc-go/modules/capability"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"

	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"

	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"cosmossdk.io/x/evidence"
	evidencetypes "cosmossdk.io/x/evidence/types"

	"cosmossdk.io/x/feegrant"
	feegrantmodule "cosmossdk.io/x/feegrant/module"

	"github.com/skip-mev/feemarket/x/feemarket"
	feemarkettypes "github.com/skip-mev/feemarket/x/feemarket/types"

	"cosmossdk.io/x/circuit"

	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"

	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	groupmodule "github.com/cosmos/cosmos-sdk/x/group/module"

	"github.com/cosmos/cosmos-sdk/x/mint"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/cosmos/cosmos-sdk/x/params"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/OmniFlix/omniflixhub/v6/x/tokenfactory"
	tokenfactorytypes "github.com/OmniFlix/omniflixhub/v6/x/tokenfactory/types"

	"cosmossdk.io/x/upgrade"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v8/modules/core"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"

	ica "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	"github.com/cosmos/ibc-go/v8/modules/apps/transfer"

	icq "github.com/cosmos/ibc-apps/modules/async-icq/v8"
	icqtypes "github.com/cosmos/ibc-apps/modules/async-icq/v8/types"

	packetforward "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/types"

	ibchooks "github.com/cosmos/ibc-apps/modules/ibc-hooks/v8"
	ibchookstypes "github.com/cosmos/ibc-apps/modules/ibc-hooks/v8/types"

	"github.com/OmniFlix/omniflixhub/v6/x/alloc"
	alloctypes "github.com/OmniFlix/omniflixhub/v6/x/alloc/types"

	"github.com/OmniFlix/omniflixhub/v6/x/onft"
	onfttypes "github.com/OmniFlix/omniflixhub/v6/x/onft/types"

	"github.com/OmniFlix/omniflixhub/v6/x/marketplace"
	marketplacetypes "github.com/OmniFlix/omniflixhub/v6/x/marketplace/types"

	"github.com/OmniFlix/omniflixhub/v6/x/itc"
	itctypes "github.com/OmniFlix/omniflixhub/v6/x/itc/types"

	"github.com/OmniFlix/omniflixhub/v6/x/medianode"
	medianodetypes "github.com/OmniFlix/omniflixhub/v6/x/medianode/types"

	"github.com/OmniFlix/streampay/v2/x/streampay"
	streampaytypes "github.com/OmniFlix/streampay/v2/x/streampay/types"
)

var (
	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(getGovProposalHandlers()),
		wasm.AppModuleBasic{},
		groupmodule.AppModuleBasic{},
		params.AppModuleBasic{},
		consensus.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		ibc.AppModuleBasic{},
		ibctm.AppModuleBasic{},
		ica.AppModuleBasic{},
		icq.AppModuleBasic{},
		packetforward.AppModuleBasic{},
		ibchooks.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		circuit.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		transfer.AppModuleBasic{},
		vesting.AppModuleBasic{},
		feemarket.AppModuleBasic{},
		tokenfactory.AppModuleBasic{},
		nfttransfer.AppModuleBasic{},

		alloc.AppModuleBasic{},
		onft.AppModuleBasic{},
		marketplace.AppModuleBasic{},
		streampay.AppModuleBasic{},
		itc.AppModuleBasic{},
		medianode.AppModuleBasic{},
		// this line is used by starport scaffolding # stargate/app/moduleBasic
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:      nil,
		distrtypes.ModuleName:           nil,
		minttypes.ModuleName:            {authtypes.Minter},
		stakingtypes.BondedPoolName:     {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName:  {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:             {authtypes.Burner},
		ibctransfertypes.ModuleName:     {authtypes.Minter, authtypes.Burner},
		ibcnfttransfertypes.ModuleName:  nil,
		icatypes.ModuleName:             nil,
		tokenfactorytypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		feemarkettypes.ModuleName:       nil,
		feemarkettypes.FeeCollectorName: nil,
		wasmtypes.ModuleName:            {authtypes.Burner},
		alloctypes.ModuleName:           {authtypes.Minter, authtypes.Burner, authtypes.Staking},
		nft.ModuleName:                  nil,
		onfttypes.ModuleName:            nil,
		marketplacetypes.ModuleName:     nil,
		streampaytypes.ModuleName:       nil,
		itctypes.ModuleName:             nil,
		medianodetypes.ModuleName:       nil,
	}
)

func appModules(
	app *OmniFlixApp,
	encodingConfig appparams.EncodingConfig,
	skipGenesisInvariants bool,
) []module.AppModule {
	appCodec := encodingConfig.Marshaler
	return []module.AppModule{
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app.BaseApp,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.AccountKeeper, nil, app.GetSubspace(authtypes.ModuleName)),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
		crisis.NewAppModule(app.CrisisKeeper, skipGenesisInvariants, app.GetSubspace(crisistypes.ModuleName)),
		gov.NewAppModule(
			appCodec,
			&app.GovKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.GetSubspace(govtypes.ModuleName),
		),
		groupmodule.NewAppModule(
			appCodec,
			app.GroupKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.interfaceRegistry,
		),
		tokenfactory.NewAppModule(
			app.AppKeepers.TokenFactoryKeeper,
			app.AppKeepers.AccountKeeper,
			app.AppKeepers.BankKeeper,
		),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper, nil, app.GetSubspace(minttypes.ModuleName)),
		slashing.NewAppModule(
			appCodec,
			app.SlashingKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.StakingKeeper,
			app.GetSubspace(slashingtypes.ModuleName),
			app.InterfaceRegistry(),
		),
		distr.NewAppModule(
			appCodec,
			app.DistrKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.StakingKeeper,
			app.GetSubspace(distrtypes.ModuleName),
		),
		staking.NewAppModule(
			appCodec,
			app.StakingKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.GetSubspace(stakingtypes.ModuleName),
		),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		circuit.NewAppModule(appCodec, app.CircuitKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper, addresscodec.NewBech32Codec(appparams.Bech32AccountAddrPrefix)),
		evidence.NewAppModule(app.EvidenceKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		params.NewAppModule(app.ParamsKeeper),
		consensus.NewAppModule(appCodec, app.ConsensusParamsKeeper),
		transfer.NewAppModule(app.TransferKeeper),
		wasm.NewAppModule(
			appCodec,
			&app.WasmKeeper,
			app.StakingKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.MsgServiceRouter(),
			app.GetSubspace(wasmtypes.ModuleName),
		),
		ica.NewAppModule(nil, &app.ICAHostKeeper),
		icq.NewAppModule(app.ICQKeeper, app.GetSubspace(icqtypes.ModuleName)),
		nfttransfer.NewAppModule(app.IBCNFTTransferKeeper),
		packetforward.NewAppModule(app.PacketForwardKeeper, app.GetSubspace(packetforwardtypes.ModuleName)),
		ibchooks.NewAppModule(app.AccountKeeper),
		feemarket.NewAppModule(appCodec, *app.FeeMarketKeeper),
		alloc.NewAppModule(appCodec, app.AllocKeeper, app.GetSubspace(alloctypes.ModuleName)),
		onft.NewAppModule(
			appCodec,
			app.ONFTKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.DistrKeeper,
			app.GetSubspace(onfttypes.ModuleName),
		),
		marketplace.NewAppModule(appCodec, app.MarketplaceKeeper, app.GetSubspace(marketplacetypes.ModuleName)),
		streampay.NewAppModule(appCodec, app.StreamPayKeeper),
		itc.NewAppModule(appCodec, app.ItcKeeper, app.GetSubspace(itctypes.ModuleName)),
		medianode.NewAppModule(appCodec, app.MedianodeKeeper),
	}
}

// simulationModules returns modules for simulation manager
// define the order of the modules for deterministic simulations
func simulationModules(
	app *OmniFlixApp,
	encodingConfig appparams.EncodingConfig,
	_ bool,
) []module.AppModuleSimulation {
	appCodec := encodingConfig.Marshaler

	return []module.AppModuleSimulation{
		auth.NewAppModule(
			appCodec,
			app.AccountKeeper,
			authsims.RandomGenesisAccounts,
			app.GetSubspace(authtypes.ModuleName),
		),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		gov.NewAppModule(appCodec, &app.GovKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper, nil, app.GetSubspace(minttypes.ModuleName)),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(distrtypes.ModuleName)),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(slashingtypes.ModuleName), app.InterfaceRegistry()),
		params.NewAppModule(app.ParamsKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		ibc.NewAppModule(app.IBCKeeper),
	}
}

/*
orderBeginBlockers tells the app's module manager how to set the order of
BeginBlockers, which are run at the beginning of every block.

During begin block slashing happens after distr.BeginBlocker so that
there is nothing left over in the validator fee pool, so to keep the
CanWithdrawInvariant invariant.
NOTE: staking module is required if HistoricalEntries param > 0
NOTE: capability module's beginBlocker must come before any modules using capabilities (e.g. IBC)
*/

func orderBeginBlockers() []string {
	return []string{
		// Changed this as part of v0.50.x migration, moved this to  preblockers
		// upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		minttypes.ModuleName,
		alloctypes.ModuleName, // must run before distribution module
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		ibcexported.ModuleName,
		vestingtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		paramstypes.ModuleName,
		consensusparamtypes.ModuleName,
		ibchookstypes.ModuleName,
		wasmtypes.ModuleName,
		ibctransfertypes.ModuleName,
		icatypes.ModuleName,
		icqtypes.ModuleName,
		packetforwardtypes.ModuleName,
		ibcnfttransfertypes.ModuleName,
		genutiltypes.ModuleName,
		authz.ModuleName,
		authtypes.ModuleName,
		crisistypes.ModuleName,
		feegrant.ModuleName,
		circuittypes.ModuleName,
		feemarkettypes.ModuleName,
		tokenfactorytypes.ModuleName,
		group.ModuleName,
		onfttypes.ModuleName,
		marketplacetypes.ModuleName,
		streampaytypes.ModuleName,
		itctypes.ModuleName,
		medianodetypes.ModuleName,
	}
}

func orderEndBlockers() []string {
	return []string{
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		capabilitytypes.ModuleName,
		genutiltypes.ModuleName,
		banktypes.ModuleName,
		upgradetypes.ModuleName,
		evidencetypes.ModuleName,
		authtypes.ModuleName,
		vestingtypes.ModuleName,
		paramstypes.ModuleName,
		consensusparamtypes.ModuleName,
		ibchookstypes.ModuleName,
		wasmtypes.ModuleName,
		ibctransfertypes.ModuleName,
		icatypes.ModuleName,
		icqtypes.ModuleName,
		packetforwardtypes.ModuleName,
		ibcnfttransfertypes.ModuleName,
		minttypes.ModuleName,
		slashingtypes.ModuleName,
		distrtypes.ModuleName,
		ibcexported.ModuleName,
		feegrant.ModuleName,
		circuittypes.ModuleName,
		feemarkettypes.ModuleName,
		group.ModuleName,
		tokenfactorytypes.ModuleName,
		authz.ModuleName,
		alloctypes.ModuleName,
		onfttypes.ModuleName,
		marketplacetypes.ModuleName,
		streampaytypes.ModuleName,
		itctypes.ModuleName,
		medianodetypes.ModuleName,
	}
}

/*
NOTE: The genutils module must occur after staking so that pools are
properly initialized with tokens from genesis accounts.
NOTE: The genutils module must also occur after auth so that it can access the params from auth.
NOTE: Capability module must occur first so that it can initialize any capabilities
so that other modules that want to create or claim capabilities after wards in InitChain
can do so safely.
*/

func orderInitGenesis() []string {
	return []string{
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		paramstypes.ModuleName,
		consensusparamtypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		feegrant.ModuleName,
		circuittypes.ModuleName,
		ibchookstypes.ModuleName,
		wasmtypes.ModuleName,
		feemarkettypes.ModuleName,
		group.ModuleName,
		tokenfactorytypes.ModuleName,
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		icatypes.ModuleName,
		icqtypes.ModuleName,
		packetforwardtypes.ModuleName,
		ibcnfttransfertypes.ModuleName,
		alloctypes.ModuleName,
		onfttypes.ModuleName,
		marketplacetypes.ModuleName,
		streampaytypes.ModuleName,
		itctypes.ModuleName,
		medianodetypes.ModuleName,
	}
}
