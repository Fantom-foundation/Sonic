package iblockproc

import (
	"crypto/sha256"
	"math/big"

	"github.com/Fantom-foundation/lachesis-base/lachesis"
	"github.com/Fantom-foundation/lachesis-base/ltypes"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/Fantom-foundation/go-opera/inter"
	"github.com/Fantom-foundation/go-opera/opera"
)

type ValidatorBlockState struct {
	LastEvent        EventInfo
	Uptime           inter.Timestamp
	LastOnlineTime   inter.Timestamp
	LastGasPowerLeft inter.GasPowerLeft
	LastBlock        ltypes.BlockID
	DirtyGasRefund   uint64
	Originated       *big.Int
}

type EventInfo struct {
	ID           ltypes.EventHash
	GasPowerLeft inter.GasPowerLeft
	Time         inter.Timestamp
}

type ValidatorEpochState struct {
	GasRefund      uint64
	PrevEpochEvent EventInfo
}

type BlockCtx struct {
	Idx     ltypes.BlockID
	Time    inter.Timestamp
	Atropos ltypes.EventHash
}

type BlockState struct {
	LastBlock          BlockCtx
	FinalizedStateRoot ltypes.Hash

	EpochGas        uint64
	EpochCheaters   lachesis.Cheaters
	CheatersWritten uint32

	ValidatorStates       []ValidatorBlockState
	NextValidatorProfiles ValidatorProfiles

	DirtyRules *opera.Rules `rlp:"nil"` // nil means that there's no changes compared to epoch rules

	AdvanceEpochs ltypes.EpochID
}

func (bs BlockState) Copy() BlockState {
	cp := bs
	cp.EpochCheaters = make(lachesis.Cheaters, len(bs.EpochCheaters))
	copy(cp.EpochCheaters, bs.EpochCheaters)
	cp.ValidatorStates = make([]ValidatorBlockState, len(bs.ValidatorStates))
	copy(cp.ValidatorStates, bs.ValidatorStates)
	for i := range cp.ValidatorStates {
		cp.ValidatorStates[i].Originated = new(big.Int).Set(cp.ValidatorStates[i].Originated)
	}
	cp.NextValidatorProfiles = bs.NextValidatorProfiles.Copy()
	if bs.DirtyRules != nil {
		rules := bs.DirtyRules.Copy()
		cp.DirtyRules = &rules
	}
	return cp
}

func (bs *BlockState) GetValidatorState(id ltypes.ValidatorID, validators *ltypes.Validators) *ValidatorBlockState {
	validatorIdx := validators.GetIdx(id)
	return &bs.ValidatorStates[validatorIdx]
}

func (bs BlockState) Hash() ltypes.Hash {
	hasher := sha256.New()
	err := rlp.Encode(hasher, &bs)
	if err != nil {
		panic("can't hash: " + err.Error())
	}
	return ltypes.BytesToHash(hasher.Sum(nil))
}

type EpochStateV1 struct {
	Epoch          ltypes.EpochID
	EpochStart     inter.Timestamp
	PrevEpochStart inter.Timestamp

	EpochStateRoot ltypes.Hash

	Validators        *ltypes.Validators
	ValidatorStates   []ValidatorEpochState
	ValidatorProfiles ValidatorProfiles

	Rules opera.Rules
}

type EpochState EpochStateV1

func (es *EpochState) GetValidatorState(id ltypes.ValidatorID, validators *ltypes.Validators) *ValidatorEpochState {
	validatorIdx := validators.GetIdx(id)
	return &es.ValidatorStates[validatorIdx]
}

func (es EpochState) Duration() inter.Timestamp {
	return es.EpochStart - es.PrevEpochStart
}

func (es EpochState) Hash() ltypes.Hash {
	var hashed interface{}
	if es.Rules.Upgrades.London {
		hashed = &es
	} else {
		es0 := EpochStateV0{
			Epoch:             es.Epoch,
			EpochStart:        es.EpochStart,
			PrevEpochStart:    es.PrevEpochStart,
			EpochStateRoot:    es.EpochStateRoot,
			Validators:        es.Validators,
			ValidatorStates:   make([]ValidatorEpochStateV0, len(es.ValidatorStates)),
			ValidatorProfiles: es.ValidatorProfiles,
			Rules:             es.Rules,
		}
		for i, v := range es.ValidatorStates {
			es0.ValidatorStates[i].GasRefund = v.GasRefund
			es0.ValidatorStates[i].PrevEpochEvent = v.PrevEpochEvent.ID
		}
		hashed = &es0
	}
	hasher := sha256.New()
	err := rlp.Encode(hasher, hashed)
	if err != nil {
		panic("can't hash: " + err.Error())
	}
	return ltypes.BytesToHash(hasher.Sum(nil))
}

func (es EpochState) Copy() EpochState {
	cp := es
	cp.ValidatorStates = make([]ValidatorEpochState, len(es.ValidatorStates))
	copy(cp.ValidatorStates, es.ValidatorStates)
	cp.ValidatorProfiles = es.ValidatorProfiles.Copy()
	if es.Rules != (opera.Rules{}) {
		cp.Rules = es.Rules.Copy()
	}
	return cp
}
