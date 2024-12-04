package vecmt

import (
	"encoding/binary"

	"github.com/Fantom-foundation/lachesis-base/ltypes"
	"github.com/Fantom-foundation/lachesis-base/vecfc"

	"github.com/Fantom-foundation/go-opera/inter"
)

/*
 * Use binary form for optimization, to avoid serialization. As a result, DB cache works as elements cache.
 */

type (
	// HighestBeforeTime is a vector of highest events (their CreationTime) which are observed by source event
	HighestBeforeTime []byte

	HighestBefore struct {
		VSeq  *vecfc.HighestBeforeSeq
		VTime *HighestBeforeTime
	}
)

// NewHighestBefore creates new HighestBefore vector.
func NewHighestBefore(size ltypes.ValidatorIdx) *HighestBefore {
	return &HighestBefore{
		VSeq:  vecfc.NewHighestBeforeSeq(size),
		VTime: NewHighestBeforeTime(size),
	}
}

// NewHighestBeforeTime creates new HighestBeforeTime vector.
func NewHighestBeforeTime(size ltypes.ValidatorIdx) *HighestBeforeTime {
	b := make(HighestBeforeTime, size*8)
	return &b
}

// Get i's position in the byte-encoded vector clock
func (b HighestBeforeTime) Get(i ltypes.ValidatorIdx) inter.Timestamp {
	for i >= b.Size() {
		return 0
	}
	return inter.Timestamp(binary.LittleEndian.Uint64(b[i*8 : (i+1)*8]))
}

// Set i's position in the byte-encoded vector clock
func (b *HighestBeforeTime) Set(i ltypes.ValidatorIdx, time inter.Timestamp) {
	for i >= b.Size() {
		// append zeros if exceeds size
		*b = append(*b, []byte{0, 0, 0, 0, 0, 0, 0, 0}...)
	}
	binary.LittleEndian.PutUint64((*b)[i*8:(i+1)*8], uint64(time))
}

// Size of the vector clock
func (b HighestBeforeTime) Size() ltypes.ValidatorIdx {
	return ltypes.ValidatorIdx(len(b) / 8)
}
