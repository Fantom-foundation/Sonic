package emitter

import (
	"fmt"
	"time"

	"github.com/Fantom-foundation/lachesis-base/common/bigendian"
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/Fantom-foundation/lachesis-base/inter/pos"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"

	"github.com/Fantom-foundation/go-opera/eventcheck/epochcheck"
	"github.com/Fantom-foundation/go-opera/inter"
	"github.com/Fantom-foundation/go-opera/utils"
	"github.com/Fantom-foundation/go-opera/utils/txtime"
)

const (
	TxTurnPeriod        = 8 * time.Second
	TxTurnPeriodLatency = 1 * time.Second
	TxTurnNonces        = 32
)

func max64(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func (em *Emitter) maxGasPowerToUse(e *inter.MutableEventPayload) uint64 {
	rules := em.world.GetRules()
	maxGasToUse := rules.Economy.Gas.MaxEventGas
	if maxGasToUse > e.GasPowerLeft().Min() {
		maxGasToUse = e.GasPowerLeft().Min()
	}
	//fmt.Printf("\t\tMax Event Gas: %d, gasPowerLeft: %d, Max Gas: %d\n", rules.Economy.Gas.MaxEventGas, e.GasPowerLeft().Min(), maxGasToUse)
	/*
		// Smooth TPS if power isn't big
		if em.config.LimitedTpsThreshold > em.config.NoTxsThreshold {
			upperThreshold := em.config.LimitedTpsThreshold
			downThreshold := em.config.NoTxsThreshold

			estimatedAlloc := gaspowercheck.CalcValidatorGasPower(e, e.CreationTime(), e.MedianTime(), 0, em.validators, gaspowercheck.Config{
				Idx:                inter.LongTermGas,
				AllocPerSec:        rules.Economy.LongGasPower.AllocPerSec * 4 / 5,
				MaxAllocPeriod:     inter.Timestamp(time.Minute),
				MinEnsuredAlloc:    0,
				StartupAllocPeriod: 0,
				MinStartupGas:      0,
			})

			gasPowerLeft := e.GasPowerLeft().Min() + estimatedAlloc
			if gasPowerLeft < downThreshold {
				return 0
			}
			newGasPowerLeft := uint64(0)
			if gasPowerLeft > maxGasToUse {
				newGasPowerLeft = gasPowerLeft - maxGasToUse
			}

			var x1, x2 = newGasPowerLeft, gasPowerLeft
			if x1 < downThreshold {
				x1 = downThreshold
			}
			if x2 > upperThreshold {
				x2 = upperThreshold
			}
			trespassingPart := uint64(0)
			if x2 > x1 {
				trespassingPart = x2 - x1
			}
			healthyPart := uint64(0)
			if gasPowerLeft > x2 {
				healthyPart = gasPowerLeft - x2
			}

			smoothGasToUse := healthyPart + trespassingPart/2
			if maxGasToUse > smoothGasToUse {
				fmt.Printf("\t\tSmooth TPS: %d -> %d\n", maxGasToUse, smoothGasToUse)
				maxGasToUse = smoothGasToUse
			}
		}
	*/
	// pendingGas should be below MaxBlockGas
	/*
		{
			maxPendingGas := max64(max64(rules.Blocks.MaxBlockGas/3, rules.Economy.Gas.MaxEventGas), 15000000)
			if maxPendingGas <= em.pendingGas {
				return 0
			}
			if maxPendingGas < em.pendingGas+maxGasToUse {
				maxGasToUse = maxPendingGas - em.pendingGas
				fmt.Printf("\t\tMax Gas reduced by pendingGas: %d - %d = %d\n", maxPendingGas, em.pendingGas, maxGasToUse)
			}
		}
	*/
	// No txs if power is low
	/*
		{
			threshold := em.config.NoTxsThreshold
			if e.GasPowerLeft().Min() <= threshold {
				return 0
			} else if e.GasPowerLeft().Min() < threshold+maxGasToUse {
				maxGasToUse = e.GasPowerLeft().Min() - threshold
				fmt.Printf("\t\tMax Gas reduced by threshold: %d - %d = %d\n", e.GasPowerLeft().Min(), threshold, maxGasToUse)
			}
		}
	*/
	return maxGasToUse
}

func getTxRoundIndex(now, txTime time.Time, validatorsNum idx.Validator) int {
	passed := now.Sub(txTime)
	if passed < 0 {
		passed = 0
	}
	return int((passed / TxTurnPeriod) % time.Duration(validatorsNum))
}

// safe for concurrent use
func (em *Emitter) isMyTxTurn(txHash common.Hash, sender common.Address, accountNonce uint64, now time.Time, validators *pos.Validators, me idx.ValidatorID, epoch idx.Epoch) bool {
	txTime := txtime.Of(txHash)

	roundIndex := getTxRoundIndex(now, txTime, validators.Len())
	if roundIndex != getTxRoundIndex(now.Add(TxTurnPeriodLatency), txTime, validators.Len()) {
		// round is about to change, avoid originating the transaction to avoid racing with another validator
		return false
	}

	// generate seed for generating the validators sequence for the tx
	roundsHash := hash.Of(sender.Bytes(), bigendian.Uint64ToBytes(accountNonce/TxTurnNonces), epoch.Bytes())

	// generate the validators sequence for the tx
	rounds := utils.WeightedPermutation(int(validators.Len()), validators.SortedWeights(), roundsHash)

	// take a validator from the sequence, skip offline validators
	for ; roundIndex < len(rounds); roundIndex++ {
		chosenValidator := validators.GetID(idx.Validator(rounds[roundIndex]))
		if chosenValidator == me {
			return true // current validator is the chosen - emit
		}
		if !em.offlineValidators[chosenValidator] {
			return false // chosen validator is online - don't emit
		}
		// otherwise try next validator in the sequence
		skippedOfflineValidatorsCounter.Inc(1)
	}
	return false
}

func (em *Emitter) addTxs(e *inter.MutableEventPayload, sorted *transactionsByPriceAndNonce) {
	const debug = false

	maxGasUsed := em.maxGasPowerToUse(e)
	if maxGasUsed <= e.GasPowerUsed() {
		return
	}

	// sort transactions by price and nonce
	rules := em.world.GetRules()
	if debug {
		fmt.Printf("Adding transactions to event - available gas power %d - maxGasUsed %d\n", e.GasPowerLeft().Min(), maxGasUsed)
	}
	for tx, _ := sorted.Peek(); tx != nil; tx, _ = sorted.Peek() {
		resolvedTx := tx.Resolve()
		sender, _ := types.Sender(em.world.TxSigner, resolvedTx)
		// check transaction epoch rules (tx type, gas price)
		if epochcheck.CheckTxs(types.Transactions{resolvedTx}, rules) != nil {
			if debug {
				fmt.Printf("\tskipped by epoch rules\n")
			}
			txsSkippedEpochRules.Inc(1)
			sorted.Pop()
			continue
		}
		// check there's enough gas power to originate the transaction
		if tx.Gas >= e.GasPowerLeft().Min() || e.GasPowerUsed()+tx.Gas >= maxGasUsed {
			txsSkippedNoValidatorGas.Inc(1)
			if params.TxGas >= e.GasPowerLeft().Min() || e.GasPowerUsed()+params.TxGas >= maxGasUsed {
				// stop if cannot originate even an empty transaction
				if debug {
					fmt.Printf("\taborting addition of transactions - out of gas\n")
				}
				break
			}
			if debug {
				fmt.Printf("\tskipping transaction - not enough gas - %d >= %d || %d+%d >= %d\n",
					tx.Gas, e.GasPowerLeft().Min(), e.GasPowerUsed(), tx.Gas, maxGasUsed,
				)
			}
			sorted.Pop()
			continue
		}
		// check not conflicted with already originated txs (in any connected event)
		if em.originatedTxs.TotalOf(sender) != 0 {
			if debug {
				fmt.Printf("\tskipped due to sender conflict\n")
			}
			txsSkippedConflictingSender.Inc(1)
			sorted.Pop()
			continue
		}
		// my turn, i.e. try to not include the same tx simultaneously by different validators
		if !em.isMyTxTurn(tx.Hash, sender, resolvedTx.Nonce(), time.Now(), em.validators, e.Creator(), em.epoch) {
			if debug {
				fmt.Printf("\tskipped due to turn conflict\n")
			}
			txsSkippedNotMyTurn.Inc(1)
			sorted.Pop()
			continue
		}
		// check transaction is not outdated
		if !em.world.TxPool.Has(tx.Hash) {
			if debug {
				fmt.Printf("\tskipped due to out-dated transaction\n")
			}
			txsSkippedOutdated.Inc(1)
			sorted.Pop()
			continue
		}
		// add
		if debug {
			fmt.Printf("\ttransaction passed all tests => included\n")
		}
		e.SetGasPowerUsed(e.GasPowerUsed() + tx.Gas)
		e.SetGasPowerLeft(e.GasPowerLeft().Sub(tx.Gas))
		e.SetTxs(append(e.Txs(), resolvedTx))
		sorted.Shift()
	}
	if debug {
		fmt.Printf("add transactions completed, total number of transactions: %d, gas power left: %d\n", len(e.Txs()), e.GasPowerLeft().Min())
	}
}
