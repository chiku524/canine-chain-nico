package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func (suite *MintTestSuite) TestBlockMint() {
	suite.SetupTest()
	app, ctx, k := suite.app, suite.ctx, suite.app.MintKeeper
	denom := k.GetParams(ctx).MintDenom
	feeAccount := app.AccountKeeper.GetModuleAccount(ctx, authtypes.FeeCollectorName)
	feeBalanceBefore := app.BankKeeper.GetBalance(ctx, feeAccount.GetAddress(), denom)
	suite.Require().True(feeBalanceBefore.Amount.IsZero())
	supplyBefore := app.BankKeeper.GetSupply(ctx, denom)

	k.BlockMint(ctx)

	feeBalanceAfter := app.BankKeeper.GetBalance(ctx, feeAccount.GetAddress(), denom)
	suite.Require().Equal(sdkmath.NewInt(3360000), feeBalanceAfter.Amount.Sub(feeBalanceBefore.Amount))
	supplyAfter := app.BankKeeper.GetSupply(ctx, denom)
	suite.Require().Equal(sdkmath.NewInt(4_200_000), supplyAfter.Amount.Sub(supplyBefore.Amount))
	// After BlockMint we now have exactly 3.6JKL in the fee collector account
}

func (suite *MintTestSuite) TestNoProviderBlockMint() {
	suite.SetupTest()
	app, ctx, k := suite.app, suite.ctx, suite.app.MintKeeper

	params := k.GetParams(ctx)
	params.StorageProviderRatio = 0
	k.SetParams(ctx, params)

	denom := k.GetParams(ctx).MintDenom

	pr := k.GetParams(ctx).StorageProviderRatio
	suite.Require().Equal(int64(0), pr)

	feeAccount := app.AccountKeeper.GetModuleAccount(ctx, authtypes.FeeCollectorName)
	feeBalanceBefore := app.BankKeeper.GetBalance(ctx, feeAccount.GetAddress(), denom)
	suite.Require().True(feeBalanceBefore.Amount.IsZero())
	supplyBefore := app.BankKeeper.GetSupply(ctx, denom)

	k.BlockMint(ctx)

	feeBalanceAfter := app.BankKeeper.GetBalance(ctx, feeAccount.GetAddress(), denom)

	suite.T().Log(params.TokensPerBlock)
	suite.Require().Equal(sdkmath.NewInt(3360000), feeBalanceAfter.Amount.Sub(feeBalanceBefore.Amount))
	supplyAfter := app.BankKeeper.GetSupply(ctx, denom)
	suite.Require().Equal(sdkmath.NewInt(4_200_000), supplyAfter.Amount.Sub(supplyBefore.Amount))
	// After BlockMint we now have exactly 3.6JKL in the fee collector account
}

func (suite *MintTestSuite) TestDecRatios() {
	suite.SetupTest()

	stakerRatio := sdkmath.LegacyNewDec(80).QuoInt64(100)

	s, err := sdkmath.LegacyNewDecFromStr("0.8")
	suite.Require().NoError(err)

	suite.Require().Equal(s, stakerRatio)

	i := stakerRatio.MulInt64(4_200_000)

	suite.Require().Equal(int64(3360000), i.TruncateInt64())
}
