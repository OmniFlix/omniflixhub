package keeper

import (
	"fmt"

	streampaytypes "github.com/OmniFlix/streampay/x/streampay/types"

	"github.com/OmniFlix/omniflixhub/x/itc/types"
	nfttypes "github.com/OmniFlix/onft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// CreateCampaign ...
func (k Keeper) CreateCampaign(ctx sdk.Context, creator sdk.AccAddress, campaign types.Campaign) error {
	// verify collection
	collection, err := k.nftKeeper.GetDenom(ctx, campaign.NftDenomId)
	if err != nil {
		return err
	}
	// Authorize
	if collection.Creator != campaign.Creator {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized,
			"nft denom id %s isn't owned by campaign creator %s", collection.Id, campaign.Creator)
	}
	if campaign.ClaimType == types.CLAIM_TYPE_FT || campaign.ClaimType == types.CLAIM_TYPE_FT_AND_NFT {
		if err := k.bankKeeper.SendCoinsFromAccountToModule(
			ctx,
			creator,
			types.ModuleName,
			sdk.NewCoins(sdk.NewCoin(campaign.TotalTokens.Fungible.Denom, campaign.TotalTokens.Fungible.Amount)),
		); err != nil {
			return err
		}
	}

	if campaign.ClaimType == types.CLAIM_TYPE_NFT || campaign.ClaimType == types.CLAIM_TYPE_FT_AND_NFT {
		mintCollection, err := k.nftKeeper.GetDenom(ctx, campaign.NftMintDetails.DenomId)
		if err != nil {
			return err
		}
		// Authorize
		if mintCollection.Creator != campaign.Creator {
			return sdkerrors.Wrapf(
				sdkerrors.ErrUnauthorized,
				"nft mint denom id %s isn't owned by campaign creator %s",
				mintCollection.Id,
				campaign.Creator,
			)
		}
	}

	k.SetCampaign(ctx, campaign)
	k.SetNextCampaignNumber(ctx, campaign.Id+1)
	k.SetCampaignWithCreator(ctx, creator, campaign.Id)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateCampaign,
			sdk.NewAttribute(types.AttributeKeyCampaignId, fmt.Sprint(campaign.GetId())),
			sdk.NewAttribute(types.AttributeKeyCreator, creator.String()),
		),
	)
	if campaign.ClaimType == types.CLAIM_TYPE_FT && campaign.TotalTokens.Fungible != nil {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCampaignDeposit,
				sdk.NewAttribute(types.AttributeKeyCampaignId, fmt.Sprint(campaign.GetId())),
				sdk.NewAttribute(types.AttributeKeyDepositor, creator.String()),
			),
		)
	}

	return nil
}

// CancelCampaign ...
func (k Keeper) CancelCampaign(ctx sdk.Context, campaignId uint64, creator sdk.AccAddress) error {
	campaign, found := k.GetCampaign(ctx, campaignId)
	if !found {
		return sdkerrors.Wrapf(types.ErrCampaignDoesNotExists, "campaign %d not exists", campaignId)
	}
	if creator.String() != campaign.Creator {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "unauthorized address %s", creator.String())
	}
	// cancel only if campaign not started
	if campaign.StartTime.Before(ctx.BlockTime()) {
		return sdkerrors.Wrapf(types.ErrCancelNotAllowed, "active campaign can not be canceled")
	}
	// return funds
	availableTokens := campaign.AvailableTokens.Fungible
	if availableTokens.IsValid() && !availableTokens.IsPositive() {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx,
			types.ModuleName, campaign.GetCreator(),
			sdk.NewCoins(sdk.NewCoin(availableTokens.Denom, availableTokens.Amount))); err != nil {
			panic(err)
		}
	}
	k.UnsetCampaignWithCreator(ctx, creator, campaignId)
	k.RemoveCampaign(ctx, campaignId)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCancelCampaign,
			sdk.NewAttribute(types.AttributeKeyCampaignId, fmt.Sprint(campaign.GetId())),
			sdk.NewAttribute(types.AttributeKeyCreator, creator.String()),
		),
	)
	return nil
}

// Claim ...
func (k Keeper) Claim(ctx sdk.Context, campaign types.Campaign, claimer sdk.AccAddress, claim types.Claim) error {
	// check nft with campaign
	nft, err := k.nftKeeper.GetONFT(ctx, campaign.NftDenomId, claim.NftId)
	if err != nil {
		return err
	}
	// check if claim exists with given nft
	if k.HasClaim(ctx, campaign.GetId(), nft.GetID()) {
		return sdkerrors.Wrapf(types.ErrClaimExists,
			"claim exists with given nft  %s", nft.GetID())
	}
	claims := k.GetClaims(ctx, campaign.GetId())
	if uint64(len(claims)) >= campaign.MaxAllowedClaims {
		return sdkerrors.Wrapf(types.ErrClaimNotAllowed,
			"max allowed claims reached for this campaign (campaign: %d, maxAllowedClaims: %d).",
			campaign.GetId(),
			campaign.GetMaxAllowedClaims(),
		)
	}

	if campaign.ClaimType == types.CLAIM_TYPE_FT {
		if campaign.AvailableTokens.Fungible.IsLT(*campaign.TokensPerClaim.Fungible) {
			return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
				"insufficient available tokens, available tokens  %s",
				campaign.AvailableTokens.Fungible.String(),
			)
		}
	}

	if campaign.Interaction == types.INTERACTION_TYPE_HOLD {
		if !claimer.Equals(nft.GetOwner()) {
			return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized,
				"nft %s isn't owned by claimer  %s", claim.NftId, claimer.String())
		}
	} else {
		// TransferOwnership to module account
		err := k.nftKeeper.TransferOwnership(ctx,
			campaign.NftDenomId,
			nft.GetID(),
			claimer,
			k.GetModuleAccountAddress(ctx),
		)
		if err != nil {
			return err
		}
		if campaign.Interaction == types.INTERACTION_TYPE_TRANSFER {
			campaign.ReceivedNftIds = append(campaign.ReceivedNftIds, nft.GetID())
		}
	}

	// Claim Claimable
	if campaign.ClaimType == types.CLAIM_TYPE_FT || campaign.ClaimType == types.CLAIM_TYPE_FT_AND_NFT {
		tokensPerClaim := campaign.TokensPerClaim.Fungible
		if campaign.Distribution != nil && campaign.Distribution.Type == types.DISTRIBUTION_TYPE_STREAM {
			if err = k.streampayKeeper.CreatePaymentStream(ctx,
				k.GetModuleAccountAddress(ctx),
				claimer, sdk.NewCoin(tokensPerClaim.Denom, tokensPerClaim.Amount),
				streampaytypes.TypeContinuous,
				ctx.BlockTime().Add(*campaign.Distribution.StreamDuration),
			); err != nil {
				return err
			}
		} else {
			if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx,
				types.ModuleName, claimer,
				sdk.NewCoins(sdk.NewCoin(tokensPerClaim.Denom, tokensPerClaim.Amount))); err != nil {
				return err
			}
		}
		availableTokensAmount := campaign.AvailableTokens.Fungible.Amount.Sub(tokensPerClaim.Amount)
		campaign.AvailableTokens.Fungible.Amount = availableTokensAmount
	} else if campaign.ClaimType == types.CLAIM_TYPE_NFT || campaign.ClaimType == types.CLAIM_TYPE_FT_AND_NFT {
		if err := k.nftKeeper.MintONFT(
			ctx,
			campaign.NftMintDetails.DenomId,
			nfttypes.GenUniqueID("onft"),
			nfttypes.Metadata{
				Name:        campaign.NftMintDetails.Name,
				Description: campaign.NftMintDetails.Description,
				MediaURI:    campaign.NftMintDetails.MediaUri,
				PreviewURI:  campaign.NftMintDetails.PreviewUri,
			},
			campaign.NftMintDetails.Data,
			campaign.NftMintDetails.Transferable,
			campaign.NftMintDetails.Extensible,
			campaign.NftMintDetails.Nsfw,
			campaign.NftMintDetails.RoyaltyShare,
			campaign.GetCreator(),
			claimer,
		); err != nil {
			// Transfer back nft if it's not hold
			if campaign.Interaction != types.INTERACTION_TYPE_HOLD {
				err := k.nftKeeper.TransferOwnership(ctx,
					campaign.NftDenomId,
					nft.GetID(),
					k.GetModuleAccountAddress(ctx),
					claimer,
				)
				if err != nil {
					panic(err)
				}
			}
			return sdkerrors.Wrapf(types.ErrClaimingNFT,
				"unable to mint nft denomId %s", campaign.NftMintDetails.DenomId)
		}
	}
	// set claim

	k.SetClaim(ctx, claim)

	// set campaign
	k.SetCampaign(ctx, campaign)

	// burn nft
	if campaign.Interaction == types.INTERACTION_TYPE_BURN {
		_ = k.nftKeeper.BurnONFT(ctx, campaign.NftDenomId, nft.GetID(), k.GetModuleAccountAddress(ctx))
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeClaim,
			sdk.NewAttribute(types.AttributeKeyCampaignId, fmt.Sprint(campaign.GetId())),
			sdk.NewAttribute(types.AttributeKeyClaimer, claimer.String()),
		),
	)

	return nil
}

func (k Keeper) CampaignDeposit(ctx sdk.Context, campaignId uint64, depositor sdk.AccAddress, amount sdk.Coin) error {
	campaign, found := k.GetCampaign(ctx, campaignId)
	if !found {
		return sdkerrors.Wrapf(types.ErrCampaignDoesNotExists, "campaign id %d not exists", campaignId)
	}
	// Transfer tokens from depositor to module account
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		depositor,
		types.ModuleName,
		sdk.NewCoins(amount),
	); err != nil {
		return err
	}
	// Update total tokens and available tokens
	campaign.TotalTokens.Fungible.Amount = campaign.TotalTokens.Fungible.Amount.Add(amount.Amount)
	campaign.AvailableTokens.Fungible.Amount = campaign.AvailableTokens.Fungible.Amount.Add(amount.Amount)

	k.SetCampaign(ctx, campaign)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCampaignDeposit,
			sdk.NewAttribute(types.AttributeKeyCampaignId, fmt.Sprint(campaign.GetId())),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyDepositor, depositor.String()),
		),
	)

	return nil
}
