package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/OmniFlix/omniflixhub/x/alloc/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type (
	Keeper struct {
		cdc      codec.BinaryCodec
		storeKey sdk.StoreKey
		memKey   sdk.StoreKey

		accountKeeper types.AccountKeeper
		bankKeeper    types.BankKeeper
		stakingKeeper types.StakingKeeper
		distrKeeper   types.DistrKeeper

		paramstore paramtypes.Subspace
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,

	accountKeeper types.AccountKeeper, bankKeeper types.BankKeeper, stakingKeeper types.StakingKeeper, distrKeeper types.DistrKeeper,
	ps paramtypes.Subspace,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		memKey:   memKey,

		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		stakingKeeper: stakingKeeper,
		distrKeeper:   distrKeeper,
		paramstore:    ps,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetModuleAccountBalance gets the coin balance of module account
func (k Keeper) GetModuleAccountAddress(ctx sdk.Context) sdk.AccAddress {
	return k.accountKeeper.GetModuleAddress(types.ModuleName)
}

// DistributeMintedCoins implements distribution of minted coins from mint to external modules.
func (k Keeper) DistributeMintedCoins(ctx sdk.Context) error {
	blockRewardsAddr := k.accountKeeper.GetModuleAccount(ctx, authtypes.FeeCollectorName).GetAddress()
	blockRewards := k.bankKeeper.GetBalance(ctx, blockRewardsAddr, k.stakingKeeper.BondDenom(ctx))
	blockRewardsDec := sdk.NewDecFromInt(blockRewards.Amount)

	params := k.GetParams(ctx)
	proportions := params.DistributionProportions

	nftIncentiveAmount := blockRewardsDec.Mul(proportions.NftIncentives)
	nftIncentiveCoin := sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), nftIncentiveAmount.TruncateInt())
	// Distribute NFT incentives to the community pool until a future update
	err := k.distrKeeper.FundCommunityPool(ctx, sdk.NewCoins(nftIncentiveCoin), blockRewardsAddr)
	if err != nil {
		return err
	}
	k.Logger(ctx).Debug("funded community pool with nft incentives", "amount", nftIncentiveCoin.String(), "from", blockRewardsAddr)

	nodeHostsIncentiveAmount := blockRewardsDec.Mul(proportions.NodeHostsIncentives)
	nodeHostsIncentiveCoin := sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), nodeHostsIncentiveAmount.TruncateInt())
	// Distribute node hosts incentives to the community pool until a future update
	err = k.distrKeeper.FundCommunityPool(ctx, sdk.NewCoins(nodeHostsIncentiveCoin), blockRewardsAddr)
	if err != nil {
		return err
	}
	k.Logger(ctx).Debug("funded community pool with node host incentives", "amount", nodeHostsIncentiveCoin.String(), "from", blockRewardsAddr)

	devRewardAmount := blockRewardsDec.Mul(proportions.DeveloperRewards)
	devRewardCoin := sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), devRewardAmount.TruncateInt())

	for _, w := range params.WeightedDeveloperRewardsReceivers {
		devRewardPortionCoins := sdk.NewCoins(k.GetProportions(ctx, devRewardCoin, w.Weight))
		if w.Address == "" {
			err := k.distrKeeper.FundCommunityPool(ctx, devRewardPortionCoins, blockRewardsAddr)
			if err != nil {
				return err
			}
		} else {
			devRewardsAddr, err := sdk.AccAddressFromBech32(w.Address)
			if err != nil {
				return err
			}
			err = k.bankKeeper.SendCoins(ctx, blockRewardsAddr, devRewardsAddr, devRewardPortionCoins)
			if err != nil {
				return err
			}
			k.Logger(ctx).Debug("sent coins to developer", "amount", devRewardPortionCoins.String(), "from", blockRewardsAddr, "to", devRewardsAddr)
		}
	}
	// calculate staking rewards
	stakingRewardsCoins := sdk.NewCoins(k.GetProportions(ctx, blockRewards, proportions.StakingRewards))

	// subtract from original provision to ensure no coins left over after the allocations
	communityPoolCoins := sdk.NewCoins(blockRewards).Sub(stakingRewardsCoins).Sub(sdk.NewCoins(nftIncentiveCoin)).Sub(sdk.NewCoins(nodeHostsIncentiveCoin)).Sub(sdk.NewCoins(devRewardCoin))
	err = k.distrKeeper.FundCommunityPool(ctx, communityPoolCoins, blockRewardsAddr)
	if err != nil {
		return err
	}

	return nil
}

// GetProportions gets the balance of the `MintedDenom` from minted coins
// and returns coins according to the `AllocationRatio`
func (k Keeper) GetProportions(ctx sdk.Context, mintedCoin sdk.Coin, ratio sdk.Dec) sdk.Coin {
	return sdk.NewCoin(mintedCoin.Denom, mintedCoin.Amount.ToDec().Mul(ratio).TruncateInt())
}
