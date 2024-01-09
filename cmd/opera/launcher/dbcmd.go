package launcher

import (
	"fmt"
	"path"

	"github.com/Fantom-foundation/lachesis-base/kvdb"
	"github.com/Fantom-foundation/lachesis-base/kvdb/cachedproducer"
	"github.com/Fantom-foundation/lachesis-base/kvdb/multidb"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"gopkg.in/urfave/cli.v1"

	"github.com/Fantom-foundation/go-opera/gossip"
	"github.com/Fantom-foundation/go-opera/integration"
	"github.com/Fantom-foundation/go-opera/utils/dbutil/compactdb"
)

var (
	experimentalFlag = cli.BoolFlag{
		Name:  "experimental",
		Usage: "Allow experimental DB fixing",
	}
	targetEpochFlag = cli.Uint64Flag{
		Name:  "target.epoch",
		Usage: "The target epoch number to revert the db state to",
	}
	dbCommand = cli.Command{
		Name:        "db",
		Usage:       "A set of commands related to leveldb database",
		Category:    "DB COMMANDS",
		Description: "",
		Subcommands: []cli.Command{
			{
				Name:      "compact",
				Usage:     "Compact all databases",
				ArgsUsage: "",
				Action:    utils.MigrateFlags(compact),
				Category:  "DB COMMANDS",
				Flags: []cli.Flag{
					utils.DataDirFlag,
				},
				Description: `
opera db compact
will compact all databases under datadir's chaindata.
`,
			},
			{
				Name:      "transform",
				Usage:     "Transform DBs layout",
				ArgsUsage: "",
				Action:    utils.MigrateFlags(dbTransform),
				Category:  "DB COMMANDS",
				Flags: []cli.Flag{
					utils.DataDirFlag,
				},
				Description: `
opera db transform
will migrate tables layout according to the configuration.
`,
			},
			{
				Name:      "revert",
				Usage:     "Experimental - revert Opera database to given epoch",
				ArgsUsage: "",
				Action:    utils.MigrateFlags(revertDb),
				Category:  "DB COMMANDS",
				Flags: []cli.Flag{
					utils.DataDirFlag,
					experimentalFlag,
					targetEpochFlag,
				},
				Description: `
opera db revert --experimental --target.epoch 123
Experimental - revert the DB state to given epoch.
If Carmen is used, its database must be replaced with appropriate older version manually.
`,
			},
		},
	}
)

func makeUncheckedDBsProducers(cfg *config) map[multidb.TypeName]kvdb.IterableDBProducer {
	dbsList, _ := integration.SupportedDBs(path.Join(cfg.Node.DataDir, "chaindata"), cfg.DBs.RuntimeCache)
	return dbsList
}

func makeUncheckedCachedDBsProducers(chaindataDir string) map[multidb.TypeName]kvdb.FullDBProducer {
	dbTypes, _ := integration.SupportedDBs(chaindataDir, integration.DBsCacheConfig{
		Table: map[string]integration.DBCacheConfig{
			"": {
				Cache:   1024 * opt.MiB,
				Fdlimit: uint64(utils.MakeDatabaseHandles() / 2),
			},
		},
	})
	wrappedDbTypes := make(map[multidb.TypeName]kvdb.FullDBProducer)
	for typ, producer := range dbTypes {
		wrappedDbTypes[typ] = cachedproducer.WrapAll(&integration.DummyScopedProducer{IterableDBProducer: producer})
	}
	return wrappedDbTypes
}

func makeCheckedDBsProducers(cfg *config) map[multidb.TypeName]kvdb.IterableDBProducer {
	if err := integration.CheckStateInitialized(path.Join(cfg.Node.DataDir, "chaindata"), cfg.DBs); err != nil {
		utils.Fatalf(err.Error())
	}
	return makeUncheckedDBsProducers(cfg)
}

func makeDirectDBsProducerFrom(dbsList map[multidb.TypeName]kvdb.IterableDBProducer, cfg *config) kvdb.FullDBProducer {
	multiRawDbs, err := integration.MakeDirectMultiProducer(dbsList, cfg.DBs.Routing)
	if err != nil {
		utils.Fatalf("Failed to initialize multi DB producer: %v", err)
	}
	return multiRawDbs
}

func makeDirectDBsProducer(cfg *config) kvdb.FullDBProducer {
	dbsList := makeCheckedDBsProducers(cfg)
	return makeDirectDBsProducerFrom(dbsList, cfg)
}

func makeGossipStore(producer kvdb.FlushableDBProducer, cfg *config) *gossip.Store {
	return gossip.NewStore(producer, cfg.OperaStore)
}

func compact(ctx *cli.Context) error {

	cfg := makeAllConfigs(ctx)

	producers := makeCheckedDBsProducers(cfg)
	for typ, p := range producers {
		for _, name := range p.Names() {
			if err := compactDB(typ, name, p); err != nil {
				return err
			}
		}
	}

	return nil
}

func compactDB(typ multidb.TypeName, name string, producer kvdb.DBProducer) error {
	humanName := path.Join(string(typ), name)
	db, err := producer.OpenDB(name)
	defer db.Close()
	if err != nil {
		log.Error("Cannot open db or db does not exists", "db", humanName)
		return err
	}

	log.Info("Stats before compaction", "db", humanName)
	showDbStats(db)

	err = compactdb.Compact(db, humanName, 64*opt.GiB)
	if err != nil {
		log.Error("Database compaction failed", "err", err)
		return err
	}

	log.Info("Stats after compaction", "db", humanName)
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
