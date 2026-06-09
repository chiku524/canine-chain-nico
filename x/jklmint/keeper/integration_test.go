package keeper_test

import (
	"encoding/json"

	abci "github.com/cometbft/cometbft/abci/types"

	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	jklapp "github.com/jackalLabs/canine-chain/v5/app"
	"github.com/jackalLabs/canine-chain/v5/x/jklmint/types"
	dbm "github.com/cometbft/cometbft-db"
)

// returns context and an app with updated mint keeper
//
//nolint:unused
func createTestApp(isCheckTx bool) (*jklapp.JackalApp, sdk.Context) {
	app := setup(isCheckTx)

	ctx := app.NewContext(isCheckTx, tmproto.Header{})

	app.MintKeeper.SetParams(ctx, types.DefaultParams())

	return app, ctx
}

func setup(isCheckTx bool) *jklapp.JackalApp {
	app, genesisState := genApp(!isCheckTx, 5)
	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: jklapp.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return app
}

func genApp(withGenesis bool, invCheckPeriod uint) (*jklapp.JackalApp, jklapp.GenesisState) {
	db := dbm.NewMemDB()
	encCdc := jklapp.MakeEncodingConfig()
	app := jklapp.NewJackalApp(
		log.NewNopLogger(),
		db,
		nil,
		true,
		map[int64]bool{},
		jklapp.DefaultNodeHome,
		invCheckPeriod,
		encCdc,
		jklapp.GetEnabledProposals(),
		jklapp.EmptyBaseAppOptions{},
		jklapp.GetWasmOpts(jklapp.EmptyBaseAppOptions{}),
	)

	if withGenesis {
		return app, jklapp.NewDefaultGenesisState()
	}

	return app, jklapp.GenesisState{}
}
