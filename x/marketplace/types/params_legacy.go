package types

import paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

// Parameter keys
var (
	ParamStoreKeySaleCommission     = []byte("SaleCommission")
	ParamStoreKeyDistribution       = []byte("MarketplaceDistribution")
	ParamStoreKeyBidCloseDuration   = []byte("BidCloseDuration")
	ParamStoreKeyMaxAuctionDuration = []byte("MaxAuctionDuration")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeySaleCommission, &p.SaleCommission, validateSaleCommission),
		paramtypes.NewParamSetPair(ParamStoreKeyDistribution, &p.Distribution, validateMarketplaceDistributionParams),
		paramtypes.NewParamSetPair(ParamStoreKeyBidCloseDuration, &p.BidCloseDuration, validateBidCloseDuration),
		paramtypes.NewParamSetPair(ParamStoreKeyMaxAuctionDuration, &p.MaxAuctionDuration, validateMaxAuctionDuration),
	}
}
