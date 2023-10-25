package v2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NFTKeeper save the denom of class
type NFTKeeper interface {
	SaveDenom(
		ctx sdk.Context,
		id,
		symbol,
		name,
		schema string,
		creator sdk.AccAddress,
		description,
		previewUri string,
		uri,
		uriHash,
		data string,
	) error
}
