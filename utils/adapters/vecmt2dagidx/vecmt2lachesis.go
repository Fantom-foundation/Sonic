package vecmt2dagidx

import (
	"github.com/Fantom-foundation/go-opera/vecclock"
	"github.com/Fantom-foundation/go-opera/vecclock/highestbefore"
	"github.com/Fantom-foundation/lachesis-base/abft"
	"github.com/Fantom-foundation/lachesis-base/abft/dagidx"
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
)

type Adapter struct {
	*vecclock.Index
}

var _ abft.DagIndex = (*Adapter)(nil)

type SeqAdapter struct {
	highestbefore.Type
}

func (s SeqAdapter) Seq() idx.Event {
	return s.Type.Seq
}

func (s SeqAdapter) MinSeq() idx.Event {
	return s.Type.MinSeq
}

func (s SeqAdapter) IsForkDetected() bool {
	return s.Type.IsForkDetected()
}

type HighestBeforeAdapter highestbefore.Types

func (h HighestBeforeAdapter) Size() int {
	return len(h)
}

func (h HighestBeforeAdapter) Get(i idx.Validator) dagidx.Seq {
	v := highestbefore.Types(h).Get(i)
	return SeqAdapter{v}
}

func (v *Adapter) GetMergedHighestBefore(id hash.Event) dagidx.HighestBeforeSeq {
	return HighestBeforeAdapter(v.Index.GetMergedHighestBefore(id))
}

func Wrap(v *vecclock.Index) *Adapter {
	return &Adapter{v}
}
