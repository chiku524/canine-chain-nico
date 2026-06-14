package app

import (
	"bytes"
	"fmt"
	"math"

	sdkmath "cosmossdk.io/math"
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	storagetypes "github.com/jackalLabs/canine-chain/v5/x/storage/types"
)

// freeStorageMessageTypeURLs are storage msgs that pay no chain fee when a tx
// contains only these message types. Ported from JackalLabs/cosmos-sdk-new
// (DeductFeeDecorator patch on SDK 0.45).
var freeStorageMessageTypeURLs = map[string]struct{}{
	sdk.MsgTypeURL(&storagetypes.MsgPostProof{}):               {},
	sdk.MsgTypeURL(&storagetypes.MsgRequestAttestationForm{}):  {},
	sdk.MsgTypeURL(&storagetypes.MsgAttest{}):                  {},
	sdk.MsgTypeURL(&storagetypes.MsgReport{}):                  {},
}

func isFreeStorageTx(tx sdk.Tx) bool {
	for _, msg := range tx.GetMsgs() {
		if _, ok := freeStorageMessageTypeURLs[sdk.MsgTypeURL(msg)]; !ok {
			return false
		}
	}
	return len(tx.GetMsgs()) > 0
}

// JackalTxFeeChecker applies the default validator min-gas-price check, then
// returns zero effective fees for free storage-only transactions.
func JackalTxFeeChecker(ctx sdk.Context, tx sdk.Tx) (sdk.Coins, int64, error) {
	fee, priority, err := checkTxFeeWithValidatorMinGasPrices(ctx, tx)
	if err != nil {
		return nil, 0, err
	}
	if isFreeStorageTx(tx) {
		return sdk.Coins{}, priority, nil
	}
	return fee, priority, nil
}

// checkTxFeeWithValidatorMinGasPrices mirrors ante.checkTxFeeWithValidatorMinGasPrices.
func checkTxFeeWithValidatorMinGasPrices(ctx sdk.Context, tx sdk.Tx) (sdk.Coins, int64, error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return nil, 0, errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	feeCoins := feeTx.GetFee()
	gas := feeTx.GetGas()

	if ctx.IsCheckTx() {
		minGasPrices := ctx.MinGasPrices()
		if !minGasPrices.IsZero() {
			requiredFees := make(sdk.Coins, len(minGasPrices))
			glDec := sdkmath.LegacyNewDec(int64(gas))
			for i, gp := range minGasPrices {
				fee := gp.Amount.Mul(glDec)
				requiredFees[i] = sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt())
			}
			if !feeCoins.IsAnyGTE(requiredFees) {
				return nil, 0, errorsmod.Wrapf(sdkerrors.ErrInsufficientFee, "insufficient fees; got: %s required: %s", feeCoins, requiredFees)
			}
		}
	}

	return feeCoins, getTxPriority(feeCoins, int64(gas)), nil
}

func getTxPriority(fee sdk.Coins, gas int64) int64 {
	var priority int64
	for _, c := range fee {
		p := int64(math.MaxInt64)
		gasPrice := c.Amount.QuoRaw(gas)
		if gasPrice.IsInt64() {
			p = gasPrice.Int64()
		}
		if priority == 0 || p < priority {
			priority = p
		}
	}
	return priority
}

// JackalDeductFeeDecorator extends the SDK deduct-fee ante step with Jackal's
// free storage message fee waiver (legacy cosmos-sdk-new fork behavior).
type JackalDeductFeeDecorator struct {
	accountKeeper  ante.AccountKeeper
	bankKeeper     authtypes.BankKeeper
	feegrantKeeper ante.FeegrantKeeper
	txFeeChecker   ante.TxFeeChecker
}

func NewJackalDeductFeeDecorator(
	ak ante.AccountKeeper,
	bk authtypes.BankKeeper,
	fk ante.FeegrantKeeper,
	tfc ante.TxFeeChecker,
) JackalDeductFeeDecorator {
	if tfc == nil {
		tfc = JackalTxFeeChecker
	}
	return JackalDeductFeeDecorator{
		accountKeeper:  ak,
		bankKeeper:     bk,
		feegrantKeeper: fk,
		txFeeChecker:   tfc,
	}
}

func (dfd JackalDeductFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	if !simulate && ctx.BlockHeight() > 0 && feeTx.GetGas() == 0 {
		return ctx, errorsmod.Wrap(sdkerrors.ErrInvalidGasLimit, "must provide positive gas")
	}

	fee := feeTx.GetFee()
	priority := int64(0)
	if !simulate {
		var err error
		fee, priority, err = dfd.txFeeChecker(ctx, tx)
		if err != nil {
			return ctx, err
		}
	} else if isFreeStorageTx(tx) {
		fee = sdk.Coins{}
	}

	if err := dfd.checkDeductFee(ctx, tx, fee); err != nil {
		return ctx, err
	}

	return next(ctx.WithPriority(priority), tx, simulate)
}

func (dfd JackalDeductFeeDecorator) checkDeductFee(ctx sdk.Context, sdkTx sdk.Tx, fee sdk.Coins) error {
	feeTx, ok := sdkTx.(sdk.FeeTx)
	if !ok {
		return errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	if addr := dfd.accountKeeper.GetModuleAddress(authtypes.FeeCollectorName); addr == nil {
		return fmt.Errorf("fee collector module account (%s) has not been set", authtypes.FeeCollectorName)
	}

	feePayer := feeTx.FeePayer()
	feeGranter := feeTx.FeeGranter()
	deductFeesFrom := feePayer

	if len(feeGranter) != 0 {
		if dfd.feegrantKeeper == nil {
			return sdkerrors.ErrInvalidRequest.Wrap("fee grants are not enabled")
		} else if !bytes.Equal(feeGranter, feePayer) {
			err := dfd.feegrantKeeper.UseGrantedFees(ctx, feeGranter, feePayer, fee, sdkTx.GetMsgs())
			if err != nil {
				return errorsmod.Wrapf(err, "%s does not allow to pay fees for %s", sdk.AccAddress(feeGranter), sdk.AccAddress(feePayer))
			}
		}
		deductFeesFrom = feeGranter
	}

	deductFeesFromAcc := dfd.accountKeeper.GetAccount(ctx, deductFeesFrom)
	if deductFeesFromAcc == nil {
		return sdkerrors.ErrUnknownAddress.Wrapf("fee payer address: %s does not exist", deductFeesFrom)
	}

	if !fee.IsZero() && !isFreeStorageTx(sdkTx) {
		if err := dfd.bankKeeper.SendCoinsFromAccountToModule(ctx, deductFeesFromAcc.GetAddress(), authtypes.FeeCollectorName, fee); err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds, "%s", err.Error())
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeTx,
			sdk.NewAttribute(sdk.AttributeKeyFee, fee.String()),
			sdk.NewAttribute(sdk.AttributeKeyFeePayer, sdk.AccAddress(deductFeesFrom).String()),
		),
	})

	return nil
}
