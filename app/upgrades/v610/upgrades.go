package v610

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	circuittypes "cosmossdk.io/x/circuit/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/jackalLabs/canine-chain/v5/app/upgrades"
)

var _ upgrades.Upgrade = &Upgrade{}

// Upgrade migrates jackal-1 from Cosmos SDK 0.47 to 0.50.
// On-chain upgrade name: "v610"
type Upgrade struct {
	mm           *module.Manager
	configurator module.Configurator
}

func NewUpgrade(mm *module.Manager, configurator module.Configurator) *Upgrade {
	return &Upgrade{mm: mm, configurator: configurator}
}

func (u *Upgrade) Name() string {
	return "v610"
}

func (u *Upgrade) Handler() upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		sdkCtx.Logger().Info("Upgrading Jackal Protocol to Cosmos SDK 0.50 (v610)...")
		return u.mm.RunMigrations(sdkCtx, u.configurator, fromVM)
	}
}

func (u *Upgrade) StoreUpgrades() *storetypes.StoreUpgrades {
	return &storetypes.StoreUpgrades{
		Added: []string{
			circuittypes.ModuleName,
		},
	}
}
