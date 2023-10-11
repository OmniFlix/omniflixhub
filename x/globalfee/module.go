package globalfee

import (
	"context"
	"encoding/json"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	abci "github.com/cometbft/cometbft/abci/types"

	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/OmniFlix/omniflixhub/v2/x/globalfee/client/cli"
	"github.com/OmniFlix/omniflixhub/v2/x/globalfee/keeper"
	"github.com/OmniFlix/omniflixhub/v2/x/globalfee/types"
)

var (
	_ module.AppModuleBasic   = AppModuleBasic{}
	_ module.AppModuleGenesis = AppModule{}
	_ module.AppModule        = AppModule{}
)

// ConsensusVersion defines the current x/globalfee module consensus version.
const ConsensusVersion = 2

// AppModuleBasic defines the basic application module used by the wasm module.
type AppModuleBasic struct {
	cdc codec.Codec
}

func (am AppModuleBasic) Name() string {
	return types.ModuleName
}

func (am AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

func (am AppModuleBasic) ValidateGenesis(marshaler codec.JSONCodec, _ client.TxEncodingConfig, message json.RawMessage) error {
	var data types.GenesisState
	err := marshaler.UnmarshalJSON(message, &data)
	if err != nil {
		return err
	}
	if err := data.Params.Validate(); err != nil {
		return errorsmod.Wrap(err, "params")
	}
	return nil
}

func (am AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))
	if err != nil {
		panic(err)
	}
}

func (am AppModuleBasic) GetTxCmd() *cobra.Command {
	return nil
}

func (am AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

func (am AppModuleBasic) RegisterInterfaces(r codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(r)
}

type AppModule struct {
	AppModuleBasic

	keeper keeper.Keeper

	// bondDenom is used solely for migration off of x/params
	bondDenom string
}

// NewAppModule constructor
func NewAppModule(
	cdc codec.Codec,
	keeper keeper.Keeper,
	debondDenom string,
) *AppModule {
	return &AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc},
		keeper:         keeper,
		bondDenom:      debondDenom,
	}
}

func (am AppModule) InitGenesis(ctx sdk.Context, marshaler codec.JSONCodec, message json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	marshaler.MustUnmarshalJSON(message, &genesisState)
	_ = am.keeper.SetParams(ctx, genesisState.Params)
	return nil
}

func (am AppModule) ExportGenesis(ctx sdk.Context, marshaler codec.JSONCodec) json.RawMessage {
	params := am.keeper.GetParams(ctx)
	genState := types.NewGenesisState(params)
	return marshaler.MustMarshalJSON(genState)
}

func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {
}

func (am AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewGrpcQuerier(am.keeper))
}

func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {
}

func (am AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return nil
}

// ConsensusVersion is a sequence number for state-breaking change of the
// module. It should be incremented on each consensus-breaking change
// introduced by the module. To avoid wrong/empty versions, the initial version
// should be set to 1.
func (am AppModule) ConsensusVersion() uint64 {
	return ConsensusVersion
}
