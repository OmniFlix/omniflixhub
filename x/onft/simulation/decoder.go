package simulation

import (
	"bytes"
	"fmt"

	"github.com/OmniFlix/omniflixhub/v3/x/onft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"
)

func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.PrefixONFT):
			var nftA, nftB types.ONFT
			cdc.MustUnmarshal(kvA.Value, &nftA)
			cdc.MustUnmarshal(kvB.Value, &nftB)
			return fmt.Sprintf("%v\n%v", nftA, nftB)
		case bytes.Equal(kvA.Key[:1], types.PrefixOwners):
			idA := types.MustUnMarshalONFTID(cdc, kvA.Value)
			idB := types.MustUnMarshalDenomID(cdc, kvB.Value)
			return fmt.Sprintf("%v\n%v", idA, idB)
		case bytes.Equal(kvA.Key[:1], types.PrefixCollection):
			supplyA := types.MustUnMarshalSupply(cdc, kvA.Value)
			supplyB := types.MustUnMarshalSupply(cdc, kvB.Value)
			return fmt.Sprintf("%d\n%d", supplyA, supplyB)
		case bytes.Equal(kvA.Key[:1], types.PrefixDenom):
			var denomA, denomB types.Denom
			cdc.MustUnmarshal(kvA.Value, &denomA)
			cdc.MustUnmarshal(kvB.Value, &denomB)
			return fmt.Sprintf("%v\n%v", denomA, denomB)

		default:
			panic(fmt.Sprintf("invalid %s key prefix %X", types.ModuleName, kvA.Key[:1]))
		}
	}
}
