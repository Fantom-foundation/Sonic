package gossip

import (
	"github.com/Fantom-foundation/lachesis-base/ltypes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"

	"github.com/Fantom-foundation/go-opera/eventcheck/gaspowercheck"
	"github.com/Fantom-foundation/go-opera/inter"
	"github.com/Fantom-foundation/go-opera/opera"
	"github.com/Fantom-foundation/go-opera/utils/concurrent"
)

type GPOBackend struct {
	store  *Store
	txpool TxPool
}

func (b *GPOBackend) GetLatestBlockIndex() ltypes.BlockID {
	return b.store.GetLatestBlockIndex()
}

func (b *GPOBackend) GetRules() opera.Rules {
	return b.store.GetRules()
}

func (b *GPOBackend) GetPendingRules() opera.Rules {
	bs, es := b.store.GetBlockEpochState()
	if bs.DirtyRules != nil {
		return *bs.DirtyRules
	}
	return es.Rules
}

func (b *GPOBackend) PendingTxs() map[common.Address]types.Transactions {
	txs, err := b.txpool.Pending(false)
	if err != nil {
		return map[common.Address]types.Transactions{}
	}
	return txs
}

func (b *GPOBackend) MinGasTip() *big.Int {
	return b.txpool.GasPrice()
}

// TotalGasPowerLeft returns a total amount of obtained gas power by the validators, according to the latest events from each validator
func (b *GPOBackend) TotalGasPowerLeft() uint64 {
	bs, es := b.store.GetBlockEpochState()
	set := b.store.GetLastEvents(es.Epoch)
	if set == nil {
		set = concurrent.WrapValidatorEventsSet(map[ltypes.ValidatorID]ltypes.EventHash{})
	}
	set.RLock()
	defer set.RUnlock()
	metValidators := map[ltypes.ValidatorID]bool{}
	total := uint64(0)
	gasPowerCheckCfg := gaspowercheck.Config{
		Idx:                inter.LongTermGas,
		AllocPerSec:        es.Rules.Economy.LongGasPower.AllocPerSec,
		MaxAllocPeriod:     es.Rules.Economy.LongGasPower.MaxAllocPeriod,
		MinEnsuredAlloc:    es.Rules.Economy.Gas.MaxEventGas,
		StartupAllocPeriod: es.Rules.Economy.LongGasPower.StartupAllocPeriod,
		MinStartupGas:      es.Rules.Economy.LongGasPower.MinStartupGas,
	}
	// count GasPowerLeft from latest events of this epoch
	for _, tip := range set.Val {
		e := b.store.GetEvent(tip)
		left := e.GasPowerLeft().Gas[inter.LongTermGas]
		left += bs.GetValidatorState(e.Creator(), es.Validators).DirtyGasRefund

		_, max, _ := gaspowercheck.CalcValidatorGasPowerPerSec(e.Creator(), es.Validators, gasPowerCheckCfg)
		if left > max {
			left = max
		}
		total += left

		metValidators[e.Creator()] = true
	}
	// count GasPowerLeft from last events of prev epoch if no event in current epoch is present
	for i := ltypes.ValidatorIdx(0); i < es.Validators.Len(); i++ {
		vid := es.Validators.GetID(i)
		if !metValidators[vid] {
			left := es.ValidatorStates[i].PrevEpochEvent.GasPowerLeft.Gas[inter.LongTermGas]
			left += es.ValidatorStates[i].GasRefund

			_, max, startup := gaspowercheck.CalcValidatorGasPowerPerSec(vid, es.Validators, gasPowerCheckCfg)
			if left > max {
				left = max
			}
			if left < startup {
				left = startup
			}
			total += left
		}
	}

	return total
}
