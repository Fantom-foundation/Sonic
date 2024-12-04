package emitter

import (
	"io"

	"github.com/Fantom-foundation/lachesis-base/ltypes"
	"github.com/ethereum/go-ethereum/log"

	"github.com/Fantom-foundation/go-opera/utils"
)

var openPrevActionFile = utils.OpenFile

func (em *Emitter) writeLastEmittedEventID(id ltypes.EventHash) {
	if em.emittedEventFile == nil {
		return
	}
	_, err := em.emittedEventFile.WriteAt(id.Bytes(), 0)
	if err != nil {
		log.Crit("Failed to write event file", "file", em.config.PrevEmittedEventFile.Path, "err", err)
	}
}

func (em *Emitter) readLastEmittedEventID() *ltypes.EventHash {
	if em.emittedEventFile == nil {
		return nil
	}
	buf := make([]byte, 32)
	_, err := em.emittedEventFile.ReadAt(buf, 0)
	if err != nil {
		if err == io.EOF {
			return nil
		}
		log.Crit("Failed to read event file", "file", em.config.PrevEmittedEventFile.Path, "err", err)
	}
	v := ltypes.BytesToEvent(buf)
	return &v
}
