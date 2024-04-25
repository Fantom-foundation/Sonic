package config

import (
	"bufio"
	"errors"
	"fmt"
	carmen "github.com/Fantom-foundation/Carmen/go/state"
	"github.com/Fantom-foundation/go-opera/config/flags"
	"github.com/Fantom-foundation/go-opera/gossip/evmstore"
	"github.com/ethereum/go-ethereum/common/fdlimit"
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
	"github.com/Fantom-foundation/go-opera/utils/memory"
	"github.com/Fantom-foundation/go-opera/vecmt"
)

const (
	// ClientIdentifier to advertise over the network.
	ClientIdentifier = "go-opera"
)

var (
	// Git SHA1 commit hash of the release (set via linker flags).
	GitCommit = ""
	GitDate   = ""
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

type Config struct {
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

func (c *Config) AppConfigs() integration.Configs {
	return integration.Configs{
		Opera:         c.Opera,
		OperaStore:    c.OperaStore,
		Lachesis:      c.Lachesis,
		LachesisStore: c.LachesisStore,
		VectorClock:   c.VectorClock,
		DBs:           c.DBs,
	}
}

func loadAllConfigs(file string, cfg *Config) error {
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
	if ctx.GlobalIsSet(flags.TxPoolLocalsFlag.Name) {
		locals := strings.Split(ctx.GlobalString(flags.TxPoolLocalsFlag.Name), ",")
		for _, account := range locals {
			if trimmed := strings.TrimSpace(account); !common.IsHexAddress(trimmed) {
				return fmt.Errorf("invalid account in --%s: %s", flags.TxPoolLocalsFlag.Name, trimmed)
			} else {
				cfg.Locals = append(cfg.Locals, common.HexToAddress(account))
			}
		}
	}
	if ctx.GlobalIsSet(flags.TxPoolNoLocalsFlag.Name) {
		cfg.NoLocals = ctx.GlobalBool(flags.TxPoolNoLocalsFlag.Name)
	}
	if ctx.GlobalIsSet(flags.TxPoolJournalFlag.Name) {
		cfg.Journal = ctx.GlobalString(flags.TxPoolJournalFlag.Name)
	}
	if ctx.GlobalIsSet(flags.TxPoolRejournalFlag.Name) {
		cfg.Rejournal = ctx.GlobalDuration(flags.TxPoolRejournalFlag.Name)
	}
	if ctx.GlobalIsSet(flags.TxPoolPriceLimitFlag.Name) {
		cfg.PriceLimit = ctx.GlobalUint64(flags.TxPoolPriceLimitFlag.Name)
	}
	if ctx.GlobalIsSet(flags.TxPoolPriceBumpFlag.Name) {
		cfg.PriceBump = ctx.GlobalUint64(flags.TxPoolPriceBumpFlag.Name)
	}
	if ctx.GlobalIsSet(flags.TxPoolAccountSlotsFlag.Name) {
		cfg.AccountSlots = ctx.GlobalUint64(flags.TxPoolAccountSlotsFlag.Name)
	}
	if ctx.GlobalIsSet(flags.TxPoolGlobalSlotsFlag.Name) {
		cfg.GlobalSlots = ctx.GlobalUint64(flags.TxPoolGlobalSlotsFlag.Name)
	}
	if ctx.GlobalIsSet(flags.TxPoolAccountQueueFlag.Name) {
		cfg.AccountQueue = ctx.GlobalUint64(flags.TxPoolAccountQueueFlag.Name)
	}
	if ctx.GlobalIsSet(flags.TxPoolGlobalQueueFlag.Name) {
		cfg.GlobalQueue = ctx.GlobalUint64(flags.TxPoolGlobalQueueFlag.Name)
	}
	if ctx.GlobalIsSet(flags.TxPoolLifetimeFlag.Name) {
		cfg.Lifetime = ctx.GlobalDuration(flags.TxPoolLifetimeFlag.Name)
	}
	return nil
}

func gossipConfigWithFlags(ctx *cli.Context, src gossip.Config) gossip.Config {
	cfg := src

	if ctx.GlobalIsSet(flags.RPCGlobalGasCapFlag.Name) {
		cfg.RPCGasCap = ctx.GlobalUint64(flags.RPCGlobalGasCapFlag.Name)
	}
	if ctx.GlobalIsSet(flags.RPCGlobalEVMTimeoutFlag.Name) {
		cfg.RPCEVMTimeout = ctx.GlobalDuration(flags.RPCGlobalEVMTimeoutFlag.Name)
	}
	if ctx.GlobalIsSet(flags.RPCGlobalTxFeeCapFlag.Name) {
		cfg.RPCTxFeeCap = ctx.GlobalFloat64(flags.RPCGlobalTxFeeCapFlag.Name)
	}
	if ctx.GlobalIsSet(flags.RPCGlobalTimeoutFlag.Name) {
		cfg.RPCTimeout = ctx.GlobalDuration(flags.RPCGlobalTimeoutFlag.Name)
	}

	return cfg
}

func setEvmStore(ctx *cli.Context, datadir string, src  evmstore.StoreConfig) (evmstore.StoreConfig, error) {
	cfg := src
	cfg.StateDb.Directory = filepath.Join(datadir, "carmen")

	if ctx.GlobalIsSet(flags.ModeFlag.Name) || ctx.IsSet(flags.ModeFlag.Name) {
		var mode string
		if ctx.IsSet(flags.ModeFlag.Name) {
			mode = ctx.String(flags.ModeFlag.Name)
		} else {
			mode = ctx.GlobalString(flags.ModeFlag.Name)
		}
		if mode != "rpc" && mode != "validator" {
			return cfg, fmt.Errorf("--%s must be 'rpc' or 'validator'", flags.ModeFlag.Name)
		}
		if mode == "validator" {
			cfg.StateDb.Archive = carmen.NoArchive
			cfg.DisableLogsIndexing = true
			cfg.DisableTxHashesIndexing = true
		}
	}
	return cfg, nil
}

// makeDatabaseHandles raises out the number of allowed file handles per process
// and returns half of the allowance to assign to the database.
func makeDatabaseHandles() (uint64, error) {
	limit, err := fdlimit.Maximum()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve file descriptor allowance: %w", err)
	}
	raised, err := fdlimit.Raise(uint64(limit))
	if err != nil {
		return 0, fmt.Errorf("failed to raise file descriptor allowance: %w", err)
	}
	return raised / 2, nil // Leave half for networking and other stuff
}

func setDBConfig(cfg Config, cacheRatio cachescale.Func) (Config, error) {
	handles, err := makeDatabaseHandles()
	if err != nil {
		return Config{}, err
	}
	cfg.DBs.RuntimeCache = integration.DBCacheConfig{
		Cache:   cacheRatio.U64(480 * opt.MiB),
		Fdlimit: handles*480/1400 + 1,
	}
	return cfg, nil
}

const (
	// DefaultCacheSize is calculated as memory consumption in a worst case scenario with default configuration
	// Average memory consumption might be 3-5 times lower than the maximum
	DefaultCacheSize  = 6 * 1024 // MB
	ConstantCacheSize = 400 // MB
)

func cacheScaler(ctx *cli.Context) cachescale.Func {
	baseSize := DefaultCacheSize
	totalMemory := int(memory.TotalMemory() / opt.MiB)
	maxCache := totalMemory * 3 / 5  // max 60% of available memory
	if maxCache < baseSize {
		maxCache = baseSize
	}
	if !ctx.GlobalIsSet(flags.CacheFlag.Name) {
		recommendedCache := totalMemory / 2
		if recommendedCache > baseSize {
			log.Warn(fmt.Sprintf("Please add '--%s %d' flag to allocate more cache for Opera. Total memory is %d MB.", flags.CacheFlag.Name, recommendedCache, totalMemory))
		}
		return cachescale.Identity
	}
	targetCache := ctx.GlobalInt(flags.CacheFlag.Name)
	if targetCache < baseSize {
		log.Crit("Invalid flag", "flag", flags.CacheFlag.Name, "err", fmt.Sprintf("minimum cache size is %d MB", baseSize))
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

func MakeAllConfigsFromFile(ctx *cli.Context, configFile string) (*Config, error) {
	// Defaults (low priority)
	cacheRatio := cacheScaler(ctx)
	cfg := Config{
		Node:          DefaultNodeConfig(),
		Opera:         gossip.DefaultConfig(cacheRatio),
		Emitter:       emitter.DefaultConfig(),
		TxPool:        evmcore.DefaultTxPoolConfig,
		OperaStore:    gossip.DefaultStoreConfig(cacheRatio),
		Lachesis:      abft.DefaultConfig(),
		LachesisStore: abft.DefaultStoreConfig(cacheRatio),
		VectorClock:   vecmt.DefaultConfig(cacheRatio),
	}

	if ctx.GlobalIsSet(FakeNetFlag.Name) {
		_, num, err := ParseFakeGen(ctx.GlobalString(FakeNetFlag.Name))
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
	err = SetNodeConfig(ctx, &cfg.Node)
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

func MakeAllConfigs(ctx *cli.Context) (*Config, error) {
	return MakeAllConfigsFromFile(ctx, ctx.GlobalString(flags.ConfigFileFlag.Name))
}

func DefaultNodeConfig() node.Config {
	cfg := NodeDefaultConfig
	cfg.Name = ClientIdentifier
	cfg.Version = params.VersionWithCommit(GitCommit, GitDate)
	cfg.HTTPModules = append(cfg.HTTPModules, "eth", "ftm", "dag", "abft", "web3")
	cfg.WSModules = append(cfg.WSModules, "eth", "ftm", "dag", "abft", "web3")
	cfg.IPCPath = "opera.ipc"
	return cfg
}
