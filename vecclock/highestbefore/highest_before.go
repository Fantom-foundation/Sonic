package highestbefore

import (
	"github.com/Fantom-foundation/go-opera/inter"
	"github.com/Fantom-foundation/go-opera/utils/swapmap"
	"github.com/Fantom-foundation/lachesis-base/common/bigendian"
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/Fantom-foundation/lachesis-base/kvdb"
	"github.com/ethereum/go-ethereum/log"
	"math"
)

type Type struct {
	Seq    idx.Event
	MinSeq idx.Event
	Time   inter.Timestamp
}

func (b Type) IsForkDetected() bool {
	return b.Seq == ForkDetectedSeq.Seq && b.MinSeq == ForkDetectedSeq.MinSeq
}

func (b *Type) SetForkDetected() {
	*b = ForkDetectedSeq
}

func (b Type) IsEmpty() bool {
	return !b.IsForkDetected() && b.Seq == 0
}

type Types []Type

func (b Types) Get(i idx.Validator) Type {
	if i >= idx.Validator(len(b)) {
		return Type{}
	}
	return b[i]
}

func (b *Types) Alloc(i idx.Validator) {
	for i >= idx.Validator(len(*b)) {
		// append zeros if exceeds size
		*b = append(*b, Type{})
	}
}

var (
	// ForkDetectedSeq is a special marker of observed fork by a creator
	ForkDetectedSeq = Type{
		Seq:    0,
		MinSeq: idx.Event(math.MaxInt32),
		Time:   0,
	}
)

func serialize(hb Types) []byte {
	bb := make([]byte, 0, len(hb)*16)
	for i := 0; i < len(hb); i++ {
		bb = append(bb, hb[i].Seq.Bytes()...)
		bb = append(bb, hb[i].MinSeq.Bytes()...)
		bb = append(bb, hb[i].Time.Bytes()...)
	}
	return bb
}

func deserialize(bb []byte) Types {
	hb := make(Types, 0, len(bb)/16)
	for i := 0; i < len(bb); i += 16 {
		h := Type{}
		h.Seq = idx.Event(bigendian.BytesToUint32(bb[i : i+4]))
		h.MinSeq = idx.Event(bigendian.BytesToUint32(bb[i+4 : i+8]))
		h.Time = inter.Timestamp(bigendian.BytesToUint64(bb[i+8 : i+16]))
		hb = append(hb, h)
	}
	return hb
}

type Map struct {
	swapmap.Map[hash.Event, Types]
}

func NewMap(maxMemSize int, getTmpDB func(name string) kvdb.Store) Map {
	sm := swapmap.New[hash.Event, Types](swapmap.Callbacks[hash.Event, Types]{
		SizeOf: func(k hash.Event, v Types) int {
			return 100 + 32 + len(v)*16
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
			return getTmpDB("hb")
		},
		Fatal: func(err error) {
			log.Crit("hb SwapMap error", "err", err)
		},
	}, maxMemSize)
	return Map{*sm}
}
