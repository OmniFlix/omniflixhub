package bindings_test

import (
	"encoding/json"
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"

	"github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/v2/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/OmniFlix/omniflixhub/v6/app"
	bindings "github.com/OmniFlix/omniflixhub/v6/x/tokenfactory/bindings/types"
	"github.com/OmniFlix/omniflixhub/v6/x/tokenfactory/types"
)

func TestCreateDenomMsg(t *testing.T) {
	creator := RandomAccountAddress()
	customApp, ctx := SetupCustomApp(t, creator)

	lucky := RandomAccountAddress()
	reflect := instantiateReflectContract(t, ctx, customApp, lucky)
	require.NotEmpty(t, reflect)

	// Fund reflect contract with 100 base denom creation fees
	reflectAmount := sdk.NewCoins(sdk.NewCoin(types.DefaultParams().DenomCreationFee[0].Denom, types.DefaultParams().DenomCreationFee[0].Amount.MulRaw(100)))
	fundAccount(t, ctx, customApp, reflect, reflectAmount)

	msg := bindings.TokenFactoryMsg{CreateDenom: &bindings.CreateDenom{
		Subdenom: "SUN",
	}}
	err := executeCustom(t, ctx, customApp, reflect, lucky, msg, sdk.Coin{})
	require.NoError(t, err)

	// query the denom and see if it matches
	query := bindings.TokenFactoryQuery{
		FullDenom: &bindings.FullDenom{
			CreatorAddr: reflect.String(),
			Subdenom:    "SUN",
		},
	}
	resp := bindings.FullDenomResponse{}
	queryCustom(t, ctx, customApp, reflect, query, &resp)

	require.Equal(t, resp.Denom, fmt.Sprintf("factory/%s/SUN", reflect.String()))
}

func TestMintMsg(t *testing.T) {
	creator := RandomAccountAddress()
	customApp, ctx := SetupCustomApp(t, creator)

	lucky := RandomAccountAddress()
	reflect := instantiateReflectContract(t, ctx, customApp, lucky)
	require.NotEmpty(t, reflect)

	// Fund reflect contract with 100 base denom creation fees
	reflectAmount := sdk.NewCoins(sdk.NewCoin(types.DefaultParams().DenomCreationFee[0].Denom, types.DefaultParams().DenomCreationFee[0].Amount.MulRaw(100)))
	fundAccount(t, ctx, customApp, reflect, reflectAmount)

	// lucky was broke
	balances := customApp.AppKeepers.BankKeeper.GetAllBalances(ctx, lucky)
	require.Empty(t, balances)

	// Create denom for minting
	msg := bindings.TokenFactoryMsg{CreateDenom: &bindings.CreateDenom{
		Subdenom: "SUN",
	}}
	err := executeCustom(t, ctx, customApp, reflect, lucky, msg, sdk.Coin{})
	require.NoError(t, err)
	sunDenom := fmt.Sprintf("factory/%s/%s", reflect.String(), msg.CreateDenom.Subdenom)

	amount, ok := sdkmath.NewIntFromString("808010808")
	require.True(t, ok)
	msg = bindings.TokenFactoryMsg{MintTokens: &bindings.MintTokens{
		Denom:         sunDenom,
		Amount:        amount,
		MintToAddress: lucky.String(),
	}}
	err = executeCustom(t, ctx, customApp, reflect, lucky, msg, sdk.Coin{})
	require.NoError(t, err)

	balances = customApp.AppKeepers.BankKeeper.GetAllBalances(ctx, lucky)
	require.Len(t, balances, 1)
	coin := balances[0]
	require.Equal(t, amount, coin.Amount)
	require.Contains(t, coin.Denom, "factory/")

	// query the denom and see if it matches
	query := bindings.TokenFactoryQuery{
		FullDenom: &bindings.FullDenom{
			CreatorAddr: reflect.String(),
			Subdenom:    "SUN",
		},
	}
	resp := bindings.FullDenomResponse{}
	queryCustom(t, ctx, customApp, reflect, query, &resp)

	require.Equal(t, resp.Denom, coin.Denom)

	// mint the same denom again
	err = executeCustom(t, ctx, customApp, reflect, lucky, msg, sdk.Coin{})
	require.NoError(t, err)

	balances = customApp.AppKeepers.BankKeeper.GetAllBalances(ctx, lucky)
	require.Len(t, balances, 1)
	coin = balances[0]
	require.Equal(t, amount.MulRaw(2), coin.Amount)
	require.Contains(t, coin.Denom, "factory/")

	// query the denom and see if it matches
	query = bindings.TokenFactoryQuery{
		FullDenom: &bindings.FullDenom{
			CreatorAddr: reflect.String(),
			Subdenom:    "SUN",
		},
	}
	resp = bindings.FullDenomResponse{}
	queryCustom(t, ctx, customApp, reflect, query, &resp)

	require.Equal(t, resp.Denom, coin.Denom)

	// now mint another amount / denom
	// create it first
	msg = bindings.TokenFactoryMsg{CreateDenom: &bindings.CreateDenom{
		Subdenom: "MOON",
	}}
	err = executeCustom(t, ctx, customApp, reflect, lucky, msg, sdk.Coin{})
	require.NoError(t, err)
	moonDenom := fmt.Sprintf("factory/%s/%s", reflect.String(), msg.CreateDenom.Subdenom)

	amount = amount.SubRaw(1)
	msg = bindings.TokenFactoryMsg{MintTokens: &bindings.MintTokens{
		Denom:         moonDenom,
		Amount:        amount,
		MintToAddress: lucky.String(),
	}}
	err = executeCustom(t, ctx, customApp, reflect, lucky, msg, sdk.Coin{})
	require.NoError(t, err)

	balances = customApp.AppKeepers.BankKeeper.GetAllBalances(ctx, lucky)
	require.Len(t, balances, 2)
	coin = balances[0]
	require.Equal(t, amount, coin.Amount)
	require.Contains(t, coin.Denom, "factory/")

	// query the denom and see if it matches
	query = bindings.TokenFactoryQuery{
		FullDenom: &bindings.FullDenom{
			CreatorAddr: reflect.String(),
			Subdenom:    "MOON",
		},
	}
	resp = bindings.FullDenomResponse{}
	queryCustom(t, ctx, customApp, reflect, query, &resp)

	require.Equal(t, resp.Denom, coin.Denom)

	// and check the first denom is unchanged
	coin = balances[1]
	require.Equal(t, amount.AddRaw(1).MulRaw(2), coin.Amount)
	require.Contains(t, coin.Denom, "factory/")

	// query the denom and see if it matches
	query = bindings.TokenFactoryQuery{
		FullDenom: &bindings.FullDenom{
			CreatorAddr: reflect.String(),
			Subdenom:    "SUN",
		},
	}
	resp = bindings.FullDenomResponse{}
	queryCustom(t, ctx, customApp, reflect, query, &resp)

	require.Equal(t, resp.Denom, coin.Denom)
}

func TestBurnMsg(t *testing.T) {
	creator := RandomAccountAddress()
	customApp, ctx := SetupCustomApp(t, creator)

	lucky := RandomAccountAddress()
	reflect := instantiateReflectContract(t, ctx, customApp, lucky)
	require.NotEmpty(t, reflect)

	// Fund reflect contract with 100 base denom creation fees
	reflectAmount := sdk.NewCoins(sdk.NewCoin(types.DefaultParams().DenomCreationFee[0].Denom, types.DefaultParams().DenomCreationFee[0].Amount.MulRaw(100)))
	fundAccount(t, ctx, customApp, reflect, reflectAmount)

	// lucky was broke
	balances := customApp.AppKeepers.BankKeeper.GetAllBalances(ctx, lucky)
	require.Empty(t, balances)

	// Create denom for minting
	msg := bindings.TokenFactoryMsg{CreateDenom: &bindings.CreateDenom{
		Subdenom: "SUN",
	}}
	err := executeCustom(t, ctx, customApp, reflect, lucky, msg, sdk.Coin{})
	require.NoError(t, err)
	sunDenom := fmt.Sprintf("factory/%s/%s", reflect.String(), msg.CreateDenom.Subdenom)

	amount, ok := sdkmath.NewIntFromString("808010809")
	require.True(t, ok)

	msg = bindings.TokenFactoryMsg{MintTokens: &bindings.MintTokens{
		Denom:         sunDenom,
		Amount:        amount,
		MintToAddress: lucky.String(),
	}}
	err = executeCustom(t, ctx, customApp, reflect, lucky, msg, sdk.Coin{})
	require.NoError(t, err)

	// can burn from different address with burnFrom
	amt, ok := sdkmath.NewIntFromString("1")
	require.True(t, ok)
	msg = bindings.TokenFactoryMsg{BurnTokens: &bindings.BurnTokens{
		Denom:           sunDenom,
		Amount:          amt,
		BurnFromAddress: lucky.String(),
	}}
	err = executeCustom(t, ctx, customApp, reflect, lucky, msg, sdk.Coin{})
	require.NoError(t, err)

	// lucky needs to send balance to reflect contract to burn it
	luckyBalance := customApp.AppKeepers.BankKeeper.GetAllBalances(ctx, lucky)
	err = customApp.AppKeepers.BankKeeper.SendCoins(ctx, lucky, reflect, luckyBalance)
	require.NoError(t, err)

	msg = bindings.TokenFactoryMsg{BurnTokens: &bindings.BurnTokens{
		Denom:           sunDenom,
		Amount:          amount.Abs().Sub(sdkmath.NewInt(1)),
		BurnFromAddress: reflect.String(),
	}}
	err = executeCustom(t, ctx, customApp, reflect, lucky, msg, sdk.Coin{})
	require.NoError(t, err)
}

type ReflectExec struct {
	ReflectMsg    *ReflectMsgs    `json:"reflect_msg,omitempty"`
	ReflectSubMsg *ReflectSubMsgs `json:"reflect_sub_msg,omitempty"`
}

type ReflectMsgs struct {
	Msgs []wasmvmtypes.CosmosMsg `json:"msgs"`
}

type ReflectSubMsgs struct {
	Msgs []wasmvmtypes.SubMsg `json:"msgs"`
}

func executeCustom(t *testing.T, ctx sdk.Context, customApp *app.OmniFlixApp, contract sdk.AccAddress, sender sdk.AccAddress, msg bindings.TokenFactoryMsg, funds sdk.Coin) error { //nolint:unparam // funds is always nil but could change in the future.
	t.Helper()

	customBz, err := json.Marshal(msg)
	require.NoError(t, err)

	reflectMsg := ReflectExec{
		ReflectMsg: &ReflectMsgs{
			Msgs: []wasmvmtypes.CosmosMsg{{
				Custom: customBz,
			}},
		},
	}
	reflectBz, err := json.Marshal(reflectMsg)
	require.NoError(t, err)

	// no funds sent if amount is 0
	var coins sdk.Coins
	if !funds.Amount.IsNil() {
		coins = sdk.Coins{funds}
	}

	contractKeeper := keeper.NewDefaultPermissionKeeper(customApp.AppKeepers.WasmKeeper)
	_, err = contractKeeper.Execute(ctx, contract, sender, reflectBz, coins)
	return err
}
