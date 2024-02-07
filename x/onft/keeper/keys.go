package keeper

import (
	nftkeeper "github.com/cosmos/cosmos-sdk/x/nft/keeper"
)

func classStoreKey(classID string) []byte {
	key := make([]byte, len(nftkeeper.ClassKey)+len(classID))
	copy(key, nftkeeper.ClassKey)
	copy(key[len(nftkeeper.ClassKey):], classID)
	return key
}
