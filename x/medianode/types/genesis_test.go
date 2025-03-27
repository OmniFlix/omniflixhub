package types_test

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/OmniFlix/omniflixhub/v6/x/medianode/types"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_ValidateGenesis(t *testing.T) {
	// generate a private/public key pair and get the respective address
	pk1 := ed25519.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pk1.Address())

	pk2 := ed25519.GenPrivKey().PubKey()
	addr2 := sdk.AccAddress(pk2.Address())

	defaultDeposit := sdk.NewCoin("uflix", sdkmath.NewInt(10000000))
	defaultPricePerHour := sdk.NewCoin("uflix", sdkmath.NewInt(1000000))

	validMediaNode := types.MediaNode{
		Id:     "mn01953d3f737572d29082215cdaa00001",
		Owner:  addr1.String(),
		Status: types.STATUS_ACTIVE,
		Leased: false,
		Url:    "https://test.medianode.com",
		Info: types.Info{
			Moniker:     "Test Node",
			Description: "Test Description",
			Contact:     "test@medianode.com",
		},
		HardwareSpecs: types.HardwareSpecs{
			Cpus:        4,
			RamInGb:     16,
			StorageInGb: 1000,
		},
		PricePerHour: defaultPricePerHour,
		Deposits: []*types.Deposit{{
			Amount:      defaultDeposit,
			Depositor:   addr1.String(),
			DepositedAt: time.Now(),
		}},
	}

	validLease := types.Lease{
		MediaNodeId:        "mn01953d3f737572d29082215cdaa11006",
		Lessee:             addr2.String(),
		Owner:              addr1.String(),
		LeasedHours:        24,
		StartTime:          time.Now(),
		TotalLeaseAmount:   sdk.NewCoin("uflix", sdkmath.NewInt(240000)),
		SettledLeaseAmount: sdk.NewCoin("uflix", sdkmath.NewInt(0)),
		LastSettledAt:      time.Now(),
	}

	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
				Nodes:       []types.MediaNode{validMediaNode},
				Leases:      []types.Lease{validLease},
				Params:      types.DefaultParams(),
				NodeCounter: uint64(1),
			},
			valid: true,
		},
		{
			desc: "invalid media node",
			genState: &types.GenesisState{
				Nodes: []types.MediaNode{
					{
						Id:     "", // Invalid: empty ID
						Owner:  addr1.String(),
						Status: types.STATUS_ACTIVE,
						Leased: false,
						Url:    "invalid-url",
						Info:   types.Info{},
						HardwareSpecs: types.HardwareSpecs{
							Cpus:        0, // Invalid: zero CPUs
							RamInGb:     0,
							StorageInGb: 0,
						},
					},
				},
				Leases:      []types.Lease{},
				Params:      types.DefaultParams(),
				NodeCounter: uint64(1),
			},
			valid: false,
		},
		{
			desc: "invalid lease",
			genState: &types.GenesisState{
				Nodes: []types.MediaNode{validMediaNode},
				Leases: []types.Lease{
					{
						MediaNodeId:  "", // Invalid: empty node ID
						Lessee:       "invalid-address",
						Owner:        addr1.String(),
						PricePerHour: sdk.NewCoin("uflix", sdkmath.NewInt(0)),
						LeasedHours:  0, // Invalid: zero hours
						StartTime:    time.Time{},

						LastSettledAt: time.Time{},
					},
				},
				Params:      types.DefaultParams(),
				NodeCounter: uint64(1),
			},
			valid: false,
		},
		{
			desc: "invalid params",
			genState: &types.GenesisState{
				Nodes:  []types.MediaNode{validMediaNode},
				Leases: []types.Lease{validLease},
				Params: types.Params{
					MinimumLeaseHours:    0,
					MaximumLeaseHours:    50,
					DepositReleasePeriod: time.Hour * 12,
				},
				NodeCounter: uint64(1),
			},
			valid: false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.ValidateGenesis()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestDefaultGenesis(t *testing.T) {
	defaultGenesis := types.DefaultGenesis()
	require.NotNil(t, defaultGenesis)
	require.Equal(t, uint64(0), defaultGenesis.NodeCounter)
	require.Empty(t, defaultGenesis.Nodes)
	require.Empty(t, defaultGenesis.Leases)
	require.Equal(t, types.DefaultParams(), defaultGenesis.Params)
}

func TestNewGenesisState(t *testing.T) {
	nodes := []types.MediaNode{}
	leases := []types.Lease{}
	nodeCounter := uint64(0)
	params := types.DefaultParams()

	genesis := types.NewGenesisState(nodes, leases, nodeCounter, params)
	require.NotNil(t, genesis)
	require.Equal(t, nodes, genesis.Nodes)
	require.Equal(t, leases, genesis.Leases)
	require.Equal(t, nodeCounter, genesis.NodeCounter)
	require.Equal(t, params, genesis.Params)
}
