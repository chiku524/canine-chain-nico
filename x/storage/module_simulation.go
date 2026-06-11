package storage

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	simappparams "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	storagesimulation "github.com/jackalLabs/canine-chain/v5/x/storage/simulation"
)

// avoid unused import issue
var (
	// _ = sample.AccAddress
	_ = storagesimulation.FindAccount
	_ = simappparams.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	//nolint:all
	opWeightMsgSetProviderIP = "op_weight_msg_set_provider_ip"
	// TODO: Determine the simulation weight value
	defaultWeightMsgSetProviderIP int = 10

	//nolint:all
	opWeightMsgSetProviderTotalSpace = "op_weight_msg_set_provider_totalspace"
	// TODO: Determine the simulation weight value
	defaultWeightMsgSetProviderTotalSpace int = 10

	//nolint:all
	opWeightMsgInitProvider = "op_weight_msg_init_provider"
	// TODO: Determine the simulation weight value
	defaultWeightMsgInitProvider int = 60

	//nolint:all
	opWeightMsgBuyStorage = "op_weight_msg_buy_storage"
	// TODO: Determine the simulation weight value
	defaultWeightMsgBuyStorage int = 100

	//nolint:all
	opWeightMsgAddProviderClaimer          = "op_weight_msg_add_provider_claimer"
	defaultWeightMsgAddProviderClaimer int = 100

	//nolint:all
	opWeightMsgRemoveProviderClaimer          = "op_weight_msg_remove_provider_claimer"
	defaultWeightMsgRemoveProviderClaimer int = 10
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
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	// Storage provider/proof sim interacts with periodic BeginBlock proof scans.
	_ = simState
	return nil
}
