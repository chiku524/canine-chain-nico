//go:build cgo

package app

import (
	"encoding/json"
	"testing"

	"github.com/cometbft/cometbft/libs/log"
	dbm "github.com/cometbft/cometbft-db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

var emptyWasmOpts []wasmkeeper.Option

func TestWasmdExport(t *testing.T) {
	setBech32ForTest()

	gapp := SetupTestingAppWithGenesis(t)

	exported, err := gapp.ExportAppStateAndValidators(false, []string{}, nil)
	require.NoError(t, err)

	db := dbm.NewMemDB()
	newGapp := NewJackalApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, DefaultNodeHome, 5,
		MakeEncodingConfig(), wasmtypes.EnableAllProposals, EmptyBaseAppOptions{}, emptyWasmOpts)

	var genesisState GenesisState
	require.NoError(t, json.Unmarshal(exported.AppState, &genesisState))

	ctx := newGapp.NewContext(true, tmproto.Header{Height: 0})
	newGapp.mm.InitGenesis(ctx, newGapp.AppCodec(), genesisState)
	newGapp.StoreConsensusParams(ctx, exported.ConsensusParams)
}

// ensure that blocked addresses are properly set in bank keeper
func TestBlockedAddrs(t *testing.T) {
	setBech32ForTest()

	gapp := SetupTestingAppWithGenesis(t)

	for acc := range maccPerms {
		t.Run(acc, func(t *testing.T) {
			require.True(t, gapp.BankKeeper.BlockedAddr(gapp.AccountKeeper.GetModuleAddress(acc)),
				"ensure that blocked addresses are properly set in bank keeper",
			)
		})
	}
}

func TestGetMaccPerms(t *testing.T) {
	dup := GetMaccPerms()
	require.Equal(t, maccPerms, dup, "duplicated module account permissions differed from actual module account permissions")
}

func TestGetEnabledProposals(t *testing.T) {
	cases := map[string]struct {
		proposalsEnabled string
		specificEnabled  string
		expected         []wasmtypes.ProposalType
	}{
		"all disabled": {
			proposalsEnabled: "false",
			expected:         wasmtypes.DisableAllProposals,
		},
		"all enabled": {
			proposalsEnabled: "true",
			expected:         wasmtypes.EnableAllProposals,
		},
		"some enabled": {
			proposalsEnabled: "okay",
			specificEnabled:  "StoreCode,InstantiateContract",
			expected:         []wasmtypes.ProposalType{wasmtypes.ProposalTypeStoreCode, wasmtypes.ProposalTypeInstantiateContract},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			ProposalsEnabled = tc.proposalsEnabled
			EnableSpecificProposals = tc.specificEnabled
			proposals := GetEnabledProposals()
			assert.Equal(t, tc.expected, proposals)
		})
	}
}
