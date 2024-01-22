package integration

import (
	"github.com/Fantom-foundation/lachesis-base/kvdb/multidb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

/*
 * pbl-1 config
 */

func DefaultDBsConfig(scale func(uint64) uint64, fdlimit uint64) DBsConfig {
	return DBsConfig{
		Routing:      Pbl1RoutingConfig(),
		RuntimeCache: Pbl1RuntimeDBsCacheConfig(scale, fdlimit),
		GenesisCache: Pbl1GenesisDBsCacheConfig(scale, fdlimit),
	}
}

func Pbl1RoutingConfig() RoutingConfig {
	return RoutingConfig{
		Table: map[string]multidb.Route{
			"": {
				Type: "pebble-fsh",
			},
		},
	}
}

func Pbl1RuntimeDBsCacheConfig(scale func(uint64) uint64, fdlimit uint64) DBsCacheConfig {
	return DBsCacheConfig{
		Table: map[string]DBCacheConfig{
			"": {
				Cache:   scale(480 * opt.MiB),
				Fdlimit: fdlimit*480/1400 + 1,
			},
		},
	}
}

func Pbl1GenesisDBsCacheConfig(scale func(uint64) uint64, fdlimit uint64) DBsCacheConfig {
	return DBsCacheConfig{
		Table: map[string]DBCacheConfig{
			"": {
				Cache:   scale(1000 * opt.MiB),
				Fdlimit: fdlimit*1000/3000 + 1,
			},
		},
	}
}
