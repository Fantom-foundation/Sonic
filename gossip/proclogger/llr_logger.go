package proclogger

import (
	"time"

	"github.com/Fantom-foundation/lachesis-base/ltypes"

	"github.com/Fantom-foundation/go-opera/inter"
	"github.com/Fantom-foundation/go-opera/logger"
	"github.com/Fantom-foundation/go-opera/utils"
)

type dagSum struct {
	connected       ltypes.EventID
	totalProcessing time.Duration
}

type llrSum struct {
	bvs ltypes.BlockID
	brs ltypes.BlockID
	evs ltypes.EpochID
	ers ltypes.EpochID
}

type Logger struct {
	// summary accumulators
	dagSum dagSum
	llrSum llrSum

	// latest logged data
	lastEpoch     ltypes.EpochID
	lastBlock     ltypes.BlockID
	lastID        ltypes.EventHash
	lastEventTime inter.Timestamp
	lastLlrTime   inter.Timestamp

	nextLogging time.Time

	emitting  bool
	noSummary bool

	logger.Instance
}

func (l *Logger) summary(now time.Time) {
	if l.noSummary {
		return
	}
	if now.After(l.nextLogging) {
		if l.llrSum != (llrSum{}) {
			age := utils.PrettyDuration(now.Sub(l.lastLlrTime.Time())).String()
			if l.lastLlrTime <= l.lastEventTime {
				age = "none"
			}
			l.Log.Info("New LLR summary", "last_epoch", l.lastEpoch, "last_block", l.lastBlock,
				"new_evs", l.llrSum.evs, "new_ers", l.llrSum.ers, "new_bvs", l.llrSum.bvs, "new_brs", l.llrSum.brs, "age", age)
		}
		if l.dagSum != (dagSum{}) {
			l.Log.Info("New DAG summary", "new", l.dagSum.connected, "last_id", l.lastID.String(),
				"age", utils.PrettyDuration(now.Sub(l.lastEventTime.Time())), "t", utils.PrettyDuration(l.dagSum.totalProcessing))
		}
		l.dagSum = dagSum{}
		l.llrSum = llrSum{}
		l.nextLogging = now.Add(8 * time.Second)
	}
}
