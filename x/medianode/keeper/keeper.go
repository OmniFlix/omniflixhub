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
)

type (
	Keeper struct {
		cdc      codec.BinaryCodec
		storeKey storetypes.StoreKey

		accountKeeper types.AccountKeeper
		bankKeeper    types.BankKeeper

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
) *Keeper {
	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
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
func (k Keeper) RegisterMediaNode(ctx sdk.Context, mediaNode types.MediaNode) error {

	k.SetMediaNode(ctx, mediaNode)
	k.SetNextMediaNodeNumber(ctx, mediaNode.Id+1)
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

	return nil
}

// LeaseMediaNode creates a new lease for a media node
func (k Keeper) LeaseMediaNode(ctx sdk.Context, mediaNodeId uint64, leaseDays uint64, lessee sdk.AccAddress) error {
	mediaNode, found := k.GetMediaNode(ctx, mediaNodeId)
	if !found {
		return errorsmod.Wrapf(types.ErrMediaNodeDoesNotExist, "media node %d does not exist", mediaNodeId)
	}

	if mediaNode.IsLeased() {
		return errorsmod.Wrapf(types.ErrMediaNodeAlreadyLeased, "media node %d is already leased", mediaNodeId)
	}

	// Calculate total lease amount
	totalLeaseAmount := sdk.NewCoin(
		mediaNode.PricePerDay.Denom,
		mediaNode.PricePerDay.Amount.MulRaw(int64(leaseDays)),
	)

	// Transfer tokens from lessee to module account
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		lessee,
		types.ModuleName,
		sdk.NewCoins(totalLeaseAmount),
	); err != nil {
		return err
	}

	// Update media node lease details
	mediaNode.Leased = true
	k.SetMediaNode(ctx, mediaNode)

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
	minDeposit := k.GetMinDeposit(ctx) // Assuming this method exists
	if totalDeposits.IsGTE(minDeposit) {
		mediaNode.Status = types.STATUS_ACTIVE // Change status to ACTIVE
	}

	k.SetMediaNode(ctx, mediaNode)

	return nil
}
