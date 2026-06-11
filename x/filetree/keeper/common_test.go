package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtime "github.com/cometbft/cometbft/types/time"

	"github.com/cosmos/cosmos-sdk/baseapp"

	storetypes "cosmossdk.io/store/types"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"

	moduletestutil "github.com/jackalLabs/canine-chain/v5/types/module/testutil" // when importing from sdk,'go mod tidy' keeps trying to import from v0.46.

	sdk "github.com/cosmos/cosmos-sdk/types"
	canineglobaltestutil "github.com/jackalLabs/canine-chain/v5/testutil"
	"github.com/jackalLabs/canine-chain/v5/x/filetree/keeper"
	types "github.com/jackalLabs/canine-chain/v5/x/filetree/types"
)

// SetupFiletreeKeeper creates a filetreeKeeper as well as all its dependencies.
func SetupFiletreeKeeper(t *testing.T) (
	*keeper.Keeper,
	moduletestutil.TestEncodingConfig,
	sdk.Context,
) {
	key := storetypes.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	testCtx := canineglobaltestutil.DefaultContextWithDB(t, storetypes.NewTransientStoreKey("transient_test"), key)
	ctx := testCtx.Ctx.WithBlockHeader(tmproto.Header{Time: tmtime.Now()})

	encCfg := moduletestutil.MakeTestEncodingConfig()
	types.RegisterInterfaces(encCfg.InterfaceRegistry)

	// Create MsgServiceRouter, but don't populate it before creating the filetree keeper.
	msr := baseapp.NewMsgServiceRouter()

	paramsSubspace := typesparams.NewSubspace(encCfg.Codec,
		types.Amino,
		key,
		memStoreKey,
		"FiletreeParams",
	)

	// filetree keeper initializations
	filetreeKeeper := keeper.NewKeeper(encCfg.Codec, key, memStoreKey, paramsSubspace)
	filetreeKeeper.SetParams(ctx, types.DefaultParams())

	// Register all handlers for the MegServiceRouter.
	msr.SetInterfaceRegistry(encCfg.InterfaceRegistry)
	types.RegisterMsgServer(msr, keeper.NewMsgServerImpl(*filetreeKeeper))

	return filetreeKeeper, encCfg, ctx
}
