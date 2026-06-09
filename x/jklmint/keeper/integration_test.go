package keeper_test

import (
	"testing"

	jklapp "github.com/jackalLabs/canine-chain/v5/app"
	"github.com/jackalLabs/canine-chain/v5/testutil"
)

func setup(t *testing.T) *jklapp.JackalApp {
	t.Helper()
	if !testutil.CgoEnabled() {
		t.Skip("integration tests require CGO for wasmvm")
	}
	return jklapp.SetupTestingAppWithGenesis(t)
}
