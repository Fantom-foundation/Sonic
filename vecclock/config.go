package vecclock

import (
	"github.com/Fantom-foundation/lachesis-base/utils/cachescale"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// IndexCacheConfig - config for cache sizes of Engine
type IndexCacheConfig struct {
	ForklessCausePairs int
	LowestAfter        int
	HighestBefore      int
	BranchID           int
}

// Config - Engine config (cache sizes)
type Config struct {
	Caches IndexCacheConfig
}

// DefaultConfig returns default index config
func DefaultConfig(scale cachescale.Func) Config {
	return Config{
		Caches: IndexCacheConfig{
			ForklessCausePairs: scale.I(20000),
			LowestAfter:        scale.I(40 * opt.MiB),
			HighestBefore:      scale.I(120 * opt.MiB),
			BranchID:           scale.I(10 * opt.MiB),
		},
	}
}

// LiteConfig returns default index config for tests
func LiteConfig() Config {
	return DefaultConfig(cachescale.Ratio{Base: 100, Target: 1})
}
