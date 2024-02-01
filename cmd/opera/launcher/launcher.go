package launcher

import (
	"fmt"
	"github.com/Fantom-foundation/go-opera/cmd/opera/launcher/diskusage"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"os"
	"os/signal"
	"path"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/console/prompt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p/discover/discfilter"
	"github.com/ethereum/go-ethereum/params"
	"gopkg.in/urfave/cli.v1"

	evmetrics "github.com/ethereum/go-ethereum/metrics"

	"github.com/Fantom-foundation/go-opera/cmd/opera/launcher/metrics"
	"github.com/Fantom-foundation/go-opera/cmd/opera/launcher/tracing"
	"github.com/Fantom-foundation/go-opera/debug"
	"github.com/Fantom-foundation/go-opera/evmcore"
	"github.com/Fantom-foundation/go-opera/flags"
	"github.com/Fantom-foundation/go-opera/gossip"
	"github.com/Fantom-foundation/go-opera/gossip/emitter"
	"github.com/Fantom-foundation/go-opera/integration"
	"github.com/Fantom-foundation/go-opera/opera/genesis"
	"github.com/Fantom-foundation/go-opera/opera/genesisstore"
	"github.com/Fantom-foundation/go-opera/utils/errlock"
	"github.com/Fantom-foundation/go-opera/valkeystore"
	_ "github.com/Fantom-foundation/go-opera/version"
)

const (
	// clientIdentifier to advertise over the network.
	clientIdentifier = "go-opera"
)

var (
	// Git SHA1 commit hash of the release (set via linker flags).
	gitCommit = ""
	gitDate   = ""
	// The app that holds all commands and flags.
	app = flags.NewApp(gitCommit, gitDate, "the go-opera command line interface")

	nodeFlags        []cli.Flag
	testFlags        []cli.Flag
	gpoFlags         []cli.Flag
	accountFlags     []cli.Flag
	performanceFlags []cli.Flag
	networkingFlags  []cli.Flag
	txpoolFlags      []cli.Flag
	operaFlags       []cli.Flag
	rpcFlags         []cli.Flag
	metricsFlags     []cli.Flag
)

func initFlags() {
	// Flags for testing purpose.
	testFlags = []cli.Flag{
		FakeNetFlag,
		FakeNetGasPowerFlag,
		JsonGenesisFlag,
	}

	// Flags that configure the node.
	gpoFlags = []cli.Flag{}
	accountFlags = []cli.Flag{
		utils.UnlockedAccountFlag,
		utils.PasswordFileFlag,
		utils.ExternalSignerFlag,
		utils.InsecureUnlockAllowedFlag,
	}
	performanceFlags = []cli.Flag{
		CacheFlag,
	}
	networkingFlags = []cli.Flag{
		utils.BootnodesFlag,
		utils.ListenPortFlag,
		utils.MaxPeersFlag,
		utils.MaxPendingPeersFlag,
		utils.NATFlag,
		utils.NoDiscoverFlag,
		utils.DiscoveryV5Flag,
		utils.NetrestrictFlag,
		utils.IPrestrictFlag,
		utils.PrivateNodeFlag,
		utils.NodeKeyFileFlag,
		utils.NodeKeyHexFlag,
	}
	txpoolFlags = []cli.Flag{
		utils.TxPoolLocalsFlag,
		utils.TxPoolNoLocalsFlag,
		utils.TxPoolJournalFlag,
		utils.TxPoolRejournalFlag,
		utils.TxPoolPriceLimitFlag,
		utils.TxPoolPriceBumpFlag,
		utils.TxPoolAccountSlotsFlag,
		utils.TxPoolGlobalSlotsFlag,
		utils.TxPoolAccountQueueFlag,
		utils.TxPoolGlobalQueueFlag,
		utils.TxPoolLifetimeFlag,
	}
	operaFlags = []cli.Flag{
		GenesisFlag,
		ExperimentalGenesisFlag,
		utils.IdentityFlag,
		DataDirFlag,
		utils.MinFreeDiskSpaceFlag,
		utils.KeyStoreDirFlag,
		utils.USBFlag,
		utils.SmartCardDaemonPathFlag,
		ExitWhenAgeFlag,
		ExitWhenEpochFlag,
		utils.LightKDFFlag,
		configFileFlag,
		validatorIDFlag,
		validatorPubkeyFlag,
		validatorPasswordFlag,
		ModeFlag,
		overrideMinGasPriceFlag,
	}

	rpcFlags = []cli.Flag{
		utils.HTTPEnabledFlag,
		utils.HTTPListenAddrFlag,
		utils.HTTPPortFlag,
		utils.HTTPCORSDomainFlag,
		utils.HTTPVirtualHostsFlag,
		utils.GraphQLEnabledFlag,
		utils.GraphQLCORSDomainFlag,
		utils.GraphQLVirtualHostsFlag,
		utils.HTTPApiFlag,
		utils.HTTPPathPrefixFlag,
		utils.WSEnabledFlag,
		utils.WSListenAddrFlag,
		utils.WSPortFlag,
		utils.WSApiFlag,
		utils.WSAllowedOriginsFlag,
		utils.WSPathPrefixFlag,
		utils.IPCDisabledFlag,
		utils.IPCPathFlag,
		RPCGlobalGasCapFlag,
		RPCGlobalEVMTimeoutFlag,
		RPCGlobalTxFeeCapFlag,
		RPCGlobalTimeoutFlag,
	}

	metricsFlags = []cli.Flag{
		utils.MetricsEnabledFlag,
		utils.MetricsEnabledExpensiveFlag,
		utils.MetricsHTTPFlag,
		utils.MetricsPortFlag,
		utils.MetricsEnableInfluxDBFlag,
		utils.MetricsInfluxDBEndpointFlag,
		utils.MetricsInfluxDBDatabaseFlag,
		utils.MetricsInfluxDBUsernameFlag,
		utils.MetricsInfluxDBPasswordFlag,
		utils.MetricsInfluxDBTagsFlag,
		utils.MetricsEnableInfluxDBV2Flag,
		utils.MetricsInfluxDBTokenFlag,
		utils.MetricsInfluxDBBucketFlag,
		utils.MetricsInfluxDBOrganizationFlag,
		tracing.EnableFlag,
	}

	nodeFlags = []cli.Flag{}
	nodeFlags = append(nodeFlags, gpoFlags...)
	nodeFlags = append(nodeFlags, accountFlags...)
	nodeFlags = append(nodeFlags, performanceFlags...)
	nodeFlags = append(nodeFlags, networkingFlags...)
	nodeFlags = append(nodeFlags, txpoolFlags...)
	nodeFlags = append(nodeFlags, operaFlags...)
}

// init the CLI app.
func init() {
	discfilter.Enable()
	overrideFlags()
	overrideParams()

	initFlags()

	// App.

	app.Action = lachesisMain
	app.Version = params.VersionWithCommit(gitCommit, gitDate)
	app.HideVersion = true // we have a command to print the version
	app.Commands = []cli.Command{
		// See accountcmd.go:
		accountCommand,
		walletCommand,
		// see validatorcmd.go:
		validatorCommand,
		// See consolecmd.go:
		consoleCommand,
		attachCommand,
		javascriptCommand,
		// See config.go:
		dumpConfigCommand,
		checkConfigCommand,
		// See misccmd.go:
		versionCommand,
		licenseCommand,
		// See chaincmd.go
		importCommand,
		exportCommand,
		checkCommand,
		// See dbcmd.go
		dbCommand,
	}
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Flags = append(app.Flags, testFlags...)
	app.Flags = append(app.Flags, nodeFlags...)
	app.Flags = append(app.Flags, rpcFlags...)
	app.Flags = append(app.Flags, consoleFlags...)
	app.Flags = append(app.Flags, debug.Flags...)
	app.Flags = append(app.Flags, metricsFlags...)

	app.Before = func(ctx *cli.Context) error {
		if err := debug.Setup(ctx); err != nil {
			return err
		}

		// Start metrics export if enabled
		utils.SetupMetrics(ctx)
		// Start system runtime metrics collection
		go evmetrics.CollectProcessMetrics(3 * time.Second)
		return nil
	}

	app.After = func(ctx *cli.Context) error {
		debug.Exit()
		prompt.Stdin.Close() // Resets terminal mode.

		return nil
	}
}

func Launch(args []string) error {
	return app.Run(args)
}

// opera is the main entry point into the system if no special subcommand is ran.
// It creates a default node based on the command line arguments and runs it in
// blocking mode, waiting for it to be shut down.
func lachesisMain(ctx *cli.Context) error {
	if args := ctx.Args(); len(args) > 0 {
		return fmt.Errorf("invalid command: %q", args[0])
	}

	// TODO: tracing flags
	//tracingStop, err := tracing.Start(ctx)
	//if err != nil {
	//	return err
	//}
	//defer tracingStop()

	cfg := makeAllConfigs(ctx)
	genesisStore := mayGetGenesisStore(ctx, cfg)

	node, _, nodeClose, err := makeNode(ctx, cfg, genesisStore)
	if err != nil {
		return fmt.Errorf("failed to initialize the node: %w", err)
	}
	defer nodeClose()

	if err := startNode(ctx, node); err != nil {
		return fmt.Errorf("failed to start the node: %w", err)
	}
	node.Wait()
	return nil
}

func makeNode(ctx *cli.Context, cfg *config, genesisStore *genesisstore.Store) (*node.Node, *gossip.Service, func(), error) {
	var success bool
	var cleanup []func()
	defer func() { // if the function fails, clean-up in defer, otherwise return cleaning function
		if !success {
			for i := len(cleanup) - 1; i >= 0 ; i-- {
				cleanup[i]()
			}
		}
	}()

	// check errlock file
	errlock.SetDefaultDatadir(cfg.Node.DataDir)
	errlock.Check()

	var g *genesis.Genesis
	if genesisStore != nil {
		gv := genesisStore.Genesis()
		g = &gv
	}

	// applies genesis
	engine, dagIndex, gdb, cdb, blockProc, closeDBs := integration.MakeEngine(path.Join(cfg.Node.DataDir, "chaindata"), g, cfg.AppConfigs())
	cleanup = append(cleanup, func() {
		if err := gdb.Close(); err != nil {
			log.Warn("Failed to close gossip store", "err", err)
		}
		if err := cdb.Close(); err != nil {
			log.Warn("Failed to close consensus database", "err", err)
		}
		if closeDBs != nil {
			if err := closeDBs(); err != nil {
				log.Warn("Failed to close databases", "err", err)
			}
		}
	})

	if genesisStore != nil {
		_ = genesisStore.Close()
	}
	metrics.SetDataDir(cfg.Node.DataDir)

	// substitute default bootnodes if requested
	networkName := ""
	if gdb.HasBlockEpochState() {
		networkName = gdb.GetRules().Name
	}
	if len(networkName) == 0 && genesisStore != nil {
		networkName = genesisStore.Header().NetworkName
	}
	if needDefaultBootnodes(cfg.Node.P2P.BootstrapNodes) {
		bootnodes := Bootnodes[networkName]
		if bootnodes == nil {
			bootnodes = []string{}
		}
		setBootnodes(ctx, bootnodes, &cfg.Node)
	}

	stack, err := makeNetworkStack(ctx, &cfg.Node)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to unlock validator key: %w", err)
	}
	cleanup = append(cleanup, func() {
		if err := stack.Close(); err != nil && err != node.ErrNodeStopped {
			log.Warn("Failed to close stack", "err", err)
		}
	})

	valKeystore := valkeystore.NewDefaultFileKeystore(path.Join(getValKeystoreDir(cfg.Node), "validator"))
	valPubkey := cfg.Emitter.Validator.PubKey
	if key := getFakeValidatorKey(ctx); key != nil && cfg.Emitter.Validator.ID != 0 {
		addFakeValidatorKey(ctx, key, valPubkey, valKeystore)
		coinbase := integration.SetAccountKey(stack.AccountManager(), key, "fakepassword")
		log.Info("Unlocked fake validator account", "address", coinbase.Address.Hex())
	}

	// unlock validator key
	if !valPubkey.Empty() {
		err := unlockValidatorKey(ctx, valPubkey, valKeystore)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to unlock validator key: %w", err)
		}
	}
	signer := valkeystore.NewSigner(valKeystore)

	// Create and register a gossip network service.
	newTxPool := func(reader evmcore.StateReader) gossip.TxPool {
		if cfg.TxPool.Journal != "" {
			cfg.TxPool.Journal = stack.ResolvePath(cfg.TxPool.Journal)
		}
		return evmcore.NewTxPool(cfg.TxPool, reader.Config(), reader)
	}
	haltCheck := func(oldEpoch, newEpoch idx.Epoch, age time.Time) bool {
		stop := ctx.GlobalIsSet(ExitWhenAgeFlag.Name) && ctx.GlobalDuration(ExitWhenAgeFlag.Name) >= time.Since(age)
		stop = stop || ctx.GlobalIsSet(ExitWhenEpochFlag.Name) && idx.Epoch(ctx.GlobalUint64(ExitWhenEpochFlag.Name)) <= newEpoch
		if stop {
			go func() {
				// do it in a separate thread to avoid deadlock
				_ = stack.Close()
			}()
			return true
		}
		return false
	}
	svc, err := gossip.NewService(stack, cfg.Opera, gdb, blockProc, engine, dagIndex, newTxPool, haltCheck)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create the gossip service: %w", err)
	}
	err = engine.Bootstrap(svc.GetConsensusCallbacks())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to bootstrap the consensus engine: %w", err)
	}
	err = engine.Reset(gdb.GetEpoch(), gdb.GetValidators())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to reset the consensus engine: %w", err)
	}
	svc.ReprocessEpochEvents()

	// commit dbs to avoid dirty state when the rest of the startup fails
	if err := gdb.Commit(); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to commit gossip store: %w", err)
	}

	if cfg.Emitter.Validator.ID != 0 {
		svc.RegisterEmitter(emitter.NewEmitter(cfg.Emitter, svc.EmitterWorld(signer)))
	}

	stack.RegisterAPIs(svc.APIs())
	stack.RegisterProtocols(svc.Protocols())
	stack.RegisterLifecycle(svc)

	success = true // skip cleanup in defer - keep it for the returned cleanup function
	return stack, svc, func() {
		for i := len(cleanup) - 1; i >= 0 ; i-- {
			cleanup[i]()
		}
	}, nil
}

func makeNetworkStack(ctx *cli.Context, cfg *node.Config) (*node.Node, error) {
	stack, err := node.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create the protocol stack: %w", err)
	}
	return stack, nil
}

// startNode boots up the system node and all registered protocols, after which
// it unlocks any requested accounts, and starts the RPC/IPC interfaces.
func startNode(ctx *cli.Context, stack *node.Node) error {
	debug.Memsize.Add("node", stack)

	// Start up the node itself
	if err := stack.Start(); err != nil {
		return fmt.Errorf("error starting protocol stack: %w", err)
	}
	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigc)

		startFreeDiskSpaceMonitor(ctx, sigc, stack.InstanceDir())

		<-sigc
		log.Info("Got interrupt, shutting down...")
		go stack.Close()
		for i := 10; i > 0; i-- {
			<-sigc
			if i > 1 {
				log.Warn("Already shutting down, interrupt more to panic.", "times", i-1)
			}
		}
		// received 10 interrupts - kill the node forcefully
		debug.Exit() // ensure trace and CPU profile data is flushed.
		debug.LoudPanic("boom")
	}()

	// Unlock any account specifically requested
	err := unlockAccounts(ctx, stack)
	if err != nil {
		_, _ = fmt.Printf("Fatal: %v\n", err) // for compatibility with tests expecting exact match
		return fmt.Errorf("failed to unlock accounts: %w", err)
	}

	// Register wallet event handlers to open and auto-derive wallets
	events := make(chan accounts.WalletEvent, 16)
	stack.AccountManager().Subscribe(events)

	// Create a client to interact with local opera node.
	rpcClient, err := stack.Attach()
	if err != nil {
		return fmt.Errorf("failed to attach to self: %w", err)
	}
	ethClient := ethclient.NewClient(rpcClient)
	go func() {
		// Open any wallets already attached
		for _, wallet := range stack.AccountManager().Wallets() {
			if err := wallet.Open(""); err != nil {
				log.Warn("Failed to open wallet", "url", wallet.URL(), "err", err)
			}
		}
		// Listen for wallet event till termination
		for event := range events {
			switch event.Kind {
			case accounts.WalletArrived:
				if err := event.Wallet.Open(""); err != nil {
					log.Warn("New wallet appeared, failed to open", "url", event.Wallet.URL(), "err", err)
				}
			case accounts.WalletOpened:
				status, _ := event.Wallet.Status()
				log.Info("New wallet appeared", "url", event.Wallet.URL(), "status", status)

				var derivationPaths []accounts.DerivationPath
				if event.Wallet.URL().Scheme == "ledger" {
					derivationPaths = append(derivationPaths, accounts.LegacyLedgerBaseDerivationPath)
				}
				derivationPaths = append(derivationPaths, accounts.DefaultBaseDerivationPath)

				event.Wallet.SelfDerive(derivationPaths, ethClient)

			case accounts.WalletDropped:
				log.Info("Old wallet dropped", "url", event.Wallet.URL())
				event.Wallet.Close()
			}
		}
	}()

	return nil
}

func startFreeDiskSpaceMonitor(ctx *cli.Context, sigc chan os.Signal, path string) {
	minFreeDiskSpace := ethconfig.Defaults.TrieDirtyCache
	if ctx.GlobalIsSet(utils.MinFreeDiskSpaceFlag.Name) {
		minFreeDiskSpace = ctx.GlobalInt(utils.MinFreeDiskSpaceFlag.Name)
	} else {
		minFreeDiskSpace = 8192
	}
	if minFreeDiskSpace > 0 {
		go diskusage.MonitorFreeDiskSpace(sigc, path, uint64(minFreeDiskSpace)*1024*1024)
	}
}

// unlockAccounts unlocks any account specifically requested.
func unlockAccounts(ctx *cli.Context, stack *node.Node) error {
	var unlocks []string
	inputs := strings.Split(ctx.GlobalString(utils.UnlockedAccountFlag.Name), ",")
	for _, input := range inputs {
		if trimmed := strings.TrimSpace(input); trimmed != "" {
			unlocks = append(unlocks, trimmed)
		}
	}
	// Short circuit if there is no account to unlock.
	if len(unlocks) == 0 {
		return nil
	}
	// If insecure account unlocking is not allowed if node's APIs are exposed to external.
	// Print warning log to user and skip unlocking.
	if !stack.Config().InsecureUnlockAllowed && stack.Config().ExtRPCEnabled() {
		return fmt.Errorf("account unlock with HTTP access is forbidden")
	}
	ks := stack.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
	passwords := utils.MakePasswordList(ctx)
	for i, account := range unlocks {
		if _, _, err := unlockAccount(ks, account, i, passwords); err != nil {
			return err
		}
	}
	return nil
}
