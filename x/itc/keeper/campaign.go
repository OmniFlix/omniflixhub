package keeper

import (
	"github.com/OmniFlix/omniflixhub/x/itc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// CreateCampaign TODO: NFT claim ownership transfers
func (k Keeper) CreateCampaign(ctx sdk.Context, creator sdk.AccAddress, campaign types.Campaign) error {
	// verify collection
	collection, err := k.nftKeeper.GetDenom(ctx, campaign.NftDenomId)
	if err != nil {
		return err
	}
	// Authorize
	if collection.Creator != campaign.Creator {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized,
			"denom id %s isn't owned by campaign creator %s", collection.Id, campaign.Creator)
	}
	if campaign.ClaimType == types.CLAIM_TYPE_FT || campaign.ClaimType == types.CLAIM_TYPE_FT_AND_NFT {
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, creator, types.ModuleName,
			sdk.Coins{*campaign.TotalTokens.Fungible}); err != nil {
			return err
		}
	}

	/*
		if campaign.ClaimType == types.CLAIM_TYPE_NFT {
			mintCollection, err := k.nftKeeper.GetDenom(ctx, campaign.NftMintDetails.DenomId)
			if err != nil {
				return err
			}
			// Authorize
			if mintCollection.Creator != campaign.Creator {
				return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized,
					"nft mint denom id %s isn't owned by campaign creator %s", mintCollection.Id, campaign.Creator)
			}
		}
	*/

	k.SetCampaign(ctx, campaign)
	k.SetNextCampaignNumber(ctx, campaign.Id+1)
	k.SetCampaignWithCreator(ctx, creator, campaign.Id)
	k.SetInactiveCampaign(ctx, campaign.Id)

	return nil
}

// CancelCampaign TODO: cancel campaign and return back funds/ nfts
func (k Keeper) CancelCampaign(ctx sdk.Context, campaign types.Campaign) error {
	return nil
}

// Claim TODO: vesting distribution
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
		if campaign.AvailableTokens.Fungible.IsLT(*campaign.ClaimableTokens.Fungible) {
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
	if campaign.ClaimType == types.CLAIM_TYPE_FT {
		claimableTokens := campaign.ClaimableTokens.Fungible
		if campaign.Distribution != nil && campaign.Distribution.Type == types.DISTRIBUTION_TYPE_VEST {
			// TODO: add vesting based on vesting periods
			if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx,
				types.ModuleName, claimer,
				sdk.NewCoins(sdk.NewCoin(claimableTokens.Denom, claimableTokens.Amount))); err != nil {
				return err
			}
		} else {
			if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx,
				types.ModuleName, claimer,
				sdk.NewCoins(sdk.NewCoin(claimableTokens.Denom, claimableTokens.Amount))); err != nil {
				return err
			}
		}
		availableTokensAmount := campaign.AvailableTokens.Fungible.Amount.Sub(claimableTokens.Amount)
		campaign.AvailableTokens.Fungible.Amount = availableTokensAmount
	}
	// set claim

	k.SetClaim(ctx, claim)
	k.SetClaimWithNft(ctx, campaign.Id, claim.NftId)

	// set campaign
	k.SetCampaign(ctx, campaign)

	// burn nft
	if campaign.Interaction == types.INTERACTION_TYPE_BURN {
		_ = k.nftKeeper.BurnONFT(ctx, campaign.NftDenomId, nft.GetID(), k.GetModuleAccountAddress(ctx))
	}

	return nil
}
