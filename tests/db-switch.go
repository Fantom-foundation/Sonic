package tests

import (
	"errors"
	carmen "github.com/Fantom-foundation/Carmen/go/state"
	"github.com/Fantom-foundation/Carmen/go/state/gostate"
	"github.com/Fantom-foundation/go-opera/gossip/evmstore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"io"
	"os"
)

type stateDb interface {
	vm.StateDB
	io.Closer

	Logs() []*types.Log
	Commit() (stateRootHash common.Hash)
}

type carmenDb struct {
	vm.StateDB
	db  *evmstore.CarmenStateDB
	st  carmen.State
	dir string
}

func newCarmenDb() (stateDb, error) {
	dir, err := os.MkdirTemp("", "eth-tests-db-*")
	if err != nil {
		return nil, err
	}

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
	return errors.Join(
		db.st.Close(),
		os.RemoveAll(db.dir),
	)
}

func (db *carmenDb) Logs() []*types.Log {
	return db.db.Logs()
}

func (db *carmenDb) Commit() common.Hash {
	db.db.Finalise()
	return db.db.GetStateHash()
}

type gethDb struct {
	vm.StateDB
	db *state.StateDB
}

func newGethDb() (stateDb, error) {
	sdb := state.NewDatabase(rawdb.NewMemoryDatabase())
	statedb, _ := state.New(types.EmptyRootHash, sdb, nil)

	return &gethDb{
		StateDB: statedb,
		db:      statedb,
	}, nil
}

func (db *gethDb) Close() error {
	return nil
}

func (db *gethDb) Logs() []*types.Log {
	return db.db.Logs()
}

func (db *gethDb) Commit() common.Hash {
	root, _ := db.db.Commit(1, true)
	return root
}
