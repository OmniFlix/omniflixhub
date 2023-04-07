package decorators

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

var MinimumInitialDepositRate = sdk.NewDecWithPrec(10, 2) // 10% of min_deposit

type MinimumInitialDepositDecorator struct {
	govKeeper govkeeper.Keeper
	cdc       codec.BinaryCodec
}

func NewMinimumInitialDepositDecorator(
	cdc codec.BinaryCodec, govKeeper govkeeper.Keeper,
) MinimumInitialDepositDecorator {
	return MinimumInitialDepositDecorator{
		govKeeper: govKeeper,
		cdc:       cdc,
	}
}

// checkProposalInitialDeposit returns error if the initial deposit amount is less than minInitialDeposit amount
func (midd MinimumInitialDepositDecorator) checkProposalInitialDeposit(ctx sdk.Context, msg sdk.Msg) error {
	if msg, ok := msg.(*govtypes.MsgSubmitProposal); ok {
		depositParams := midd.govKeeper.GetDepositParams(ctx)
		minimumInitialDeposit := midd.calculateMinimumInitialDeposit(depositParams.MinDeposit)
		if msg.InitialDeposit.IsAllLT(minimumInitialDeposit) {
			return errorsmod.Wrapf(sdkerrors.ErrUnauthorized,
				"initial deposit amount is not enough. required: %v", minimumInitialDeposit)
		}
	}
	return nil
}

// Validate validates messages
func (midd MinimumInitialDepositDecorator) Validate(ctx sdk.Context, msgs []sdk.Msg) error {
	// Check every msg in the tx, if it's a MsgExec, check the inner msgs.
	// If it's a MsgSubmitProposal, check the initial deposit is enough.
	for _, m := range msgs {
		var innerMsg sdk.Msg
		if msg, ok := m.(*authz.MsgExec); ok {
			for _, v := range msg.Msgs {
				err := midd.cdc.UnpackAny(v, &innerMsg)
				if err != nil {
					return errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "cannot unmarshal authz exec msgs")
				}

				err = midd.checkProposalInitialDeposit(ctx, innerMsg)
				if err != nil {
					return err
				}
			}
		} else {
			err := midd.checkProposalInitialDeposit(ctx, m)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (midd MinimumInitialDepositDecorator) AnteHandle(
	ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {
	msgs := tx.GetMsgs()
	if err := midd.Validate(ctx, msgs); err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate)
}

func (midd MinimumInitialDepositDecorator) calculateMinimumInitialDeposit(
	minDeposit sdk.Coins,
) (minimumInitialDeposit sdk.Coins) {
	for _, coin := range minDeposit {
		minimumInitialCoin := MinimumInitialDepositRate.MulInt(coin.Amount).RoundInt()
		minimumInitialDeposit = minimumInitialDeposit.Add(sdk.NewCoin(coin.Denom, minimumInitialCoin))
	}

	return
}
