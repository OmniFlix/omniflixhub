package marketplace

import (
	"fmt"

	"github.com/OmniFlix/omniflixhub/v6/x/marketplace/keeper"
	"github.com/OmniFlix/omniflixhub/v6/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	if err := genState.ValidateGenesis(); err != nil {
		panic(err.Error())
	}
	for _, l := range genState.Listings {
		k.SetListing(ctx, l)
		k.SetWithOwner(ctx, l.GetOwner(), l.GetId())
		k.SetWithNFTID(ctx, l.GetNftId(), l.GetId())
		k.SetWithPriceDenom(ctx, l.Price.GetDenom(), l.GetId())
	}
	k.SetListingCount(ctx, genState.ListingCount)
	err := k.SetParams(ctx, genState.Params)
	if err != nil {
		panic(err)
	}

	for _, al := range genState.Auctions {
		k.SetAuctionListing(ctx, al)
		k.SetAuctionListingWithOwner(ctx, al.GetOwner(), al.GetId())
		k.SetAuctionListingWithNFTID(ctx, al.GetNftId(), al.GetId())
		k.SetAuctionListingWithPriceDenom(ctx, al.StartPrice.GetDenom(), al.GetId())
	}

	for _, b := range genState.Bids {
		k.SetBid(ctx, b)
	}
	k.SetNextAuctionNumber(ctx, genState.NextAuctionNumber)

	// check if the module account exists
	moduleAcc := k.GetMarketplaceAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return types.NewGenesisState(
		k.GetAllListings(ctx),
		k.GetListingCount(ctx),
		k.GetParams(ctx),
		k.GetAllAuctionListings(ctx),
		k.GetAllBids(ctx),
		k.GetNextAuctionNumber(ctx),
	)
}

func DefaultGenesisState() *types.GenesisState {
	return types.NewGenesisState([]types.Listing{}, 0, types.DefaultParams(), []types.AuctionListing{}, []types.Bid{}, 1)
}
