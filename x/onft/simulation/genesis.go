package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/OmniFlix/omniflixhub/v5/x/onft/types"
)

const (
	denomName1 = "denom1"
	denomName2 = "denom2"
	denomId1   = "onftdenom1"
	denomId2   = "onftdenom2"
)

// RandomizedGenState generates a random GenesisState for nft
func RandomizedGenState(simState *module.SimulationState) {
	collections := types.NewCollections(
		types.NewCollection(
			types.Denom{
				Id:          denomId1,
				Name:        denomName1,
				Schema:      "{}",
				Creator:     "",
				Symbol:      "denom1",
				Description: simtypes.RandStringOfLength(simState.Rand, 45),
				PreviewURI:  simtypes.RandStringOfLength(simState.Rand, 45),
			},
			types.ONFTs{},
		),
		types.NewCollection(
			types.Denom{
				Id:          denomId2,
				Name:        denomName2,
				Schema:      "{}",
				Creator:     "",
				Symbol:      "denom2",
				Description: simtypes.RandStringOfLength(simState.Rand, 45),
				PreviewURI:  simtypes.RandStringOfLength(simState.Rand, 45),
			},
			types.ONFTs{}),
	)
	for i, acc := range simState.Accounts {
		if simState.Rand.Intn(100) < 10 {
			oNFT := types.NewONFT(
				RandID(simState.Rand, "onft", 10),
				RandMetadata(simState.Rand),
				"{}",
				genRandomBool(simState.Rand),
				genRandomBool(simState.Rand),
				acc.Address,
				time.Time{},
				genRandomBool(simState.Rand),
				RandRoyaltyShare(simState.Rand),
			)

			if i < 50 {
				collections[0].Denom.Creator = oNFT.Owner
				collections[0] = collections[0].AddONFT(oNFT)
			} else {
				collections[1].Denom.Creator = oNFT.Owner
				collections[1] = collections[1].AddONFT(oNFT)
			}
		}
	}

	nftGenesis := types.NewGenesisState(collections, types.DefaultParams())

	bz, err := json.MarshalIndent(nftGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated %s parameters:\n%s\n", types.ModuleName, bz)

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(nftGenesis)
}

func RandID(r *rand.Rand, prefix string, n int) string {
	id := simtypes.RandStringOfLength(r, n)
	return strings.ToLower(prefix + id)
}

func RandMetadata(r *rand.Rand) types.Metadata {
	return types.Metadata{
		Name:        simtypes.RandStringOfLength(r, 10),
		Description: simtypes.RandStringOfLength(r, 45),
		PreviewURI:  simtypes.RandStringOfLength(r, 45),
		MediaURI:    simtypes.RandStringOfLength(r, 45),
	}
}

func RandRoyaltyShare(r *rand.Rand) sdk.Dec {
	return simtypes.RandomDecAmount(r, sdk.NewDecWithPrec(999999999999999999, 18))
}
