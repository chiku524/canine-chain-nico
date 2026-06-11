package storage

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	storagesimulation "github.com/jackalLabs/canine-chain/v5/x/storage/simulation"
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	storagesimulation.RandomizedGenState(simState)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RegisterStoreDecoder registers a decoder
func (AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (AppModule) WeightedOperations(_ module.SimulationState) []simtypes.WeightedOperation {
	// Storage provider/proof sim interacts with periodic BeginBlock proof scans.
	return nil
}
