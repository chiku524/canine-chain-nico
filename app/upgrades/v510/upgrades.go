package v510

import (
	"context"
	_ "embed"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	return "v510"
}

func (u *Upgrade) Handler() upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)

		sdkCtx.Logger().Info("Updating the Jackal Protocol to v5.1.0...")

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
