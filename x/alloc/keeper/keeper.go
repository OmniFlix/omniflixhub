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

	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	distrKeeper types.DistrKeeper,
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

// GetModuleAccountAddress gets the address of module account
func (k Keeper) GetModuleAccountAddress() sdk.AccAddress {
	return k.accountKeeper.GetModuleAddress(types.ModuleName)
}

// DistributeMintedCoins implements distribution of minted coins from mint to external modules.
func (k Keeper) DistributeMintedCoins(ctx sdk.Context) error {
	blockRewardsAddr := k.accountKeeper.GetModuleAccount(ctx, authtypes.FeeCollectorName).GetAddress()
	blockRewards := k.bankKeeper.GetBalance(ctx, blockRewardsAddr, k.stakingKeeper.BondDenom(ctx))
	blockRewardsAmountDec := blockRewards.Amount.ToDec()

	params := k.GetParams(ctx)
	proportions := params.DistributionProportions

	nftIncentiveAmount := blockRewardsAmountDec.Mul(proportions.NftIncentives).TruncateInt()
	nftIncentiveCoin := sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), nftIncentiveAmount)

	k.Logger(ctx).Debug(
		"distributing minted coins to nft incentives receivers",
		"amount", nftIncentiveCoin.String(), "from", blockRewardsAddr,
	)
	if err := k.distributeCoinToWeightedAddresses(
		ctx,
		params.WeightedNftIncentivesReceivers,
		nftIncentiveCoin,
		blockRewardsAddr,
	); err != nil {
		return err
	}

	nodeHostsIncentiveAmount := blockRewardsAmountDec.Mul(proportions.NodeHostsIncentives)
	nodeHostsIncentiveCoin := sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), nodeHostsIncentiveAmount.TruncateInt())

	k.Logger(ctx).Debug(
		"distributing minted coins to node host incentives receivers",
		"amount", nodeHostsIncentiveCoin.String(), "from", blockRewardsAddr,
	)
	if err := k.distributeCoinToWeightedAddresses(
		ctx,
		params.WeightedNodeHostsIncentivesReceivers,
		nodeHostsIncentiveCoin,
		blockRewardsAddr,
	); err != nil {
		return err
	}

	devRewardAmount := blockRewardsAmountDec.Mul(proportions.DeveloperRewards)
	devRewardCoin := sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), devRewardAmount.TruncateInt())

	k.Logger(ctx).Debug(
		"distributing minted coins to developer rewards receivers",
		"amount", devRewardCoin.String(), "from", blockRewardsAddr,
	)

	if err := k.distributeCoinToWeightedAddresses(
		ctx,
		params.WeightedDeveloperRewardsReceivers,
		devRewardCoin,
		blockRewardsAddr,
	); err != nil {
		return err
	}

	// calculate staking rewards
	stakingRewardAmount := blockRewardsAmountDec.Mul(proportions.StakingRewards)
	stakingRewardCoin := sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), stakingRewardAmount.TruncateInt())

	// subtract from original provision to ensure no coins left over after the allocations
	communityPoolCoin := blockRewards.
		Sub(stakingRewardCoin).
		Sub(nftIncentiveCoin).
		Sub(nodeHostsIncentiveCoin).
		Sub(devRewardCoin)

	if err := k.distrKeeper.FundCommunityPool(ctx, sdk.NewCoins(communityPoolCoin), blockRewardsAddr); err != nil {
		return err
	}

	return nil
}

func (k Keeper) distributeCoinToWeightedAddresses(
	ctx sdk.Context,
	weightedAddresses []types.WeightedAddress,
	totalAmount sdk.Coin,
	fromAddress sdk.AccAddress,
) error {
	totalAmountDec := totalAmount.Amount.ToDec()
	for _, w := range weightedAddresses {
		amount := sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), totalAmountDec.Mul(w.Weight).TruncateInt())
		if w.Address == "" {
			err := k.distrKeeper.FundCommunityPool(ctx, sdk.NewCoins(amount), fromAddress)
			if err != nil {
				return err
			}
		} else {
			toAddress, err := sdk.AccAddressFromBech32(w.Address)
			if err != nil {
				return err
			}
			err = k.bankKeeper.SendCoins(ctx, fromAddress, toAddress, sdk.NewCoins(amount))
			if err != nil {
				return err
			}
			k.Logger(ctx).Debug(
				"sent coins to address", "amount", amount.String(), "from", fromAddress, "to", toAddress)
		}
	}
	return nil
}
