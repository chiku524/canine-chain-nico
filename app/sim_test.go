//go:build simulation

package app

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/stretchr/testify/require"

	"cosmossdk.io/log/v2"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	simcli "github.com/cosmos/cosmos-sdk/x/simulation/client/cli"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

// SimAppChainID hardcoded chainID for simulation.
const SimAppChainID = "simulation-app"

var simWasmProposals = wasmtypes.DisableAllProposals

func init() {
	simcli.GetSimulatorFlags()
}

func fauxMerkleModeOpt(bapp *baseapp.BaseApp) {
	bapp.SetFauxMerkleMode()
}

func TestFullAppSimulation(t *testing.T) {
	config, db, _, app := setupSimulationApp(t, "skipping application simulation")

	_, simParams, simErr := simulation.SimulateFromSeed(
		t,
		os.Stdout,
		app.BaseApp,
		simtestutil.AppStateFn(app.AppCodec(), app.SimulationManager(), app.DefaultGenesis()),
		simtypes.RandomAccounts,
		simtestutil.BuildSimulationOperations(app, app.AppCodec(), config, app.TxConfig()),
		BlockedAddresses(),
		config,
		app.AppCodec(),
	)

	err := simtestutil.CheckExportSimulation(app, config, simParams)
	require.NoError(t, err)
	require.NoError(t, simErr)

	if config.Commit {
		simtestutil.PrintStats(db)
	}
}

func TestAppImportExport(t *testing.T) {
	config, db, appOptions, app := setupSimulationApp(t, "skipping application import/export simulation")

	_, simParams, simErr := simulation.SimulateFromSeed(
		t,
		os.Stdout,
		app.BaseApp,
		simtestutil.AppStateFn(app.AppCodec(), app.SimulationManager(), app.DefaultGenesis()),
		simtypes.RandomAccounts,
		simtestutil.BuildSimulationOperations(app, app.AppCodec(), config, app.TxConfig()),
		BlockedAddresses(),
		config,
		app.AppCodec(),
	)

	err := simtestutil.CheckExportSimulation(app, config, simParams)
	require.NoError(t, err)
	require.NoError(t, simErr)

	if config.Commit {
		simtestutil.PrintStats(db)
	}

	t.Log("exporting genesis...")

	exported, err := app.ExportAppStateAndValidators(false, []string{}, nil)
	require.NoError(t, err)

	t.Log("importing genesis...")

	newDB, newDir, _, _, err := simtestutil.SetupSimulation(config, "leveldb-app-sim-2", "Simulation-2", simcli.FlagVerboseValue, true)
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		require.NoError(t, newDB.Close())
		require.NoError(t, os.RemoveAll(newDir))
	}()

	appOptions[flags.FlagHome] = newDir

	encConf := MakeEncodingConfig()
	newApp := NewJackalApp(log.NewNopLogger(), newDB, nil, true, map[int64]bool{}, newDir, simcli.FlagPeriodValue,
		encConf, simWasmProposals, appOptions, nil, fauxMerkleModeOpt, baseapp.SetChainID(SimAppChainID))
	require.Equal(t, appName, newApp.Name())

	initReq := &abci.RequestInitChain{
		AppStateBytes: exported.AppState,
	}

	ctxB := newApp.NewContextLegacy(true, cmtproto.Header{Height: app.LastBlockHeight()})
	_, err = newApp.InitChainer(ctxB, initReq)
	if err != nil {
		if strings.Contains(err.Error(), "validator set is empty after InitGenesis") {
			t.Log("Skipping import/export compare: all validators unbonded")
			t.Logf("err: %s stacktrace: %s\n", err, string(debug.Stack()))
			return
		}
		require.NoError(t, err)
	}

	err = newApp.StoreConsensusParams(ctxB, exported.ConsensusParams)
	require.NoError(t, err)

	t.Log("re-exporting genesis from imported app...")
	reExported, err := newApp.ExportAppStateAndValidators(false, []string{}, nil)
	require.NoError(t, err)
	require.JSONEq(t, string(exported.AppState), string(reExported.AppState))
}

func BenchmarkFullAppSimulation(b *testing.B) {
	config, db, _, app := setupSimulationApp(b, "skipping application simulation")

	_, simParams, simErr := simulation.SimulateFromSeed(
		b,
		os.Stdout,
		app.BaseApp,
		simtestutil.AppStateFn(app.AppCodec(), app.SimulationManager(), app.DefaultGenesis()),
		simtypes.RandomAccounts,
		simtestutil.BuildSimulationOperations(app, app.AppCodec(), config, app.TxConfig()),
		BlockedAddresses(),
		config,
		app.AppCodec(),
	)

	err := simtestutil.CheckExportSimulation(app, config, simParams)
	require.NoError(b, err)
	require.NoError(b, simErr)

	if config.Commit {
		simtestutil.PrintStats(db)
	}
}

func setupSimulationApp(t testing.TB, skipMsg string) (simtypes.Config, dbm.DB, simtestutil.AppOptionsMap, *JackalApp) {
	t.Helper()
	SetBech32ForTest()

	config := simcli.NewConfigFromFlags()
	config.ChainID = SimAppChainID

	db, dir, logger, skip, err := simtestutil.SetupSimulation(config, "leveldb-app-sim", "Simulation", simcli.FlagVerboseValue, simcli.FlagEnabledValue)
	if skip {
		t.Skip(skipMsg)
	}
	require.NoError(t, err, "simulation setup failed")

	t.Cleanup(func() {
		require.NoError(t, db.Close())
		require.NoError(t, os.RemoveAll(dir))
	})

	appOptions := make(simtestutil.AppOptionsMap)
	appOptions[flags.FlagHome] = dir
	appOptions[server.FlagInvCheckPeriod] = simcli.FlagPeriodValue

	encConf := MakeEncodingConfig()
	app := NewJackalApp(logger, db, nil, true, map[int64]bool{}, dir, simcli.FlagPeriodValue,
		encConf, simWasmProposals, appOptions, nil, fauxMerkleModeOpt, baseapp.SetChainID(SimAppChainID))
	require.Equal(t, appName, app.Name())

	return config, db, appOptions, app
}

func TestAppStateDeterminism(t *testing.T) {
	if !simcli.FlagEnabledValue {
		t.Skip("skipping application simulation")
	}

	config := simcli.NewConfigFromFlags()
	config.InitialBlockHeight = 1
	config.ExportParamsPath = ""
	config.OnOperation = false
	config.AllInvariants = false
	config.ChainID = SimAppChainID

	numSeeds := 3
	numTimesToRunPerSeed := 5
	appHashList := make([]json.RawMessage, numTimesToRunPerSeed)

	appOptions := make(simtestutil.AppOptionsMap)
	appOptions[flags.FlagHome] = t.TempDir()
	appOptions[server.FlagInvCheckPeriod] = simcli.FlagPeriodValue

	for i := 0; i < numSeeds; i++ {
		config.Seed += int64(i)

		for j := 0; j < numTimesToRunPerSeed; j++ {
			var logger log.Logger
			if simcli.FlagVerboseValue {
				logger = log.NewTestLogger(t)
			} else {
				logger = log.NewNopLogger()
			}

			db := dbm.NewMemDB()
			encConf := MakeEncodingConfig()
			app := NewJackalApp(logger, db, nil, true, map[int64]bool{}, t.TempDir(), simcli.FlagPeriodValue,
				encConf, simWasmProposals, appOptions, nil, fauxMerkleModeOpt, baseapp.SetChainID(SimAppChainID))

			fmt.Printf(
				"running non-determinism simulation; seed %d: %d/%d, attempt: %d/%d\n",
				config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
			)

			_, _, err := simulation.SimulateFromSeed(
				t,
				os.Stdout,
				app.BaseApp,
				simtestutil.AppStateFn(app.AppCodec(), app.SimulationManager(), app.DefaultGenesis()),
				simtypes.RandomAccounts,
				simtestutil.BuildSimulationOperations(app, app.AppCodec(), config, app.TxConfig()),
				BlockedAddresses(),
				config,
				app.AppCodec(),
			)
			require.NoError(t, err)

			if config.Commit {
				simtestutil.PrintStats(db)
			}

			appHash := app.LastCommitID().Hash
			appHashList[j] = appHash

			if j != 0 {
				require.Equal(
					t, string(appHashList[0]), string(appHashList[j]),
					"non-determinism in seed %d: %d/%d, attempt: %d/%d\n", config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
				)
			}
		}
	}
}
