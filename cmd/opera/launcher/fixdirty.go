package launcher

import (
	"fmt"
	"time"

	"github.com/Fantom-foundation/lachesis-base/abft"
	"github.com/Fantom-foundation/lachesis-base/common/bigendian"
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/Fantom-foundation/lachesis-base/kvdb"
	"github.com/Fantom-foundation/lachesis-base/kvdb/batched"
	"github.com/Fantom-foundation/lachesis-base/kvdb/flushable"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"gopkg.in/urfave/cli.v1"

	"github.com/Fantom-foundation/go-opera/integration"
	"github.com/Fantom-foundation/go-opera/inter/iblockproc"
)

// revertDb is the 'db revert' command.
func revertDb(ctx *cli.Context) error {
	if !ctx.Bool(experimentalFlag.Name) {
		utils.Fatalf("Add --experimental flag")
	}
	var targetEpoch = idx.Epoch(ctx.Uint64(targetEpochFlag.Name))
	if targetEpoch == 0 {
		utils.Fatalf("Add --target.epoch flag")
	}

	cfg := makeAllConfigs(ctx)

	log.Info("Opening databases")
	dbTypes := makeUncheckedDBsProducers(cfg)
	multiProducer := makeDirectDBsProducerFrom(dbTypes, cfg)

	// reverts the gossip database state
	epochState, oldTopEpoch, err := revertGossipDb(multiProducer, cfg, targetEpoch)
	if err != nil {
		return err
	}

	// drop epoch-related databases and consensus database
	log.Info("Removing epoch DBs - will be recreated on next start")
	for _, name := range []string{
		fmt.Sprintf("gossip-%d", oldTopEpoch),
		fmt.Sprintf("lachesis-%d", oldTopEpoch),
		"lachesis",
	} {
		err = eraseTable(name, multiProducer)
		if err != nil {
			return err
		}
	}

	// prepare consensus database from epochState
	log.Info("Recreating lachesis DB")
	cMainDb := mustOpenDB(multiProducer, "lachesis")
	cGetEpochDB := func(epoch idx.Epoch) kvdb.Store {
		return mustOpenDB(multiProducer, fmt.Sprintf("lachesis-%d", epoch))
	}
	cdb := abft.NewStore(cMainDb, cGetEpochDB, panics("Lachesis store"), cfg.LachesisStore)
	err = cdb.ApplyGenesis(&abft.Genesis{
		Epoch:      epochState.Epoch,
		Validators: epochState.Validators,
	})
	if err != nil {
		return fmt.Errorf("failed to init consensus database: %v", err)
	}
	_ = cdb.Close()

	log.Info("Clearing DBs dirty flags")
	id := bigendian.Uint64ToBytes(uint64(time.Now().UnixNano()))
	for typ, producer := range dbTypes {
		err := clearDirtyFlags(id, producer)
		if err != nil {
			log.Crit("Failed to write clean FlushID", "type", typ, "err", err)
		}
	}

	log.Info("Revert is complete")
	return nil
}

// revertGossipDb reverts the gossip database into state, when was one of last epochs sealed
func revertGossipDb(producer kvdb.FlushableDBProducer, cfg *config, targetEpoch idx.Epoch) (
	epochState *iblockproc.EpochState, oldTopEpoch idx.Epoch, err error) {
	gdb := makeGossipStore(producer, cfg) // requires FlushIDKey present (not clean) in all dbs
	defer gdb.Close()
	oldTopEpoch = gdb.GetEpoch()

	// get target epoch/block state
	blockState, epochState := gdb.GetHistoryBlockEpochState(targetEpoch)
	if blockState == nil || epochState == nil {
		return nil, 0, fmt.Errorf("target epoch not available")
	}

	if err := gdb.StateDbManager.Open(); err != nil {
		return nil, 0, fmt.Errorf("failed to open StateDbManager; %w", err)
	}
	log.Warn("Carmen db has must be replace manually",
		"block", blockState.LastBlock.Idx, "stateRoot", blockState.FinalizedStateRoot)

	// set the historic state to be the current
	log.Info("Reverting to epoch state", "epoch", epochState.Epoch, "block", blockState.LastBlock.Idx)
	gdb.SetBlockEpochState(*blockState, *epochState)
	gdb.FlushBlockEpochState()

	// Service.switchEpochTo
	gdb.SetHighestLamport(0)
	gdb.FlushHighestLamport()

	// removing excessive events (event epoch >= closed epoch)
	log.Info("Removing excessive events")
	gdb.ForEachEventRLP(epochState.Epoch.Bytes(), func(id hash.Event, _ rlp.RawValue) bool {
		gdb.DelEvent(id)
		return true
	})

	return epochState, oldTopEpoch, nil
}

func eraseTable(name string, producer kvdb.IterableDBProducer) error {
	log.Info("Cleaning table", "name", name)
	db, err := producer.OpenDB(name)
	if err != nil {
		return fmt.Errorf("unable to open DB %s; %s", name, err)
	}
	db = batched.Wrap(db)
	defer db.Close()
	it := db.NewIterator(nil, nil)
	defer it.Release()
	for it.Next() {
		err := db.Delete(it.Key())
		if err != nil {
			return err
		}
	}
	return nil
}

// clearDirtyFlags - writes the CleanPrefix into all databases
func clearDirtyFlags(id []byte, rawProducer kvdb.IterableDBProducer) error {
	names := rawProducer.Names()
	for _, name := range names {
		db, err := rawProducer.OpenDB(name)
		if err != nil {
			return err
		}

		err = db.Put(integration.FlushIDKey, append([]byte{flushable.CleanPrefix}, id...))
		if err != nil {
			log.Crit("Failed to write CleanPrefix", "name", name)
			return err
		}
		log.Info("Database set clean", "name", name)
		_ = db.Close()
	}
	return nil
}

func mustOpenDB(producer kvdb.DBProducer, name string) kvdb.Store {
	db, err := producer.OpenDB(name)
	if err != nil {
		utils.Fatalf("Failed to open '%s' database: %v", name, err)
	}
	return db
}

func panics(name string) func(error) {
	return func(err error) {
		log.Crit(fmt.Sprintf("%s error", name), "err", err)
	}
}
