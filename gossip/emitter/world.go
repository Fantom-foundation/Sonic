package emitter

import (
	"errors"
	"sync"

	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/Fantom-foundation/lachesis-base/ltypes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/Fantom-foundation/go-opera/inter"
	"github.com/Fantom-foundation/go-opera/inter/state"
	"github.com/Fantom-foundation/go-opera/opera"
	"github.com/Fantom-foundation/go-opera/valkeystore"
	"github.com/Fantom-foundation/go-opera/vecmt"
)

var (
	ErrNotEnoughGasPower = errors.New("not enough gas power")
)

type (
	// External world
	External interface {
		sync.Locker
		Reader

		Check(e *inter.EventPayload, parents inter.Events) error
		Process(*inter.EventPayload) error
		Broadcast(*inter.EventPayload)
		Build(*inter.MutableEventPayload, func()) error
		DagIndex() *vecmt.Index

		IsBusy() bool
		IsSynced() bool
		PeersNum() int

		StateDB() state.StateDB
	}

	// aliases for mock generator
	Signer   valkeystore.SignerI
	TxSigner types.Signer

	// World is an emitter's environment
	World struct {
		External
		TxPool   TxPool
		Signer   valkeystore.SignerI
		TxSigner types.Signer
	}
)

// Reader is a callback for getting events from an external storage.
type Reader interface {
	GetLatestBlockIndex() idx.BlockID
	GetEpochValidators() (*ltypes.Validators, idx.EpochID)
	GetEvent(hash.EventHash) *inter.Event
	GetEventPayload(hash.EventHash) *inter.EventPayload
	GetLastEvent(epoch idx.EpochID, from idx.ValidatorID) *hash.EventHash
	GetHeads(idx.EpochID) hash.EventHashes
	GetGenesisTime() inter.Timestamp
	GetRules() opera.Rules
}

type TxPool interface {
	// Has returns an indicator whether txpool has a transaction cached with the
	// given hash.
	Has(hash common.Hash) bool
	// Pending should return pending transactions.
	// The slice should be modifiable by the caller.
	Pending(enforceTips bool) (map[common.Address]types.Transactions, error)

	// Count returns the total number of transactions
	Count() int
}
