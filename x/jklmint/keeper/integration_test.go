package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	jklapp "github.com/jackalLabs/canine-chain/v5/app"
	"github.com/jackalLabs/canine-chain/v5/testutil"
	"github.com/jackalLabs/canine-chain/v5/x/jklmint/types"
)

// returns context and an app with updated mint keeper
func createTestApp(t *testing.T, isCheckTx bool) (*jklapp.JackalApp, sdk.Context) {
	t.Helper()
	app := setup(t, isCheckTx)

	ctx := app.NewContext(isCheckTx, tmproto.Header{})
	app.MintKeeper.SetParams(ctx, types.DefaultParams())

	return app, ctx
}

func setup(t *testing.T, isCheckTx bool) *jklapp.JackalApp {
	t.Helper()
	if !testutil.CgoEnabled() {
		t.Skip("integration tests require CGO for wasmvm")
	}
	return jklapp.SetupTestingAppWithGenesis(t)
}
