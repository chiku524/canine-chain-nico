package v450

import (
	"context"
	_ "embed"

	types2 "github.com/jackalLabs/canine-chain/v5/types"
	"github.com/jackalLabs/canine-chain/v5/x/storage/types"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/jackalLabs/canine-chain/v5/app/upgrades"
	storageKeeper "github.com/jackalLabs/canine-chain/v5/x/storage/keeper"
)

var _ upgrades.Upgrade = &Upgrade{}

//go:embed upgrade_data
var UpgradeData string

type Upgrade struct {
	mm           *module.Manager
	configurator module.Configurator
	sk           *storageKeeper.Keeper
	bk           types.BankKeeper
	ak           types.AccountKeeper
}

// NewUpgrade returns a new Upgrade instance
func NewUpgrade(mm *module.Manager, configurator module.Configurator, sk *storageKeeper.Keeper, bk types.BankKeeper, ak types.AccountKeeper) *Upgrade {
	return &Upgrade{
		mm:           mm,
		configurator: configurator,
		sk:           sk,
		bk:           bk,
		ak:           ak,
	}
}

// Name implements upgrades.Upgrade
func (u *Upgrade) Name() string {
	return "v450"
}

func (u *Upgrade) fixPayment(ctx sdk.Context) { // fixing payment mismatch from burned files
	usedCounter := make(map[string]int64)

	u.sk.IterateAndParseFilesByMerkle(ctx, false, func(_ []byte, val types.UnifiedFile) bool {
		owner := val.Owner
		size := val.FileSize

		if val.Expires > 0 {
			return false
		}

		usedCounter[owner] += size

		return false
	})

	for owner, size := range usedCounter {
		payInfo, found := u.sk.GetStoragePaymentInfo(ctx, owner)
		if !found {
			continue
		}

		payInfo.SpaceUsed = size

		u.sk.SetStoragePaymentInfo(ctx, payInfo)
	}
}

// Handler implements upgrades.Upgrade
func (u *Upgrade) Handler() upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)

		err := upgrades.RecoverFiles(sdkCtx, u.sk, UpgradeData, plan.Height, "v4.5.0")
		if err != nil {
			return nil, err
		}

		pol, err := types2.GetPOLAccount()
		if err != nil {
			return nil, err
		}

		storageAccount := u.ak.GetModuleAddress(types.ModuleName)

		bal := u.bk.GetBalance(sdkCtx, storageAccount, "ujkl")

		err = u.bk.SendCoinsFromModuleToAccount(sdkCtx, types.ModuleName, pol, sdk.NewCoins(bal))
		if err != nil {
			return nil, err
		}

		u.fixPayment(sdkCtx)

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
