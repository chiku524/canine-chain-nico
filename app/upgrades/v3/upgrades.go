package v3

import (
	"context"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/jackalLabs/canine-chain/v5/app/upgrades"
	storagekeeper "github.com/jackalLabs/canine-chain/v5/x/storage/keeper"

	storagemoduletypes "github.com/jackalLabs/canine-chain/v5/x/storage/types"
)

var _ upgrades.Upgrade = &Upgrade{}

// Upgrade represents the v3 upgrade
type Upgrade struct {
	mm           *module.Manager
	configurator module.Configurator
	sk           storagekeeper.Keeper
}

// NewUpgrade returns a new Upgrade instance
func NewUpgrade(mm *module.Manager, configurator module.Configurator, sk storagekeeper.Keeper) *Upgrade {
	return &Upgrade{
		mm:           mm,
		configurator: configurator,
		sk:           sk,
	}
}

// Name implements upgrades.Upgrade
func (u *Upgrade) Name() string {
	return "v3"
}

// Handler implements upgrades.Upgrade
func (u *Upgrade) Handler() upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)

		fromVM[storagemoduletypes.ModuleName] = 4

		newVM, err := u.mm.RunMigrations(sdkCtx, u.configurator, fromVM)
		if err != nil {
			return newVM, err
		}

		return newVM, err
	}
}

// StoreUpgrades implements upgrades.Upgrade
func (u *Upgrade) StoreUpgrades() *storetypes.StoreUpgrades {
	return &storetypes.StoreUpgrades{
		Added: []string{
			"feeibc",
		},
		Deleted: []string{
			"intertx", // legacy interchain-accounts inter-tx module store
		},
	}
}
