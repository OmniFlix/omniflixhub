package keeper

import (
	"fmt"

	onfttypes "github.com/OmniFlix/omniflixhub/v2/x/onft/types"

	errorsmod "cosmossdk.io/errors"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/OmniFlix/omniflixhub/v2/x/marketplace/types"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec

	accountKeeper      types.AccountKeeper
	bankKeeper         types.BankKeeper
	nftKeeper          types.NftKeeper
	distributionKeeper types.DistributionKeeper

	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string
}

func NewKeeper(
	cdc codec.BinaryCodec,
	key storetypes.StoreKey,

	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	nftKeeper types.NftKeeper,
	distrKeeper types.DistributionKeeper,
	authority string,
) Keeper {
	// ensure marketplace module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		storeKey:           key,
		cdc:                cdc,
		accountKeeper:      accountKeeper,
		bankKeeper:         bankKeeper,
		nftKeeper:          nftKeeper,
		distributionKeeper: distrKeeper,
		authority:          authority,
	}
}

// GetAuthority returns the x/marketplace module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("OmniFlix/%s", types.ModuleName))
}

// AddListing adds a listing in the store and set owner to listing and updates the count
func (k Keeper) AddListing(ctx sdk.Context, listing types.Listing) error {
	// check listing already exists
	if k.HasListing(ctx, listing.GetId()) {
		return errorsmod.Wrapf(types.ErrListingAlreadyExists, "listing already exists: %s", listing.GetId())
	}

	err := k.nftKeeper.TransferOwnership(ctx,
		listing.GetDenomId(), listing.GetNftId(), listing.GetOwner(),
		k.accountKeeper.GetModuleAddress(types.ModuleName))
	if err != nil {
		return err
	}
	// set listing
	k.SetListing(ctx, listing)
	if len(listing.GetOwner()) != 0 {
		// set listing id with owner prefix
		k.SetWithOwner(ctx, listing.GetOwner(), listing.GetId())
	}
	// Update listing count
	count := k.GetListingCount(ctx)
	k.SetListingCount(ctx, count+1)
	k.SetWithNFTID(ctx, listing.NftId, listing.Id)

	if len(listing.Price.Denom) > 0 {
		k.SetWithPriceDenom(ctx, listing.Price.Denom, listing.Id)
	}
	return nil
}

func (k Keeper) DeleteListing(ctx sdk.Context, listing types.Listing) {
	k.RemoveListing(ctx, listing.GetId())
	k.UnsetWithOwner(ctx, listing.GetOwner(), listing.GetId())
	k.UnsetWithNFTID(ctx, listing.GetNftId())
	k.UnsetWithPriceDenom(ctx, listing.Price.Denom, listing.GetId())
}

func (k Keeper) Buy(ctx sdk.Context, listing types.Listing, buyer sdk.AccAddress) error {
	owner, err := sdk.AccAddressFromBech32(listing.Owner)
	if err != nil {
		return err
	}
	denom, err := k.nftKeeper.GetDenomInfo(ctx, listing.DenomId)
	if err != nil {
		return err
	}
	listingPriceCoin := listing.Price
	listingSaleAmountCoin := listingPriceCoin
	nft, err := k.nftKeeper.GetONFT(ctx, listing.DenomId, listing.NftId)
	if err != nil {
		return err
	}

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, buyer, types.ModuleName, sdk.NewCoins(listingPriceCoin))
	if err != nil {
		return err
	}
	err = k.nftKeeper.TransferOwnership(ctx, listing.GetDenomId(), listing.GetNftId(),
		k.accountKeeper.GetModuleAddress(types.ModuleName), buyer)
	if err != nil {
		return err
	}
	saleCommission := k.GetSaleCommission(ctx)
	marketplaceCoin := k.GetProportions(listing.Price, saleCommission)
	if marketplaceCoin.Amount.GTE(sdk.OneInt()) {
		err = k.DistributeCommission(ctx, marketplaceCoin)
		if err != nil {
			return err
		}
		listingSaleAmountCoin = listingPriceCoin.Sub(marketplaceCoin)
	}
	// check if it is a valid royalty share
	if nft.GetRoyaltyShare().GT(sdk.ZeroDec()) && nft.GetRoyaltyShare().LTE(sdk.OneDec()) {
		nftRoyaltyShareCoin := k.GetProportions(listingSaleAmountCoin, nft.GetRoyaltyShare())
		creator, err := sdk.AccAddressFromBech32(denom.Creator)
		if err != nil {
			return err
		}
		if err := k.TransferRoyalty(ctx, nftRoyaltyShareCoin, denom.RoyaltyReceivers, creator); err != nil {
			return err
		}
		listingSaleAmountCoin = listingSaleAmountCoin.Sub(nftRoyaltyShareCoin)
	}
	remaining := listingSaleAmountCoin

	if len(listing.SplitShares) > 0 {
		for _, share := range listing.SplitShares {
			sharePortionCoin := k.GetProportions(listingSaleAmountCoin, share.Weight)
			sharePortionCoins := sdk.NewCoins(sharePortionCoin)
			if share.Address == "" {
				err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, owner, sharePortionCoins)
				if err != nil {
					return err
				}
			} else {
				saleSplitAddr, err := sdk.AccAddressFromBech32(share.Address)
				if err != nil {
					return err
				}
				err = k.bankKeeper.SendCoinsFromModuleToAccount(
					ctx, types.ModuleName, saleSplitAddr, sharePortionCoins)
				if err != nil {
					return err
				}
				k.createSplitShareTransferEvent(ctx, k.accountKeeper.GetModuleAddress(types.ModuleName), saleSplitAddr, sharePortionCoin)
			}
			remaining = remaining.Sub(sharePortionCoin)
		}
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, owner, sdk.NewCoins(remaining))
		if err != nil {
			return err
		}
	} else {
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, owner, sdk.NewCoins(remaining))
		if err != nil {
			return err
		}
	}

	k.DeleteListing(ctx, listing)
	return nil
}

func (k Keeper) GetProportions(totalCoin sdk.Coin, ratio sdk.Dec) sdk.Coin {
	return sdk.NewCoin(totalCoin.Denom, sdk.NewDecFromInt(totalCoin.Amount).Mul(ratio).TruncateInt())
}

func (k Keeper) DistributeCommission(ctx sdk.Context, marketplaceCoin sdk.Coin) error {
	distrParams := k.GetMarketplaceDistributionParams(ctx)
	stakingCommissionCoin := k.GetProportions(marketplaceCoin, distrParams.Staking)
	if distrParams.Staking.GT(sdk.ZeroDec()) && stakingCommissionCoin.Amount.GT(sdk.ZeroInt()) {
		err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, authtypes.FeeCollectorName, sdk.NewCoins(stakingCommissionCoin))
		if err != nil {
			return err
		}
		k.createSaleCommissionTransferEvent(ctx,
			k.accountKeeper.GetModuleAddress(types.ModuleName),
			k.accountKeeper.GetModuleAddress(authtypes.FeeCollectorName),
			stakingCommissionCoin,
		)
		marketplaceCoin = marketplaceCoin.Sub(stakingCommissionCoin)
	}
	communityPoolCommissionCoin := marketplaceCoin

	err := k.distributionKeeper.FundCommunityPool(
		ctx,
		sdk.NewCoins(communityPoolCommissionCoin),
		k.accountKeeper.GetModuleAddress(types.ModuleName),
	)
	if err != nil {
		return err
	}
	k.createSaleCommissionTransferEvent(ctx,
		k.accountKeeper.GetModuleAddress(types.ModuleName),
		k.accountKeeper.GetModuleAddress("distribution"),
		communityPoolCommissionCoin,
	)

	return nil
}

// CreateAuctionListing creates a auction in the store and set owner to auction and updates the next auction number
func (k Keeper) CreateAuctionListing(ctx sdk.Context, auction types.AuctionListing) error {
	// check auction already exists or not
	if k.HasAuctionListing(ctx, auction.GetId()) {
		return errorsmod.Wrapf(types.ErrListingAlreadyExists, "auction listing already exists: %d", auction.GetId())
	}

	err := k.nftKeeper.TransferOwnership(ctx,
		auction.GetDenomId(), auction.GetNftId(), auction.GetOwner(),
		k.accountKeeper.GetModuleAddress(types.ModuleName))
	if err != nil {
		return err
	}
	// set auction listing
	k.SetAuctionListing(ctx, auction)

	if len(auction.GetOwner()) != 0 {
		// set auction listing id with owner prefix
		k.SetAuctionListingWithOwner(ctx, auction.GetOwner(), auction.GetId())
	}
	// Update auction listing next number
	auctionId := k.GetNextAuctionNumber(ctx)
	k.SetNextAuctionNumber(ctx, auctionId+1)
	k.SetAuctionListingWithNFTID(ctx, auction.NftId, auction.Id)

	if len(auction.StartPrice.Denom) > 0 {
		k.SetAuctionListingWithPriceDenom(ctx, auction.StartPrice.Denom, auction.Id)
	}
	return nil
}

func (k Keeper) CancelAuctionListing(ctx sdk.Context, auction types.AuctionListing) error {
	// Check bid Exists or Not
	if k.HasBid(ctx, auction.Id) {
		return errorsmod.Wrapf(types.ErrBidExists, "cannot cancel auction %d, bid exists ", auction.Id)
	}

	// Transfer Back NFT ownership to auction owner
	err := k.nftKeeper.TransferOwnership(ctx, auction.GetDenomId(), auction.GetNftId(),
		k.accountKeeper.GetModuleAddress(types.ModuleName), auction.GetOwner())
	if err != nil {
		return err
	}
	k.RemoveAuctionListing(ctx, auction.GetId())
	k.UnsetAuctionListingWithOwner(ctx, auction.GetOwner(), auction.GetId())
	k.UnsetAuctionListingWithNFTID(ctx, auction.GetNftId())
	k.UnsetAuctionListingWithPriceDenom(ctx, auction.StartPrice.Denom, auction.GetId())

	return nil
}

func (k Keeper) PlaceBid(ctx sdk.Context, auction types.AuctionListing, newBid types.Bid) error {
	// Check bids of auction
	newBidPrice := auction.StartPrice
	prevBid, bidExists := k.GetBid(ctx, auction.Id)
	if bidExists {
		newBidPrice = k.GetNewBidPrice(auction.StartPrice.Denom, prevBid.Amount, auction.IncrementPercentage)
	}
	if newBid.Amount.IsLT(newBidPrice) {
		return errorsmod.Wrapf(types.ErrBidAmountNotEnough,
			"cannot place bid for given auction %d, required amount to bid is %s", auction.Id, newBidPrice.String())
	}

	// Transfer amount from bidder to module account
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, newBid.GetBidder(), types.ModuleName, sdk.NewCoins(newBid.Amount))
	if err != nil {
		return err
	}
	// Release previous Bid
	if bidExists {
		_ = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, prevBid.GetBidder(), sdk.NewCoins(prevBid.Amount))
		k.RemoveBid(ctx, prevBid.AuctionId)
	}
	// Set new bid
	k.SetBid(ctx, newBid)

	return nil
}

func (k Keeper) GetNewBidPrice(denom string, amount sdk.Coin, increment sdk.Dec) sdk.Coin {
	return sdk.NewCoin(denom, amount.Amount.Add(sdk.NewDecFromInt(amount.Amount).Mul(increment).TruncateInt()))
}

func (k Keeper) TransferRoyalty(
	ctx sdk.Context,
	nftRoyaltyShareCoin sdk.Coin,
	royaltyReceivers []*onfttypes.WeightedAddress,
	creator sdk.AccAddress,
) error {
	// if royalty splits configured, distributing royalty
	// else sending royalty to collection creator
	if len(royaltyReceivers) > 0 {
		remaining := nftRoyaltyShareCoin
		for _, share := range royaltyReceivers {
			sharePortionCoin := k.GetProportions(nftRoyaltyShareCoin, share.Weight)
			sharePortionCoins := sdk.NewCoins(sharePortionCoin)
			royaltySplitAddr, err := sdk.AccAddressFromBech32(share.Address)
			if err != nil {
				// ignoring error and sending royalty to creator
				if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creator, sharePortionCoins); err != nil {
					return err
				}
				k.createRoyaltyShareTransferEvent(ctx, k.accountKeeper.GetModuleAddress(types.ModuleName), creator, sharePortionCoin)
			} else {
				err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, royaltySplitAddr, sharePortionCoins)
				if err != nil {
					return err
				}
				k.createRoyaltyShareTransferEvent(ctx, k.accountKeeper.GetModuleAddress(types.ModuleName), royaltySplitAddr, sharePortionCoin)
			}

			remaining = remaining.Sub(sharePortionCoin)
		}
		// sending remaining to creator
		if remaining.Amount.GT(sdk.ZeroInt()) {
			err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creator, sdk.NewCoins(remaining))
			if err != nil {
				return err
			}
			k.createRoyaltyShareTransferEvent(ctx, k.accountKeeper.GetModuleAddress(types.ModuleName), creator, nftRoyaltyShareCoin)
		}
	} else {
		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creator, sdk.NewCoins(nftRoyaltyShareCoin))
		if err != nil {
			return err
		}
		k.createRoyaltyShareTransferEvent(ctx, k.accountKeeper.GetModuleAddress(types.ModuleName), creator, nftRoyaltyShareCoin)
	}
	return nil
}
