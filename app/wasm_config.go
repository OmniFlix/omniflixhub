package app

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
)

const (
	// DefaultOmniFlixInstanceCost is initially set the same as in wasmd
	DefaultOmniFlixInstanceCost uint64 = 60_000
	// DefaultOmniFlixCompileCost set to a large number for testing
	DefaultOmniFlixCompileCost uint64 = 100
)

// OmniFlixGasRegisterConfig is defaults plus a custom compile amount
func OmniFlixGasRegisterConfig() wasmkeeper.WasmGasRegisterConfig {
	gasConfig := wasmkeeper.DefaultGasRegisterConfig()
	gasConfig.InstanceCost = DefaultOmniFlixInstanceCost
	gasConfig.CompileCost = DefaultOmniFlixCompileCost

	return gasConfig
}

func NewOmniFlixWasmGasRegister() wasmkeeper.WasmGasRegister {
	return wasmkeeper.NewWasmGasRegister(OmniFlixGasRegisterConfig())
}
