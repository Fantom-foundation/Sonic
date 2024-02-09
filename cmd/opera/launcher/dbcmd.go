package launcher

import (
	"fmt"
	"path"

	"github.com/Fantom-foundation/lachesis-base/kvdb"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"gopkg.in/urfave/cli.v1"

	"github.com/Fantom-foundation/go-opera/integration"
	"github.com/Fantom-foundation/go-opera/utils/dbutil/compactdb"
)

var (
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
		},
	}
)

func makeDBsProducer(cfg *config) kvdb.FullDBProducer {
	if err := integration.CheckStateInitialized(path.Join(cfg.Node.DataDir, "chaindata"), cfg.DBs); err != nil {
		utils.Fatalf(err.Error())
	}
	producer, err := integration.GetDbProducer(path.Join(cfg.Node.DataDir, "chaindata"), cfg.DBs.RuntimeCache)
	if err != nil {
		utils.Fatalf("Failed to initialize DB producer: %v", err)
	}
	return producer
}

func compact(ctx *cli.Context) error {

	cfg := makeAllConfigs(ctx)

	producer := makeDBsProducer(cfg)
	for _, name := range producer.Names() {
		if err := compactDB(name, producer); err != nil {
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
