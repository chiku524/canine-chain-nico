package v600

import (
	"context"
	"github.com/cosmos/cosmos-sdk/baseapp"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/jackalLabs/canine-chain/v5/app/upgrades"
)

var _ upgrades.Upgrade = &Upgrade{}

// Upgrade migrates jackal-1 from Cosmos SDK 0.45 to 0.47.
// On-chain upgrade name: "v600"
type Upgrade struct {
	mm              *module.Manager
	configurator    module.Configurator
	paramsKeeper    paramskeeper.Keeper
	consensusKeeper keeper.Keeper
}

func NewUpgrade(
	mm *module.Manager,
	configurator module.Configurator,
	paramsKeeper paramskeeper.Keeper,
	consensusKeeper keeper.Keeper,
) *Upgrade {
	return &Upgrade{
		mm:              mm,
		configurator:    configurator,
		paramsKeeper:    paramsKeeper,
		consensusKeeper: consensusKeeper,
	}
}

func (u *Upgrade) Name() string {
	return "v600"
}

func (u *Upgrade) Handler() upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)

		sdkCtx.Logger().Info("Upgrading Jackal Protocol to Cosmos SDK 0.47 (v600)...")

		baseAppLegacySS := u.paramsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable())
		baseapp.MigrateParams(sdkCtx, baseAppLegacySS, u.consensusKeeper.ParamsStore)

		return u.mm.RunMigrations(sdkCtx, u.configurator, fromVM)
	}
}

func (u *Upgrade) StoreUpgrades() *storetypes.StoreUpgrades {
	return &storetypes.StoreUpgrades{
		Added: []string{
			"consensus",
			"crisis",
		},
	}
}
