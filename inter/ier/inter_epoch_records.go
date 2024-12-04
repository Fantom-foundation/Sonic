package ier

import (
	"github.com/Fantom-foundation/lachesis-base/ltypes"

	"github.com/Fantom-foundation/go-opera/inter/iblockproc"
)

type LlrFullEpochRecord struct {
	BlockState iblockproc.BlockState
	EpochState iblockproc.EpochState
}

type LlrIdxFullEpochRecord struct {
	LlrFullEpochRecord
	Idx ltypes.EpochID
}

func (er LlrFullEpochRecord) Hash() ltypes.Hash {
	return ltypes.Of(er.BlockState.Hash().Bytes(), er.EpochState.Hash().Bytes())
}
