package types_test

import (
	"testing"
	"time"

	"github.com/OmniFlix/omniflixhub/v2/x/itc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func TestGenesisState_ValidateGenesis(t *testing.T) {
	// generate a private/public key pair and get the respective address
	pk1 := ed25519.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pk1.Address())

	pk2 := ed25519.GenPrivKey().PubKey()
	addr2 := sdk.AccAddress(pk2.Address())

	defaultAmount := sdk.NewInt64Coin(types.DefaultCampaignCreationFee.Denom, 100_000_000)

	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc: "default is valid",
			genState: &types.GenesisState{
				Campaigns:          []types.Campaign{},
				Claims:             []types.Claim{},
				Params:             types.DefaultParams(),
				NextCampaignNumber: uint64(1),
			},
			valid: true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
				Campaigns: []types.Campaign{
					types.NewCampaign(
						1,
						"test",
						"test",
						time.Time{},
						time.Time{},
						addr1.String(),
						"onftdenom0001",
						10,
						types.INTERACTION_TYPE_BURN,
						types.CLAIM_TYPE_FT,
						defaultAmount,
						sdk.NewInt64Coin(types.DefaultCampaignCreationFee.Denom, 1_000_000_000),
						sdk.NewInt64Coin(types.DefaultCampaignCreationFee.Denom, 100_000_000),
						&types.NFTDetails{},
						&types.Distribution{},
					),
				},
				Claims: []types.Claim{
					types.NewClaim(
						1,
						addr2.String(),
						"onft0001",
						types.INTERACTION_TYPE_BURN,
					),
				},
				Params:             types.DefaultParams(),
				NextCampaignNumber: uint64(2),
			},
			valid: true,
		},
	} {
		{
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
}
