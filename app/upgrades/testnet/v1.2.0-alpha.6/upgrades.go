package v120alpha6

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/jackalLabs/canine-chain/v5/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)

		logger := sdkCtx.Logger().With("upgrade", UpgradeName)

		if types.IsTestnet(sdkCtx.ChainID()) {
			logger.Debug("Updating to 1.2.0-alpha.6")
		}

		if types.IsMainnet(sdkCtx.ChainID()) {
			logger.Debug("Ignoring Infra & Storage Deals for mainnet")
		}

		return mm.RunMigrations(sdkCtx, configurator, vm)
	}
}
