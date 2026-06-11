//go:build !simulation

package app

import "github.com/cosmos/cosmos-sdk/types/module"

func simulationOverrides(base map[string]module.AppModuleSimulation, _ map[string]interface{}) map[string]module.AppModuleSimulation {
	return base
}
