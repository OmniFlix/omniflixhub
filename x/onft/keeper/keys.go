package keeper

import (
	nftkeeper "cosmossdk.io/x/nft/keeper"
)

func classStoreKey(classID string) []byte {
	key := make([]byte, len(nftkeeper.ClassKey)+len(classID))
	copy(key, nftkeeper.ClassKey)
	copy(key[len(nftkeeper.ClassKey):], classID)
	return key
}
