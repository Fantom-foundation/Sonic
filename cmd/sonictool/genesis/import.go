package genesis

import (
	"fmt"
	"github.com/Fantom-foundation/go-opera/cmd/sonictool/db"
	"github.com/Fantom-foundation/go-opera/config/flags"
	"github.com/Fantom-foundation/go-opera/opera/genesis"
	"github.com/Fantom-foundation/go-opera/opera/genesisstore"
	"github.com/Fantom-foundation/lachesis-base/abft"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/Fantom-foundation/lachesis-base/kvdb"
	"github.com/Fantom-foundation/lachesis-base/utils/cachescale"
	"github.com/ethereum/go-ethereum/log"
	"gopkg.in/urfave/cli.v1"
	"path/filepath"
)

func ImportGenesisStore(ctx *cli.Context, genesisStore *genesisstore.Store, dataDir string, validatorMode bool, cacheRatio cachescale.Func) error {
	if err := db.AssertDatabaseNotInitialized(dataDir); err != nil {
		return fmt.Errorf("database in datadir is already initialized: %w", err)
	}
	if err := db.RemoveDatabase(dataDir); err != nil {
		return fmt.Errorf("failed to remove existing data from the datadir: %w", err)
	}

	chaindataDir := filepath.Join(dataDir, "chaindata")
	dbs, err := db.MakeDbProducer(chaindataDir, cacheRatio)
	if err != nil {
		return err
	}
	defer dbs.Close()
	setGenesisProcessing(chaindataDir)

	gdb, err := db.MakeGossipDb(db.GossipDbParameters{
		Dbs:           dbs,
		DataDir:       dataDir,
		ValidatorMode: validatorMode,
		CacheRatio:    cacheRatio,
		LiveDbCache:   ctx.Int64(flags.LiveDbCacheFlag.Name),
		ArchiveCache:  ctx.Int64(flags.ArchiveCacheFlag.Name),
	})
	if err != nil {
		return err
	}
	defer gdb.Close()

	err = gdb.ApplyGenesis(genesisStore.Genesis())
	if err != nil {
		return fmt.Errorf("failed to write Gossip genesis state: %v", err)
	}

	cMainDb, err := dbs.OpenDB("lachesis")
	if err != nil {
		return err
	}
	cGetEpochDB := func(epoch idx.Epoch) kvdb.Store {
		db, err := dbs.OpenDB(fmt.Sprintf("lachesis-%d", epoch))
		if err != nil {
			panic(fmt.Errorf("failed to open epoch db: %w", err))
		}
		return db
	}
	abftCrit := func(err error) {
		panic(fmt.Errorf("lachesis store error: %w", err))
	}
	cdb := abft.NewStore(cMainDb, cGetEpochDB, abftCrit, abft.DefaultStoreConfig(cacheRatio))
	defer cdb.Close()

	err = cdb.ApplyGenesis(&abft.Genesis{
		Epoch:      gdb.GetEpoch(),
		Validators: gdb.GetValidators(),
	})
	if err != nil {
		return fmt.Errorf("failed to write lachesis genesis state: %w", err)
	}

	err = gdb.Commit()
	if err != nil {
		return err
	}
	setGenesisComplete(chaindataDir)
	log.Info("Successfully imported genesis file")
	return nil
}

func IsGenesisTrusted(genesisStore *genesisstore.Store, genesisHashes genesis.Hashes) error {
	g := genesisStore.Genesis()

	// try trusted hashes first
	for _, allowed := range allowedGenesis {
		if allowed.Hashes.Equal(genesisHashes) && allowed.Header.Equal(g.Header) {
			return nil
		}
	}

	// try using SignedMetadata section
	hash, _, err := GetGenesisMetadata(g.Header, genesisHashes)
	if err != nil {
		return fmt.Errorf("failed to calculate hash of genesis: %w", err)
	}
	signature, err := g.SignatureSection.GetSignature()
	if err != nil {
		return fmt.Errorf("genesis file doesn't refer to any trusted preset, signature not found: %w", err)
	}
	if err := CheckGenesisSignature(hash, signature); err != nil {
		return fmt.Errorf("genesis file doesn't refer to any trusted preset: %w", err)
	}
	return nil
}
