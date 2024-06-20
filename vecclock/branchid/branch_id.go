package branchid

import (
	"github.com/Fantom-foundation/go-opera/utils/swapmap"
	"github.com/Fantom-foundation/lachesis-base/common/bigendian"
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/Fantom-foundation/lachesis-base/kvdb"
	"github.com/ethereum/go-ethereum/log"
)

func deserialize(bb []byte) idx.Validator {
	return idx.Validator(bigendian.BytesToUint32(bb))
}

type Map struct {
	swapmap.Map[hash.Event, idx.Validator]
}

func NewMap(maxMemSize int, getTmpDB func(name string) kvdb.Store) Map {
	sm := swapmap.New[hash.Event, idx.Validator](swapmap.Callbacks[hash.Event, idx.Validator]{
		SizeOf: func(k hash.Event, v idx.Validator) int {
			return 100
		},
		SerializeK: func(k hash.Event) []byte {
			return k.Bytes()
		},
		SerializeV: func(v idx.Validator) []byte {
			return v.Bytes()
		},
		DeserializeV: func(bb []byte) idx.Validator {
			return deserialize(bb)
		},
		GetSwapDB: func() kvdb.Store {
			return getTmpDB("branchid")
		},
		Fatal: func(err error) {
			log.Crit("branchid SwapMap error", "err", err)
		},
	}, maxMemSize)
	return Map{*sm}
}
