package keeper

import (
	"github.com/OmniFlix/omniflixhub/x/itc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO: create campaign and deduct funds/nfts
func (k Keeper) CreateCampaign(ctx sdk.Context, campaign types.Campaign) error {
	return nil
}

// TODO: cancel campaign and return back funds/nfts
func (k Keeper) CancelCampaign(ctx sdk.Context, campaign types.Campaign) error {
	return nil
}

// TODO: claim tokens from campaign
func (k Keeper) Claim(ctx sdk.Context, campaign types.Campaign, claim types.Claim) error {
	return nil
}
