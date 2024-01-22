package integration

import (
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

/*
 * pbl-1 config
 */

func DefaultDBsConfig(scale func(uint64) uint64) DBsConfig {
	return DBsConfig{
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
