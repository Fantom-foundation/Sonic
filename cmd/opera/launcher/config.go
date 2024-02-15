package launcher

import (
	"bufio"
	"errors"
	"fmt"
	carmen "github.com/Fantom-foundation/Carmen/go/state"
	"github.com/Fantom-foundation/go-opera/cmd/opera/launcher/utils"
	"github.com/Fantom-foundation/go-opera/gossip/evmstore"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/Fantom-foundation/lachesis-base/abft"
	"github.com/Fantom-foundation/lachesis-base/utils/cachescale"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/params"
	"github.com/naoina/toml"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"gopkg.in/urfave/cli.v1"

	"github.com/Fantom-foundation/go-opera/evmcore"
	"github.com/Fantom-foundation/go-opera/gossip"
	"github.com/Fantom-foundation/go-opera/gossip/emitter"
	"github.com/Fantom-foundation/go-opera/integration"
	"github.com/Fantom-foundation/go-opera/opera/genesis"
	"github.com/Fantom-foundation/go-opera/utils/memory"
	"github.com/Fantom-foundation/go-opera/vecmt"
)

var (
	configFileFlag = cli.StringFlag{
		Name:  "config",
		Usage: "TOML configuration file",
	}
	CacheFlag = cli.IntFlag{
		Name:  "cache",
		Usage: "Megabytes of memory allocated to internal caching",
		Value: DefaultCacheSize,
	}
	RPCGlobalGasCapFlag = cli.Uint64Flag{
		Name:  "rpc.gascap",
		Usage: "Sets a cap on gas that can be used in ftm_call/estimateGas (0=infinite)",
		Value: gossip.DefaultConfig(cachescale.Identity).RPCGasCap,
	}
	RPCGlobalEVMTimeoutFlag = &cli.DurationFlag{
		Name:  "rpc.evmtimeout",
		Usage: "Sets a timeout used for eth_call (0=infinite)",
		Value: gossip.DefaultConfig(cachescale.Identity).RPCEVMTimeout,
	}
	RPCGlobalTxFeeCapFlag = cli.Float64Flag{
		Name:  "rpc.txfeecap",
		Usage: "Sets a cap on transaction fee (in FTM) that can be sent via the RPC APIs (0 = no cap)",
		Value: gossip.DefaultConfig(cachescale.Identity).RPCTxFeeCap,
	}
	RPCGlobalTimeoutFlag = cli.DurationFlag{
		Name:  "rpc.timeout",
		Usage: "Time limit for RPC calls execution",
		Value: gossip.DefaultConfig(cachescale.Identity).RPCTimeout,
	}
	ModeFlag = cli.StringFlag{
		Name:  "mode",
		Usage: `Mode of the node ("rpc" or "validator")`,
		Value: "rpc",
	}
	ExitWhenAgeFlag = cli.DurationFlag{
		Name:  "exitwhensynced.age",
		Usage: "Exits after synchronisation reaches the required age",
	}
	ExitWhenEpochFlag = cli.Uint64Flag{
		Name:  "exitwhensynced.epoch",
		Usage: "Exits after synchronisation reaches the required epoch",
	}
)

type GenesisTemplate struct {
	Name   string
	Header genesis.Header
	Hashes genesis.Hashes
}

const (
	// DefaultCacheSize is calculated as memory consumption in a worst case scenario with default configuration
	// Average memory consumption might be 3-5 times lower than the maximum
	DefaultCacheSize  = 3600
	ConstantCacheSize = 400
)

// These settings ensure that TOML keys use the same names as Go struct fields.
var TomlSettings = toml.Config{
	NormFieldName: func(rt reflect.Type, key string) string {
		return key
	},
	FieldToKey: func(rt reflect.Type, field string) string {
		return field
	},
	MissingField: func(rt reflect.Type, field string) error {
		return fmt.Errorf("field '%s' is not defined in %s", field, rt.String())
	},
}

type config struct {
	Node          node.Config
	Opera         gossip.Config
	Emitter       emitter.Config
	TxPool        evmcore.TxPoolConfig
	OperaStore    gossip.StoreConfig
	Lachesis      abft.Config
	LachesisStore abft.StoreConfig
	VectorClock   vecmt.IndexConfig
	DBs           integration.DBsConfig
}

func (c *config) AppConfigs() integration.Configs {
	return integration.Configs{
		Opera:         c.Opera,
		OperaStore:    c.OperaStore,
		Lachesis:      c.Lachesis,
		LachesisStore: c.LachesisStore,
		VectorClock:   c.VectorClock,
		DBs:           c.DBs,
	}
}

func loadAllConfigs(file string, cfg *config) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	err = TomlSettings.NewDecoder(bufio.NewReader(f)).Decode(cfg)
	// Add file name to errors that have a line number.
	if _, ok := err.(*toml.LineError); ok {
		err = errors.New(file + ", " + err.Error())
	}
	if err != nil {
		return errors.New(fmt.Sprintf("TOML config file error: %v.\n"+
			"Use 'dumpconfig' command to get an example config file.\n"+
			"If node was recently upgraded and a previous network config file is used, then check updates for the config file.", err))
	}
	return err
}

func setBootnodes(ctx *cli.Context, urls []string, cfg *node.Config) {
	cfg.P2P.BootstrapNodesV5 = []*enode.Node{}
	for _, url := range urls {
		if url != "" {
			node, err := enode.Parse(enode.ValidSchemes, url)
			if err != nil {
				log.Error("Bootstrap URL invalid", "enode", url, "err", err)
				continue
			}
			cfg.P2P.BootstrapNodesV5 = append(cfg.P2P.BootstrapNodesV5, node)
		}
	}
	cfg.P2P.BootstrapNodes = cfg.P2P.BootstrapNodesV5
}

func setTxPool(ctx *cli.Context, cfg *evmcore.TxPoolConfig) error {
	if ctx.GlobalIsSet(utils.TxPoolLocalsFlag.Name) {
		locals := strings.Split(ctx.GlobalString(utils.TxPoolLocalsFlag.Name), ",")
		for _, account := range locals {
			if trimmed := strings.TrimSpace(account); !common.IsHexAddress(trimmed) {
				return fmt.Errorf("invalid account in --%s: %s", utils.TxPoolLocalsFlag.Name, trimmed)
			} else {
				cfg.Locals = append(cfg.Locals, common.HexToAddress(account))
			}
		}
	}
	if ctx.GlobalIsSet(utils.TxPoolNoLocalsFlag.Name) {
		cfg.NoLocals = ctx.GlobalBool(utils.TxPoolNoLocalsFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolJournalFlag.Name) {
		cfg.Journal = ctx.GlobalString(utils.TxPoolJournalFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolRejournalFlag.Name) {
		cfg.Rejournal = ctx.GlobalDuration(utils.TxPoolRejournalFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolPriceLimitFlag.Name) {
		cfg.PriceLimit = ctx.GlobalUint64(utils.TxPoolPriceLimitFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolPriceBumpFlag.Name) {
		cfg.PriceBump = ctx.GlobalUint64(utils.TxPoolPriceBumpFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolAccountSlotsFlag.Name) {
		cfg.AccountSlots = ctx.GlobalUint64(utils.TxPoolAccountSlotsFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolGlobalSlotsFlag.Name) {
		cfg.GlobalSlots = ctx.GlobalUint64(utils.TxPoolGlobalSlotsFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolAccountQueueFlag.Name) {
		cfg.AccountQueue = ctx.GlobalUint64(utils.TxPoolAccountQueueFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolGlobalQueueFlag.Name) {
		cfg.GlobalQueue = ctx.GlobalUint64(utils.TxPoolGlobalQueueFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolLifetimeFlag.Name) {
		cfg.Lifetime = ctx.GlobalDuration(utils.TxPoolLifetimeFlag.Name)
	}
	return nil
}

func gossipConfigWithFlags(ctx *cli.Context, src gossip.Config) gossip.Config {
	cfg := src

	if ctx.GlobalIsSet(RPCGlobalGasCapFlag.Name) {
		cfg.RPCGasCap = ctx.GlobalUint64(RPCGlobalGasCapFlag.Name)
	}
	if ctx.GlobalIsSet(RPCGlobalEVMTimeoutFlag.Name) {
		cfg.RPCEVMTimeout = ctx.GlobalDuration(RPCGlobalEVMTimeoutFlag.Name)
	}
	if ctx.GlobalIsSet(RPCGlobalTxFeeCapFlag.Name) {
		cfg.RPCTxFeeCap = ctx.GlobalFloat64(RPCGlobalTxFeeCapFlag.Name)
	}
	if ctx.GlobalIsSet(RPCGlobalTimeoutFlag.Name) {
		cfg.RPCTimeout = ctx.GlobalDuration(RPCGlobalTimeoutFlag.Name)
	}

	return cfg
}

func setEvmStore(ctx *cli.Context, datadir string, src  evmstore.StoreConfig) (evmstore.StoreConfig, error) {
	cfg := src
	cfg.StateDb.Directory = filepath.Join(datadir, "carmen")

	if ctx.GlobalIsSet(ModeFlag.Name) {
		mode := ctx.GlobalString(ModeFlag.Name)
		if mode != "rpc" && mode != "validator" {
			return cfg, fmt.Errorf("--%s must be 'rpc' or 'validator'", ModeFlag.Name)
		}
		if mode == "validator" {
			cfg.StateDb.Archive = carmen.NoArchive
			cfg.DisableLogsIndexing = true
			cfg.DisableTxHashesIndexing = true
		}
	}
	return cfg, nil
}

func setDBConfig(cfg config, cacheRatio cachescale.Func) (config, error) {
	handles, err := utils.MakeDatabaseHandles()
	if err != nil {
		return config{}, err
	}
	cfg.DBs.RuntimeCache = integration.DBCacheConfig{
		Cache:   cacheRatio.U64(480 * opt.MiB),
		Fdlimit: handles*480/1400 + 1,
	}
	return cfg, nil
}

func cacheScaler(ctx *cli.Context) cachescale.Func {
	targetCache := ctx.GlobalInt(CacheFlag.Name)
	baseSize := DefaultCacheSize
	totalMemory := int(memory.TotalMemory() / opt.MiB)
	maxCache := totalMemory * 3 / 5
	if maxCache < baseSize {
		maxCache = baseSize
	}
	if !ctx.GlobalIsSet(CacheFlag.Name) {
		recommendedCache := totalMemory / 2
		if recommendedCache > baseSize {
			log.Warn(fmt.Sprintf("Please add '--%s %d' flag to allocate more cache for Opera. Total memory is %d MB.", CacheFlag.Name, recommendedCache, totalMemory))
		}
		return cachescale.Identity
	}
	if targetCache < baseSize {
		log.Crit("Invalid flag", "flag", CacheFlag.Name, "err", fmt.Sprintf("minimum cache size is %d MB", baseSize))
	}
	if totalMemory != 0 && targetCache > maxCache {
		log.Warn(fmt.Sprintf("Requested cache size exceeds 60%% of available memory. Reducing cache size to %d MB.", maxCache))
		targetCache = maxCache
	}

	return cachescale.Ratio{
		Base:   uint64(baseSize - ConstantCacheSize),
		Target: uint64(targetCache - ConstantCacheSize),
	}
}

func MakeAllConfigsFromFile(ctx *cli.Context, configFile string) (*config, error) {
	// Defaults (low priority)
	cacheRatio := cacheScaler(ctx)
	cfg := config{
		Node:          defaultNodeConfig(),
		Opera:         gossip.DefaultConfig(cacheRatio),
		Emitter:       emitter.DefaultConfig(),
		TxPool:        evmcore.DefaultTxPoolConfig,
		OperaStore:    gossip.DefaultStoreConfig(cacheRatio),
		Lachesis:      abft.DefaultConfig(),
		LachesisStore: abft.DefaultStoreConfig(cacheRatio),
		VectorClock:   vecmt.DefaultConfig(cacheRatio),
	}

	if ctx.GlobalIsSet(FakeNetFlag.Name) {
		_, num, err := parseFakeGen(ctx.GlobalString(FakeNetFlag.Name))
		if err != nil {
			return nil, fmt.Errorf("invalid fakenet flag")
		}
		cfg.Emitter = emitter.FakeConfig(num)
		setBootnodes(ctx, []string{}, &cfg.Node)
	} else {
		// "asDefault" means set network defaults
		cfg.Node.P2P.BootstrapNodes = asDefault
		cfg.Node.P2P.BootstrapNodesV5 = asDefault
	}

	// Load config file (medium priority)
	if configFile != "" {
		if err := loadAllConfigs(configFile, &cfg); err != nil {
			return &cfg, err
		}
	}

	// Apply flags (high priority)
	var err error
	cfg.Opera = gossipConfigWithFlags(ctx, cfg.Opera)
	err = utils.SetNodeConfig(ctx, &cfg.Node)
	if err != nil {
		return nil, err
	}
	cfg.OperaStore.EVM, err = setEvmStore(ctx, cfg.Node.DataDir, cfg.OperaStore.EVM)
	if err != nil {
		return nil, err
	}

	err = setValidator(ctx, &cfg.Emitter)
	if err != nil {
		return nil, err
	}
	if cfg.Emitter.Validator.ID != 0 && len(cfg.Emitter.PrevEmittedEventFile.Path) == 0 {
		cfg.Emitter.PrevEmittedEventFile.Path = cfg.Node.ResolvePath(path.Join("emitter", fmt.Sprintf("last-%d", cfg.Emitter.Validator.ID)))
	}
	if err := setTxPool(ctx, &cfg.TxPool); err != nil {
		return nil, err
	}

	// Process DBs defaults in the end because they are applied only in absence of config or flags
	cfg, err = setDBConfig(cfg, cacheRatio)
	if err != nil {
		return nil, err
	}

	if err := cfg.Opera.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func MakeAllConfigs(ctx *cli.Context) (*config, error) {
	return MakeAllConfigsFromFile(ctx, ctx.GlobalString(configFileFlag.Name))
}

func defaultNodeConfig() node.Config {
	cfg := NodeDefaultConfig
	cfg.Name = clientIdentifier
	cfg.Version = params.VersionWithCommit(gitCommit, gitDate)
	cfg.HTTPModules = append(cfg.HTTPModules, "eth", "ftm", "dag", "abft", "web3")
	cfg.WSModules = append(cfg.WSModules, "eth", "ftm", "dag", "abft", "web3")
	cfg.IPCPath = "opera.ipc"
	return cfg
}
