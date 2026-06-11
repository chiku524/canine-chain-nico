package v610_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	v610 "github.com/jackalLabs/canine-chain/v5/app/upgrades/v610"
)

func TestUpgradeNameAndStores(t *testing.T) {
	u := v610.NewUpgrade(nil, nil)
	require.Equal(t, "v610", u.Name())

	stores := u.StoreUpgrades()
	require.NotNil(t, stores)
	require.Contains(t, stores.Added, "circuit")
}
