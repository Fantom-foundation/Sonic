package integration

import (
	"fmt"
	"github.com/Fantom-foundation/go-opera/gossip"
	"github.com/Fantom-foundation/go-opera/utils/dbutil/dbcounter"
	"github.com/Fantom-foundation/go-opera/utils/dbutil/threads"
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/dag"
	"github.com/Fantom-foundation/lachesis-base/kvdb"
	"github.com/Fantom-foundation/lachesis-base/kvdb/cachedproducer"
	"github.com/Fantom-foundation/lachesis-base/kvdb/flaggedproducer"
	"github.com/Fantom-foundation/lachesis-base/kvdb/leveldb"
	"github.com/Fantom-foundation/lachesis-base/kvdb/pebble"
	"github.com/Fantom-foundation/lachesis-base/kvdb/skipkeys"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
	"golang.org/x/exp/rand"
	"io"
	"os"
	"path"
	"sync/atomic"
)

type DBsConfig struct {
	RuntimeCache DBCacheConfig
	TmpDBCache   DBCacheConfig
}

type DBCacheConfig struct {
	Cache   uint64
	Fdlimit uint64
}

func GetRawDbProducer(chaindataDir string, cfg DBCacheConfig) kvdb.IterableDBProducer {
	if chaindataDir == "inmemory" || chaindataDir == "" {
		chaindataDir, _ = os.MkdirTemp("", "opera-tmp")
	}
	cacher := func(name string) (int, int) {
		return int(cfg.Cache), int(cfg.Fdlimit)
	}

	rawProducer := dbcounter.Wrap(pebble.NewProducer(chaindataDir, cacher), true)

	if metrics.Enabled {
		rawProducer = WrapDatabaseWithMetrics(rawProducer)
	}
	return rawProducer
}

func GetDbProducer(chaindataDir string, cfg DBCacheConfig) (kvdb.FullDBProducer, error) {
	rawProducer := GetRawDbProducer(chaindataDir, cfg)
	scopedProducer := flaggedproducer.Wrap(rawProducer, FlushIDKey) // pebble-flg
	_, err := scopedProducer.Initialize(rawProducer.Names(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open existing databases: %v", err)
	}
	cachedProducer := cachedproducer.WrapAll(scopedProducer)
	skippingProducer := skipkeys.WrapAllProducer(cachedProducer, MetadataPrefix)
	return threads.CountedFullDBProducer(skippingProducer), nil
}

func isEmpty(dir string) bool {
	f, err := os.Open(dir)
	if err != nil {
		return true
	}
	defer f.Close()
	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true
	}
	return false
}

type GossipStoreAdapter struct {
	*gossip.Store
}

func (g *GossipStoreAdapter) GetEvent(id hash.Event) dag.Event {
	e := g.Store.GetEvent(id)
	if e == nil {
		return nil
	}
	return e
}

func tmpDBsMaker(dir string, cache, handles uint64) func(string) kvdb.Store {
	dir = path.Join(dir, fmt.Sprintf("%x", rand.Uint64()))
	ldbProducer := leveldb.NewProducer(dir, func(_ string) (int, int) {
		return int(cache), int(handles)
	})
	counter := atomic.Uint64{}

	return func(name string) kvdb.Store {
		seq := counter.Add(1)
		name += fmt.Sprintf("%s-%d", name, seq)
		tmpDB, err := ldbProducer.OpenDB(name)
		if err != nil {
			log.Crit("Failed to create tmp database", "name", name, "err", err)
		}
		return tmpDB
	}
}
