package v600_test

import (
	"testing"

	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	"github.com/stretchr/testify/require"

	v600 "github.com/jackalLabs/canine-chain/v5/app/upgrades/v600"
)

func TestUpgradeNameAndStores(t *testing.T) {
	u := v600.NewUpgrade(nil, nil, paramskeeper.Keeper{}, consensusparamkeeper.Keeper{})
	require.Equal(t, "v600", u.Name())

	stores := u.StoreUpgrades()
	require.NotNil(t, stores)
	require.Contains(t, stores.Added, "consensus")
	require.Contains(t, stores.Added, "crisis")
}
