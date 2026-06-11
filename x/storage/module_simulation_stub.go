//go:build !simulation

package storage

import (
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func (AppModule) GenerateGenesisState(_ *module.SimulationState) {}

func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

func (AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

func (AppModule) WeightedOperations(_ module.SimulationState) []simtypes.WeightedOperation {
	return nil
}
