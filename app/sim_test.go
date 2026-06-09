//go:build simulation

package app

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	simcli "github.com/cosmos/cosmos-sdk/x/simulation/client/cli"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	minttypes "github.com/jackalLabs/canine-chain/v5/x/jklmint/types"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	filetreemoduletypes "github.com/jackalLabs/canine-chain/v5/x/filetree/types"
	oraclemoduletypes "github.com/jackalLabs/canine-chain/v5/x/oracle/types"
	rnsmoduletypes "github.com/jackalLabs/canine-chain/v5/x/rns/types"
	storagemoduletypes "github.com/jackalLabs/canine-chain/v5/x/storage/types"
)

// SimAppChainID hardcoded chainID for simulation.
const SimAppChainID = "simulation-app"

func init() {
	simcli.GetSimulatorFlags()
}

type StoreKeysPrefixes struct {
	A        storetypes.StoreKey
	B        storetypes.StoreKey
	Prefixes [][]byte
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
		simtestutil.SimulationOperations(app, app.AppCodec(), config),
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
	config, db, _, app := setupSimulationApp(t, "skipping application import/export simulation")

	_, simParams, simErr := simulation.SimulateFromSeed(
		t,
		os.Stdout,
		app.BaseApp,
		simtestutil.AppStateFn(app.AppCodec(), app.SimulationManager(), app.DefaultGenesis()),
		simtypes.RandomAccounts,
		simtestutil.SimulationOperations(app, app.AppCodec(), config),
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

	newDB, newDir, _, _, err := simtestutil.SetupSimulation(config, "leveldb-app-sim-2", "Simulation-2", simcli.FlagVerboseValue, simcli.FlagEnabledValue)
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		require.NoError(t, newDB.Close())
		require.NoError(t, os.RemoveAll(newDir))
	}()

	encConf := MakeEncodingConfig()
	appOptions := make(simtestutil.AppOptionsMap)
	appOptions[flags.FlagHome] = newDir
	appOptions[server.FlagInvCheckPeriod] = simcli.FlagPeriodValue

	newApp := NewJackalApp(log.NewNopLogger(), newDB, nil, true, map[int64]bool{}, newDir, simcli.FlagPeriodValue,
		encConf, wasmtypes.EnableAllProposals, appOptions, nil, fauxMerkleModeOpt, baseapp.SetChainID(SimAppChainID))
	require.Equal(t, appName, newApp.Name())

	var genesisState GenesisState
	err = json.Unmarshal(exported.AppState, &genesisState)
	require.NoError(t, err)

	ctxA := app.NewContext(true, tmproto.Header{Height: app.LastBlockHeight()})
	ctxB := newApp.NewContext(true, tmproto.Header{Height: app.LastBlockHeight()})
	newApp.mm.InitGenesis(ctxB, app.AppCodec(), genesisState)
	newApp.StoreConsensusParams(ctxB, exported.ConsensusParams)

	t.Log("comparing stores...")

	storeKeysPrefixes := []StoreKeysPrefixes{
		{app.keys[authtypes.StoreKey], newApp.keys[authtypes.StoreKey], [][]byte{}},
		{
			app.keys[stakingtypes.StoreKey], newApp.keys[stakingtypes.StoreKey],
			[][]byte{
				stakingtypes.UnbondingQueueKey, stakingtypes.RedelegationQueueKey, stakingtypes.ValidatorQueueKey,
				stakingtypes.HistoricalInfoKey,
			},
		},
		{app.keys[slashingtypes.StoreKey], newApp.keys[slashingtypes.StoreKey], [][]byte{}},
		{app.keys[minttypes.StoreKey], newApp.keys[minttypes.StoreKey], [][]byte{}},
		{app.keys[distrtypes.StoreKey], newApp.keys[distrtypes.StoreKey], [][]byte{}},
		{app.keys[banktypes.StoreKey], newApp.keys[banktypes.StoreKey], [][]byte{banktypes.BalancesPrefix}},
		{app.keys[paramstypes.StoreKey], newApp.keys[paramstypes.StoreKey], [][]byte{}},
		{app.keys[govtypes.StoreKey], newApp.keys[govtypes.StoreKey], [][]byte{}},
		{app.keys[evidencetypes.StoreKey], newApp.keys[evidencetypes.StoreKey], [][]byte{}},
		{app.keys[capabilitytypes.StoreKey], newApp.keys[capabilitytypes.StoreKey], [][]byte{}},
		{app.keys[ibcexported.StoreKey], newApp.keys[ibcexported.StoreKey], [][]byte{}},
		{app.keys[ibctransfertypes.StoreKey], newApp.keys[ibctransfertypes.StoreKey], [][]byte{}},
		{app.keys[authzkeeper.StoreKey], newApp.keys[authzkeeper.StoreKey], [][]byte{}},
		{app.keys[feegrant.StoreKey], newApp.keys[feegrant.StoreKey], [][]byte{}},
		{app.keys[wasm.StoreKey], newApp.keys[wasm.StoreKey], [][]byte{}},
		{app.keys[oraclemoduletypes.StoreKey], newApp.keys[oraclemoduletypes.StoreKey], [][]byte{}},
		{app.keys[storagemoduletypes.StoreKey], newApp.keys[storagemoduletypes.StoreKey], [][]byte{}},
		{app.keys[filetreemoduletypes.StoreKey], newApp.keys[filetreemoduletypes.StoreKey], [][]byte{}},
		{app.keys[rnsmoduletypes.StoreKey], newApp.keys[rnsmoduletypes.StoreKey], [][]byte{}},
	}

	ctxA.KVStore(app.keys[wasm.StoreKey]).Delete(wasmtypes.TXCounterPrefix)

	dropContractHistory := func(s sdk.KVStore, keys ...[]byte) {
		for _, key := range keys {
			prefixStore := prefix.NewStore(s, key)
			iter := prefixStore.Iterator(nil, nil)
			for ; iter.Valid(); iter.Next() {
				prefixStore.Delete(iter.Key())
			}
			iter.Close()
		}
	}
	prefixes := [][]byte{wasmtypes.ContractCodeHistoryElementPrefix, wasmtypes.ContractByCodeIDAndCreatedSecondaryIndexPrefix}
	dropContractHistory(ctxA.KVStore(app.keys[wasm.StoreKey]), prefixes...)
	dropContractHistory(ctxB.KVStore(newApp.keys[wasm.StoreKey]), prefixes...)

	normalizeContractInfo := func(ctx sdk.Context, jackalApp *JackalApp) {
		var index uint64
		jackalApp.wasmKeeper.IterateContractInfo(ctx, func(address sdk.AccAddress, info wasmtypes.ContractInfo) bool {
			created := &wasmtypes.AbsoluteTxPosition{
				BlockHeight: uint64(0),
				TxIndex:     index,
			}
			info.Created = created
			kvStore := ctx.KVStore(jackalApp.keys[wasm.StoreKey])
			kvStore.Set(wasmtypes.GetContractAddressKey(address), jackalApp.appCodec.MustMarshal(&info))
			index++
			return false
		})
	}
	normalizeContractInfo(ctxA, app)
	normalizeContractInfo(ctxB, newApp)

	for _, skp := range storeKeysPrefixes {
		storeA := ctxA.KVStore(skp.A)
		storeB := ctxB.KVStore(skp.B)

		failedKVAs, failedKVBs := sdk.DiffKVStores(storeA, storeB, skp.Prefixes)
		require.Equal(t, len(failedKVAs), len(failedKVBs), "unequal sets of key-values to compare")

		t.Logf("compared %d different key/value pairs between %s and %s\n", len(failedKVAs), skp.A, skp.B)
		require.Len(t, failedKVAs, 0, simtestutil.GetSimulationLog(skp.A.Name(), app.SimulationManager().StoreDecoders, failedKVAs, failedKVBs))
	}
}

func BenchmarkFullAppSimulation(b *testing.B) {
	config, db, _, app := setupSimulationApp(b, "skipping application simulation")

	_, simParams, simErr := simulation.SimulateFromSeed(
		b,
		os.Stdout,
		app.BaseApp,
		simtestutil.AppStateFn(app.AppCodec(), app.SimulationManager(), app.DefaultGenesis()),
		simtypes.RandomAccounts,
		simtestutil.SimulationOperations(app, app.AppCodec(), config),
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
	setBech32ForTest()

	config := simcli.NewConfigFromFlags()
	config.ChainID = SimAppChainID
	config.Commit = true

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
		encConf, wasmtypes.EnableAllProposals, appOptions, nil, fauxMerkleModeOpt, baseapp.SetChainID(SimAppChainID))
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
				logger = log.TestingLogger()
			} else {
				logger = log.NewNopLogger()
			}

			db := dbm.NewMemDB()
			encConf := MakeEncodingConfig()
			app := NewJackalApp(logger, db, nil, true, map[int64]bool{}, t.TempDir(), simcli.FlagPeriodValue,
				encConf, wasmtypes.EnableAllProposals, appOptions, nil, fauxMerkleModeOpt, baseapp.SetChainID(SimAppChainID))

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
				simtestutil.SimulationOperations(app, app.AppCodec(), config),
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
