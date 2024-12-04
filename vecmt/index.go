package vecmt

import (
	"github.com/Fantom-foundation/lachesis-base/kvdb"
	"github.com/Fantom-foundation/lachesis-base/kvdb/table"
	"github.com/Fantom-foundation/lachesis-base/ltypes"
	"github.com/Fantom-foundation/lachesis-base/utils/cachescale"
	"github.com/Fantom-foundation/lachesis-base/utils/wlru"
	"github.com/Fantom-foundation/lachesis-base/vecengine"
	"github.com/Fantom-foundation/lachesis-base/vecfc"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// IndexCacheConfig - config for cache sizes of Engine
type IndexCacheConfig struct {
	HighestBeforeTimeSize uint
	DBCache               int
}

// IndexConfig - Engine config (cache sizes)
type IndexConfig struct {
	Fc     vecfc.IndexConfig
	Caches IndexCacheConfig
}

// Index is a data to detect forkless-cause condition, calculate median timestamp, detect forks.
type Index struct {
	*vecfc.Index
	Base          *vecfc.Index
	baseCallbacks vecengine.Callbacks

	crit          func(error)
	validators    *ltypes.Validators
	validatorIdxs map[ltypes.ValidatorID]ltypes.ValidatorIdx

	getEvent func(ltypes.EventHash) ltypes.Event

	vecDb kvdb.Store
	table struct {
		HighestBeforeTime kvdb.Store `table:"T"`
	}

	cache struct {
		HighestBeforeTime *wlru.Cache
	}

	cfg IndexConfig
}

// DefaultConfig returns default index config
func DefaultConfig(scale cachescale.Func) IndexConfig {
	return IndexConfig{
		Fc: vecfc.DefaultConfig(scale),
		Caches: IndexCacheConfig{
			HighestBeforeTimeSize: scale.U(160 * 1024),
			DBCache:               scale.I(10 * opt.MiB),
		},
	}
}

// LiteConfig returns default index config for tests
func LiteConfig() IndexConfig {
	return IndexConfig{
		Fc: vecfc.LiteConfig(),
		Caches: IndexCacheConfig{
			HighestBeforeTimeSize: 4 * 1024,
		},
	}
}

// NewIndex creates Index instance.
func NewIndex(crit func(error), config IndexConfig) *Index {
	vi := &Index{
		cfg:  config,
		crit: crit,
	}
	engine := vecengine.NewIndex(crit, vi.GetEngineCallbacks())

	vi.Base = vecfc.NewIndexWithEngine(crit, config.Fc, engine)
	vi.Index = vi.Base
	vi.baseCallbacks = vi.Base.GetEngineCallbacks()
	vi.initCaches()

	return vi
}

func (vi *Index) initCaches() {
	vi.cache.HighestBeforeTime, _ = wlru.New(vi.cfg.Caches.HighestBeforeTimeSize, int(vi.cfg.Caches.HighestBeforeTimeSize))
}

// Reset resets buffers.
func (vi *Index) Reset(validators *ltypes.Validators, db kvdb.Store, getEvent func(ltypes.EventHash) ltypes.Event) {
	fdb := WrapByVecFlushable(db, vi.cfg.Caches.DBCache)
	vi.vecDb = fdb
	vi.Base.Reset(validators, fdb, getEvent)
	vi.getEvent = getEvent
	vi.validators = validators
	vi.validatorIdxs = validators.Idxs()
	vi.onDropNotFlushed()

	table.MigrateTables(&vi.table, vi.vecDb)
}

func (vi *Index) Close() error {
	return vi.vecDb.Close()
}

func (vi *Index) GetEngineCallbacks() vecengine.Callbacks {
	return vecengine.Callbacks{
		GetHighestBefore: func(event ltypes.EventHash) vecengine.HighestBeforeI {
			return vi.GetHighestBefore(event)
		},
		GetLowestAfter: func(event ltypes.EventHash) vecengine.LowestAfterI {
			return vi.baseCallbacks.GetLowestAfter(event)
		},
		SetHighestBefore: func(event ltypes.EventHash, b vecengine.HighestBeforeI) {
			vi.SetHighestBefore(event, b.(*HighestBefore))
		},
		SetLowestAfter: func(event ltypes.EventHash, i vecengine.LowestAfterI) {
			vi.baseCallbacks.SetLowestAfter(event, i)
		},
		NewHighestBefore: func(size ltypes.ValidatorIdx) vecengine.HighestBeforeI {
			return NewHighestBefore(size)
		},
		NewLowestAfter: func(size ltypes.ValidatorIdx) vecengine.LowestAfterI {
			return vi.baseCallbacks.NewLowestAfter(size)
		},
		OnDropNotFlushed: func() {
			vi.baseCallbacks.OnDropNotFlushed()
			vi.onDropNotFlushed()
		},
	}
}

func (vi *Index) onDropNotFlushed() {
	vi.cache.HighestBeforeTime.Purge()
}

// GetMergedHighestBefore returns HighestBefore vector clock without branches, where branches are merged into one
func (vi *Index) GetMergedHighestBefore(id ltypes.EventHash) *HighestBefore {
	return vi.Engine.GetMergedHighestBefore(id).(*HighestBefore)
}
