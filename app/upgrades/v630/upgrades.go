package v630

import (
	"context"

	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/jackalLabs/canine-chain/v5/app/upgrades"
)

var _ upgrades.Upgrade = &Upgrade{}

// Upgrade migrates jackal-1 from Cosmos SDK 0.53 to 0.54 (store/v2, ibc-go v11, wasmd 0.70).
// On-chain upgrade name: "v630"
type Upgrade struct {
	mm           *module.Manager
	configurator module.Configurator
}

func NewUpgrade(mm *module.Manager, configurator module.Configurator) *Upgrade {
	return &Upgrade{mm: mm, configurator: configurator}
}

func (u *Upgrade) Name() string {
	return "v630"
}

func (u *Upgrade) Handler() upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		sdkCtx.Logger().Info("Upgrading Jackal Protocol to Cosmos SDK 0.54 (v630)...")
		return u.mm.RunMigrations(sdkCtx, u.configurator, fromVM)
	}
}

func (u *Upgrade) StoreUpgrades() *storetypes.StoreUpgrades {
	return &storetypes.StoreUpgrades{
		Deleted: []string{
			"crisis",
			"circuit",
		},
	}
}
