package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QuerySupply     = "supply"
	QueryOwner      = "owner"
	QueryCollection = "collection"
	QueryDenoms     = "denoms"
	QueryDenom      = "denom"
	QueryONFT       = "onft"
)

type QuerySupplyParams struct {
	Denom string
	Owner sdk.AccAddress
}

func NewQuerySupplyParams(denom string, owner sdk.AccAddress) QuerySupplyParams {
	return QuerySupplyParams{
		Denom: denom,
		Owner: owner,
	}
}

func (q QuerySupplyParams) Bytes() []byte {
	return []byte(q.Denom)
}

type QueryOwnerParams struct {
	Denom string
	Owner sdk.AccAddress
}

func NewQueryOwnerParams(denom string, owner sdk.AccAddress) QueryOwnerParams {
	return QueryOwnerParams{
		Denom: denom,
		Owner: owner,
	}
}

type QueryCollectionParams struct {
	Denom string
}

func NewQueryCollectionParams(denom string) QueryCollectionParams {
	return QueryCollectionParams{
		Denom: denom,
	}
}

type QueryDenomParams struct {
	ID string
}

func NewQueryDenomParams(id string) QueryDenomParams {
	return QueryDenomParams{
		ID: id,
	}
}

type QueryONFTParams struct {
	Denom  string
	ONFTID string
}

func NewQueryONFTParams(denom, id string) QueryONFTParams {
	return QueryONFTParams{
		Denom:  denom,
		ONFTID: id,
	}
}
