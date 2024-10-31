package check

import (
	"fmt"
	"os"
	"path/filepath"

	carmen "github.com/Fantom-foundation/Carmen/go/state"
	"github.com/Fantom-foundation/go-opera/gossip"
	"github.com/Fantom-foundation/go-opera/integration"
	"github.com/Fantom-foundation/lachesis-base/kvdb"
	"github.com/Fantom-foundation/lachesis-base/utils/cachescale"
	"github.com/ethereum/go-ethereum/log"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

func createGdb(dataDir string, cacheRatio cachescale.Func, archive carmen.ArchiveType, skipArchiveCheck bool) (*gossip.Store, kvdb.FullDBProducer, error) {
	chaindataDir := filepath.Join(dataDir, "chaindata")
	carmenDir := filepath.Join(dataDir, "carmen")

	if stat, err := os.Stat(chaindataDir); err != nil || !stat.IsDir() {
		return nil, nil, fmt.Errorf("unable to validate: datadir does not contain chandata")
	}
	if stat, err := os.Stat(carmenDir); err != nil || !stat.IsDir() {
		return nil, nil, fmt.Errorf("unable to validate: datadir does not contain carmen")
	}

	dbs, err := integration.GetDbProducer(chaindataDir, integration.DBCacheConfig{
		Cache:   cacheRatio.U64(480 * opt.MiB),
		Fdlimit: 100,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to make DB producer: %v", err)
	}

	gdbConfig := gossip.DefaultStoreConfig(cacheRatio)
	gdbConfig.EVM.StateDb.Directory = carmenDir
	gdbConfig.EVM.StateDb.Archive = archive
	gdbConfig.EVM.SkipArchiveCheck = skipArchiveCheck // skip archive mode check (allow "check live" to run with archive enabled)

	gdb, err := gossip.NewStore(dbs, gdbConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create gossip store: %w", err)
	}

	err = gdb.EvmStore().Open()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open EvmStore: %v", err)
	}

	return gdb, dbs, nil
}

type verificationObserver struct{}

func (o verificationObserver) StartVerification() {}

func (o verificationObserver) Progress(msg string) {
	log.Info(msg)
}

func (o verificationObserver) EndVerification(res error) {}
