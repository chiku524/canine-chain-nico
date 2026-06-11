package v620

import (
	"context"

	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/jackalLabs/canine-chain/v5/app/upgrades"
)

var _ upgrades.Upgrade = &Upgrade{}

// Upgrade migrates jackal-1 from Cosmos SDK 0.50 to 0.53 (ibc-go v10, wasmd 0.60).
// On-chain upgrade name: "v620"
type Upgrade struct {
	mm           *module.Manager
	configurator module.Configurator
}

func NewUpgrade(mm *module.Manager, configurator module.Configurator) *Upgrade {
	return &Upgrade{mm: mm, configurator: configurator}
}

func (u *Upgrade) Name() string {
	return "v620"
}

func (u *Upgrade) Handler() upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		sdkCtx.Logger().Info("Upgrading Jackal Protocol to Cosmos SDK 0.53 (v620)...")
		return u.mm.RunMigrations(sdkCtx, u.configurator, fromVM)
	}
}

func (u *Upgrade) StoreUpgrades() *storetypes.StoreUpgrades {
	return &storetypes.StoreUpgrades{
		Deleted: []string{
			"capability",
			"feeibc",
		},
	}
}
