package ics721nft

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	// AccountKeeper defines the contract required for account APIs.
	AccountKeeper interface {
		NewAccountWithAddress(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
		SetAccount(context.Context, sdk.AccountI)
		GetModuleAddress(name string) sdk.AccAddress
	}
	BankKeeper interface {
		BlockedAddr(addr sdk.AccAddress) bool
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
