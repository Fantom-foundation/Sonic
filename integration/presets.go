package integration

import (
	"github.com/Fantom-foundation/lachesis-base/kvdb/multidb"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

/*
 * pbl-1 config
 */

func DefaultDBsConfig(scale func(uint64) uint64) DBsConfig {
	return DBsConfig{
		Routing: RoutingConfig{
			Table: map[string]multidb.Route{
				"": {
					Type: "pebble-fsh",
				},
			},
		},
		RuntimeCache: DBCacheConfig{
			Cache:   scale(480 * opt.MiB),
			Fdlimit: uint64(utils.MakeDatabaseHandles())*480/1400 + 1,
		},
		GenesisCache: DBCacheConfig{
			Cache:   scale(1000 * opt.MiB),
			Fdlimit: uint64(utils.MakeDatabaseHandles())*1000/3000 + 1,
		},
	}
}
