package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackalLabs/canine-chain/v5/docs"
	"github.com/jackalLabs/canine-chain/v5/docs/openapiconsole"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	reflectionv1 "cosmossdk.io/api/cosmos/reflection/v1"
	"cosmossdk.io/log/v2"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cast"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/posthandler"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ica "github.com/cosmos/ibc-go/v11/modules/apps/27-interchain-accounts"
	icacontroller "github.com/cosmos/ibc-go/v11/modules/apps/27-interchain-accounts/controller"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v11/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v11/modules/apps/27-interchain-accounts/controller/types"
	icahost "github.com/cosmos/ibc-go/v11/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v11/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v11/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v11/modules/apps/27-interchain-accounts/types"
	ibccallbacks "github.com/cosmos/ibc-go/v11/modules/apps/callbacks"
	"github.com/cosmos/ibc-go/v11/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v11/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v11/modules/apps/transfer/types"
	transferv2 "github.com/cosmos/ibc-go/v11/modules/apps/transfer/v2"
	ibc "github.com/cosmos/ibc-go/v11/modules/core"
	porttypes "github.com/cosmos/ibc-go/v11/modules/core/05-port/types"
	ibcapi "github.com/cosmos/ibc-go/v11/modules/core/api"
	ibcexported "github.com/cosmos/ibc-go/v11/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v11/modules/core/keeper"
	ibctm "github.com/cosmos/ibc-go/v11/modules/light-clients/07-tendermint"

	mint "github.com/jackalLabs/canine-chain/v5/x/jklmint"
	mintkeeper "github.com/jackalLabs/canine-chain/v5/x/jklmint/keeper"
	minttypes "github.com/jackalLabs/canine-chain/v5/x/jklmint/types"
	rnsmodule "github.com/jackalLabs/canine-chain/v5/x/rns"
	rnsmodulekeeper "github.com/jackalLabs/canine-chain/v5/x/rns/keeper"
	rnsmoduletypes "github.com/jackalLabs/canine-chain/v5/x/rns/types"
	oraclemodule "github.com/jackalLabs/canine-chain/v5/x/oracle"
	oraclemodulekeeper "github.com/jackalLabs/canine-chain/v5/x/oracle/keeper"
	oraclemoduletypes "github.com/jackalLabs/canine-chain/v5/x/oracle/types"
	storagemodule "github.com/jackalLabs/canine-chain/v5/x/storage"
	storagemodulekeeper "github.com/jackalLabs/canine-chain/v5/x/storage/keeper"
	storagemoduletypes "github.com/jackalLabs/canine-chain/v5/x/storage/types"
	filetreemodule "github.com/jackalLabs/canine-chain/v5/x/filetree"
	filetreemodulekeeper "github.com/jackalLabs/canine-chain/v5/x/filetree/keeper"
	filetreemoduletypes "github.com/jackalLabs/canine-chain/v5/x/filetree/types"
	notificationsmodule "github.com/jackalLabs/canine-chain/v5/x/notifications"
	notificationsmodulekeeper "github.com/jackalLabs/canine-chain/v5/x/notifications/keeper"
	notificationsmoduletypes "github.com/jackalLabs/canine-chain/v5/x/notifications/types"

	wasmappparams "github.com/jackalLabs/canine-chain/v5/app/params"
	owasm "github.com/jackalLabs/canine-chain/v5/wasmbinding"
)

const appName = "JackalApp"

var (
	NodeDir      = ".canine"
	Bech32Prefix = "jkl"

	ProposalsEnabled        = "false"
	EnableSpecificProposals = ""
)

func GetEnabledProposals() []wasmtypes.ProposalType {
	if EnableSpecificProposals == "" {
		if ProposalsEnabled == "true" {
			return wasmtypes.EnableAllProposals
		}
		return wasmtypes.DisableAllProposals
	}
	chunks := strings.Split(EnableSpecificProposals, ",")
	proposals := make([]wasmtypes.ProposalType, 0, len(chunks))
	for _, c := range chunks {
		proposals = append(proposals, wasmtypes.ProposalType(strings.TrimSpace(c)))
	}
	return proposals
}

var (
	DefaultNodeHome = os.ExpandEnv("$HOME/") + NodeDir

	Bech32PrefixAccAddr  = Bech32Prefix
	Bech32PrefixAccPub   = Bech32Prefix + sdk.PrefixPublic
	Bech32PrefixValAddr  = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator
	Bech32PrefixValPub   = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	Bech32PrefixConsAddr = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus
	Bech32PrefixConsPub  = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
)

var (
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		distr.AppModuleBasic{},
		mint.AppModuleBasic{},
		gov.NewAppModuleBasic([]govclient.ProposalHandler{
			paramsclient.ProposalHandler,
		}),
		params.AppModuleBasic{},
		slashing.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		ibc.AppModuleBasic{},
		ibctm.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		transfer.AppModuleBasic{},
		vesting.AppModuleBasic{},
		consensus.AppModuleBasic{},
		wasm.AppModuleBasic{},
		ica.AppModuleBasic{},
		rnsmodule.AppModuleBasic{},
		storagemodule.AppModuleBasic{},
		filetreemodule.AppModuleBasic{},
		oraclemodule.AppModuleBasic{},
		notificationsmodule.AppModuleBasic{},
	)

	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:                 nil,
		distrtypes.ModuleName:                      nil,
		minttypes.ModuleName:                       {authtypes.Minter},
		stakingtypes.BondedPoolName:                {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName:             {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:                        {authtypes.Burner},
		ibctransfertypes.ModuleName:                {authtypes.Minter, authtypes.Burner},
		wasmtypes.ModuleName:                       {authtypes.Burner},
		rnsmoduletypes.ModuleName:                  {authtypes.Minter, authtypes.Burner},
		storagemoduletypes.ModuleName:              {authtypes.Minter, authtypes.Burner},
		oraclemoduletypes.ModuleName:               nil,
		notificationsmoduletypes.ModuleName:        nil,
		icatypes.ModuleName:                        nil,
		storagemoduletypes.CollateralCollectorName: nil,
	}
)

var (
	_ runtime.AppI            = (*JackalApp)(nil)
	_ servertypes.Application = (*JackalApp)(nil)
)

type JackalApp struct {
	*baseapp.BaseApp
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	InterfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	AccountKeeper         authkeeper.AccountKeeper
	BankKeeper            bankkeeper.BaseKeeper
	stakingKeeper         *stakingkeeper.Keeper
	slashingKeeper        slashingkeeper.Keeper
	MintKeeper            mintkeeper.Keeper
	distrKeeper           distrkeeper.Keeper
	govKeeper             govkeeper.Keeper
	upgradeKeeper         *upgradekeeper.Keeper
	paramsKeeper          paramskeeper.Keeper
	ConsensusParamsKeeper consensusparamkeeper.Keeper
	evidenceKeeper        evidencekeeper.Keeper
	ibcKeeper             *ibckeeper.Keeper
	transferKeeper        *ibctransferkeeper.Keeper
	feeGrantKeeper        feegrantkeeper.Keeper
	authzKeeper           authzkeeper.Keeper
	wasmKeeper            wasmkeeper.Keeper

	ICAControllerKeeper *icacontrollerkeeper.Keeper
	ICAHostKeeper       *icahostkeeper.Keeper

	RnsKeeper           rnsmodulekeeper.Keeper
	OracleKeeper        oraclemodulekeeper.Keeper
	StorageKeeper       storagemodulekeeper.Keeper
	FileTreeKeeper      filetreemodulekeeper.Keeper
	NotificationsKeeper notificationsmodulekeeper.Keeper

	mm           *module.Manager
	sm           *module.SimulationManager
	configurator module.Configurator
}

func NewJackalApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	encodingConfig wasmappparams.EncodingConfig,
	_ []wasmtypes.ProposalType,
	appOpts servertypes.AppOptions,
	wasmOpts []wasmkeeper.Option,
	baseAppOptions ...func(*baseapp.BaseApp),
) *JackalApp {
	appCodec := encodingConfig.Marshaler
	legacyAmino := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry
	txConfig := encodingConfig.TxConfig

	bApp := baseapp.NewBaseApp(appName, logger, db, txConfig.TxDecoder(), baseAppOptions...)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)
	bApp.SetTxEncoder(txConfig.TxEncoder())

	keys := storetypes.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		minttypes.StoreKey, distrtypes.StoreKey, slashingtypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, consensusparamtypes.StoreKey, upgradetypes.StoreKey,
		evidencetypes.StoreKey, feegrant.StoreKey, authzkeeper.StoreKey,
		ibcexported.StoreKey, ibctransfertypes.StoreKey,
		wasmtypes.StoreKey, icacontrollertypes.StoreKey, icahosttypes.StoreKey,
		rnsmoduletypes.StoreKey, storagemoduletypes.StoreKey, filetreemoduletypes.StoreKey,
		oraclemoduletypes.StoreKey, notificationsmoduletypes.StoreKey,
	)
	tkeys := storetypes.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := storetypes.NewMemoryStoreKeys(
		oraclemoduletypes.MemStoreKey,
		storagemoduletypes.MemStoreKey,
		rnsmoduletypes.MemStoreKey,
		filetreemoduletypes.MemStoreKey,
		notificationsmoduletypes.MemStoreKey,
	)

	if err := bApp.RegisterStreamingServices(appOpts, keys); err != nil {
		panic(err)
	}

	app := &JackalApp{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		txConfig:          txConfig,
		InterfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	app.paramsKeeper = initParamsKeeper(appCodec, legacyAmino, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])
	app.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[consensusparamtypes.StoreKey]),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		runtime.EventService{},
	)
	bApp.SetParamStore(app.ConsensusParamsKeeper.ParamsStore)

	govAuthority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		maccPerms,
		authcodec.NewBech32Codec(Bech32Prefix),
		Bech32Prefix,
		govAuthority,
	)
	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		app.AccountKeeper,
		BlockedAddresses(),
		govAuthority,
		logger,
	)
	app.authzKeeper = authzkeeper.NewKeeper(
		runtime.NewKVStoreService(keys[authzkeeper.StoreKey]),
		appCodec,
		app.MsgServiceRouter(),
		app.AccountKeeper,
	)
	app.feeGrantKeeper = feegrantkeeper.NewKeeper(appCodec, runtime.NewKVStoreService(keys[feegrant.StoreKey]), app.AccountKeeper)

	app.stakingKeeper = stakingkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[stakingtypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		govAuthority,
		authcodec.NewBech32Codec(Bech32PrefixValAddr),
		authcodec.NewBech32Codec(Bech32PrefixConsAddr),
	)

	app.MintKeeper = mintkeeper.NewKeeper(
		appCodec,
		keys[minttypes.StoreKey],
		app.GetSubspace(minttypes.ModuleName),
		app.stakingKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.FeeCollectorName,
		storagemoduletypes.ModuleName,
	)
	mintModule := mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(minttypes.ModuleName))

	app.distrKeeper = distrkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[distrtypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		app.stakingKeeper,
		authtypes.FeeCollectorName,
		govAuthority,
	)
	app.slashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		legacyAmino,
		runtime.NewKVStoreService(keys[slashingtypes.StoreKey]),
		app.stakingKeeper,
		govAuthority,
	)

	app.stakingKeeper.SetHooks(stakingtypes.NewMultiStakingHooks(
		app.distrKeeper.Hooks(),
		app.slashingKeeper.Hooks(),
	))

	app.upgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		runtime.NewKVStoreService(keys[upgradetypes.StoreKey]),
		appCodec,
		homePath,
		app.BaseApp,
		govAuthority,
	)

	app.ibcKeeper = ibckeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[ibcexported.StoreKey]),
		app.upgradeKeeper,
		govAuthority,
	)

	app.transferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		app.AccountKeeper.AddressCodec(),
		runtime.NewKVStoreService(keys[ibctransfertypes.StoreKey]),
		app.ibcKeeper.ChannelKeeper,
		app.MsgServiceRouter(),
		app.AccountKeeper,
		app.BankKeeper,
		govAuthority,
	)
	transferModule := transfer.NewAppModule(app.transferKeeper)

	app.ICAHostKeeper = icahostkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[icahosttypes.StoreKey]),
		app.ibcKeeper.ChannelKeeper,
		app.AccountKeeper,
		app.MsgServiceRouter(),
		app.GRPCQueryRouter(),
		govAuthority,
	)
	app.ICAControllerKeeper = icacontrollerkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[icacontrollertypes.StoreKey]),
		app.ibcKeeper.ChannelKeeper,
		app.MsgServiceRouter(),
		govAuthority,
	)

	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[evidencetypes.StoreKey]),
		app.stakingKeeper,
		app.slashingKeeper,
		app.AccountKeeper.AddressCodec(),
		runtime.ProvideCometInfoService(),
	)
	app.evidenceKeeper = *evidenceKeeper

	app.RnsKeeper = *rnsmodulekeeper.NewKeeper(
		appCodec, keys[rnsmoduletypes.StoreKey], app.GetSubspace(rnsmoduletypes.ModuleName), app.BankKeeper,
	)
	rnsModule := rnsmodule.NewAppModule(appCodec, app.RnsKeeper, app.AccountKeeper, app.BankKeeper)

	app.OracleKeeper = *oraclemodulekeeper.NewKeeper(
		appCodec, keys[oraclemoduletypes.StoreKey], app.GetSubspace(oraclemoduletypes.ModuleName), app.BankKeeper,
	)
	oracleModule := oraclemodule.NewAppModule(appCodec, app.OracleKeeper, app.AccountKeeper, app.BankKeeper)

	app.StorageKeeper = *storagemodulekeeper.NewKeeper(
		appCodec, keys[storagemoduletypes.StoreKey], app.GetSubspace(storagemoduletypes.ModuleName),
		app.BankKeeper, app.AccountKeeper, app.OracleKeeper, app.RnsKeeper, authtypes.FeeCollectorName,
	)
	storageModule := storagemodule.NewAppModule(
		appCodec, app.StorageKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(storagemoduletypes.ModuleName),
	)

	app.FileTreeKeeper = *filetreemodulekeeper.NewKeeper(
		appCodec, keys[filetreemoduletypes.StoreKey], memKeys[filetreemoduletypes.MemStoreKey],
		app.GetSubspace(filetreemoduletypes.ModuleName),
	)
	filetreeModule := filetreemodule.NewAppModule(appCodec, app.FileTreeKeeper, app.AccountKeeper, app.BankKeeper)

	app.NotificationsKeeper = *notificationsmodulekeeper.NewKeeper(
		appCodec, keys[notificationsmoduletypes.StoreKey], memKeys[notificationsmoduletypes.MemStoreKey],
		app.GetSubspace(notificationsmoduletypes.ModuleName), app.RnsKeeper,
	)
	notificationsModule := notificationsmodule.NewAppModule(appCodec, app.NotificationsKeeper, app.AccountKeeper, app.BankKeeper)

	wasmDir := filepath.Join(homePath, "wasm")
	nodeConfig, err := wasm.ReadNodeConfig(appOpts)
	if err != nil {
		panic(fmt.Sprintf("error while reading wasm config: %s", err))
	}

	wasmOpts = append(owasm.RegisterCustomPlugins(&app.FileTreeKeeper, &app.StorageKeeper, &app.NotificationsKeeper), wasmOpts...)

	app.wasmKeeper = wasmkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[wasmtypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		app.stakingKeeper,
		distrkeeper.NewQuerier(app.distrKeeper),
		app.ibcKeeper.ChannelKeeper,
		app.ibcKeeper.ChannelKeeper,
		app.ibcKeeper.ChannelKeeperV2,
		app.transferKeeper,
		app.MsgServiceRouter(),
		app.GRPCQueryRouter(),
		wasmDir,
		nodeConfig,
		wasmtypes.VMConfig{},
		wasmkeeper.BuiltInCapabilities(),
		govAuthority,
		wasmOpts...,
	)

	govConfig := govtypes.DefaultConfig()
	govKeeper := govkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[govtypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		app.distrKeeper,
		app.MsgServiceRouter(),
		govConfig,
		govAuthority,
		govkeeper.NewDefaultCalculateVoteResultsAndVotingPower(app.stakingKeeper),
	)
	app.govKeeper = *govKeeper.SetHooks(govtypes.NewMultiGovHooks())

	wasmStackIBCHandler := wasm.NewIBCHandler(app.wasmKeeper, app.ibcKeeper.ChannelKeeper, app.transferKeeper, app.ibcKeeper.ChannelKeeper)

	var icaControllerStack porttypes.IBCModule
	var noAuthzModule porttypes.IBCModule
	icaStackBuilder := porttypes.NewIBCStackBuilder(app.ibcKeeper.ChannelKeeper)
	icaStackBuilder.Base(
		icacontroller.NewIBCMiddlewareWithAuth(noAuthzModule, app.ICAControllerKeeper)).Next(
		icacontroller.NewIBCMiddlewareWithAuth(icaControllerStack, app.ICAControllerKeeper)).Next(
		ibccallbacks.NewIBCMiddleware(wasmStackIBCHandler, wasm.DefaultMaxIBCCallbackGas),
	)
	icaControllerStack = icaStackBuilder.Build()
	icaICS4Wrapper := icaControllerStack.(porttypes.ICS4Wrapper)
	app.ICAControllerKeeper.WithICS4Wrapper(icaICS4Wrapper)

	icaHostStack := icahost.NewIBCModule(app.ICAHostKeeper)

	var transferStack porttypes.IBCModule
	transferStackBuilder := porttypes.NewIBCStackBuilder(app.ibcKeeper.ChannelKeeper)
	transferStackBuilder.Base(
		transfer.NewIBCModule(app.transferKeeper)).Next(
		ibccallbacks.NewIBCMiddleware(wasmStackIBCHandler, wasm.DefaultMaxIBCCallbackGas),
	)
	transferStack = transferStackBuilder.Build()
	transferICS4Wrapper := transferStack.(porttypes.ICS4Wrapper)
	app.transferKeeper.WithICS4Wrapper(transferICS4Wrapper)

	ibcRouter := porttypes.NewRouter().
		AddRoute(ibctransfertypes.ModuleName, transferStack).
		AddRoute(wasmtypes.ModuleName, wasmStackIBCHandler).
		AddRoute(icacontrollertypes.SubModuleName, icaControllerStack).
		AddRoute(icahosttypes.SubModuleName, icaHostStack)
	app.ibcKeeper.SetRouter(ibcRouter)

	ibcRouterV2 := ibcapi.NewRouter()
	ibcRouterV2 = ibcRouterV2.
		AddRoute(ibctransfertypes.PortID, transferv2.NewIBCModule(app.transferKeeper)).
		AddPrefixRoute(wasmkeeper.PortIDPrefixV2, wasmkeeper.NewIBC2Handler(app.wasmKeeper))
	app.ibcKeeper.SetRouterV2(ibcRouterV2)

	clientKeeper := app.ibcKeeper.ClientKeeper
	storeProvider := app.ibcKeeper.ClientKeeper.GetStoreProvider()
	tmLightClientModule := ibctm.NewLightClientModule(appCodec, storeProvider)
	clientKeeper.AddRoute(ibctm.ModuleName, &tmLightClientModule)

	icaModule := ica.NewAppModule(app.ICAControllerKeeper, app.ICAHostKeeper)

	app.mm = module.NewManager(
		genutil.NewAppModule(app.AccountKeeper, app.stakingKeeper, app, txConfig),
		auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts, app.GetSubspace(authtypes.ModuleName)),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.feeGrantKeeper, app.InterfaceRegistry),
		gov.NewAppModule(appCodec, &app.govKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
		mintModule,
		slashing.NewAppModule(appCodec, app.slashingKeeper, app.AccountKeeper, app.BankKeeper, app.stakingKeeper, app.GetSubspace(slashingtypes.ModuleName), app.InterfaceRegistry),
		distr.NewAppModule(appCodec, app.distrKeeper, app.AccountKeeper, app.BankKeeper, app.stakingKeeper, app.GetSubspace(distrtypes.ModuleName)),
		staking.NewAppModule(appCodec, app.stakingKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		upgrade.NewAppModule(app.upgradeKeeper, app.AccountKeeper.AddressCodec()),
		evidence.NewAppModule(app.evidenceKeeper),
		params.NewAppModule(app.paramsKeeper),
		authzmodule.NewAppModule(appCodec, app.authzKeeper, app.AccountKeeper, app.BankKeeper, app.InterfaceRegistry),
		consensus.NewAppModule(appCodec, app.ConsensusParamsKeeper),
		wasm.NewAppModule(appCodec, &app.wasmKeeper, app.stakingKeeper, app.AccountKeeper, app.BankKeeper, app.MsgServiceRouter(), app.GetSubspace(wasmtypes.ModuleName)),
		ibc.NewAppModule(app.ibcKeeper),
		transferModule,
		icaModule,
		ibctm.NewAppModule(tmLightClientModule),
		rnsModule,
		storageModule,
		filetreeModule,
		oracleModule,
		notificationsModule,
	)

	app.mm.SetOrderPreBlockers(upgradetypes.ModuleName, authtypes.ModuleName)

	app.mm.SetOrderBeginBlockers(
		upgradetypes.ModuleName, minttypes.ModuleName, distrtypes.ModuleName,
		slashingtypes.ModuleName, evidencetypes.ModuleName, stakingtypes.ModuleName,
		authtypes.ModuleName, banktypes.ModuleName, govtypes.ModuleName,
		genutiltypes.ModuleName, authz.ModuleName, feegrant.ModuleName, paramstypes.ModuleName,
		vestingtypes.ModuleName, consensusparamtypes.ModuleName,
		icatypes.ModuleName, ibctransfertypes.ModuleName, ibcexported.ModuleName,
		rnsmoduletypes.ModuleName, storagemoduletypes.ModuleName, filetreemoduletypes.ModuleName,
		oraclemoduletypes.ModuleName, notificationsmoduletypes.ModuleName, wasmtypes.ModuleName,
	)

	app.mm.SetOrderEndBlockers(
		govtypes.ModuleName, stakingtypes.ModuleName,
		authtypes.ModuleName, banktypes.ModuleName, distrtypes.ModuleName, slashingtypes.ModuleName,
		minttypes.ModuleName, genutiltypes.ModuleName, evidencetypes.ModuleName, authz.ModuleName,
		feegrant.ModuleName, paramstypes.ModuleName, vestingtypes.ModuleName, consensusparamtypes.ModuleName,
		upgradetypes.ModuleName, icatypes.ModuleName, ibctransfertypes.ModuleName,
		ibcexported.ModuleName, rnsmoduletypes.ModuleName, storagemoduletypes.ModuleName,
		filetreemoduletypes.ModuleName, oraclemoduletypes.ModuleName, notificationsmoduletypes.ModuleName,
		wasmtypes.ModuleName,
	)

	genesisModuleOrder := []string{
		authtypes.ModuleName, banktypes.ModuleName, distrtypes.ModuleName,
		stakingtypes.ModuleName, slashingtypes.ModuleName, govtypes.ModuleName, minttypes.ModuleName,
		genutiltypes.ModuleName, evidencetypes.ModuleName, authz.ModuleName,
		feegrant.ModuleName, paramstypes.ModuleName, upgradetypes.ModuleName, vestingtypes.ModuleName,
		consensusparamtypes.ModuleName, ibctransfertypes.ModuleName, ibcexported.ModuleName,
		icatypes.ModuleName, rnsmoduletypes.ModuleName, storagemoduletypes.ModuleName,
		filetreemoduletypes.ModuleName, oraclemoduletypes.ModuleName, notificationsmoduletypes.ModuleName,
		wasmtypes.ModuleName,
	}
	app.mm.SetOrderInitGenesis(genesisModuleOrder...)
	app.mm.SetOrderExportGenesis(genesisModuleOrder...)

	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	if err := app.mm.RegisterServices(app.configurator); err != nil {
		panic(err)
	}

	app.registerTestnetUpgradeHandlers()
	app.registerMainnetUpgradeHandlers()

	overrideModules := simulationOverrides(map[string]module.AppModuleSimulation{
		authtypes.ModuleName: auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts, app.GetSubspace(authtypes.ModuleName)),
	}, app.mm.Modules)
	app.sm = module.NewSimulationManagerFromAppModules(app.mm.Modules, overrideModules)
	app.sm.RegisterStoreDecoders()

	autocliv1.RegisterQueryServer(app.GRPCQueryRouter(), runtimeservices.NewAutoCLIQueryService(app.mm.Modules))
	reflectionSvc, err := runtimeservices.NewReflectionService()
	if err != nil {
		panic(err)
	}
	reflectionv1.RegisterReflectionServiceServer(app.GRPCQueryRouter(), reflectionSvc)

	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	app.setAnteHandler(txConfig, nodeConfig, keys[wasmtypes.StoreKey])
	app.setPostHandler()

	app.SetInitChainer(app.InitChainer)
	app.SetPreBlocker(app.PreBlocker)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	if manager := app.SnapshotManager(); manager != nil {
		if err := manager.RegisterExtensions(
			wasmkeeper.NewWasmSnapshotter(app.CommitMultiStore(), &app.wasmKeeper),
		); err != nil {
			panic(fmt.Errorf("failed to register snapshot extension: %w", err))
		}
	}

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			panic(fmt.Errorf("failed to load latest version: %w", err))
		}
		ctx := app.NewContextLegacy(true, tmproto.Header{})
		if err := app.wasmKeeper.InitializePinnedCodes(ctx); err != nil {
			panic(fmt.Errorf("failed initialize pinned codes: %w", err))
		}
	}

	return app
}

func (app *JackalApp) setAnteHandler(txConfig client.TxConfig, nodeConfig wasmtypes.NodeConfig, txCounterStoreKey *storetypes.KVStoreKey) {
	anteHandler, err := NewAnteHandler(HandlerOptions{
		HandlerOptions: ante.HandlerOptions{
			AccountKeeper:   app.AccountKeeper,
			BankKeeper:      app.BankKeeper,
			FeegrantKeeper:  app.feeGrantKeeper,
			SignModeHandler: txConfig.SignModeHandler(),
			SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
			TxFeeChecker:    JackalTxFeeChecker,
		},
		IBCKeeper:             app.ibcKeeper,
		WasmKeeper:            &app.wasmKeeper,
		NodeConfig:            &nodeConfig,
		TXCounterStoreService: runtime.NewKVStoreService(txCounterStoreKey),
	})
	if err != nil {
		panic(fmt.Errorf("failed to create AnteHandler: %w", err))
	}
	app.SetAnteHandler(anteHandler)
}

func (app *JackalApp) setPostHandler() {
	postHandler, err := posthandler.NewPostHandler(posthandler.HandlerOptions{})
	if err != nil {
		panic(err)
	}
	app.SetPostHandler(postHandler)
}

func (app *JackalApp) Name() string { return app.BaseApp.Name() }

func (app *JackalApp) PreBlocker(ctx sdk.Context, _ *abci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
	return app.mm.PreBlock(ctx)
}

func (app *JackalApp) BeginBlocker(ctx sdk.Context) (sdk.BeginBlock, error) {
	return app.mm.BeginBlock(ctx)
}

func (app *JackalApp) EndBlocker(ctx sdk.Context) (sdk.EndBlock, error) {
	return app.mm.EndBlock(ctx)
}

func (app *JackalApp) InitChainer(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
	var genesisState GenesisState
	if err := json.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}
	if err := app.upgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap()); err != nil {
		panic(err)
	}
	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

func (app *JackalApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

func (app *JackalApp) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

func (app *JackalApp) AppCodec() codec.Codec {
	return app.appCodec
}

func (app *JackalApp) TxConfig() client.TxConfig {
	return app.txConfig
}

func (app *JackalApp) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.paramsKeeper.GetSubspace(moduleName)
	return subspace
}

func (app *JackalApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

func (app *JackalApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	cmtservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	nodeservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	if apiConfig.Swagger {
		apiSvr.Router.Handle("/swagger-ui/swagger.yaml", http.FileServer(http.FS(docs.Docs)))
		apiSvr.Router.HandleFunc("/", openapiconsole.Handler(appName, "/swagger-ui/swagger.yaml"))
	}
}

func (app *JackalApp) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.GRPCQueryRouter(), clientCtx, app.Simulate, app.InterfaceRegistry)
}

func (app *JackalApp) RegisterTendermintService(clientCtx client.Context) {
	cmtApp := server.NewCometABCIWrapper(app)
	cmtservice.RegisterTendermintService(clientCtx, app.GRPCQueryRouter(), app.InterfaceRegistry, cmtApp.Query)
}

func (app *JackalApp) RegisterNodeService(clientCtx client.Context, cfg config.Config) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter(), cfg, func() int64 {
		return app.CommitMultiStore().EarliestVersion()
	})
}

func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}
	return dupMaccPerms
}

func BlockedAddresses() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range GetMaccPerms() {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}
	delete(modAccAddrs, authtypes.NewModuleAddress(govtypes.ModuleName).String())
	return modAccAddrs
}

func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govv1.ParamKeyTable()) //nolint:staticcheck
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibcexported.ModuleName)
	paramsKeeper.Subspace(icacontrollertypes.SubModuleName)
	paramsKeeper.Subspace(icahosttypes.SubModuleName)
	paramsKeeper.Subspace(rnsmoduletypes.ModuleName)
	paramsKeeper.Subspace(oraclemoduletypes.ModuleName)
	paramsKeeper.Subspace(storagemoduletypes.ModuleName)
	paramsKeeper.Subspace(filetreemoduletypes.ModuleName)
	paramsKeeper.Subspace(notificationsmoduletypes.ModuleName)
	paramsKeeper.Subspace(wasmtypes.ModuleName)

	return paramsKeeper
}

func GetWasmOpts(appOpts servertypes.AppOptions) []wasmkeeper.Option {
	var wasmOpts []wasmkeeper.Option
	if cast.ToBool(appOpts.Get("telemetry.enabled")) {
		wasmOpts = append(wasmOpts, wasmkeeper.WithVMCacheMetrics(prometheus.DefaultRegisterer))
	}
	wasmOpts = append(wasmOpts, wasmkeeper.WithGasRegister(NewJackalWasmGasRegister()))
	return wasmOpts
}
