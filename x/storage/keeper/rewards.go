package keeper

import (
	sdkmath "cosmossdk.io/math"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/jackalLabs/canine-chain/v5/x/storage/types"
)

func (k Keeper) pullTokensFromGauges(ctx sdk.Context) sdk.Coins {
	currentTime := ctx.BlockTime()
	coinsToDistribute := make(sdk.Coins, 0)

	k.IterateGauges(ctx, func(pg types.PaymentGauge) { // check every gauge
		if pg.End.Before(currentTime) || pg.End.Before(pg.Start) || pg.End.Equal(pg.Start) { // if the gauge is expired or has an invalid end time, we remove it
			k.RemoveGauge(ctx, pg.Id)
			return
		}

		gaugeWallet, err := types.GetGaugeAccount(pg)
		if err != nil {
			ctx.Logger().Error(err.Error())
			return
		}

		gaugeBalance := k.bankKeeper.GetAllBalances(ctx, gaugeWallet) // delete empty gauges
		if gaugeBalance.Empty() {
			k.RemoveGauge(ctx, pg.Id)
			return
		}

		allGaugeCoins := pg.Coins

		totalTime := pg.End.Sub(pg.Start)
		timeLeft := pg.End.Sub(currentTime)

		totalTimeDec := sdkmath.LegacyNewDec(totalTime.Microseconds())
		timeLeftDec := sdkmath.LegacyNewDec(timeLeft.Microseconds())

		timeRatio := sdkmath.LegacyNewDec(1).Sub(timeLeftDec.Quo(totalTimeDec))

		for _, coin := range allGaugeCoins {
			coinAmountDec := sdkmath.LegacyNewDecFromInt(coin.Amount)
			bal := gaugeBalance.AmountOf(coin.Denom)

			b := sdkmath.LegacyNewDecFromInt(coin.Amount.Sub(bal))
			wouldBeBalance := timeRatio.Mul(coinAmountDec)
			newBalance := wouldBeBalance.Sub(b)

			amt64 := newBalance.TruncateInt64()
			if amt64 == 0 {
				continue
			}

			c := sdk.NewInt64Coin(coin.Denom, amt64)
			coinsToDistribute = coinsToDistribute.Add(c)
			err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, gaugeWallet, types.ModuleName, sdk.NewCoins(c))
			if err != nil {
				ctx.Logger().Error(errorsmod.Wrapf(err, "cannot send tokens from gauge to storage account").Error())
				continue
			}
		}
	})

	return coinsToDistribute
}

func (k Keeper) rewardAllProviders(ctx sdk.Context, totalSize int64, trackers []types.RewardTracker) {
	coins := k.pullTokensFromGauges(ctx)
	networkValue := sdkmath.LegacyNewDec(totalSize)

	for _, tracker := range trackers { // loop through a sorted list of providers
		worth := tracker.Size_
		prover := tracker.Provider
		providerValue := sdkmath.LegacyNewDec(worth)

		networkPercentage := providerValue.Quo(networkValue)
		pAddress, err := sdk.AccAddressFromBech32(prover)
		if err != nil {
			ctx.Logger().Error(errorsmod.Wrapf(err, "failed to convert prover address %s to bech32", prover).Error())
			continue
		}

		for _, coin := range coins {
			tokensValueOwed := networkPercentage.Mul(sdkmath.LegacyNewDecFromInt(coin.Amount)).TruncateInt()
			c := sdk.NewCoin(coin.Denom, tokensValueOwed)

			err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, pAddress, sdk.NewCoins(c))
			if err != nil {
				ctx.Logger().Error(errorsmod.Wrapf(err, "failed to send %s to %s", coins.String(), prover).Error())
				continue
			}
		}

	}
}

// ManageRewards pays out providers based on their proofs
func (k Keeper) ManageRewards(ctx sdk.Context) {
	var totalSize int64
	trackers := k.GetAllRewardTrackers(ctx)
	for _, tracker := range trackers {
		totalSize += tracker.Size_
	}

	if totalSize == 0 { // if there are no proofs to check, skip it all
		return
	}

	k.rewardAllProviders(ctx, totalSize, trackers)

	for _, tracker := range trackers {
		tracker.Size_ = 0
		k.SetRewardTracker(ctx, tracker) // reset tracker
	}
}

func (k Keeper) RunRewardBlock(ctx sdk.Context) {
	DayBlocks := k.GetParams(ctx).CheckWindow // checks more often than proofs take to catch them more frequently

	if ctx.BlockHeight()%DayBlocks > 0 { // runs once each window (usually a full days worth of blocks)
		ctx.Logger().Debug("skipping reward handling for this block")
		return
	}

	k.ManageRewards(ctx)
}
