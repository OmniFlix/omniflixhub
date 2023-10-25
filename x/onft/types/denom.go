package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewDenom(id, symbol, name, schema string, creator sdk.AccAddress, description, previewURI string) Denom {
	return Denom{
		Id:          id,
		Symbol:      symbol,
		Name:        name,
		Schema:      schema,
		Creator:     creator.String(),
		Description: description,
		PreviewURI:  previewURI,
	}
}
