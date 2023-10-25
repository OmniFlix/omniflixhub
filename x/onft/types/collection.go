package types

import (
	"github.com/OmniFlix/omniflixhub/v2/x/onft/exported"
)

func NewCollection(denom Denom, onfts []exported.ONFTI) (c Collection) {
	c.Denom = denom
	for _, onft := range onfts {
		c = c.AddONFT(onft.(ONFT))
	}
	return c
}

func (c Collection) AddONFT(nft ONFT) Collection {
	c.ONFTs = append(c.ONFTs, nft)
	return c
}

func (c Collection) Supply() int {
	return len(c.ONFTs)
}

func NewCollections(c ...Collection) []Collection {
	return append([]Collection{}, c...)
}
