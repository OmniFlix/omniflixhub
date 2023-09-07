package keeper

import (
	"errors"
	"fmt"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/cometbft/cometbft/libs/log"

	sdkmath "cosmossdk.io/math"
	"github.com/OmniFlix/omniflixhub/v2/x/alloc/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type (
	Keeper struct {
		cdc      codec.BinaryCodec
		storeKey storetypes.StoreKey
		memKey   storetypes.StoreKey

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
	memKey storetypes.StoreKey,

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

	params := k.GetParams(ctx)
	proportions := params.DistributionProportions

	nftIncentiveCoin, err := getProportionAmount(blockRewards, proportions.NftIncentives)
	if err != nil {
		return err
	}

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

	nodeHostsIncentiveCoin, err := getProportionAmount(blockRewards, proportions.NodeHostsIncentives)
	if err != nil {
		return err
	}

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

	devRewardCoin, err := getProportionAmount(blockRewards, proportions.DeveloperRewards)
	if err != nil {
		return err
	}

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
	stakingRewardCoin, err := getProportionAmount(blockRewards, proportions.StakingRewards)
	if err != nil {
		return err
	}

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
	totalCoin sdk.Coin,
	fromAddress sdk.AccAddress,
) error {
	for _, w := range weightedAddresses {
		amount, err := getProportionAmount(totalCoin, w.Weight)
		if err != nil {
			return err
		}
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

func getProportionAmount(totalCoin sdk.Coin, ratio sdk.Dec) (sdk.Coin, error) {
	if ratio.GT(sdkmath.LegacyOneDec()) {
		return sdk.Coin{}, errors.New("ratio cannot be greater than 1")
	}
	return sdk.NewCoin(totalCoin.Denom, sdkmath.LegacyNewDecFromInt(totalCoin.Amount).Mul(ratio).TruncateInt()), nil
}
