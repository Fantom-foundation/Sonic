package main

import (
	"fmt"
	"github.com/Fantom-foundation/go-opera/config/flags"
	"github.com/Fantom-foundation/go-opera/integration"
	"github.com/Fantom-foundation/go-opera/utils/dbutil/compactdb"
	"github.com/Fantom-foundation/lachesis-base/kvdb"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"gopkg.in/urfave/cli.v1"
	"path/filepath"
)

func compactDbs(ctx *cli.Context) error {
	dataDir := ctx.GlobalString(flags.DataDirFlag.Name)
	if dataDir == "" {
		return fmt.Errorf("--%s need to be set", flags.DataDirFlag.Name)
	}
	cacheRatio, err := cacheScaler(ctx)
	if err != nil {
		return err
	}
	chaindataDir := filepath.Join(dataDir, "chaindata")
	dbs, err := integration.GetDbProducer(chaindataDir, integration.DBCacheConfig{
		Cache:   cacheRatio.U64(480 * opt.MiB),
		Fdlimit: 100,
	})
	if err != nil {
		return fmt.Errorf("failed to make DB producer: %v", err)
	}
	defer dbs.Close()

	for _, name := range dbs.Names() {
		if err := compactDB(name, dbs); err != nil {
			return err
		}
	}
	return nil
}

func compactDB(name string, producer kvdb.DBProducer) error {
	db, err := producer.OpenDB(name)
	defer db.Close()
	if err != nil {
		log.Error("Cannot open db or db does not exists", "db", name)
		return err
	}

	log.Info("Stats before compaction", "db", name)
	showDbStats(db)

	err = compactdb.Compact(db, name, 64*opt.GiB)
	if err != nil {
		log.Error("Database compaction failed", "err", err)
		return err
	}

	log.Info("Stats after compaction", "db", name)
	showDbStats(db)

	return nil
}

func showDbStats(db ethdb.Stater) {
	if stats, err := db.Stat("stats"); err != nil {
		log.Warn("Failed to read database stats", "error", err)
	} else {
		fmt.Println(stats)
	}
	if ioStats, err := db.Stat("iostats"); err != nil {
		log.Warn("Failed to read database iostats", "error", err)
	} else {
		fmt.Println(ioStats)
	}
}
