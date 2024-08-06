package apptesting

import (
	"fmt"
	"os"
	"time"

	"github.com/cosmos/cosmos-sdk/x/consensus/testutil"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	coreheader "cosmossdk.io/core/header"
	sdkmath "cosmossdk.io/math"

	"github.com/OmniFlix/omniflixhub/v5/app"

	"cosmossdk.io/log"
	"cosmossdk.io/store/rootmulti"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto/ed25519"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/testutil"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"
)

var (
	SecondaryDenom       = "uxflx"
	SecondaryAmount      = sdkmath.NewInt(100000000)
	baseTestAccts        = []sdk.AccAddress{}
	defaultTestStartTime = time.Now().UTC()
	testDescription      = stakingtypes.NewDescription("test_moniker", "test_identity", "test_website", "test_security_contact", "test_details")
)

func init() {
	baseTestAccts = CreateRandomAccounts(3)
}

type KeeperTestHelper struct {
	suite.Suite

	// defaults to false,
	// set to true if any method that potentially alters baseapp/abci is used.
	// this controls whether we can reuse the app instance, or have to set a new one.
	hasUsedAbci bool
	// defaults to false, set to true if we want to use a new app instance with caching enabled.
	// then on new setup test call, we just drop the current cache.
	// this is not always enabled, because some tests may take a painful performance hit due to CacheKv.
	withCaching bool

	App         *app.OmniFlixApp
	Ctx         sdk.Context
	QueryHelper *baseapp.QueryServiceTestHelper
	TestAccs    []sdk.AccAddress
}

// Setup sets up basic environment for suite (App, Ctx, and test accounts)
func (s *KeeperTestHelper) Setup() {
	dir, err := os.MkdirTemp("", "omniflixhub-test-home")
	if err != nil {
		panic(fmt.Sprintf("failed creating temporary directory: %v", err))
	}
	s.T().Cleanup(func() { os.RemoveAll(dir); s.withCaching = false })
	s.App = app.SetupWithCustomHome(false, dir)
	s.setupGeneral()

	// Manually set validator signing info, otherwise we panic
	vals, err := s.App.StakingKeeper.GetAllValidators(s.Ctx)
	if err != nil {
		panic(err)
	}
	for _, val := range vals {
		consAddr, _ := val.GetConsAddr()
		signingInfo := slashingtypes.NewValidatorSigningInfo(
			consAddr,
			s.Ctx.BlockHeight(),
			s.Ctx.BlockTime().Unix(),
			time.Unix(0, 0),
			false,
			0,
		)
		err := s.App.SlashingKeeper.SetValidatorSigningInfo(s.Ctx, consAddr, signingInfo)
		if err != nil {
			panic(err)
		}
	}
}

func (s *KeeperTestHelper) SetupWithCustomChainId(chainId string) {
	dir, err := os.MkdirTemp("", "omniflixhub-test-home")
	if err != nil {
		panic(fmt.Sprintf("failed creating temporary directory: %v", err))
	}
	s.T().Cleanup(func() { os.RemoveAll(dir); s.withCaching = false })
	s.App = app.SetupWithCustomHomeAndChainId(false, dir, chainId)
	s.setupGeneralCustomChainId(chainId)

	// Manually set validator signing info, otherwise we panic
	vals, err := s.App.StakingKeeper.GetAllValidators(s.Ctx)
	if err != nil {
		panic(err)
	}
	for _, val := range vals {
		consAddr, _ := val.GetConsAddr()
		signingInfo := slashingtypes.NewValidatorSigningInfo(
			consAddr,
			s.Ctx.BlockHeight(),
			0,
			time.Unix(0, 0),
			false,
			0,
		)
		err := s.App.SlashingKeeper.SetValidatorSigningInfo(s.Ctx, consAddr, signingInfo)
		if err != nil {
			panic(err)
		}
	}
}

func (s *KeeperTestHelper) Reset() {
	if s.hasUsedAbci || !s.withCaching {
		s.withCaching = true
		s.Setup()
	} else {
		s.setupGeneral()
	}
}

func (s *KeeperTestHelper) SetupWithLevelDb() func() {
	app, cleanup := app.SetupTestingAppWithLevelDb(false)
	s.App = app
	s.setupGeneral()
	return cleanup
}

func (s *KeeperTestHelper) setupGeneral() {
	s.setupGeneralCustomChainId("omniflixhub-1")
}

func (s *KeeperTestHelper) setupGeneralCustomChainId(chainId string) {
	s.Ctx = s.App.BaseApp.NewContextLegacy(false, cmtproto.Header{Height: 1, ChainID: chainId, Time: defaultTestStartTime})
	if s.withCaching {
		s.Ctx, _ = s.Ctx.CacheContext()
	}
	s.QueryHelper = &baseapp.QueryServiceTestHelper{
		GRPCQueryRouter: s.App.GRPCQueryRouter(),
		Ctx:             s.Ctx,
	}

	s.TestAccs = []sdk.AccAddress{}
	s.TestAccs = append(s.TestAccs, baseTestAccts...)
	s.hasUsedAbci = false
}

func (s *KeeperTestHelper) SetupTestForInitGenesis() {
	// Setting to True, leads to init genesis not running
	s.App = app.Setup(true)
	s.Ctx = s.App.BaseApp.NewContextLegacy(true, cmtproto.Header{})
	s.hasUsedAbci = true
}

// CreateTestContext creates a test context.
func (s *KeeperTestHelper) CreateTestContext() sdk.Context {
	ctx := s.CreateTestContextWithMultiStore()
	return ctx
}

// CreateTestContextWithMultiStore creates a test context and returns it together with multi store.
func (s *KeeperTestHelper) CreateTestContextWithMultiStore() sdk.Context {
	db := dbm.NewMemDB()
	logger := log.NewNopLogger()

	ms := rootmulti.NewStore(db, logger, nil)

	return sdk.NewContext(ms, cmtproto.Header{}, false, logger)
}

func (s *KeeperTestHelper) Commit() {
	// oldHeight := s.Ctx.BlockHeight()
	// oldHeader := s.Ctx.BlockHeader()
	// commit, _ := s.App.Commit()
	// newHeader := cmtproto.Header{Height: oldHeight + 1, ChainID: oldHeader.ChainID, Time: oldHeader.Time.Add(time.Second)}
	// s.App.PreBlocker(s.Ctx, abci.Request_Commit{Commit: commit})
	// s.Ctx = s.App.BaseApp.NewContextLegacy(false, newHeader)

	_, err := s.App.FinalizeBlock(&abci.RequestFinalizeBlock{Height: s.Ctx.BlockHeight(), Time: s.Ctx.BlockTime()})
	if err != nil {
		panic(err)
	}
	_, err = s.App.Commit()
	if err != nil {
		panic(err)
	}

	newBlockTime := s.Ctx.BlockTime().Add(time.Second)

	header := s.Ctx.BlockHeader()
	header.Time = newBlockTime
	header.Height++

	s.Ctx = s.App.BaseApp.NewUncachedContext(false, header).WithHeaderInfo(coreheader.Info{
		Height: header.Height,
		Time:   header.Time,
	})

	s.hasUsedAbci = true
}

// FundAcc funds target address with specified amount.
func (s *KeeperTestHelper) FundAcc(acc sdk.AccAddress, amounts sdk.Coins) {
	err := banktestutil.FundAccount(s.Ctx, s.App.BankKeeper, acc, amounts)
	s.Require().NoError(err)
}

// FundModuleAcc funds target modules with specified amount.
func (s *KeeperTestHelper) FundModuleAcc(moduleName string, amounts sdk.Coins) {
	err := banktestutil.FundModuleAccount(s.Ctx, s.App.BankKeeper, moduleName, amounts)
	s.Require().NoError(err)
}

func (s *KeeperTestHelper) MintCoins(coins sdk.Coins) {
	err := s.App.BankKeeper.MintCoins(s.Ctx, minttypes.ModuleName, coins)
	s.Require().NoError(err)
}

// SetupValidator sets up a validator and returns the ValAddress.
func (s *KeeperTestHelper) SetupValidator(bondStatus stakingtypes.BondStatus) sdk.ValAddress {
	valPub := secp256k1.GenPrivKey().PubKey()
	valAddr := sdk.ValAddress(valPub.Address())
	stakingParams, err := s.App.StakingKeeper.GetParams(s.Ctx)
	s.Require().NoError(err)
	bondDenom := stakingParams.BondDenom
	bondAmt := sdk.DefaultPowerReduction
	selfBond := sdk.NewCoins(sdk.Coin{Amount: bondAmt, Denom: bondDenom})

	s.FundAcc(sdk.AccAddress(valAddr), selfBond)

	stakingCoin := sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: selfBond[0].Amount}
	ZeroCommission := stakingtypes.NewCommissionRates(sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec())
	valCreateMsg, err := stakingtypes.NewMsgCreateValidator(valAddr.String(), valPub, stakingCoin, testDescription, ZeroCommission, sdkmath.OneInt())
	s.Require().NoError(err)
	stakingMsgSvr := stakingkeeper.NewMsgServerImpl(s.App.StakingKeeper)
	res, err := stakingMsgSvr.CreateValidator(s.Ctx, valCreateMsg)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	val, err := s.App.StakingKeeper.GetValidator(s.Ctx, valAddr)
	s.Require().NoError(err)

	val = val.UpdateStatus(bondStatus)
	err = s.App.StakingKeeper.SetValidator(s.Ctx, val)
	s.Require().NoError(err)

	consAddr, err := val.GetConsAddr()
	s.Suite.Require().NoError(err)

	signingInfo := slashingtypes.NewValidatorSigningInfo(
		consAddr,
		s.Ctx.BlockHeight(),
		0,
		time.Unix(0, 0),
		false,
		0,
	)
	err = s.App.SlashingKeeper.SetValidatorSigningInfo(s.Ctx, consAddr, signingInfo)
	s.Require().NoError(err)

	return valAddr
}

// SetupMultipleValidators setups "numValidator" validators and returns their address in string
func (s *KeeperTestHelper) SetupMultipleValidators(numValidator int) []string {
	valAddrs := []string{}
	for i := 0; i < numValidator; i++ {
		valAddr := s.SetupValidator(stakingtypes.Bonded)
		valAddrs = append(valAddrs, valAddr.String())
	}
	return valAddrs
}

// BeginNewBlock starts a new block.
func (s *KeeperTestHelper) BeginNewBlock(executeNextEpoch bool) {
	var valAddr []byte

	validators, err := s.App.StakingKeeper.GetAllValidators(s.Ctx)
	s.Require().NoError(err)
	if len(validators) >= 1 {
		valAddrFancy, err := validators[0].GetConsAddr()
		s.Require().NoError(err)
		valAddr = valAddrFancy
	} else {
		valAddrFancy := s.SetupValidator(stakingtypes.Bonded)
		validator, _ := s.App.StakingKeeper.GetValidator(s.Ctx, valAddrFancy)
		valAddr2, _ := validator.GetConsAddr()
		valAddr = valAddr2
	}

	s.BeginNewBlockWithProposer(executeNextEpoch, valAddr)
}

// BeginNewBlockWithProposer begins a new block with a proposer.
func (s *KeeperTestHelper) BeginNewBlockWithProposer(executeNextEpoch bool, proposer sdk.ValAddress) {
	// validator, err := s.App.StakingKeeper.GetValidator(s.Ctx, proposer)
	// s.Assert().NoError(err)

	// valConsAddr, err := validator.GetConsAddr()
	// s.Require().NoError(err)

	// valAddr := valConsAddr

	validator, err := s.App.StakingKeeper.GetValidator(s.Ctx, proposer)
	s.Assert().NoError(err)

	valAddr, err := validator.GetConsAddr()
	s.Require().NoError(err)

	newBlockTime := s.Ctx.BlockTime().Add(5 * time.Second)

	header := cmtproto.Header{Height: s.Ctx.BlockHeight() + 1, Time: newBlockTime}
	newCtx := s.Ctx.WithBlockTime(newBlockTime).WithBlockHeight(s.Ctx.BlockHeight() + 1)
	s.Ctx = newCtx

	voteInfos := []abci.VoteInfo{
		{
			Validator:   abci.Validator{Address: valAddr, Power: 1000},
			BlockIdFlag: cmtproto.BlockIDFlagCommit,
		},
	}
	s.Ctx = s.Ctx.WithVoteInfos(voteInfos)

	fmt.Println("beginning block ", s.Ctx.BlockHeight())

	_, err = s.App.BeginBlocker(s.Ctx)
	s.Require().NoError(err)

	s.Ctx = s.App.NewContextLegacy(false, header)
	s.hasUsedAbci = true
}

// EndBlock ends the block, and runs commit
func (s *KeeperTestHelper) EndBlock() {
	_, err := s.App.EndBlocker(s.Ctx)
	s.Require().NoError(err)
	s.hasUsedAbci = true
}

func (s *KeeperTestHelper) RunMsg(msg sdk.Msg) (*sdk.Result, error) {
	// cursed that we have to copy this internal logic from SDK
	router := s.App.BaseApp.MsgServiceRouter()
	if handler := router.Handler(msg); handler != nil {
		// ADR 031 request type routing
		return handler(s.Ctx, msg)
	}
	s.FailNow("msg %v could not be ran", msg)
	return nil, fmt.Errorf("msg %v could not be ran", msg)
}

// AllocateRewardsToValidator allocates reward tokens to a distribution module then allocates rewards to the validator address.
func (s *KeeperTestHelper) AllocateRewardsToValidator(valAddr sdk.ValAddress, rewardAmt sdkmath.Int) {
	validator, err := s.App.StakingKeeper.GetValidator(s.Ctx, valAddr)
	s.Require().NoError(err)

	// allocate reward tokens to distribution module
	coins := sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, rewardAmt)}
	err = testutil.FundModuleAccount(s.App.BankKeeper, s.Ctx, distrtypes.ModuleName, coins)
	s.Require().NoError(err)

	// allocate rewards to validator
	s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
	decTokens := sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: sdkmath.LegacyNewDec(20000)}}
	err = s.App.DistrKeeper.AllocateTokensToValidator(s.Ctx, validator, decTokens)
	s.Require().NoError(err)
}

// BuildTx builds a transaction.
func (s *KeeperTestHelper) BuildTx(
	txBuilder client.TxBuilder,
	msgs []sdk.Msg,
	sigV2 signing.SignatureV2,
	memo string, txFee sdk.Coins,
	gasLimit uint64,
) authsigning.Tx {
	err := txBuilder.SetMsgs(msgs[0])
	s.Require().NoError(err)

	err = txBuilder.SetSignatures(sigV2)
	s.Require().NoError(err)

	txBuilder.SetMemo(memo)
	txBuilder.SetFeeAmount(txFee)
	txBuilder.SetGasLimit(gasLimit)

	return txBuilder.GetTx()
}

// CreateRandomAccounts is a function return a list of randomly generated AccAddresses
func CreateRandomAccounts(numAccts int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, numAccts)
	for i := 0; i < numAccts; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
	}

	return testAddrs
}

func GenerateTestAddrs() (string, string) {
	pk1 := ed25519.GenPrivKey().PubKey()
	validAddr := sdk.AccAddress(pk1.Address()).String()
	invalidAddr := sdk.AccAddress("invalid").String()
	return validAddr, invalidAddr
}
