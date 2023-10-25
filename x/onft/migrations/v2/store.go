package v2

import (
	"time"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/OmniFlix/omniflixhub/v2/x/onft/types"
)

// MigrateCollections is used to migrate nft data from onft to x/nft
func MigrateCollections(ctx sdk.Context,
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	logger log.Logger,
	k keeper,
) error {
	logger.Info("migrate store data from version 1 to 2")
	startTime := time.Now()

	store := ctx.KVStore(storeKey)
	iterator := sdk.KVStorePrefixIterator(store, KeyDenomID(""))
	defer iterator.Close()

	var (
		totkalDenoms int64
		totalNFTs    int64
	)
	for ; iterator.Valid(); iterator.Next() {
		var denom types.Denom
		cdc.MustUnmarshal(iterator.Value(), &denom)

		// delete unused key
		store.Delete(KeyDenomID(denom.Id))
		store.Delete(KeyDenomSymbol(denom.Name))
		store.Delete(KeyCollection(denom.Id))

		creator, err := sdk.AccAddressFromBech32(denom.Creator)
		if err != nil {
			return err
		}

		if err := k.nftKeeper.SaveDenom(
			ctx,
			denom.Id,
			denom.Symbol,
			denom.Name,
			denom.Schema,
			creator,
			denom.Description,
			denom.PreviewURI,
			denom.Uri,
			denom.UriHash,
			denom.Data,
		); err != nil {
			return err
		}

		totalNFTsInDenom, err := migrateONFT(ctx, k, logger, denom.Id)
		if err != nil {
			return err
		}
		totkalDenoms++
		totalNFTs += totalNFTsInDenom

	}
	logger.Info("migrate store data success",
		"Total Denoms", totkalDenoms,
		"total NFTs", totalNFTs,
		"time taken", time.Since(startTime).String(),
	)
	return nil
}

func migrateONFT(
	ctx sdk.Context,
	k keeper,
	logger log.Logger,
	denomID string,
) (int64, error) {
	var iterator sdk.Iterator
	defer func() {
		if iterator != nil {
			_ = iterator.Close()
		}
	}()

	store := ctx.KVStore(k.storeKey)

	total := int64(0)
	iterator = sdk.KVStorePrefixIterator(store, KeyONFT(denomID, ""))
	for ; iterator.Valid(); iterator.Next() {
		var oNFT types.ONFT
		k.cdc.MustUnmarshal(iterator.Value(), &oNFT)

		owner, err := sdk.AccAddressFromBech32(oNFT.Owner)
		if err != nil {
			return 0, err
		}

		// delete unused key
		store.Delete(KeyONFT(denomID, oNFT.Id))
		store.Delete(KeyOwner(owner, denomID, oNFT.Id))

		if err := k.saveNFT(
			ctx,
			denomID,
			oNFT.Id,
			oNFT.Metadata.Name,
			oNFT.Metadata.Description,
			oNFT.Metadata.MediaURI,
			oNFT.Metadata.PreviewURI,
			oNFT.Metadata.UriHash,
			oNFT.Data,
			oNFT.Extensible,
			oNFT.Transferable,
			oNFT.Nsfw,
			oNFT.CreatedAt,
			oNFT.RoyaltyShare,
			owner,
		); err != nil {
			return 0, err
		}
		total++
	}
	logger.Info("migrate onft collection success", "DenomID", denomID, "TotalNFTs", total)
	return total, nil
}
