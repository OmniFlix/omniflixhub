package v2_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/OmniFlix/omniflixhub/v5/app/apptesting"
	"github.com/OmniFlix/omniflixhub/v5/x/onft/keeper"
	v2 "github.com/OmniFlix/omniflixhub/v5/x/onft/migrations/v2"
	"github.com/OmniFlix/omniflixhub/v5/x/onft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/cosmos/gogoproto/types"
)

func generateCollectionsData(
	ctx sdk.Context,
	storeKey storetypes.StoreKey,
	cdc codec.Codec,
) (collections []types.Collection) {
	addrs := apptesting.CreateRandomAccounts(10)
	randomIpfsPaths := []string{"ipfs://testhash", "https://ipfs.some.url/ipfs/testhash", "http://url/ipfs/", "https://ipfs.some.url/ipfs/testhash?1.jpg"}
	for i := 1; i <= 10; i++ {
		denom := types.Denom{
			Id:          fmt.Sprintf("onftdenom%d", i),
			Name:        fmt.Sprintf("Test Denom%d", i),
			Schema:      fmt.Sprintf("{\"id\":%d}", i),
			Creator:     addrs[rand.Intn(len(addrs))].String(),
			PreviewURI:  randomIpfsPaths[rand.Intn(len(randomIpfsPaths))],
			Symbol:      fmt.Sprintf("Testd%d", i),
			Description: fmt.Sprintf("Test Description %d", i),
			Uri:         fmt.Sprintf("ipfs://testuri%d", i),
			UriHash:     fmt.Sprintf("urihash%d", i),
			Data:        fmt.Sprintf("data%d", i),
		}
		setDenom(ctx, storeKey, cdc, denom)

		var nfts []types.ONFT
		for j := 1; j <= 100; j++ {
			metadata := types.Metadata{
				Name:        fmt.Sprintf("Test NFT #%d", j),
				Description: fmt.Sprintf("Test description %d", j),
				MediaURI:    randomIpfsPaths[rand.Intn(len(randomIpfsPaths))],
				PreviewURI:  randomIpfsPaths[rand.Intn(len(randomIpfsPaths))],
			}
			_nft := types.ONFT{
				Id:           fmt.Sprintf("onft%d", j),
				Metadata:     metadata,
				Transferable: []bool{true, false}[rand.Intn(2)],
				Extensible:   []bool{true, false}[rand.Intn(2)],
				Nsfw:         []bool{true, false}[rand.Intn(2)],
				RoyaltyShare: sdk.NewDecWithPrec(5, 2),
				Data:         fmt.Sprintf("nftData%d", j),
				Owner:        addrs[rand.Intn(len(addrs))].String(),
				CreatedAt:    time.Time{},
			}
			nfts = append(nfts, _nft)
			mintONFT(ctx, storeKey, cdc, denom.Id, _nft)
		}
		collections = append(collections, types.Collection{
			Denom: denom,
			ONFTs: nfts,
		})
	}
	return collections
}

func check(t *testing.T, ctx sdk.Context, k keeper.Keeper, collections []types.Collection) {
	t.Helper()

	for _, collection := range collections {
		denom := collection.Denom
		count := k.GetTotalSupply(ctx, denom.Id)
		require.Equal(t, count, uint64(len(collection.ONFTs)))
	}
	keeper.SupplyInvariant(k)
}

func setDenom(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.Codec, denom types.Denom) {
	store := ctx.KVStore(storeKey)
	bz := cdc.MustMarshal(&denom)
	store.Set(v2.KeyDenomID(denom.Id), bz)
	store.Set(v2.KeyDenomSymbol(denom.Name), []byte(denom.Id))
}

func mintONFT(
	ctx sdk.Context,
	storeKey storetypes.StoreKey,
	cdc codec.Codec,
	denomID string,
	onft types.ONFT,
) {
	setONFT(ctx, storeKey, cdc, denomID, onft)
	setOwner(ctx, storeKey, cdc, denomID, onft.Id, onft.Owner)
	increaseSupply(ctx, storeKey, cdc, denomID)
}

func setONFT(
	ctx sdk.Context,
	storeKey storetypes.StoreKey,
	cdc codec.Codec,
	denomID string,
	onft types.ONFT,
) {
	store := ctx.KVStore(storeKey)

	bz := cdc.MustMarshal(&onft)
	store.Set(v2.KeyONFT(denomID, onft.Id), bz)
}

func setOwner(
	ctx sdk.Context,
	storeKey storetypes.StoreKey,
	cdc codec.Codec,
	denomID, onftID, owner string,
) {
	store := ctx.KVStore(storeKey)
	bz := mustMarshalOnftID(cdc, onftID)
	ownerAddr := sdk.MustAccAddressFromBech32(owner)
	store.Set(v2.KeyOwner(ownerAddr, denomID, onftID), bz)
}

func increaseSupply(
	ctx sdk.Context,
	storeKey storetypes.StoreKey,
	cdc codec.Codec,
	denomID string,
) {
	supply := getTotalSupply(ctx, storeKey, cdc, denomID)
	supply++

	store := ctx.KVStore(storeKey)
	bz := mustMarshalSupply(cdc, supply)
	store.Set(v2.KeyCollection(denomID), bz)
}

func getTotalSupply(
	ctx sdk.Context,
	storeKey storetypes.StoreKey,
	cdc codec.Codec,
	denomID string,
) uint64 {
	store := ctx.KVStore(storeKey)
	bz := store.Get(v2.KeyCollection(denomID))
	if len(bz) == 0 {
		return 0
	}
	return mustUnMarshalSupply(cdc, bz)
}

func mustMarshalSupply(cdc codec.Codec, supply uint64) []byte {
	supplyWrap := gogotypes.UInt64Value{Value: supply}
	return cdc.MustMarshal(&supplyWrap)
}

func mustUnMarshalSupply(cdc codec.Codec, value []byte) uint64 {
	var supplyWrap gogotypes.UInt64Value
	cdc.MustUnmarshal(value, &supplyWrap)
	return supplyWrap.Value
}

func mustMarshalOnftID(cdc codec.Codec, onftID string) []byte {
	onftIDWrap := gogotypes.StringValue{Value: onftID}
	return cdc.MustMarshal(&onftIDWrap)
}
