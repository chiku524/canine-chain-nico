//go:build simulation

package app

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	filetreemoduletypes "github.com/jackalLabs/canine-chain/v5/x/filetree/types"
	notificationsmoduletypes "github.com/jackalLabs/canine-chain/v5/x/notifications/types"
	oraclemoduletypes "github.com/jackalLabs/canine-chain/v5/x/oracle/types"
	rnsmoduletypes "github.com/jackalLabs/canine-chain/v5/x/rns/types"
	storagemoduletypes "github.com/jackalLabs/canine-chain/v5/x/storage/types"
)

// noopSimModule disables WeightedOperations while preserving other simulation hooks.
type noopSimModule struct {
	module.AppModuleSimulation
}

func (noopSimModule) WeightedOperations(_ module.SimulationState) []simtypes.WeightedOperation {
	return nil
}

func simulationOverrides(base map[string]module.AppModuleSimulation, modules map[string]interface{}) map[string]module.AppModuleSimulation {
	if base == nil {
		base = make(map[string]module.AppModuleSimulation)
	}

	for _, name := range []string{
		wasmtypes.ModuleName,
		filetreemoduletypes.ModuleName,
		storagemoduletypes.ModuleName,
		rnsmoduletypes.ModuleName,
		oraclemoduletypes.ModuleName,
		notificationsmoduletypes.ModuleName,
	} {
		if m, ok := modules[name].(module.AppModuleSimulation); ok {
			base[name] = noopSimModule{m}
		}
	}

	return base
}
