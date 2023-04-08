package keeper

import (
	"context"
	"time"

	errorsmod "cosmossdk.io/errors"

	"github.com/OmniFlix/omniflixhub/v2/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"golang.org/x/exp/slices"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the token MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) ListNFT(goCtx context.Context, msg *types.MsgListNFT) (*types.MsgListNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	nft, err := m.nftKeeper.GetONFT(ctx, msg.DenomId, msg.NftId)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrNftNotExists,
			"invalid nft and or denomId, nftId %s, denomId %s", msg.NftId, msg.DenomId)
	}
	if owner.String() != nft.GetOwner().String() {
		return nil, errorsmod.Wrapf(types.ErrUnauthorized, "unauthorized address %s", owner)
	}
	if !nft.IsTransferable() {
		return nil, errorsmod.Wrapf(
			types.ErrNftNonTransferable, "non-transferable nfts not allowed to list in marketplace")
	}

	listing := types.NewListing(msg.Id, msg.NftId, msg.DenomId, msg.Price, owner, msg.SplitShares)
	err = m.Keeper.AddListing(ctx, listing)
	if err != nil {
		return nil, err
	}

	m.Keeper.createListNftEvent(ctx, owner, listing.Id, listing.DenomId, listing.NftId, listing.Price)

	return &types.MsgListNFTResponse{}, nil
}

func (m msgServer) EditListing(goCtx context.Context,
	msg *types.MsgEditListing,
) (*types.MsgEditListingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	listing, found := m.Keeper.GetListing(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrListingDoesNotExists, "listing id %s not exists", msg.Id)
	}
	if owner.String() != listing.Owner {
		return nil, errorsmod.Wrapf(types.ErrUnauthorized, "unauthorized address %s", owner)
	}
	if err := types.ValidatePrice(msg.Price); err != nil {
		return nil, err
	}
	listing.Price = msg.Price
	m.Keeper.SetListing(ctx, listing)

	m.Keeper.createEditListingEvent(ctx, owner, listing.Id, listing.Price)

	return &types.MsgEditListingResponse{}, nil
}

func (m msgServer) DeListNFT(goCtx context.Context,
	msg *types.MsgDeListNFT,
) (*types.MsgDeListNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}
	listing, found := m.Keeper.GetListing(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrListingDoesNotExists, "listing id %s not exists", msg.Id)
	}
	if owner.String() != listing.Owner {
		return nil, errorsmod.Wrapf(types.ErrUnauthorized, "unauthorized address %s", owner)
	}
	err = m.nftKeeper.TransferOwnership(ctx, listing.GetDenomId(), listing.GetNftId(),
		m.accountKeeper.GetModuleAddress(types.ModuleName), listing.GetOwner())
	if err != nil {
		return nil, err
	}
	m.Keeper.DeleteListing(ctx, listing)

	m.Keeper.createDeListNftEvent(ctx, owner, listing.Id)

	return &types.MsgDeListNFTResponse{}, nil
}

func (m msgServer) BuyNFT(goCtx context.Context, msg *types.MsgBuyNFT) (*types.MsgBuyNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	buyer, err := sdk.AccAddressFromBech32(msg.Buyer)
	if err != nil {
		return nil, err
	}

	listing, found := m.Keeper.GetListing(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrListingDoesNotExists, "listing id %s not exists", msg.Id)
	}
	if err := types.ValidatePrice(msg.Price); err != nil {
		return nil, err
	}
	if !msg.Price.Equal(listing.Price) {
		if msg.Price.Denom != listing.Price.Denom {
			return nil, errorsmod.Wrapf(types.ErrInvalidPriceDenom, "invalid price denom %s", msg.Price.Denom)
		}
		if msg.Price.Amount.LT(listing.Price.Amount) {
			return nil, errorsmod.Wrapf(types.ErrNotEnoughAmount,
				"%s is not enough, to buy %s required", msg.Price.String(), listing.Price.String())
		}
		return nil, errorsmod.Wrapf(types.ErrInvalidPrice,
			"price %s not matched with listing price", msg.Price.String())
	}
	err = m.Keeper.Buy(ctx, listing, buyer)
	if err != nil {
		return nil, err
	}

	m.Keeper.createBuyNftEvent(ctx, buyer, listing.Id, listing.NftId, listing.Price)

	return &types.MsgBuyNFTResponse{}, nil
}

// CreateAuction
func (m msgServer) CreateAuction(goCtx context.Context, msg *types.MsgCreateAuction) (*types.MsgCreateAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}
	if err := msg.Validate(ctx.BlockTime()); err != nil {
		return nil, err
	}

	nft, err := m.nftKeeper.GetONFT(ctx, msg.DenomId, msg.NftId)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrNftNotExists,
			"invalid nft and or denomId, nftId %s, denomId %s", msg.NftId, msg.DenomId)
	}
	if owner.String() != nft.GetOwner().String() {
		return nil, errorsmod.Wrapf(types.ErrUnauthorized, "unauthorized address %s", owner)
	}
	if !nft.IsTransferable() {
		return nil, errorsmod.Wrapf(
			types.ErrNftNonTransferable, "non-transferable nfts not allowed to list in marketplace")
	}
	var endTime *time.Time
	if msg.Duration != nil {
		endAt := msg.StartTime.Add(*msg.Duration)
		endTime = &endAt
		if endTime.Before(msg.StartTime) || endTime.Equal(msg.StartTime) {
			return nil, errorsmod.Wrapf(types.ErrInvalidDuration, "duration must be positive or nil")
		}
	}
	auctionNumber := m.Keeper.GetNextAuctionNumber(ctx)
	auction := types.NewAuctionListing(auctionNumber, msg.NftId, msg.DenomId,
		msg.StartTime, endTime, msg.StartPrice,
		msg.IncrementPercentage, owner, msg.WhitelistAccounts, msg.SplitShares)
	err = m.Keeper.CreateAuctionListing(ctx, auction)
	if err != nil {
		return nil, err
	}

	m.Keeper.createAuctionEvent(ctx, auction)

	return &types.MsgCreateAuctionResponse{
		Auction: &auction,
	}, nil
}

// CancelAuction
func (m msgServer) CancelAuction(goCtx context.Context, msg *types.MsgCancelAuction) (*types.MsgCancelAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	auction, found := m.Keeper.GetAuctionListing(ctx, msg.AuctionId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrAuctionDoesNotExists, "auction id %d not exists", msg.AuctionId)
	}
	if owner.String() != auction.Owner {
		return nil, errorsmod.Wrapf(types.ErrUnauthorized, "unauthorized address %s", owner.String())
	}

	err = m.Keeper.CancelAuctionListing(ctx, auction)
	if err != nil {
		return nil, err
	}

	m.Keeper.cancelAuctionEvent(ctx, auction)

	return &types.MsgCancelAuctionResponse{}, nil
}

// PlaceBid
func (m msgServer) PlaceBid(goCtx context.Context, msg *types.MsgPlaceBid) (*types.MsgPlaceBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	bidder, err := sdk.AccAddressFromBech32(msg.Bidder)
	if err != nil {
		return nil, err
	}

	auction, found := m.Keeper.GetAuctionListing(ctx, msg.AuctionId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrAuctionDoesNotExists, "auction id %d not exists", msg.AuctionId)
	}
	if !auction.StartTime.Before(ctx.BlockTime()) {
		return nil, errorsmod.Wrapf(types.ErrInActiveAuction, "cannot place a bid for inactive auction %d, ", auction.Id)
	}
	if len(auction.WhitelistAccounts) > 0 && !slices.Contains(auction.WhitelistAccounts, bidder.String()) {
		return nil, errorsmod.Wrapf(types.ErrUnauthorized, "cannot place a bid for this auction %d, only whitelisted accounts allowed to bid", auction.Id)
	}
	if msg.Amount.GetDenom() != auction.StartPrice.GetDenom() {
		return nil, errorsmod.Wrapf(types.ErrInvalidPriceDenom, "given auction only accepts bids in %s, ", auction.StartPrice.GetDenom())
	}

	bid := types.NewBid(auction.Id, msg.Amount, ctx.BlockTime(), bidder)

	err = m.Keeper.PlaceBid(ctx, auction, bid)
	if err != nil {
		return nil, err
	}

	m.Keeper.placeBidEvent(ctx, auction, bid)

	return &types.MsgPlaceBidResponse{}, nil
}
