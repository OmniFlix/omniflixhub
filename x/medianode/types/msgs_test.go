package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/OmniFlix/omniflixhub/v6/x/medianode/types"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var (
	defaultNodeName        = "test node"
	defaultNodeDescription = "test description"
	defaultContact         = "test@medianode.com"
	defaultNodeURL         = "https://test.medianode.com"
	defaultPricePerHour    = sdk.NewCoin("uflix", sdkmath.NewInt(10000))
	defaultDeposit         = sdk.NewCoin("uflix", sdkmath.NewInt(1000000))
	defaultHardwareSpecs   = types.HardwareSpecs{
		Cpus:        4,
		RamInGb:     16,
		StorageInGb: 1000,
	}
	defaultInfo = types.Info{
		Moniker:     defaultNodeName,
		Description: defaultNodeDescription,
		Contact:     defaultContact,
	}
)

func TestNewMsgRegisterMediaNode(t *testing.T) {
	pk1 := ed25519.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pk1.Address())

	msg, err := types.NewMsgRegisterMediaNode(
		defaultNodeURL,
		defaultInfo,
		defaultHardwareSpecs,
		defaultPricePerHour,
		defaultDeposit,
		addr1.String(),
	)
	require.NoError(t, err)

	require.Equal(t, msg.Route(), types.MsgRoute)
	require.Equal(t, msg.Type(), types.TypeMsgRegisterMediaNode)
	signers := msg.GetSigners()
	require.Equal(t, len(signers), 1)
	require.Equal(t, signers[0].String(), addr1.String())

	testCases := []struct {
		desc  string
		msg   types.MsgRegisterMediaNode
		valid bool
	}{
		{
			desc:  "valid message",
			msg:   *msg,
			valid: true,
		},
		{
			desc: "invalid sender",
			msg: types.MsgRegisterMediaNode{
				Sender: "invalid-address",
				Url:    defaultNodeURL,
				Info:   defaultInfo,
			},
			valid: false,
		},
		{
			desc: "empty URL",
			msg: types.MsgRegisterMediaNode{
				Sender: addr1.String(),
				Url:    "",
				Info:   defaultInfo,
			},
			valid: false,
		},
		{
			desc: "invalid price per hour",
			msg: types.MsgRegisterMediaNode{
				Sender:       addr1.String(),
				Url:          defaultNodeURL,
				Info:         defaultInfo,
				PricePerHour: sdk.NewCoin("uflix", sdkmath.NewInt(0)),
			},
			valid: false,
		},
	}

	for _, tc := range testCases {
		if tc.valid {
			require.NoError(t, tc.msg.ValidateBasic(), "test: %v", tc.desc)
		} else {
			require.Error(t, tc.msg.ValidateBasic(), "test: %v", tc.desc)
		}
	}
}

func TestNewMsgUpdateMediaNode(t *testing.T) {
	pk1 := ed25519.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pk1.Address())

	updatedInfo := types.Info{
		Moniker:     "updated name",
		Description: "updated description",
		Contact:     "updated@medianode.com",
	}
	updatedPrice := sdk.NewCoin("uflix", sdkmath.NewInt(20000))

	msg := types.NewMsgUpdateMediaNode(
		"testnode1",
		&updatedInfo,
		&defaultHardwareSpecs,
		&updatedPrice,
		addr1.String(),
	)

	require.Equal(t, msg.Route(), types.MsgRoute)
	require.Equal(t, msg.Type(), types.TypeMsgUpdateMediaNode)

	testCases := []struct {
		desc  string
		msg   types.MsgUpdateMediaNode
		valid bool
	}{
		{
			desc:  "valid message",
			msg:   *msg,
			valid: true,
		},
		{
			desc: "invalid sender",
			msg: types.MsgUpdateMediaNode{
				Sender: "invalid-address",
				Id:     "testnode1",
			},
			valid: false,
		},
		{
			desc: "no updates",
			msg: types.MsgUpdateMediaNode{
				Sender: addr1.String(),
				Id:     "testnode1",
			},
			valid: false,
		},
	}

	for _, tc := range testCases {
		if tc.valid {
			require.NoError(t, tc.msg.ValidateBasic(), "test: %v", tc.desc)
		} else {
			require.Error(t, tc.msg.ValidateBasic(), "test: %v", tc.desc)
		}
	}
}

func TestNewMsgLeaseMediaNode(t *testing.T) {
	pk1 := ed25519.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pk1.Address())

	msg := types.NewMsgLeaseMediaNode(
		"testnode1",
		24, // 24 hours
		sdk.NewCoin("uflix", sdkmath.NewInt(240000)),
		addr1.String(),
	)

	require.Equal(t, msg.Route(), types.MsgRoute)
	require.Equal(t, msg.Type(), types.TypeMsgLeaseMediaNode)

	testCases := []struct {
		desc  string
		msg   types.MsgLeaseMediaNode
		valid bool
	}{
		{
			desc:  "valid message",
			msg:   *msg,
			valid: true,
		},
		{
			desc: "invalid sender",
			msg: types.MsgLeaseMediaNode{
				Sender:      "invalid-address",
				MediaNodeId: "testnode1",
				LeaseHours:  24,
			},
			valid: false,
		},
		{
			desc: "zero lease hours",
			msg: types.MsgLeaseMediaNode{
				Sender:      addr1.String(),
				MediaNodeId: "testnode1",
				LeaseHours:  0,
			},
			valid: false,
		},
		{
			desc: "zero amount",
			msg: types.MsgLeaseMediaNode{
				Sender:      addr1.String(),
				MediaNodeId: "testnode1",
				LeaseHours:  24,
				Amount:      sdk.NewCoin("uflix", sdkmath.NewInt(0)),
			},
			valid: false,
		},
	}

	for _, tc := range testCases {
		if tc.valid {
			require.NoError(t, tc.msg.ValidateBasic(), "test: %v", tc.desc)
		} else {
			require.Error(t, tc.msg.ValidateBasic(), "test: %v", tc.desc)
		}
	}
}

func TestNewMsgExtendLease(t *testing.T) {
	pk1 := ed25519.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pk1.Address())

	msg := types.NewMsgExtendLease(
		"testnode1",
		12, // 12 additional hours
		sdk.NewCoin("uflix", sdkmath.NewInt(120000)),
		addr1.String(),
	)

	require.Equal(t, msg.Route(), types.MsgRoute)
	require.Equal(t, msg.Type(), types.TypeMsgExtendLease)

	testCases := []struct {
		desc  string
		msg   types.MsgExtendLease
		valid bool
	}{
		{
			desc:  "valid message",
			msg:   *msg,
			valid: true,
		},
		{
			desc: "invalid sender",
			msg: types.MsgExtendLease{
				Sender:      "invalid-address",
				MediaNodeId: "testnode1",
			},
			valid: false,
		},
		{
			desc: "zero lease hours",
			msg: types.MsgExtendLease{
				Sender:      addr1.String(),
				MediaNodeId: "testnode1",
				LeaseHours:  0,
			},
			valid: false,
		},
	}

	for _, tc := range testCases {
		if tc.valid {
			require.NoError(t, tc.msg.ValidateBasic(), "test: %v", tc.desc)
		} else {
			require.Error(t, tc.msg.ValidateBasic(), "test: %v", tc.desc)
		}
	}
}

func TestNewMsgCancelLease(t *testing.T) {
	pk1 := ed25519.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pk1.Address())

	msg := types.NewMsgCancelLease("testnode1", addr1.String())

	require.Equal(t, msg.Route(), types.MsgRoute)
	require.Equal(t, msg.Type(), types.TypeMsgCancelLease)

	testCases := []struct {
		desc  string
		msg   types.MsgCancelLease
		valid bool
	}{
		{
			desc:  "valid message",
			msg:   *msg,
			valid: true,
		},
		{
			desc: "invalid sender",
			msg: types.MsgCancelLease{
				Sender:      "invalid-address",
				MediaNodeId: "testnode1",
			},
			valid: false,
		},
	}

	for _, tc := range testCases {
		if tc.valid {
			require.NoError(t, tc.msg.ValidateBasic(), "test: %v", tc.desc)
		} else {
			require.Error(t, tc.msg.ValidateBasic(), "test: %v", tc.desc)
		}
	}
}

func TestNewMsgDepositMediaNode(t *testing.T) {
	pk1 := ed25519.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pk1.Address())

	msg := types.NewMsgDepositMediaNode(
		"testnode1",
		sdk.NewCoin("uflix", sdkmath.NewInt(1000000)),
		addr1.String(),
	)

	require.Equal(t, msg.Route(), types.MsgRoute)
	require.Equal(t, msg.Type(), types.TypeMsgDepositMediaNode)

	testCases := []struct {
		desc  string
		msg   types.MsgDepositMediaNode
		valid bool
	}{
		{
			desc:  "valid message",
			msg:   *msg,
			valid: true,
		},
		{
			desc: "invalid sender",
			msg: types.MsgDepositMediaNode{
				Sender:      "invalid-address",
				MediaNodeId: "testnode1",
			},
			valid: false,
		},
		{
			desc: "zero amount",
			msg: types.MsgDepositMediaNode{
				Sender:      addr1.String(),
				MediaNodeId: "testnode1",
				Amount:      sdk.NewCoin("uflix", sdkmath.NewInt(0)),
			},
			valid: false,
		},
	}

	for _, tc := range testCases {
		if tc.valid {
			require.NoError(t, tc.msg.ValidateBasic(), "test: %v", tc.desc)
		} else {
			require.Error(t, tc.msg.ValidateBasic(), "test: %v", tc.desc)
		}
	}
}

func TestNewMsgCloseMediaNode(t *testing.T) {
	pk1 := ed25519.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pk1.Address())

	msg := types.NewMsgCloseMediaNode("testnode1", addr1.String())

	require.Equal(t, msg.Route(), types.MsgRoute)
	require.Equal(t, msg.Type(), types.TypeMsgCloseMediaNode)

	testCases := []struct {
		desc  string
		msg   types.MsgCloseMediaNode
		valid bool
	}{
		{
			desc:  "valid message",
			msg:   *msg,
			valid: true,
		},
		{
			desc: "invalid sender",
			msg: types.MsgCloseMediaNode{
				Sender:      "invalid-address",
				MediaNodeId: "testnode1",
			},
			valid: false,
		},
	}

	for _, tc := range testCases {
		if tc.valid {
			require.NoError(t, tc.msg.ValidateBasic(), "test: %v", tc.desc)
		} else {
			require.Error(t, tc.msg.ValidateBasic(), "test: %v", tc.desc)
		}
	}
}
