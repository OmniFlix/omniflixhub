package keeper_test

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/OmniFlix/omniflixhub/v6/app/apptesting"
	"github.com/OmniFlix/omniflixhub/v6/x/medianode/keeper"
	"github.com/OmniFlix/omniflixhub/v6/x/medianode/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper

	queryClient types.QueryClient
	msgServer   types.MsgServer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

var (
	defaultMediaNodeURL = "https://default-medianode.omniflx.studio"
	defaultMediaNodeId  = "mn01953d3f737572d20082215cdaa12001"

	defaultInfo = types.Info{
		Moniker:     "Test Media Node",
		Description: "Test Media Node description",
		Contact:     "test-medianode@omniflix.network",
	}
	defaultHardwareSpecs = types.HardwareSpecs{
		Cpus:        4,
		RamInGb:     100,
		StorageInGb: 100,
	}
	defaultPricePerHour  = sdk.NewCoin("uflix", sdkmath.NewInt(10_000_000))
	defaultDepositAmount = sdk.NewCoin("uflix", sdkmath.NewInt(100_000_000))
)

func (suite *KeeperTestSuite) SetupTest() {
	suite.Setup()

	// Fund test accounts
	fundAccsAmount := sdk.NewCoins(sdk.NewCoin(apptesting.PrimaryDenom, sdkmath.NewInt(100000000000)))
	for _, acc := range suite.TestAccs {
		suite.FundAcc(acc, fundAccsAmount)
	}

	suite.queryClient = types.NewQueryClient(suite.QueryHelper)
	suite.msgServer = keeper.NewMsgServerImpl(suite.App.MedianodeKeeper)
}

func (suite *KeeperTestSuite) TestRegisterMediaNode() {
	suite.SetupTest()
	keeper := suite.App.MedianodeKeeper

	testCases := []struct {
		name      string
		msg       *types.MsgRegisterMediaNode
		expectErr bool
	}{
		{
			name: "valid media node",
			msg: &types.MsgRegisterMediaNode{
				Id:            defaultMediaNodeId,
				Url:           defaultMediaNodeURL,
				HardwareSpecs: defaultHardwareSpecs,
				PricePerHour:  defaultPricePerHour,
				Sender:        suite.TestAccs[0].String(),
				Deposit:       &defaultDepositAmount,
				Info:          defaultInfo,
			},
			expectErr: false,
		},
		{
			name: "duplicate medianode ID",
			msg: &types.MsgRegisterMediaNode{
				Id:            defaultMediaNodeId,
				Url:           "https://localhost:8188/medianode",
				HardwareSpecs: defaultHardwareSpecs,
				PricePerHour:  defaultPricePerHour,
				Sender:        suite.TestAccs[0].String(),
				Deposit:       &defaultDepositAmount,
				Info:          types.Info{Moniker: "another name", Description: "another desc", Contact: "another@media.com"},
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			_, err := suite.msgServer.RegisterMediaNode(suite.Ctx, tc.msg)
			if tc.expectErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				// Verify node was created
				node, found := keeper.GetMediaNode(suite.Ctx, tc.msg.Id)
				suite.Require().True(found)
				suite.Require().Equal(tc.msg.Info.Moniker, node.Info.Moniker)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGetAllMediaNodes() {
	suite.SetupTest()
	keeper := suite.App.MedianodeKeeper

	// Initially should be empty
	nodes := keeper.GetAllMediaNodes(suite.Ctx)
	suite.Require().Empty(nodes)

	// register a node
	msg, _ := types.NewMsgRegisterMediaNode(
		defaultMediaNodeURL,
		defaultInfo,
		defaultHardwareSpecs,
		defaultPricePerHour,
		defaultDepositAmount,
		suite.TestAccs[0].String(),
	)
	_, err := suite.msgServer.RegisterMediaNode(suite.Ctx, msg)
	suite.Require().NoError(err)

	// Should now have one node
	nodes = keeper.GetAllMediaNodes(suite.Ctx)
	suite.Require().Len(nodes, 1)
}

func (suite *KeeperTestSuite) TestUpdateMediaNode() {
	suite.SetupTest()
	keeper := suite.App.MedianodeKeeper

	// Create initial node
	createMsg, _ := types.NewMsgRegisterMediaNode(
		defaultMediaNodeURL,
		defaultInfo,
		defaultHardwareSpecs,
		defaultPricePerHour,
		defaultDepositAmount,
		suite.TestAccs[0].String(),
	)
	_, err := suite.msgServer.RegisterMediaNode(suite.Ctx, createMsg)
	suite.Require().NoError(err)

	// Update the node
	updateMsg := types.NewMsgUpdateMediaNode(
		createMsg.Id,
		&types.Info{Moniker: "updated medianode name", Description: "updated description", Contact: "contact@medianode.com"},
		nil,
		nil,
		suite.TestAccs[0].String(),
	)
	_, err = suite.msgServer.UpdateMediaNode(suite.Ctx, updateMsg)
	suite.Require().NoError(err)

	// Verify updates
	node, found := keeper.GetMediaNode(suite.Ctx, createMsg.Id)
	suite.Require().True(found)
	suite.Require().Equal(updateMsg.Info.Moniker, node.Info.Moniker)
	suite.Require().Equal(updateMsg.Info.Description, node.Info.Description)
	suite.Require().Equal(createMsg.HardwareSpecs.Cpus, node.HardwareSpecs.Cpus)
}

func (suite *KeeperTestSuite) TestLeaseMediaNode() {
	suite.SetupTest()
	keeper := suite.App.MedianodeKeeper

	// First, register a new media node
	createMsg, err := types.NewMsgRegisterMediaNode(
		defaultMediaNodeURL,
		defaultInfo,
		defaultHardwareSpecs,
		defaultPricePerHour,
		defaultDepositAmount,
		suite.TestAccs[0].String(),
	)
	suite.Require().NoError(err)
	_, err = suite.msgServer.RegisterMediaNode(suite.Ctx, createMsg)
	suite.Require().NoError(err)

	// Lease the media node using a different account (e.g. TestAccs[1])
	leaseHours := uint64(24)
	leaseMsg := types.NewMsgLeaseMediaNode(
		createMsg.Id,
		leaseHours,
		sdk.NewCoin(defaultPricePerHour.Denom, defaultPricePerHour.Amount.Mul(sdkmath.NewIntFromUint64(leaseHours))),
		suite.TestAccs[1].String(),
	)
	_, err = suite.msgServer.LeaseMediaNode(suite.Ctx, leaseMsg)
	suite.Require().NoError(err)

	// Verify that the media node now reflects the lease
	node, found := keeper.GetMediaNode(suite.Ctx, createMsg.Id)
	suite.Require().True(found)
	suite.Require().True(node.Leased)

	lease, found := keeper.GetMediaNodeLease(suite.Ctx, createMsg.Id)
	suite.Require().True(found)
	suite.Require().Equal(suite.TestAccs[1].String(), lease.Lessee, "media node should be leased by the leasing account")
	suite.Require().Equal(leaseHours, lease.LeasedHours, "lease duration should be recorded correctly")
}

// Start Generation Here
func (suite *KeeperTestSuite) TestExtendLeaseMediaNode() {
	suite.SetupTest()
	keeper := suite.App.MedianodeKeeper

	// Register a new media node
	createMsg, err := types.NewMsgRegisterMediaNode(
		defaultMediaNodeURL,
		defaultInfo,
		defaultHardwareSpecs,
		defaultPricePerHour,
		defaultDepositAmount,
		suite.TestAccs[0].String(),
	)
	suite.Require().NoError(err)
	_, err = suite.msgServer.RegisterMediaNode(suite.Ctx, createMsg)
	suite.Require().NoError(err)

	// Lease the media node using a different account (e.g. TestAccs[1])
	leaseHours := uint64(24)
	leaseMsg := types.NewMsgLeaseMediaNode(
		createMsg.Id,
		leaseHours,
		sdk.NewCoin(defaultPricePerHour.Denom, defaultPricePerHour.Amount.Mul(sdkmath.NewIntFromUint64(leaseHours))),
		suite.TestAccs[1].String(),
	)
	_, err = suite.msgServer.LeaseMediaNode(suite.Ctx, leaseMsg)
	suite.Require().NoError(err)

	// Extend the lease by an additional duration (e.g. 12 more hours) using the same leasing account
	additionalHours := uint64(12)
	extendMsg := types.NewMsgExtendLease(
		createMsg.Id,
		additionalHours,
		sdk.NewCoin(defaultPricePerHour.Denom, defaultPricePerHour.Amount.Mul(sdkmath.NewIntFromUint64(additionalHours))),
		suite.TestAccs[1].String(),
	)
	_, err = suite.msgServer.ExtendLease(suite.Ctx, extendMsg)
	suite.Require().NoError(err)

	// Verify that the lease duration has been extended
	lease, found := keeper.GetMediaNodeLease(suite.Ctx, createMsg.Id)
	suite.Require().True(found)
	expectedTotalHours := leaseHours + additionalHours
	suite.Require().Equal(expectedTotalHours, lease.LeasedHours, "lease duration should be extended")
}

func (suite *KeeperTestSuite) TestDepositMediaNode() {
	suite.SetupTest()
	keeper := suite.App.MedianodeKeeper

	// Register a new media node
	initialDepositAmount := sdk.NewCoin(defaultDepositAmount.Denom, sdkmath.NewInt(5_000_000))
	createMsg, err := types.NewMsgRegisterMediaNode(
		defaultMediaNodeURL,
		defaultInfo,
		defaultHardwareSpecs,
		defaultPricePerHour,
		initialDepositAmount,
		suite.TestAccs[0].String(),
	)
	suite.Require().NoError(err)
	_, err = suite.msgServer.RegisterMediaNode(suite.Ctx, createMsg)
	suite.Require().NoError(err)

	// Deposit additional funds into the media node with the owner's account
	additionalDeposit := sdk.NewCoin(defaultDepositAmount.Denom, sdkmath.NewInt(5_000_000))
	depositMsg := types.NewMsgDepositMediaNode(
		createMsg.Id,
		additionalDeposit,
		suite.TestAccs[0].String(),
	)
	_, err = suite.msgServer.DepositMediaNode(suite.Ctx, depositMsg)
	suite.Require().NoError(err)

	// Verify that the deposit has been updated
	node, found := keeper.GetMediaNode(suite.Ctx, createMsg.Id)
	suite.Require().True(found)
	expectedDeposit := sdk.NewCoin(
		defaultDepositAmount.Denom,
		initialDepositAmount.Amount.Add(additionalDeposit.Amount),
	)
	suite.Require().Equal(len(node.Deposits), 1)
	suite.Require().Equal(expectedDeposit.Amount, node.Deposits[0].Amount.Amount, "deposit should be updated after additional deposit")
}

func (suite *KeeperTestSuite) TestCancelLeaseMediaNode() {
	suite.SetupTest()
	keeper := suite.App.MedianodeKeeper

	// Register a new media node
	createMsg, err := types.NewMsgRegisterMediaNode(
		defaultMediaNodeURL,
		defaultInfo,
		defaultHardwareSpecs,
		defaultPricePerHour,
		defaultDepositAmount,
		suite.TestAccs[0].String(),
	)
	suite.Require().NoError(err)
	_, err = suite.msgServer.RegisterMediaNode(suite.Ctx, createMsg)
	suite.Require().NoError(err)

	// Lease the media node using a different account (e.g. TestAccs[1])
	leaseHours := uint64(24)
	leaseMsg := types.NewMsgLeaseMediaNode(
		createMsg.Id,
		leaseHours,
		sdk.NewCoin(defaultPricePerHour.Denom, defaultPricePerHour.Amount.Mul(sdkmath.NewIntFromUint64(leaseHours))),
		suite.TestAccs[1].String(),
	)
	_, err = suite.msgServer.LeaseMediaNode(suite.Ctx, leaseMsg)
	suite.Require().NoError(err)

	// Cancel the lease using the leasing account
	cancelMsg := types.NewMsgCancelLease(
		createMsg.Id,
		suite.TestAccs[1].String(),
	)
	_, err = suite.msgServer.CancelLease(suite.Ctx, cancelMsg)
	suite.Require().NoError(err)

	// Verify that the lease has been cancelled
	node, found := keeper.GetMediaNode(suite.Ctx, createMsg.Id)
	suite.Require().True(found)
	suite.Require().False(node.Leased, "media node lease should be cancelled")

	_, found = keeper.GetMediaNodeLease(suite.Ctx, createMsg.Id)
	suite.Require().False(found, "media node lease record should be removed after cancellation")
}

func (suite *KeeperTestSuite) TestCloseMediaNode() {
	suite.SetupTest()
	keeper := suite.App.MedianodeKeeper

	// Register a new media node
	createMsg, err := types.NewMsgRegisterMediaNode(
		defaultMediaNodeURL,
		defaultInfo,
		defaultHardwareSpecs,
		defaultPricePerHour,
		defaultDepositAmount,
		suite.TestAccs[0].String(),
	)
	suite.Require().NoError(err)
	_, err = suite.msgServer.RegisterMediaNode(suite.Ctx, createMsg)
	suite.Require().NoError(err)

	// Close the media node using the owner's account
	closeMsg := types.NewMsgCloseMediaNode(
		createMsg.Id,
		suite.TestAccs[0].String(),
	)
	_, err = suite.msgServer.CloseMediaNode(suite.Ctx, closeMsg)
	suite.Require().NoError(err)

	// Verify that the media node is no longer retrievable (i.e., it has been closed)
	node, found := keeper.GetMediaNode(suite.Ctx, createMsg.Id)
	suite.Require().True(found)
	suite.Require().Equal(node.Status.String(), types.STATUS_CLOSED.String(), "closed media node status should be CLOSED")
}

func (suite *KeeperTestSuite) TestSettleActiveLeases() {
	suite.SetupTest()
	keeper := suite.App.MedianodeKeeper

	// Register a new media node
	createMsg, err := types.NewMsgRegisterMediaNode(
		defaultMediaNodeURL,
		defaultInfo,
		defaultHardwareSpecs,
		defaultPricePerHour,
		defaultDepositAmount,
		suite.TestAccs[0].String(),
	)
	suite.Require().NoError(err)
	_, err = suite.msgServer.RegisterMediaNode(suite.Ctx, createMsg)
	suite.Require().NoError(err)

	// Manually create an active lease for the media node
	lease := types.Lease{
		MediaNodeId:  createMsg.Id,
		Lessee:       suite.TestAccs[1].String(),
		Owner:        suite.TestAccs[0].String(),
		LeasedHours:  2,
		PricePerHour: defaultPricePerHour,
		// Set the start and last settled time to 1 hours ago to simulate an ongoing lease
		StartTime:          suite.Ctx.BlockTime().Add(-1 * time.Hour),
		TotalLeaseAmount:   sdk.NewCoin(defaultPricePerHour.Denom, defaultPricePerHour.Amount.Mul(sdkmath.NewInt(2))),
		SettledLeaseAmount: sdk.NewCoin(defaultPricePerHour.Denom, sdkmath.NewInt(0)),
		LastSettledAt:      suite.Ctx.BlockTime().Add(-1 * time.Hour),
	}
	// Assume SetMediaNodeLease sets the lease record in the keeper's store.
	keeper.SetLease(suite.Ctx, lease)

	// Call the SettleActiveLeases function (typically invoked in EndBlock)
	err = keeper.SettleActiveLeases(suite.Ctx)
	suite.Require().NoError(err)

	// Retrieve the updated lease and verify that settlement has occurred
	settledLease, found := keeper.GetMediaNodeLease(suite.Ctx, createMsg.Id)
	suite.Require().True(found, "active lease record should exist after settlement")
	suite.Require().Equal(suite.Ctx.BlockTime(), settledLease.LastSettledAt, "lease LastSettledAt should be updated to current block time")
	suite.Require().True(settledLease.SettledLeaseAmount.Amount.Equal(defaultPricePerHour.Amount), "settled lease amount should be equal to pricePerHour")
}

func (suite *KeeperTestSuite) TestReleaseDeposits() {
	suite.SetupTest()
	keeper := suite.App.MedianodeKeeper

	// Register a new media node which includes a deposit
	createMsg, err := types.NewMsgRegisterMediaNode(
		defaultMediaNodeURL,
		defaultInfo,
		defaultHardwareSpecs,
		defaultPricePerHour,
		defaultDepositAmount,
		suite.TestAccs[0].String(),
	)
	suite.Require().NoError(err)
	_, err = suite.msgServer.RegisterMediaNode(suite.Ctx, createMsg)
	suite.Require().NoError(err)

	// Verify that the media node has deposits before release
	node, found := keeper.GetMediaNode(suite.Ctx, createMsg.Id)
	suite.Require().True(found, "media node should exist")
	suite.Require().NotEmpty(node.Deposits, "media node should have deposits before release")

	// Close the media node to make it eligible for deposit release
	closeMsg := types.NewMsgCloseMediaNode(
		createMsg.Id,
		suite.TestAccs[0].String(),
	)
	_, err = suite.msgServer.CloseMediaNode(suite.Ctx, closeMsg)
	suite.Require().NoError(err)

	// Advance block time to simulate deposit lock expiry (if applicable)
	suite.Ctx = suite.Ctx.WithBlockTime(suite.Ctx.BlockTime().Add(7 * 24 * time.Hour))

	// Call the ReleaseDeposits function (typically invoked in EndBlock)
	err = keeper.ReleaseDeposits(suite.Ctx)
	suite.Require().NoError(err)

	// Verify that the deposits have been released (cleared) from the media node record
	node, found = keeper.GetMediaNode(suite.Ctx, createMsg.Id)
	suite.Require().True(found, "media node should still exist")
	suite.Require().Empty(node.Deposits, "media node deposits should be released after calling ReleaseDeposits")
}
