package keeper

import (
	"crypto/sha256"
	"fmt"

	errorsmod "cosmossdk.io/errors"

	"github.com/OmniFlix/omniflixhub/v5/x/itc/types"
	nfttypes "github.com/OmniFlix/omniflixhub/v5/x/onft/types"
	streampaytypes "github.com/OmniFlix/streampay/v2/x/streampay/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// CreateCampaign ...
func (k Keeper) CreateCampaign(
	ctx sdk.Context,
	creator sdk.AccAddress,
	campaign types.Campaign,
) error {
	_, err := k.nftKeeper.GetDenomInfo(ctx, campaign.NftDenomId)
	if err != nil {
		return err
	}
	if campaign.ClaimType == types.CLAIM_TYPE_FT || campaign.ClaimType == types.CLAIM_TYPE_FT_AND_NFT {
		if err := k.bankKeeper.SendCoinsFromAccountToModule(
			ctx,
			creator,
			types.ModuleName,
			sdk.NewCoins(campaign.TotalTokens),
		); err != nil {
			return err
		}
	}

	if campaign.ClaimType == types.CLAIM_TYPE_NFT || campaign.ClaimType == types.CLAIM_TYPE_FT_AND_NFT {
		// check if the mint collection exists
		mintCollection, err := k.nftKeeper.GetDenomInfo(ctx, campaign.NftMintDetails.DenomId)
		if err != nil {
			return err
		}
		// Authorize
		if mintCollection.Creator != campaign.Creator {
			return errorsmod.Wrapf(
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
	k.emitCreateCampaignEvent(ctx, campaign)

	return nil
}

// CancelCampaign ...
func (k Keeper) CancelCampaign(ctx sdk.Context, campaignId uint64, creator sdk.AccAddress) error {
	campaign, found := k.GetCampaign(ctx, campaignId)
	if !found {
		return errorsmod.Wrapf(types.ErrCampaignDoesNotExists, "campaign %d not exists", campaignId)
	}
	if creator.String() != campaign.Creator {
		return errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "unauthorized address %s", creator.String())
	}

	// return funds
	if campaign.AvailableTokens.IsValid() && campaign.AvailableTokens.Amount.GT(sdk.ZeroInt()) {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(
			ctx,
			types.ModuleName,
			campaign.GetCreator(),
			sdk.NewCoins(campaign.AvailableTokens),
		); err != nil {
			panic(err)
		}
	}
	// return NFTs if received
	if len(campaign.ReceivedNftIds) > 0 {
		for _, nftId := range campaign.ReceivedNftIds {
			if err := k.nftKeeper.TransferOwnership(ctx,
				campaign.NftDenomId,
				nftId,
				k.GetModuleAccountAddress(ctx),
				campaign.GetCreator(),
			); err != nil {
				return err
			}
		}
	}
	k.UnsetCampaignWithCreator(ctx, creator, campaignId)
	k.RemoveCampaign(ctx, campaignId)
	k.emitCancelCampaignEvent(ctx, campaignId, creator.String())

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
		return errorsmod.Wrapf(types.ErrClaimExists,
			"claim exists with given nft  %s", nft.GetID())
	}
	if (campaign.MaxAllowedClaims - campaign.ClaimCount) <= 0 {
		return errorsmod.Wrapf(types.ErrClaimNotAllowed,
			"max allowed claims reached for this campaign (campaign: %d, maxAllowedClaims: %d).",
			campaign.GetId(),
			campaign.GetMaxAllowedClaims(),
		)
	}

	if campaign.ClaimType == types.CLAIM_TYPE_FT || campaign.ClaimType == types.CLAIM_TYPE_FT_AND_NFT {
		if campaign.AvailableTokens.IsLT(campaign.TokensPerClaim) {
			return errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds,
				"insufficient available tokens, available tokens  %s",
				campaign.AvailableTokens.String(),
			)
		}
	}
	if !claimer.Equals(nft.GetOwner()) {
		return errorsmod.Wrapf(sdkerrors.ErrUnauthorized,
			"nft %s isn't owned by address  %s", claim.NftId, claimer.String())
	}

	if campaign.Interaction == types.INTERACTION_TYPE_TRANSFER {
		err := k.nftKeeper.TransferOwnership(ctx,
			campaign.NftDenomId,
			nft.GetID(),
			claimer,
			k.GetModuleAccountAddress(ctx),
		)
		if err != nil {
			return err
		}
		campaign.ReceivedNftIds = append(campaign.ReceivedNftIds, nft.GetID())
	} else if campaign.Interaction == types.INTERACTION_TYPE_BURN {
		err := k.nftKeeper.BurnONFT(ctx,
			campaign.NftDenomId,
			nft.GetID(),
			claimer,
		)
		if err != nil {
			return err
		}
	}

	// Claim Claimable
	if campaign.ClaimType == types.CLAIM_TYPE_FT || campaign.ClaimType == types.CLAIM_TYPE_FT_AND_NFT {
		claimAmount := campaign.TokensPerClaim
		if campaign.Distribution != nil && campaign.Distribution.Type == types.DISTRIBUTION_TYPE_STREAM {
			if _, err := k.streampayKeeper.CreateStreamPayment(
				ctx,
				k.GetModuleAccountAddress(ctx),
				claimer, claimAmount,
				streampaytypes.TypeContinuous,
				campaign.Distribution.StreamDuration,
				nil,
				false,
			); err != nil {
				return err
			}
		} else {
			if err := k.bankKeeper.SendCoinsFromModuleToAccount(
				ctx,
				types.ModuleName,
				claimer,
				sdk.NewCoins(claimAmount),
			); err != nil {
				return err
			}
		}
		campaign.AvailableTokens.Amount = campaign.AvailableTokens.Amount.Sub(claimAmount.Amount)
	}

	if campaign.ClaimType == types.CLAIM_TYPE_NFT || campaign.ClaimType == types.CLAIM_TYPE_FT_AND_NFT {
		nftIndex := campaign.NftMintDetails.StartIndex + campaign.MintCount
		nftTitle := fmt.Sprintf(
			"%s %s%d",
			campaign.NftMintDetails.Name,
			campaign.NftMintDetails.NameDelimiter,
			nftIndex,
		)
		if err := k.nftKeeper.MintONFT(
			ctx,
			campaign.NftMintDetails.DenomId,
			generateClaimNftId(ctx, campaign.Id, nftIndex),
			nftTitle,
			campaign.NftMintDetails.Description,
			campaign.NftMintDetails.MediaUri,
			campaign.NftMintDetails.UriHash,
			campaign.NftMintDetails.PreviewUri,
			campaign.NftMintDetails.Data,
			ctx.BlockTime().UTC(),
			campaign.NftMintDetails.Transferable,
			campaign.NftMintDetails.Extensible,
			campaign.NftMintDetails.Nsfw,
			campaign.NftMintDetails.RoyaltyShare,
			claimer,
		); err != nil {
			return errorsmod.Wrapf(types.ErrClaimingNFT,
				"unable to mint nft denomId %s", campaign.NftMintDetails.DenomId)
		}
		// set campaign mint count
		campaign.MintCount += 1
	}
	// set claim
	k.SetClaim(ctx, claim)

	// set campaign
	campaign.ClaimCount += 1
	k.SetCampaign(ctx, campaign)

	// emit events
	k.emitClaimEvent(ctx, campaign.Id, claimer.String(), nft.GetID())

	return nil
}

func (k Keeper) DepositCampaign(ctx sdk.Context, campaignId uint64, depositor sdk.AccAddress, amount sdk.Coin) error {
	campaign, found := k.GetCampaign(ctx, campaignId)
	if !found {
		return errorsmod.Wrapf(types.ErrCampaignDoesNotExists, "campaign id %d not exists", campaignId)
	}
	if depositor.String() != campaign.Creator {
		return errorsmod.Wrapf(types.ErrDepositNotAllowed, "deposit not allowed from address %s"+
			"only creator of the campaign is allowed to deposit", depositor.String())
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
	campaign.TotalTokens = campaign.TotalTokens.Add(amount)
	campaign.AvailableTokens = campaign.AvailableTokens.Add(amount)

	k.SetCampaign(ctx, campaign)
	k.emitDepositCampaignEvent(ctx, campaignId, depositor.String(), amount)

	return nil
}

func generateClaimNftId(ctx sdk.Context, campaignId, mintCount uint64) string {
	lastBlockHash := fmt.Sprintf("%x", ctx.BlockHeader().LastBlockId.Hash)
	hash := sha256.Sum256([]byte(fmt.Sprintf("%sitc%dmc%d", lastBlockHash, campaignId, mintCount)))
	return nfttypes.IDPrefix + fmt.Sprintf("%x", hash)[:32]
}
