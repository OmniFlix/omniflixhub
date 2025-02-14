package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/OmniFlix/omniflixhub/v6/x/medianode/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

type (
	Keeper struct {
		cdc      codec.BinaryCodec
		storeKey storetypes.StoreKey

		accountKeeper      types.AccountKeeper
		bankKeeper         types.BankKeeper
		distributionKeeper types.DistributionKeeper

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	distributionKeeper types.DistributionKeeper,
	authority string,
) *Keeper {
	return &Keeper{
		cdc:                cdc,
		storeKey:           storeKey,
		accountKeeper:      accountKeeper,
		bankKeeper:         bankKeeper,
		distributionKeeper: distributionKeeper,
		authority:          authority,
	}
}

// GetAuthority returns the x/itc module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetModuleAccountAddress(ctx sdk.Context) sdk.AccAddress {
	return k.accountKeeper.GetModuleAddress(types.ModuleName)
}

// RegisterMediaNode creates a new media node entry
func (k Keeper) RegisterMediaNode(ctx sdk.Context, mediaNode types.MediaNode, depositAmount sdk.Coin) error {
	// Create a deposit object
	deposit := types.Deposit{
		Depositor:   mediaNode.Owner,
		DepositedAt: ctx.BlockTime(),
		Amount:      depositAmount,
	}

	mediaNodeid := k.GetNextMediaNodeNumber(ctx)
	mediaNode.Id = mediaNodeid
	mediaNode.RegisteredAt = ctx.BlockTime()
	// Update the media node's deposits
	mediaNode.Deposits = append(mediaNode.Deposits, &deposit)

	k.SetMediaNode(ctx, mediaNode)
	k.SetNextMediaNodeNumber(ctx, mediaNodeid+1)

	k.registerMediaNodeEvent(ctx, mediaNode.Owner, mediaNode.Id, mediaNode.Url, mediaNode.PricePerDay, mediaNode.Status)
	return nil
}

// UpdateMediaNode updates an existing media node
func (k Keeper) UpdateMediaNode(ctx sdk.Context, mediaNode types.MediaNode, owner sdk.AccAddress) error {
	existingNode, found := k.GetMediaNode(ctx, mediaNode.Id)
	if !found {
		return errorsmod.Wrapf(types.ErrMediaNodeDoesNotExist, "media node %d does not exist", mediaNode.Id)
	}

	existingOwner, err := sdk.AccAddressFromBech32(existingNode.Owner)
	if err != nil {
		return err
	}

	if !owner.Equals(existingOwner) {
		return errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "unauthorized address %s", owner.String())
	}

	k.SetMediaNode(ctx, mediaNode)
	k.updateMediaNodeEvent(ctx, mediaNode.Owner, mediaNode.Id)
	return nil
}

// LeaseMediaNode creates a new lease for a media node
func (k Keeper) LeaseMediaNode(ctx sdk.Context, mediaNode types.MediaNode, leaseDays uint64, lessee sdk.AccAddress, leaseAmount sdk.Coin) error {
	// Create a new lease object
	lease := types.Lease{
		MediaNodeId:      mediaNode.Id,
		Leasee:           lessee.String(),
		LeasedDays:       leaseDays,
		PricePerDay:      mediaNode.PricePerDay,
		StartTime:        ctx.BlockTime(),
		Status:           types.LEASE_STATUS_ACTIVE,
		Owner:            mediaNode.Owner,
		TotalLeaseAmount: leaseAmount,
	}

	// Transfer tokens from lessee to module account
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		lessee,
		types.ModuleName,
		sdk.NewCoins(leaseAmount),
	); err != nil {
		return err
	}

	lease.LastSettledAt = ctx.BlockTime()
	lease.SettledAmount = sdk.NewCoin(leaseAmount.Denom, sdkmath.ZeroInt())
	// Set the lease using the SetLease method
	k.SetLease(ctx, lease)

	// Update media node lease details
	mediaNode.Leased = true
	k.SetMediaNode(ctx, mediaNode)

	k.leaseMediaNodeEvent(ctx, lease.Leasee, lease.MediaNodeId, leaseAmount)

	return nil
}

// DepositMediaNode allows a user to deposit a media node
func (k Keeper) DepositMediaNode(ctx sdk.Context, mediaNodeId uint64, amount sdk.Coin, sender sdk.AccAddress) error {
	mediaNode, found := k.GetMediaNode(ctx, mediaNodeId)
	if !found {
		return errorsmod.Wrapf(types.ErrMediaNodeDoesNotExist, "media node %d does not exist", mediaNodeId)
	}

	// Allow deposit only if the media node status is PENDING
	if mediaNode.Status != types.STATUS_PENDING {
		return errorsmod.Wrapf(types.ErrInvalidMediaNodeStatus, "media node %d is not in PENDING status", mediaNodeId)
	}

	minDeposit := k.GetMinDeposit(ctx)
	initialDepositPerc := k.GetInitialDepositPercentage(ctx)
	minInitialDeposit := sdk.NewCoin(minDeposit.Denom, sdkmath.LegacyNewDecFromInt(minDeposit.Amount).Mul(initialDepositPerc).TruncateInt())
	if !amount.IsGTE(minInitialDeposit) {
		return errorsmod.Wrapf(types.ErrInsufficientDeposit, "%s of deposit is required", minInitialDeposit.String())
	}

	// Create a deposit object
	deposit := types.Deposit{
		Depositor:   sender.String(),
		DepositedAt: ctx.BlockTime(),
		Amount:      amount,
	}

	// Update the media node's deposits
	mediaNode.Deposits = append(mediaNode.Deposits, &deposit)

	// Calculate total deposits
	totalDeposits := sdk.NewCoin(mediaNode.PricePerDay.Denom, sdkmath.ZeroInt())
	for _, d := range mediaNode.Deposits {
		totalDeposits = totalDeposits.Add(d.Amount)
	}

	// Check if total deposits meet or exceed the required minimum deposit
	if totalDeposits.IsGTE(minDeposit) {
		mediaNode.Status = types.STATUS_ACTIVE // Change status to ACTIVE
	}

	k.SetMediaNode(ctx, mediaNode)
	k.depositMediaNodeEvent(ctx, deposit.Depositor, mediaNodeId, deposit.Amount)

	return nil
}

// CancelLease cancels an existing lease for a media node
func (k Keeper) CancelLease(ctx sdk.Context, mediaNodeId uint64, sender sdk.AccAddress) error {
	mediaNode, found := k.GetMediaNode(ctx, mediaNodeId)
	if !found {
		return errorsmod.Wrapf(types.ErrMediaNodeDoesNotExist, "media node %d does not exist", mediaNodeId)
	}
	lease, found := k.GetMediaNodeLease(ctx, mediaNodeId)
	if !found {
		return errorsmod.Wrapf(types.ErrLeaseNotFound, "lease for media node %d does not exist", mediaNodeId)
	}

	if lease.Status != types.LEASE_STATUS_ACTIVE {
		return errorsmod.Wrapf(types.ErrLeaseNotActive, "lease for media node %d is not active", mediaNodeId)
	}

	if sender.String() != lease.Leasee {
		return errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "unauthorized address %s", sender.String())
	}

	// Calculate remaining lease days and refund amount //TODO: change to 24
	remainingDays := lease.LeasedDays - uint64(ctx.BlockTime().Sub(lease.StartTime).Minutes()/2)
	if remainingDays > 0 {
		refundAmount := sdk.NewCoin(
			mediaNode.PricePerDay.Denom,
			mediaNode.PricePerDay.Amount.MulRaw(int64(remainingDays)),
		)

		// Return remaining funds to lessee
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(
			ctx,
			types.ModuleName,
			sender,
			sdk.NewCoins(refundAmount),
		); err != nil {
			return err
		}
	}

	// Clear lease
	mediaNode.Leased = false
	k.SetMediaNode(ctx, mediaNode)
	k.RemoveLease(ctx, lease.MediaNodeId)

	k.cancelLeaseMediaNodeEvent(ctx, lease.Leasee, mediaNodeId)
	return nil
}

// CloseMediaNode closes an existing media node
func (k Keeper) CloseMediaNode(ctx sdk.Context, mediaNodeId uint64, owner sdk.AccAddress) error {
	mediaNode, found := k.GetMediaNode(ctx, mediaNodeId)
	if !found {
		return errorsmod.Wrapf(types.ErrMediaNodeDoesNotExist, "media node %d does not exist", mediaNodeId)
	}

	existingOwner, err := sdk.AccAddressFromBech32(mediaNode.Owner)
	if err != nil {
		return err
	}

	if !owner.Equals(existingOwner) {
		return errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "unauthorized address %s", owner.String())
	}

	// Update media node status to CLOSED
	mediaNode.Status = types.STATUS_CLOSED
	mediaNode.ClosedAt = ctx.BlockTime()
	k.SetMediaNode(ctx, mediaNode)

	k.closeMediaNodeEvent(ctx, mediaNode.Owner, mediaNodeId)
	return nil
}

// SettleActiveLeases iterates through all active leases and settles payment if 24 hours have passed
func (k Keeper) SettleActiveLeases(ctx sdk.Context) error {
	leaseCommissionPercentage := k.GetLeaseCommission(ctx)
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.PrefixLease)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var lease types.Lease
		k.cdc.MustUnmarshal(iterator.Value(), &lease)
		if lease.Status == types.LEASE_STATUS_ACTIVE &&
			ctx.BlockTime().Sub(lease.StartTime).Minutes() >= 2 &&
			ctx.BlockTime().Sub(lease.LastSettledAt).Minutes() >= 2 {

			// Calculate payment amount
			paymentAmount := sdk.NewCoin(
				lease.PricePerDay.Denom,
				lease.PricePerDay.Amount,
			)

			owner, err := sdk.AccAddressFromBech32(lease.Owner)
			if err != nil {
				return err
			}

			remainingAmount := lease.TotalLeaseAmount.Sub(lease.SettledAmount)

			if remainingAmount.Amount.GT(paymentAmount.Amount) {

				leaseCommissionCoin := k.GetProportions(paymentAmount, leaseCommissionPercentage)

				// distribute lease commission
				err := k.DistributeLeaseCommission(ctx, leaseCommissionCoin)
				if err != nil {
					return err
				}

				// Transfer remaining payment to media node owner
				mediaNodeOwnerAmount := paymentAmount.Sub(leaseCommissionCoin)

				if err := k.bankKeeper.SendCoinsFromModuleToAccount(
					ctx,
					types.ModuleName,
					owner,
					sdk.NewCoins(mediaNodeOwnerAmount),
				); err != nil {
					return err
				}
				k.createLeasePaymentTransferEvent(ctx, k.accountKeeper.GetModuleAddress(types.ModuleName), owner, mediaNodeOwnerAmount)
				// Update last settled time
				lease.SettledAmount = lease.SettledAmount.AddAmount(paymentAmount.Amount)
				lease.LastSettledAt = ctx.BlockTime()
				k.SetLease(ctx, lease)
			} else {

				leaseCommissionCoin := k.GetProportions(remainingAmount, leaseCommissionPercentage)

				// distribute lease commission
				err := k.DistributeLeaseCommission(ctx, leaseCommissionCoin)
				if err != nil {
					return err
				}

				// Transfer remaining amount to the medianode owner
				mediaNodeOwnerAmount := remainingAmount.Sub(leaseCommissionCoin)

				if err := k.bankKeeper.SendCoinsFromModuleToAccount(
					ctx,
					types.ModuleName,
					owner,
					sdk.NewCoins(mediaNodeOwnerAmount),
				); err != nil {
					return err
				}

				k.createLeasePaymentTransferEvent(ctx, k.accountKeeper.GetModuleAddress(types.ModuleName), owner, mediaNodeOwnerAmount)

				mediaNode, found := k.GetMediaNode(ctx, lease.MediaNodeId)
				if found {
					mediaNode.Leased = false
					k.SetMediaNode(ctx, mediaNode)
				}
				k.RemoveLease(ctx, lease.MediaNodeId)
			}
		}
	}

	return nil
}

// ReleaseDeposits iterates through CLOSED media nodes and returns deposits to depositors
func (k Keeper) ReleaseDeposits(ctx sdk.Context) error {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.PrefixMediaNode)

	defer iterator.Close()

	// Retrieve the deposit release period from params
	params := k.GetParams(ctx) // Assuming this method exists to get the current parameters
	depositReleasePeriod := params.DepositReleasePeriod

	for ; iterator.Valid(); iterator.Next() {
		var mediaNode types.MediaNode
		k.cdc.MustUnmarshal(iterator.Value(), &mediaNode)

		// Check if the media node is CLOSED
		if mediaNode.Status == types.STATUS_CLOSED && len(mediaNode.Deposits) > 0 {
			// Check if the time since closed exceeds the deposit release period
			if ctx.BlockTime().Sub(mediaNode.ClosedAt) > depositReleasePeriod {
				k.Logger(ctx).Info("Releasing Deposits ..")
				for _, deposit := range mediaNode.Deposits {
					// Return deposit to the depositor
					depositorAddr, err := sdk.AccAddressFromBech32(deposit.Depositor)
					if err != nil {
						return err
					}

					if err := k.bankKeeper.SendCoinsFromModuleToAccount(
						ctx,
						types.ModuleName,
						depositorAddr,
						sdk.NewCoins(deposit.Amount),
					); err != nil {
						return err
					}
					k.createMediaNodeDepositRefundEvent(ctx, k.accountKeeper.GetModuleAddress(types.ModuleName), depositorAddr, deposit.Amount)
				}
				mediaNode.Deposits = []*types.Deposit{}
				k.SetMediaNode(ctx, mediaNode)
			}
		}
	}

	return nil
}

func (k Keeper) GetProportions(totalCoin sdk.Coin, ratio sdkmath.LegacyDec) sdk.Coin {
	return sdk.NewCoin(totalCoin.Denom, sdkmath.LegacyNewDecFromInt(totalCoin.Amount).Mul(ratio).TruncateInt())
}

func (k Keeper) DistributeLeaseCommission(ctx sdk.Context, leaseCommissionCoin sdk.Coin) error {
	distrParams := k.GetCommissionDistribution(ctx)
	stakingCommissionCoin := k.GetProportions(leaseCommissionCoin, distrParams.Staking)
	moduleAccAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	feeCollectorAddr := k.accountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	if distrParams.Staking.GT(sdkmath.LegacyZeroDec()) && stakingCommissionCoin.Amount.GT(sdkmath.ZeroInt()) {
		err := k.bankKeeper.SendCoins(ctx, moduleAccAddr, feeCollectorAddr, sdk.NewCoins(stakingCommissionCoin))
		if err != nil {
			return err
		}
		k.createLeaseCommissionTransferEvent(ctx,
			moduleAccAddr,
			feeCollectorAddr,
			stakingCommissionCoin,
		)
		leaseCommissionCoin = leaseCommissionCoin.Sub(stakingCommissionCoin)
	}
	communityPoolCommissionCoin := leaseCommissionCoin

	err := k.distributionKeeper.FundCommunityPool(
		ctx,
		sdk.NewCoins(communityPoolCommissionCoin),
		moduleAccAddr,
	)
	if err != nil {
		return err
	}
	k.createLeaseCommissionTransferEvent(ctx,
		moduleAccAddr,
		k.accountKeeper.GetModuleAddress("distribution"),
		communityPoolCommissionCoin,
	)

	return nil
}
