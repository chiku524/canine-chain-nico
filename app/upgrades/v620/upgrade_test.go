package v620_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	v620 "github.com/jackalLabs/canine-chain/v5/app/upgrades/v620"
)

func TestUpgradeNameAndStores(t *testing.T) {
	u := v620.NewUpgrade(nil, nil)
	require.Equal(t, "v620", u.Name())

	stores := u.StoreUpgrades()
	require.NotNil(t, stores)
	require.Contains(t, stores.Deleted, "capability")
	require.Contains(t, stores.Deleted, "feeibc")
}
