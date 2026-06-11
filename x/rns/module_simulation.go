package rns

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	simappparams "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	//	"github.com/jackalLabs/canine-chain/testutil/sample"
	rnssimulation "github.com/jackalLabs/canine-chain/v5/x/rns/simulation"
)

// avoid unused import issue
var (
	//	_ = sample.AccAddress
	_ = rnssimulation.FindAccount
	_ = simappparams.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

//nolint:gosec // these aren't hard-coded credentials
const (
	opWeightMsgRegister = "op_weight_msg_register"
	// TODO: Determine the simulation weight value
	defaultWeightMsgRegister int = 100

	opWeightMsgBid = "op_weight_msg_bid"
	// TODO: Determine the simulation weight value
	defaultWeightMsgBid int = 100

	opWeightMsgCancelBid = "op_weight_msg_cancel_bid"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCancelBid int = 10

	opWeightMsgAcceptBid = "op_weight_msg_accept_bid"
	// TODO: Determine the simulation weight value
	defaultWeightMsgAcceptBid int = 100

	opWeightMsgList = "op_weight_msg_list"
	// TODO: Determine the simulation weight value
	defaultWeightMsgList int = 100

	opWeightMsgBuy = "op_weight_msg_buy"
	// TODO: Determine the simulation weight value
	defaultWeightMsgBuy int = 100

	opWeightMsgDelist = "op_weight_msg_delist"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDelist int = 100

	opWeightMsgTransfer = "op_weight_msg_transfer"
	// TODO: Determine the simulation weight value
	defaultWeightMsgTransfer int = 100

	opWeightMsgAddRecord = "op_weight_msg_add_record"
	// TODO: Determine the simulation weight value
	defaultWeightMsgAddRecord int = 60

	opWeightMsgDelRecord = "op_weight_msg_del_record"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDelRecord int = 40
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	rnssimulation.RandomizedGenState(simState)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}
// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the rns module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	// RNS sim mints/funds ujkl and schedules chained future ops that can destabilize
	// SDK 0.47 import/export simulation on some seeds. Core modules are still covered.
	_ = am
	_ = simState
	return nil
}
