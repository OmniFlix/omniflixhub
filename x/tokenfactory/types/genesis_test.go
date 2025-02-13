package types_test

import (
	"testing"

	"github.com/OmniFlix/omniflixhub/v6/app"

	"github.com/stretchr/testify/require"

	"github.com/OmniFlix/omniflixhub/v6/x/tokenfactory/types"
)

func TestGenesisState_Validate(t *testing.T) {
	app.SetConfig()
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
				FactoryDenoms: []types.GenesisDenom{
					{
						Denom: "factory/omniflix1t7egva48prqmzl59x5ngv4zx0dtrwewcs7ut5e/bitcoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "omniflix1t7egva48prqmzl59x5ngv4zx0dtrwewcs7ut5e",
						},
					},
				},
			},
			valid: true,
		},
		{
			desc: "different admin from creator",
			genState: &types.GenesisState{
				FactoryDenoms: []types.GenesisDenom{
					{
						Denom: "factory/omniflix1t7egva48prqmzl59x5ngv4zx0dtrwewcs7ut5e/bitcoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "omniflix1yca6mlzq9xvnmmylw0asfjhz0q483pw57pj79v",
						},
					},
				},
			},
			valid: true,
		},
		{
			desc: "empty admin",
			genState: &types.GenesisState{
				FactoryDenoms: []types.GenesisDenom{
					{
						Denom: "factory/omniflix1t7egva48prqmzl59x5ngv4zx0dtrwewcs7ut5e/bitcoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "",
						},
					},
				},
			},
			valid: true,
		},
		{
			desc: "no admin",
			genState: &types.GenesisState{
				FactoryDenoms: []types.GenesisDenom{
					{
						Denom: "factory/omniflix1t7egva48prqmzl59x5ngv4zx0dtrwewcs7ut5e/bitcoin",
					},
				},
			},
			valid: true,
		},
		{
			desc: "invalid admin",
			genState: &types.GenesisState{
				FactoryDenoms: []types.GenesisDenom{
					{
						Denom: "factory/omniflix1t7egva48prqmzl59x5ngv4zx0dtrwewcs7ut5e/bitcoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "moose",
						},
					},
				},
			},
			valid: false,
		},
		{
			desc: "multiple denoms",
			genState: &types.GenesisState{
				FactoryDenoms: []types.GenesisDenom{
					{
						Denom: "factory/omniflix1t7egva48prqmzl59x5ngv4zx0dtrwewcs7ut5e/bitcoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "",
						},
					},
					{
						Denom: "factory/omniflix1t7egva48prqmzl59x5ngv4zx0dtrwewcs7ut5e/litecoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "",
						},
					},
				},
			},
			valid: true,
		},
		{
			desc: "duplicate denoms",
			genState: &types.GenesisState{
				FactoryDenoms: []types.GenesisDenom{
					{
						Denom: "factory/omniflix1t7egva48prqmzl59x5ngv4zx0dtrwewcs7ut5e/bitcoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "",
						},
					},
					{
						Denom: "factory/omniflix1t7egva48prqmzl59x5ngv4zx0dtrwewcs7ut5e/bitcoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "",
						},
					},
				},
			},
			valid: false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
