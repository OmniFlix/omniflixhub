package ics721nft

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	nftkeeper "github.com/cosmos/cosmos-sdk/x/nft/keeper"
)

type (
	// AccountKeeper defines the contract required for account APIs.
	AccountKeeper interface {
		NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
		// Set an account in the store.
		SetAccount(sdk.Context, authtypes.AccountI)
		GetModuleAddress(name string) sdk.AccAddress
	}
	// ICS721NftKeeper defines the ICS721 Keeper
	ICS721NftKeeper struct {
		nk  nftkeeper.Keeper
		cdc codec.Codec
		ak  AccountKeeper
	}

	ICS721Class struct {
		ID   string
		URI  string
		Data string
	}

	ICS721Token struct {
		ClassID string
		ID      string
		URI     string
		Data    string
	}
)

func (c ICS721Class) GetID() string      { return c.ID }
func (c ICS721Class) GetURI() string     { return c.URI }
func (c ICS721Class) GetData() string    { return c.Data }
func (t ICS721Token) GetClassID() string { return t.ClassID }
func (t ICS721Token) GetID() string      { return t.ID }
func (t ICS721Token) GetURI() string     { return t.URI }
func (t ICS721Token) GetData() string    { return t.Data }
