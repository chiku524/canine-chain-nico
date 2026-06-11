package v460

import (
	"context"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/jackalLabs/canine-chain/v5/app/upgrades"
)

var _ upgrades.Upgrade = &Upgrade{}

type Upgrade struct {
	mm           *module.Manager
	configurator module.Configurator
}

// NewUpgrade returns a new Upgrade instance
func NewUpgrade(mm *module.Manager, configurator module.Configurator) *Upgrade {
	return &Upgrade{
		mm:           mm,
		configurator: configurator,
	}
}

// Name implements upgrades.Upgrade
func (u *Upgrade) Name() string {
	return "v460"
}

func (u *Upgrade) Handler() upgradetypes.UpgradeHandler {
	return func(_ context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		return fromVM, nil
	}
}

// StoreUpgrades implements upgrades.Upgrade
func (u *Upgrade) StoreUpgrades() *storetypes.StoreUpgrades {
	return &storetypes.StoreUpgrades{
		Added:   []string{},
		Deleted: []string{},
	}
}
