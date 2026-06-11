package v630_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	v630 "github.com/jackalLabs/canine-chain/v5/app/upgrades/v630"
)

func TestUpgradeNameAndStores(t *testing.T) {
	u := v630.NewUpgrade(nil, nil)
	require.Equal(t, "v630", u.Name())

	stores := u.StoreUpgrades()
	require.NotNil(t, stores)
	require.Contains(t, stores.Deleted, "crisis")
	require.Contains(t, stores.Deleted, "circuit")
}
