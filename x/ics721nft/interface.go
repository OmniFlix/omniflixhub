package ics721nft

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

type (
	// AccountKeeper defines the contract required for account APIs.
	AccountKeeper interface {
		NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
		SetAccount(sdk.Context, authtypes.AccountI)
		GetModuleAddress(name string) sdk.AccAddress
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
