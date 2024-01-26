package gossip

import (
	"github.com/ethereum/go-ethereum/metrics"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

var syncStatusGauge = metrics.GetOrRegisterGauge("chain/syncStage", nil)
var isMaybeSyncedGauge = metrics.GetOrRegisterGauge("chain/maybeSynced", nil)

type syncStage uint32

type syncStatus struct {
	stage       uint32
	maybeSynced uint32
}

const (
	ssUnknown syncStage = iota
	ssSnaps
	ssEvmSnapGen
	ssEvents
)

const (
	snapsyncMinEndAge   = 14 * 24 * time.Hour
	snapsyncMaxStartAge = 6 * time.Hour
)

func (ss *syncStatus) Is(s ...syncStage) bool {
	self := &ss.stage
	for _, v := range s {
		if atomic.LoadUint32(self) == uint32(v) {
			return true
		}
	}
	return false
}

func (ss *syncStatus) Set(s syncStage) {
	atomic.StoreUint32(&ss.stage, uint32(s))
	syncStatusGauge.Update(int64(s))
}

func (ss *syncStatus) MaybeSynced() bool {
	return atomic.LoadUint32(&ss.maybeSynced) != 0
}

func (ss *syncStatus) MarkMaybeSynced() {
	atomic.StoreUint32(&ss.maybeSynced, uint32(1))
	isMaybeSyncedGauge.Update(int64(1))
}

func (ss *syncStatus) AcceptEvents() bool {
	return ss.Is(ssEvents)
}

func (ss *syncStatus) AcceptBlockRecords() bool {
	return !ss.Is(ssEvents)
}

func (ss *syncStatus) AcceptTxs() bool {
	return ss.MaybeSynced() && ss.Is(ssEvents)
}

func (ss *syncStatus) RequestLLR() bool {
	return !ss.Is(ssEvents) || ss.MaybeSynced()
}

type txsync struct {
	p     *peer
	txids []common.Hash
}

// syncTransactions starts sending all currently pending transactions to the given peer.
func (h *handler) syncTransactions(p *peer, txids []common.Hash) {
	if len(txids) == 0 {
		return
	}
	select {
	case h.txsyncCh <- &txsync{p, txids}:
	case <-h.quitSync:
	}
}

// txsyncLoop takes care of the initial transaction sync for each new
// connection. When a new peer appears, we relay all currently pending
// transactions. In order to minimise egress bandwidth usage, we send
// the transactions in small packs to one peer at a time.
func (h *handler) txsyncLoop() {
	var (
		pending = make(map[enode.ID]*txsync)
		sending = false               // whether a send is active
		pack    = new(txsync)         // the pack that is being sent
		done    = make(chan error, 1) // result of the send
	)

	// send starts a sending a pack of transactions from the sync.
	send := func(s *txsync) {
		// Fill pack with transactions up to the target size.
		pack.p = s.p
		pack.txids = pack.txids[:0]
		for i := 0; i < len(s.txids) && len(pack.txids) < softLimitItems; i++ {
			pack.txids = append(pack.txids, s.txids[i])
		}
		// Remove the transactions that will be sent.
		s.txids = s.txids[len(pack.txids):]
		if len(s.txids) == 0 {
			delete(pending, s.p.ID())
		}
		// Send the pack in the background.
		s.p.Log().Trace("Sending batch of transaction hashes", "count", len(pack.txids))
		sending = true
		go func() {
			if len(pack.txids) != 0 {
				done <- pack.p.SendTransactionHashes(pack.txids)
			} else {
				done <- nil
			}
		}()
	}

	// pick chooses the next pending sync.
	pick := func() *txsync {
		if len(pending) == 0 {
			return nil
		}
		n := rand.Intn(len(pending)) + 1
		for _, s := range pending {
			if n--; n == 0 {
				return s
			}
		}
		return nil
	}

	for {
		select {
		case s := <-h.txsyncCh:
			pending[s.p.ID()] = s
			if !sending {
				send(s)
			}
		case err := <-done:
			sending = false
			// Stop tracking peers that cause send failures.
			if err != nil {
				pack.p.Log().Debug("Transaction send failed", "err", err)
				delete(pending, pack.p.ID())
			}
			// Schedule the next send.
			if s := pick(); s != nil {
				send(s)
			}
		case <-h.quitSync:
			return
		}
	}
}
