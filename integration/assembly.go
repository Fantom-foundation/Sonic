package integration

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/Fantom-foundation/go-opera/gossip"
	"github.com/Fantom-foundation/go-opera/utils/adapters/vecmt2dagidx"
	"github.com/Fantom-foundation/go-opera/vecmt"
	"github.com/Fantom-foundation/lachesis-base/abft"
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/Fantom-foundation/lachesis-base/kvdb"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/status-im/keycard-go/hexutils"
)

var (
	MetadataPrefix = hexutils.HexToBytes("0068c2927bf842c3e9e2f1364494a33a752db334b9a819534bc9f17d2c3b4e5970008ff519d35a86f29fcaa5aae706b75dee871f65f174fcea1747f2915fc92158f6bfbf5eb79f65d16225738594bffb")
	FlushIDKey     = append(common.CopyBytes(MetadataPrefix), 0x0c)
	TablesKey      = append(common.CopyBytes(MetadataPrefix), 0x0d)
)

// GenesisMismatchError is raised when trying to overwrite an existing
// genesis block with an incompatible one.
type GenesisMismatchError struct {
	Stored, New hash.Hash
}

// Error implements error interface.
func (e *GenesisMismatchError) Error() string {
	return fmt.Sprintf("database contains incompatible genesis (have %s, new %s)", e.Stored.String(), e.New.String())
}

type Configs struct {
	Opera         gossip.Config
	OperaStore    gossip.StoreConfig
	Lachesis      abft.Config
	LachesisStore abft.StoreConfig
	VectorClock   vecmt.IndexConfig
	DBs           DBsConfig
}

func panics(name string) func(error) {
	return func(err error) {
		log.Crit(fmt.Sprintf("%s error", name), "err", err)
	}
}

func mustOpenDB(producer kvdb.DBProducer, name string) kvdb.Store {
	db, err := producer.OpenDB(name)
	if err != nil {
		utils.Fatalf("Failed to open '%s' database: %v", name, err)
	}
	return db
}

func getStores(producer kvdb.FlushableDBProducer, cfg Configs) (*gossip.Store, *abft.Store, error) {
	gdb, err := gossip.NewStore(producer, cfg.OperaStore)
	if err != nil {
		return nil, nil, err
	}

	cMainDb := mustOpenDB(producer, "lachesis")
	cGetEpochDB := func(epoch idx.Epoch) kvdb.Store {
		return mustOpenDB(producer, fmt.Sprintf("lachesis-%d", epoch))
	}
	cdb := abft.NewStore(cMainDb, cGetEpochDB, panics("Lachesis store"), cfg.LachesisStore)
	return gdb, cdb, err
}

func rawMakeEngine(gdb *gossip.Store, cdb *abft.Store, cfg Configs) (*abft.Lachesis, *vecmt.Index, gossip.BlockProc, error) {
	blockProc := gossip.DefaultBlockProc()
	// create consensus
	vecClock := vecmt.NewIndex(panics("Vector clock"), cfg.VectorClock)
	engine := abft.NewLachesis(cdb, &GossipStoreAdapter{gdb}, vecmt2dagidx.Wrap(vecClock), panics("Lachesis"), cfg.Lachesis)
	return engine, vecClock, blockProc, nil
}

func CheckStateInitialized(chaindataDir string, cfg DBsConfig) error {
	if isInterrupted(chaindataDir) {
		return errors.New("genesis processing isn't finished")
	}
	dbs, err := GetDbProducer(chaindataDir, cfg.RuntimeCache)
	if err != nil {
		return err
	}
	return dbs.Close()
}

func makeEngine(chaindataDir string, cfg Configs) (*abft.Lachesis, *vecmt.Index, *gossip.Store, *abft.Store, gossip.BlockProc, func() error, error) {
	dbs, err := GetDbProducer(chaindataDir, cfg.DBs.RuntimeCache)
	if err != nil {
		return nil, nil, nil, nil, gossip.BlockProc{}, nil, err
	}

	gdb, cdb, err := getStores(dbs, cfg)
	if err != nil {
		err = fmt.Errorf("failed to get stores: %w", err)
		return nil, nil, nil, nil, gossip.BlockProc{}, nil, err
	}
	defer func() {
		if err != nil {
			gdb.Close()
			cdb.Close()
			dbs.Close()
		}
	}()

	err = gdb.EvmStore().Open()
	if err != nil {
		err = fmt.Errorf("failed to open EvmStore: %v", err)
		return nil, nil, nil, nil, gossip.BlockProc{}, dbs.Close, err
	}

	engine, vecClock, blockProc, err := rawMakeEngine(gdb, cdb, cfg)
	if err != nil {
		err = fmt.Errorf("failed to make engine: %v", err)
		return nil, nil, nil, nil, gossip.BlockProc{}, dbs.Close, err
	}

	return engine, vecClock, gdb, cdb, blockProc, dbs.Close, nil
}

// MakeEngine makes consensus engine from config.
func MakeEngine(chaindataDir string, cfg Configs) (*abft.Lachesis, *vecmt.Index, *gossip.Store, *abft.Store, gossip.BlockProc, func() error, error) {
	if isEmpty(chaindataDir) || isInterrupted(chaindataDir) {
		return nil, nil, nil, nil, gossip.BlockProc{}, nil, fmt.Errorf("database is empty or the genesis import interrupted")
	}

	engine, vecClock, gdb, cdb, blockProc, closeDBs, err := makeEngine(chaindataDir, cfg)
	if err != nil {
		return nil, nil, nil, nil, gossip.BlockProc{}, nil, err
	}

	rules := gdb.GetRules()
	genesisID := gdb.GetGenesisID()
	log.Info("Genesis is written", "name", rules.Name, "id", rules.NetworkID, "genesis", genesisID.String())

	return engine, vecClock, gdb, cdb, blockProc, closeDBs, nil
}

// SetAccountKey sets key into accounts manager and unlocks it with pswd.
func SetAccountKey(
	am *accounts.Manager, key *ecdsa.PrivateKey, pswd string,
) (
	acc accounts.Account,
) {
	kss := am.Backends(keystore.KeyStoreType)
	if len(kss) < 1 {
		log.Crit("Keystore is not found")
		return
	}
	ks := kss[0].(*keystore.KeyStore)

	acc = accounts.Account{
		Address: crypto.PubkeyToAddress(key.PublicKey),
	}

	imported, err := ks.ImportECDSA(key, pswd)
	if err == nil {
		acc = imported
	} else if err.Error() != "account already exists" {
		log.Crit("Failed to import key", "err", err)
	}

	err = ks.Unlock(acc, pswd)
	if err != nil {
		log.Crit("failed to unlock key", "err", err)
	}

	return
}
