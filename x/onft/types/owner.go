package types

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewIDCollection creates a new IDCollection instance
func NewIDCollection(denomId string, onftIds []string) IDCollection {
	return IDCollection{
		DenomId: denomId,
		OnftIds: onftIds,
	}
}

// Supply return the amount of the denom
func (idc IDCollection) Supply() int {
	return len(idc.OnftIds)
}

// AddID adds an nftID to the idCollection
func (idc IDCollection) AddId(onftId string) IDCollection {
	idc.OnftIds = append(idc.OnftIds, onftId)
	return idc
}

// ----------------------------------------------------------------------------
// IDCollections is an array of ID Collections
type IDCollections []IDCollection

// Add adds an ID to the idCollection
func (idcs IDCollections) Add(denomId, onftId string) IDCollections {
	for i, idc := range idcs {
		if idc.DenomId == denomId {
			idcs[i] = idc.AddId(onftId)
			return idcs
		}
	}
	return append(idcs, IDCollection{
		DenomId: denomId,
		OnftIds: []string{onftId},
	})
}

// String follows stringer interface
func (idcs IDCollections) String() string {
	if len(idcs) == 0 {
		return ""
	}

	var buf bytes.Buffer
	for _, idCollection := range idcs {
		if buf.Len() > 0 {
			buf.WriteString("\n")
		}
		buf.WriteString(idCollection.String())
	}
	return buf.String()
}

// NewOwner creates a new Owner
func NewOwner(owner sdk.AccAddress, idCollections ...IDCollection) Owner {
	return Owner{
		Address:       owner.String(),
		IDCollections: idCollections,
	}
}

type Owners []Owner

// NewOwners creates a new Owner
func NewOwners(owner ...Owner) Owners {
	return append([]Owner{}, owner...)
}

// String follows stringer interface
func (owners Owners) String() string {
	var buf bytes.Buffer
	for _, owner := range owners {
		if buf.Len() > 0 {
			buf.WriteString("\n")
		}
		buf.WriteString(owner.String())
	}
	return buf.String()
}
