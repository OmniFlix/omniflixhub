package v2

import (
	"reflect"
	"time"
	"unsafe"

	sdkmath "cosmossdk.io/math"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/nft"
	nftkeeper "cosmossdk.io/x/nft/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"

	"github.com/OmniFlix/omniflixhub/v6/x/onft/types"
)

type keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec

	nftKeeper NFTKeeper
}

func (k keeper) saveNFT(
	ctx sdk.Context,
	denomID,
	nftID,
	nftName,
	nftDescription,
	nftMediaURI,
	nftPreviewURI,
	nftURIHash,
	nftData string,
	extensible,
	transferable,
	nsfw bool,
	nftCreatedAt time.Time,
	nftRoyaltyShare sdkmath.LegacyDec,
	recipient sdk.AccAddress,
) error {
	nftMetadata := &types.ONFTMetadata{
		Name:         nftName,
		Description:  nftDescription,
		PreviewURI:   nftPreviewURI,
		Extensible:   extensible,
		Transferable: transferable,
		Nsfw:         nsfw,
		Data:         nftData,
		CreatedAt:    nftCreatedAt,
		RoyaltyShare: nftRoyaltyShare,
		UriHash:      nftURIHash,
	}
	data, err := codectypes.NewAnyWithValue(nftMetadata)
	if err != nil {
		return err
	}

	nftObj := nft.NFT{
		ClassId: denomID,
		Id:      nftID,
		Uri:     nftMediaURI,
		UriHash: nftURIHash,
		Data:    data,
	}
	k.setNFT(ctx, nftObj)
	k.setOwner(ctx, nftObj.ClassId, nftObj.Id, recipient)
	k.incrTotalSupply(ctx, nftObj.ClassId)
	return nil
}

func (k keeper) setNFT(ctx sdk.Context, token nft.NFT) {
	nftStore := k.getNFTStore(ctx, token.ClassId)
	bz := k.cdc.MustMarshal(&token)
	nftStore.Set([]byte(token.Id), bz)
}

func (k keeper) setOwner(ctx sdk.Context, classID, nftID string, owner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(ownerStoreKey(classID, nftID), owner.Bytes())

	ownerStore := k.getClassStoreByOwner(ctx, owner, classID)
	ownerStore.Set([]byte(nftID), nftkeeper.Placeholder)
}

func (k keeper) incrTotalSupply(ctx sdk.Context, classID string) {
	supply := k.GetTotalSupply(ctx, classID) + 1
	k.updateTotalSupply(ctx, classID, supply)
}

// GetTotalSupply returns the number of all nfts under the specified classID
func (k keeper) GetTotalSupply(ctx sdk.Context, classID string) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(classTotalSupply(classID))
	return sdk.BigEndianToUint64(bz)
}

func (k keeper) updateTotalSupply(ctx sdk.Context, classID string, supply uint64) {
	store := ctx.KVStore(k.storeKey)
	supplyKey := classTotalSupply(classID)
	store.Set(supplyKey, sdk.Uint64ToBigEndian(supply))
}

func (k keeper) getClassStoreByOwner(ctx sdk.Context, owner sdk.AccAddress, classID string) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	key := nftOfClassByOwnerStoreKey(owner, classID)
	return prefix.NewStore(store, key)
}

func (k keeper) getNFTStore(ctx sdk.Context, classID string) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, nftStoreKey(classID))
}

// classTotalSupply returns the byte representation of the ClassTotalSupply
func classTotalSupply(classID string) []byte {
	key := make([]byte, len(nftkeeper.ClassTotalSupply)+len(classID))
	copy(key, nftkeeper.ClassTotalSupply)
	copy(key[len(nftkeeper.ClassTotalSupply):], classID)
	return key
}

// nftStoreKey returns the byte representation of the nft
func nftStoreKey(classID string) []byte {
	key := make([]byte, len(nftkeeper.NFTKey)+len(classID)+len(nftkeeper.Delimiter))
	copy(key, nftkeeper.NFTKey)
	copy(key[len(nftkeeper.NFTKey):], classID)
	copy(key[len(nftkeeper.NFTKey)+len(classID):], nftkeeper.Delimiter)
	return key
}

// ownerStoreKey returns the byte representation of the nft owner
// Items are stored with the following key: values
// 0x04<classID><Delimiter(1 Byte)><nftID>
func ownerStoreKey(classID, nftID string) []byte {
	// key is of format:
	classIDBz := UnsafeStrToBytes(classID)
	nftIDBz := UnsafeStrToBytes(nftID)

	key := make([]byte, len(nftkeeper.OwnerKey)+len(classIDBz)+len(nftkeeper.Delimiter)+len(nftIDBz))
	copy(key, nftkeeper.OwnerKey)
	copy(key[len(nftkeeper.OwnerKey):], classIDBz)
	copy(key[len(nftkeeper.OwnerKey)+len(classIDBz):], nftkeeper.Delimiter)
	copy(key[len(nftkeeper.OwnerKey)+len(classIDBz)+len(nftkeeper.Delimiter):], nftIDBz)
	return key
}

// nftOfClassByOwnerStoreKey returns the byte representation of the nft owner
// Items are stored with the following key: values
// 0x03<owner><Delimiter(1 Byte)><classID><Delimiter(1 Byte)>
func nftOfClassByOwnerStoreKey(owner sdk.AccAddress, classID string) []byte {
	owner = address.MustLengthPrefix(owner)
	classIDBz := UnsafeStrToBytes(classID)

	key := make([]byte, len(nftkeeper.NFTOfClassByOwnerKey)+len(owner)+len(nftkeeper.Delimiter)+len(classIDBz)+len(nftkeeper.Delimiter))
	copy(key, nftkeeper.NFTOfClassByOwnerKey)
	copy(key[len(nftkeeper.NFTOfClassByOwnerKey):], owner)
	copy(key[len(nftkeeper.NFTOfClassByOwnerKey)+len(owner):], nftkeeper.Delimiter)
	copy(key[len(nftkeeper.NFTOfClassByOwnerKey)+len(owner)+len(nftkeeper.Delimiter):], classIDBz)
	copy(key[len(nftkeeper.NFTOfClassByOwnerKey)+len(owner)+len(nftkeeper.Delimiter)+len(classIDBz):], nftkeeper.Delimiter)
	return key
}

// UnsafeStrToBytes uses unsafe to convert string into byte array. Returned bytes
// must not be altered after this function is called as it will cause a segmentation fault.
func UnsafeStrToBytes(s string) []byte {
	var buf []byte
	sHdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bufHdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	bufHdr.Data = sHdr.Data
	bufHdr.Cap = sHdr.Len
	bufHdr.Len = sHdr.Len
	return buf
}

// UnsafeBytesToStr is meant to make a zero allocation conversion
// from []byte -> string to speed up operations, it is not meant
// to be used generally, but for a specific pattern to delete keys
// from a map.
func UnsafeBytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
