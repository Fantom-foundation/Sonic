package lowestafter

import (
	"github.com/Fantom-foundation/go-opera/utils/swapmap"
	"github.com/Fantom-foundation/lachesis-base/common/bigendian"
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/Fantom-foundation/lachesis-base/kvdb"
	"github.com/ethereum/go-ethereum/log"
)

type Types []idx.Event

func (b Types) Get(i idx.Validator) idx.Event {
	if i >= idx.Validator(len(b)) {
		return 0
	}
	return b[i]
}

func (b *Types) Alloc(i idx.Validator) {
	for i >= idx.Validator(len(*b)) {
		// append zeros if exceeds size
		*b = append(*b, 0)
	}
}

func serialize(la Types) []byte {
	bb := make([]byte, 0, len(la)*4)
	for i := 0; i < len(la); i++ {
		bb = append(bb, la[i].Bytes()...)
	}
	return bb
}

func deserialize(bb []byte) Types {
	la := make(Types, 0, len(bb)/4)
	for i := 0; i < len(bb); i += 4 {
		la = append(la, idx.Event(bigendian.BytesToUint32(bb[i:i+4])))
	}
	return la
}

type Map struct {
	swapmap.Map[hash.Event, Types]
}

func NewMap(maxMemSize int, getTmpDB func(name string) kvdb.Store) Map {
	sm := swapmap.New[hash.Event, Types](swapmap.Callbacks[hash.Event, Types]{
		SizeOf: func(k hash.Event, v Types) int {
			return 100 + 32 + len(v)*4
		},
		SerializeK: func(k hash.Event) []byte {
			return k.Bytes()
		},
		SerializeV: func(v Types) []byte {
			return serialize(v)
		},
		DeserializeV: func(bb []byte) Types {
			return deserialize(bb)
		},
		GetSwapDB: func() kvdb.Store {
			return getTmpDB("la")
		},
		Fatal: func(err error) {
			log.Crit("la SwapMap error", "err", err)
		},
	}, maxMemSize)
	return Map{*sm}
}
