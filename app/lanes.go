package app

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	signerextractionadapter "github.com/skip-mev/block-sdk/v2/adapters/signer_extraction_adapter"
	blocksdkbase "github.com/skip-mev/block-sdk/v2/block/base"
	baselane "github.com/skip-mev/block-sdk/v2/lanes/base"
)

const (
	MaxTxsForDefaultLane = 3000 // maximal number of txs that can be stored in this lane at any point in time
)

var MaxBlockSpaceForDefaultLane = sdkmath.LegacyMustNewDecFromStr("1") // maximal fraction of blockMaxBytes / gas that can be used by this lane at any point in time (100%)

// CreateLanes creates a LaneMempool containing MEV, default lanes (in that order)
func CreateLanes(app *OmniFlixApp, txConfig client.TxConfig) *blocksdkbase.BaseLane {
	// initialize the default lane
	baseCfg := blocksdkbase.LaneConfig{
		Logger:          app.Logger(),
		TxDecoder:       txConfig.TxDecoder(),
		TxEncoder:       txConfig.TxEncoder(),
		SignerExtractor: signerextractionadapter.NewDefaultAdapter(),
		MaxBlockSpace:   MaxBlockSpaceForDefaultLane,
		MaxTxs:          MaxTxsForDefaultLane,
	}

	// BaseLane (DefaultLane) is intended to hold all txs that are not matched by any lanes ordered before this lane
	baseLane := baselane.NewDefaultLane(baseCfg, blocksdkbase.DefaultMatchHandler())
	baseLane.LaneMempool = blocksdkbase.NewMempool(
		blocksdkbase.NewDefaultTxPriority(),
		baseCfg.SignerExtractor,
		baseCfg.MaxTxs,
	)

	return baseLane
}
