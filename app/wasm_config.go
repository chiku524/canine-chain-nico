package app

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

const (
	// DefaultJackalInstanceCost is initially set the same as in wasmd
	DefaultJackalInstanceCost uint64 = 60_000
	// DefaultJackalCompileCost set to a large number for testing
	DefaultJackalCompileCost uint64 = 3
)

// JackalGasRegisterConfig is defaults plus a custom compile amount
func JackalGasRegisterConfig() wasmtypes.WasmGasRegisterConfig {
	gasConfig := wasmtypes.DefaultGasRegisterConfig()
	gasConfig.InstanceCost = DefaultJackalInstanceCost
	gasConfig.CompileCost = DefaultJackalCompileCost

	return gasConfig
}

func NewJackalWasmGasRegister() wasmtypes.GasRegister {
	return wasmtypes.NewWasmGasRegister(JackalGasRegisterConfig())
}
