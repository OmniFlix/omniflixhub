package cmd

import (
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	tmcfg "github.com/tendermint/tendermint/config"
)

type AppConfig struct {
	serverconfig.Config
}

func initAppConfig() (string, interface{}) {
	srvCfg := serverconfig.DefaultConfig()

	srvCfg.MinGasPrices = "0.001uflix"
	srvCfg.IAVLDisableFastNode = false

	simAppConfig := AppConfig{
		Config: *srvCfg,
	}

	simAppTemplate := serverconfig.DefaultConfigTemplate

	return simAppTemplate, simAppConfig
}

func initTMConfig() *tmcfg.Config {
	cfg := tmcfg.DefaultConfig()

	// to put a higher strain on node memory, use these values:
	// cfg.P2P.MaxNumInboundPeers = 100
	// cfg.P2P.MaxNumOutboundPeers = 40

	return cfg
}
