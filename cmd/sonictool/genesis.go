package main

import (
	"fmt"
	"github.com/Fantom-foundation/go-opera/cmd/sonictool/genesis"
	"github.com/Fantom-foundation/go-opera/cmd/sonictool/utils"
	"github.com/Fantom-foundation/go-opera/integration/makefakegenesis"
	"github.com/Fantom-foundation/go-opera/opera/genesisstore"
	futils "github.com/Fantom-foundation/go-opera/utils"
	"github.com/Fantom-foundation/go-opera/utils/memory"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/Fantom-foundation/lachesis-base/utils/cachescale"
	"github.com/ethereum/go-ethereum/log"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"gopkg.in/urfave/cli.v1"
	"os"
	"strconv"
	"strings"
)

var (
	GenesisFlag = cli.StringFlag{
		Name:  "genesis",
		Usage: "The genesis file path",
	}
	ModeFlag = cli.StringFlag{
		Name:  "mode",
		Usage: `Mode of the node ("rpc" or "validator")`,
		Value: "rpc",
	}
	FakeNetFlag = cli.StringFlag{
		Name:  "fakenet",
		Usage: "'n/N' - sets coinbase as fake n-th key from genesis of N validators.",
	}
	ExperimentalFlag = cli.BoolFlag{
		Name:  "experimental",
		Usage: "Allow experimental features",
	}
	CacheFlag = cli.IntFlag{
		Name:  "cache",
		Usage: "Megabytes of memory allocated to internal pebble caching",
		Value: utils.DefaultCacheSize,
	}
)

func sonicGenesisImport(ctx *cli.Context) error {
	dataDir := ctx.String(DataDirFlag.Name)
	if dataDir == "" {
		return fmt.Errorf("--%s need to be set", DataDirFlag.Name)
	}
	genesisPath := ctx.String(GenesisFlag.Name)
	if genesisPath == "" {
		return fmt.Errorf("--%s need to be set", GenesisFlag.Name)
	}
	genesisReader, err := os.Open(genesisPath)
	if err != nil {
		return err
	}
	defer genesisReader.Close()
	return genesis.SonicImport(dataDir, genesisReader)
}

func legacyGenesisImport(ctx *cli.Context) error {
	dataDir := ctx.String(DataDirFlag.Name)
	if dataDir == "" {
		return fmt.Errorf("--%s need to be set", DataDirFlag.Name)
	}
	genesisPath := ctx.String(GenesisFlag.Name)
	if genesisPath == "" {
		return fmt.Errorf("--%s need to be set", GenesisFlag.Name)
	}
	validatorMode, err := isValidatorModeSet(ctx)
	if err != nil {
		return err
	}
	cacheRatio, err := cacheScaler(ctx)
	if err != nil {
		return err
	}

	genesisReader, err := os.Open(genesisPath)
	if err != nil {
		return fmt.Errorf("failed to open the legacy genesis file: %w", err)
	}
	defer genesisReader.Close()

	genesisStore, genesisHashes, err := genesisstore.OpenGenesisStore(genesisReader)
	if err != nil {
		return fmt.Errorf("failed to read genesis file: %w", err)
	}
	defer genesisStore.Close()
	if err := genesis.IsGenesisTrusted(genesisStore, genesisHashes); err != nil {
		if ctx.IsSet(ExperimentalFlag.Name) {
			log.Warn("Experimental genesis file is used", "err", err)
		} else {
			return fmt.Errorf("genesis file check failed: %w", err)
		}
	}
	return genesis.ImportGenesisStore(genesisStore, dataDir, validatorMode, cacheRatio)
}

func jsonGenesisImport(ctx *cli.Context) error {
	if !ctx.IsSet(ExperimentalFlag.Name) {
		return fmt.Errorf("using JSON genesis is for experimental usage only and requires --%s flag", ExperimentalFlag.Name)
	}
	dataDir := ctx.String(DataDirFlag.Name)
	if dataDir == "" {
		return fmt.Errorf("--%s need to be set", DataDirFlag.Name)
	}
	genesisPath := ctx.String(GenesisFlag.Name)
	if genesisPath == "" {
		return fmt.Errorf("--%s need to be set", GenesisFlag.Name)
	}
	validatorMode, err := isValidatorModeSet(ctx)
	if err != nil {
		return err
	}
	cacheRatio, err := cacheScaler(ctx)
	if err != nil {
		return err
	}

	genesisJson, err := makefakegenesis.LoadGenesisJson(genesisPath)
	if err != nil {
		return fmt.Errorf("failed to load JSON genesis: %w", err)
	}
	genesisStore, err := makefakegenesis.ApplyGenesisJson(genesisJson)
	if err != nil {
		return fmt.Errorf("failed to prepare JSON genesis: %w", err)
	}
	defer genesisStore.Close()
	return genesis.ImportGenesisStore(genesisStore, dataDir, validatorMode, cacheRatio)
}

func fakeGenesisImport(ctx *cli.Context) error {
	dataDir := ctx.String(DataDirFlag.Name)
	if dataDir == "" {
		return fmt.Errorf("--%s need to be set", DataDirFlag.Name)
	}
	_, num, err := parseFakeGen(ctx.String(FakeNetFlag.Name))
	if err != nil {
		return fmt.Errorf("--%s invalid: %w", FakeNetFlag.Name, err)
	}
	validatorMode, err := isValidatorModeSet(ctx)
	if err != nil {
		return err
	}
	cacheRatio, err := cacheScaler(ctx)
	if err != nil {
		return err
	}

	genesisStore := makefakegenesis.FakeGenesisStore(num, futils.ToFtm(1000000000), futils.ToFtm(5000000))
	defer genesisStore.Close()
	return genesis.ImportGenesisStore(genesisStore, dataDir, validatorMode, cacheRatio)
}

func isValidatorModeSet(ctx *cli.Context) (bool, error) {
	if ctx.IsSet(ModeFlag.Name) {
		mode := ctx.String(ModeFlag.Name)
		if mode != "rpc" && mode != "validator" {
			return false, fmt.Errorf("--%s must be 'rpc' or 'validator'", ModeFlag.Name)
		}
		if mode == "validator" {
			return true, nil
		}
	}
	return false, nil
}

func parseFakeGen(s string) (id idx.ValidatorID, num idx.Validator, err error) {
	parts := strings.SplitN(s, "/", 2)
	if len(parts) != 2 {
		err = fmt.Errorf("use %%d/%%d format")
		return
	}

	var u32 uint64
	u32, err = strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		return
	}
	id = idx.ValidatorID(u32)

	u32, err = strconv.ParseUint(parts[1], 10, 32)
	num = idx.Validator(u32)
	if num < 0 || idx.Validator(id) > num {
		err = fmt.Errorf("key-num should be in range from 1 to validators (<key-num>/<validators>), or should be zero for non-validator node")
		return
	}

	return
}

func cacheScaler(ctx *cli.Context) (cachescale.Func, error) {
	targetCache := ctx.Int(CacheFlag.Name)
	baseSize := utils.DefaultCacheSize
	totalMemory := int(memory.TotalMemory() / opt.MiB)
	maxCache := totalMemory * 3 / 5
	if maxCache < baseSize {
		maxCache = baseSize
	}
	if !ctx.IsSet(CacheFlag.Name) {
		recommendedCache := totalMemory / 2
		if recommendedCache > baseSize {
			log.Warn(fmt.Sprintf("Please add '--%s %d' flag to allocate more cache for the database. Total memory is %d MB.", CacheFlag.Name, recommendedCache, totalMemory))
		}
		return cachescale.Identity, nil
	}
	if targetCache < baseSize {
		return nil, fmt.Errorf("invalid flag %s - minimum cache size is %d MB", CacheFlag.Name, baseSize)
	}
	if totalMemory != 0 && targetCache > maxCache {
		log.Warn(fmt.Sprintf("Requested cache size exceeds 60%% of available memory. Reducing cache size to %d MB.", maxCache))
		targetCache = maxCache
	}

	return cachescale.Ratio{
		Base:   uint64(baseSize - utils.ConstantCacheSize),
		Target: uint64(targetCache - utils.ConstantCacheSize),
	}, nil
}
