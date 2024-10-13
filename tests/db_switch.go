package tests

import (
	carmen "github.com/Fantom-foundation/Carmen/go/state"
	"github.com/Fantom-foundation/Carmen/go/state/gostate"
	"github.com/Fantom-foundation/go-opera/gossip/evmstore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"io"
)

type stateDb interface {
	vm.StateDB
	io.Closer

	Logs() []*types.Log
	Commit() (stateRootHash common.Hash)
	Reopen()
}

type carmenDb struct {
	vm.StateDB
	db   *evmstore.CarmenStateDB
	st   carmen.State
	dir  string
	logs []*types.Log
}

func newCarmenDb(dir string) (stateDb, error) {
	parameters := carmen.Parameters{
		Variant:   gostate.VariantGoMemory,
		Schema:    carmen.Schema(5),
		Archive:   carmen.NoArchive,
		Directory: dir,
	}

	st, err := carmen.NewState(parameters)
	if err != nil {
		return nil, err
	}

	carmenstatedb := carmen.CreateCustomStateDBUsing(st, 1024)
	statedb := evmstore.CreateCarmenStateDb(carmenstatedb)

	return &carmenDb{
		StateDB: statedb,
		db:      statedb.(*evmstore.CarmenStateDB),
		st:      st,
		dir:     dir,
	}, nil
}

func (db *carmenDb) Close() error {
	return db.st.Close()
}

func (db *carmenDb) Logs() []*types.Log {
	return db.logs
}

func (db *carmenDb) Commit() common.Hash {
	// save logs before carmen clears them on tx and block end
	db.logs = db.db.Logs()

	db.db.Finalise() // ends transaction
	db.db.EndBlock(0)
	return db.db.GetStateHash()
}

func (db *carmenDb) Reopen() {
	carmenstatedb := carmen.CreateCustomStateDBUsing(db.st, 1024)
	statedb := evmstore.CreateCarmenStateDb(carmenstatedb)

	db.StateDB = statedb
	db.db = statedb.(*evmstore.CarmenStateDB)
}

type gethDb struct {
	vm.StateDB
	db  *state.StateDB
	sdb state.Database
}

func newGethDb() (stateDb, error) {
	sdb := state.NewDatabase(rawdb.NewMemoryDatabase())
	statedb, _ := state.New(types.EmptyRootHash, sdb, nil)

	return &gethDb{
		StateDB: statedb,
		db:      statedb,
		sdb:     sdb,
	}, nil
}

func (db *gethDb) Close() error {
	return nil
}

func (db *gethDb) Logs() []*types.Log {
	return db.db.Logs()
}

func (db *gethDb) Commit() common.Hash {
	root, _ := db.db.Commit(0, true)
	return root
}

func (db *gethDb) Reopen() {
	root := db.Commit()
	statedb, _ := state.New(root, db.sdb, nil)

	db.StateDB = statedb
	db.db = statedb
}
