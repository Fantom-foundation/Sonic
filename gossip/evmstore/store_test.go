package evmstore

import (
	"github.com/Fantom-foundation/lachesis-base/kvdb/memorydb"
)

func cachedStore() *Store {
	cfg := LiteStoreConfig()

	return NewStore(memorydb.NewProducer(""), cfg, nil)
}

func nonCachedStore() *Store {
	cfg := StoreConfig{}

	return NewStore(memorydb.NewProducer(""), cfg, nil)
}
