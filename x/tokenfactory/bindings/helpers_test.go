package bindings_test

import (
	"os"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/stretchr/testify/require"

	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/testutil"

	"github.com/OmniFlix/omniflixhub/v5/app"
)

func CreateTestInput(t *testing.T) (*app.OmniFlixApp, sdk.Context) {
	t.Helper()

	omniflix := app.Setup(t)
	ctx := omniflix.BaseApp.NewContext(false)
	return omniflix, ctx
}

func FundAccount(t *testing.T, ctx sdk.Context, customApp *app.OmniFlixApp, acct sdk.AccAddress) {
	t.Helper()

	err := banktestutil.FundAccount(ctx, customApp.AppKeepers.BankKeeper, acct, sdk.NewCoins(
		sdk.NewCoin("uflix", sdkmath.NewInt(10000000000)),
	))
	require.NoError(t, err)
}

// we need to make this deterministic (same every test run), as content might affect gas costs
func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := ed25519.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

func RandomAccountAddress() sdk.AccAddress {
	_, _, addr := keyPubAddr()
	return addr
}

func RandomBech32AccountAddress() string {
	return RandomAccountAddress().String()
}

func storeReflectCode(t *testing.T, ctx sdk.Context, customApp *app.OmniFlixApp, addr sdk.AccAddress) uint64 {
	t.Helper()

	wasmCode, err := os.ReadFile("./testdata/token_reflect.wasm")
	require.NoError(t, err)

	contractKeeper := keeper.NewDefaultPermissionKeeper(customApp.AppKeepers.WasmKeeper)
	codeID, _, err := contractKeeper.Create(ctx, addr, wasmCode, nil)
	require.NoError(t, err)

	return codeID
}

func instantiateReflectContract(t *testing.T, ctx sdk.Context, customApp *app.OmniFlixApp, funder sdk.AccAddress) sdk.AccAddress {
	t.Helper()

	initMsgBz := []byte("{}")
	contractKeeper := keeper.NewDefaultPermissionKeeper(customApp.AppKeepers.WasmKeeper)
	codeID := uint64(1)
	addr, _, err := contractKeeper.Instantiate(ctx, codeID, funder, funder, initMsgBz, "demo contract", nil)
	require.NoError(t, err)

	return addr
}

func fundAccount(t *testing.T, ctx sdk.Context, customApp *app.OmniFlixApp, addr sdk.AccAddress, coins sdk.Coins) {
	t.Helper()

	err := banktestutil.FundAccount(
		ctx,
		customApp.AppKeepers.BankKeeper,
		addr,
		coins,
	)
	require.NoError(t, err)
}

func SetupCustomApp(t *testing.T, addr sdk.AccAddress) (*app.OmniFlixApp, sdk.Context) {
	t.Helper()

	customApp, ctx := CreateTestInput(t)
	wasmKeeper := customApp.AppKeepers.WasmKeeper

	storeReflectCode(t, ctx, customApp, addr)

	cInfo := wasmKeeper.GetCodeInfo(ctx, 1)
	require.NotNil(t, cInfo)

	return customApp, ctx
}
