package marketplace

import (
	"github.com/OmniFlix/omniflixhub/v2/x/marketplace/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

func EndBlock(ctx sdk.Context, k keeper.Keeper) []abcitypes.ValidatorUpdate {
	log := k.Logger(ctx)
	err := k.UpdateAuctionStatusesAndProcessBids(ctx)
	if err != nil {
		panic(err)
	}
	log.Info("Updated Auctions and Processed bids.. ")
	return []abcitypes.ValidatorUpdate{}
}
