package gossip

import (
	"fmt"

	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/kvdb"
	"github.com/ethereum/go-ethereum/common"

	"github.com/Fantom-foundation/go-opera/inter"
	"github.com/Fantom-foundation/go-opera/utils/migration"
)

func isEmptyDB(db kvdb.Iteratee) bool {
	it := db.NewIterator(nil, nil)
	defer it.Release()
	return !it.Next()
}

func (s *Store) migrateData() error {
	versions := migration.NewKvdbIDStore(s.table.Version)
	if isEmptyDB(s.table.Version) {
		// short circuit if empty DB
		versions.SetID(s.migrations().ID())
		return nil
	}

	err := s.migrations().Exec(versions, s.flushDBs)
	return err
}

func (s *Store) migrations() *migration.Migration {
	return migration.
		Begin("opera-gossip-store").
		Next("used gas recovery", unsupportedMigration).
		Next("tx hashes recovery", unsupportedMigration).
		Next("DAG heads recovery", unsupportedMigration).
		Next("DAG last events recovery", unsupportedMigration).
		Next("BlockState recovery", unsupportedMigration).
		Next("LlrState recovery", unsupportedMigration).
		Next("erase gossip-async db", unsupportedMigration).
		Next("erase SFC API table", unsupportedMigration).
		Next("erase legacy genesis DB", unsupportedMigration).
		Next("calculate upgrade heights", unsupportedMigration)
}

func unsupportedMigration() error {
	return fmt.Errorf("DB version isn't supported, please restart from scratch")
}

var (
	fixTxHash1  = common.HexToHash("0xb6840d4c0eb562b0b1731760223d91b36edc6c160958e23e773e6058eea30458")
	fixTxEvent1 = hash.HexToEventHash("0x00001718000003d4d3955bf592e12fb80a60574fa4b18bd5805b4c010d75e86d")
	fixTxHash2  = common.HexToHash("0x3aeede91740093cb8feb1296a34cf70d86d2f802cff860edd798978e94a40534")
	fixTxEvent2 = hash.HexToEventHash("0x0000179e00000c464d756a7614d0ca067fcb37ee4452004bf308c9df561e85e8")
)

const (
	fixTxEventPos1 = 2
	fixTxBlock1    = 4738821
	fixTxEventPos2 = 0
	fixTxBlock2    = 4801307
)

func fixEventTxHashes(e *inter.EventPayload) {
	if e.ID() == fixTxEvent1 {
		e.Txs()[fixTxEventPos1].SetHash(fixTxHash1)
	}
	if e.ID() == fixTxEvent2 {
		e.Txs()[fixTxEventPos2].SetHash(fixTxHash2)
	}
}
