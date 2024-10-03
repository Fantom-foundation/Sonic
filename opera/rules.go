package opera

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	ethparams "github.com/ethereum/go-ethereum/params"

	"github.com/Fantom-foundation/Tosca/go/geth_adapter"
	"github.com/Fantom-foundation/Tosca/go/interpreter/lfvm"
	"github.com/Fantom-foundation/go-opera/inter"
	"github.com/Fantom-foundation/go-opera/opera/contracts/evmwriter"
)

const (
	MainNetworkID   uint64 = 0xfa
	TestNetworkID   uint64 = 0xfa2
	FakeNetworkID   uint64 = 0xfa3
	DefaultEventGas uint64 = 28000
	berlinBit              = 1 << 0
	londonBit              = 1 << 1
	llrBit                 = 1 << 2
	sonicBit               = 1 << 3
)

var DefaultVMConfig = func() vm.Config {

	// For transaction processing, Tosca's LFVM is used.
	interpreter, err := lfvm.NewInterpreter(lfvm.Config{})
	if err != nil {
		panic(err)
	}
	lfvmFactory := geth_adapter.NewGethInterpreterFactory(interpreter)

	// For tracing, Geth's EVM is used.
	gethFactory := func(evm *vm.EVM) vm.Interpreter {
		return vm.NewEVMInterpreter(evm)
	}

	return vm.Config{
		StatePrecompiles: map[common.Address]vm.PrecompiledStateContract{
			evmwriter.ContractAddress: &evmwriter.PreCompiledContract{},
		},
		Interpreter:           lfvmFactory,
		InterpreterForTracing: gethFactory,

		// Fantom/Sonic modifications
		ChargeExcessGas:                 true,
		IgnoreGasFeeCap:                 true,
		InsufficientBalanceIsNotAnError: true,
		SkipTipPaymentToCoinbase:        true,
	}
}()

type RulesRLP struct {
	Name      string
	NetworkID uint64

	// Graph options
	Dag DagRules

	// Epochs options
	Epochs EpochsRules

	// Blockchain options
	Blocks BlocksRules

	// Economy options
	Economy EconomyRules

	Upgrades Upgrades `rlp:"-"`
}

// Rules describes opera net.
// Note keep track of all the non-copiable variables in Copy()
type Rules RulesRLP

// GasPowerRules defines gas power rules in the consensus.
type GasPowerRules struct {
	AllocPerSec        uint64
	MaxAllocPeriod     inter.Timestamp
	StartupAllocPeriod inter.Timestamp
	MinStartupGas      uint64
}

type GasRulesRLPV1 struct {
	MaxEventGas  uint64
	EventGas     uint64
	ParentGas    uint64
	ExtraDataGas uint64
	// Post-LLR fields
	BlockVotesBaseGas    uint64
	BlockVoteGas         uint64
	EpochVoteGas         uint64
	MisbehaviourProofGas uint64
}

type GasRules GasRulesRLPV1

type EpochsRules struct {
	MaxEpochGas      uint64
	MaxEpochDuration inter.Timestamp
}

// DagRules of Lachesis DAG (directed acyclic graph).
type DagRules struct {
	MaxParents     idx.Event
	MaxFreeParents idx.Event // maximum number of parents with no gas cost
	MaxExtraData   uint32
}

// BlocksMissed is information about missed blocks from a staker
type BlocksMissed struct {
	BlocksNum idx.Block
	Period    inter.Timestamp
}

// EconomyRules contains economy constants
type EconomyRules struct {
	BlockMissedSlack idx.Block

	Gas GasRules

	MinGasPrice *big.Int

	ShortGasPower GasPowerRules
	LongGasPower  GasPowerRules
}

// BlocksRules contains blocks constants
type BlocksRules struct {
	MaxBlockGas             uint64 // technical hard limit, gas is mostly governed by gas power allocation
	MaxEmptyBlockSkipPeriod inter.Timestamp
}

type Upgrades struct {
	Berlin bool
	London bool
	Llr    bool
	Sonic  bool
}

type UpgradeHeight struct {
	Upgrades Upgrades
	Height   idx.Block
}

var BaseChainConfig = ethparams.ChainConfig{
	ChainID:                       big.NewInt(1337),
	HomesteadBlock:                big.NewInt(0),
	DAOForkBlock:                  nil,
	DAOForkSupport:                false,
	EIP150Block:                   big.NewInt(0),
	EIP155Block:                   big.NewInt(0),
	EIP158Block:                   big.NewInt(0),
	ByzantiumBlock:                big.NewInt(0),
	ConstantinopleBlock:           big.NewInt(0),
	PetersburgBlock:               big.NewInt(0),
	IstanbulBlock:                 big.NewInt(0),
	MuirGlacierBlock:              big.NewInt(0), // EIP-2384: Muir Glacier Difficulty Bomb Delay
	BerlinBlock:                   nil, // to be overwritten in EvmChainConfig
	LondonBlock:                   nil, // to be overwritten in EvmChainConfig
	ArrowGlacierBlock:             nil, // EIP-4345: Difficulty Bomb Delay
	GrayGlacierBlock:              nil, // EIP-5133: Delaying Difficulty Bomb
	MergeNetsplitBlock:            nil,
	ShanghaiTime:                  nil,
	CancunTime:                    nil,
	PragueTime:                    nil,
	VerkleTime:                    nil,
	TerminalTotalDifficulty:       nil,
	TerminalTotalDifficultyPassed: true,
	Ethash:                        new(ethparams.EthashConfig),
	Clique:                        nil,
}

// EvmChainConfig returns ChainConfig for transactions signing and execution
func (r Rules) EvmChainConfig(hh []UpgradeHeight) *ethparams.ChainConfig {
	cfg := BaseChainConfig
	zero := new(uint64)
	cfg.ChainID = new(big.Int).SetUint64(r.NetworkID)
	for i, h := range hh {
		height := new(big.Int)
		if i > 0 {
			height.SetUint64(uint64(h.Height))
		}
		if cfg.BerlinBlock == nil && h.Upgrades.Berlin {
			cfg.BerlinBlock = height
		}
		if !h.Upgrades.Berlin {
			// disabling upgrade breaks the history replay - should be never used
			cfg.BerlinBlock = nil
		}

		if cfg.LondonBlock == nil && h.Upgrades.London {
			cfg.LondonBlock = height
		}
		if !h.Upgrades.London {
			// disabling upgrade breaks the history replay - should be never used
			cfg.LondonBlock = nil
		}

		if cfg.CancunTime == nil && h.Upgrades.Sonic {
			cfg.ArrowGlacierBlock = height
			cfg.GrayGlacierBlock = height
			// enabling with no specified block/time breaks the history replay - should be used only at start of a chain
			cfg.ShanghaiTime = zero
			cfg.CancunTime = zero
		}
		if !h.Upgrades.Sonic {
			// disabling upgrade breaks the history replay - should be never used
			cfg.ArrowGlacierBlock = nil
			cfg.GrayGlacierBlock = nil
			cfg.ShanghaiTime = nil
			cfg.CancunTime = nil
		}
	}
	return &cfg
}

func MainNetRules() Rules {
	return Rules{
		Name:      "main",
		NetworkID: MainNetworkID,
		Dag:       DefaultDagRules(),
		Epochs:    DefaultEpochsRules(),
		Economy:   DefaultEconomyRules(),
		Blocks: BlocksRules{
			MaxBlockGas:             20500000,
			MaxEmptyBlockSkipPeriod: inter.Timestamp(1 * time.Minute),
		},
	}
}

func FakeNetRules() Rules {
	return Rules{
		Name:      "fake",
		NetworkID: FakeNetworkID,
		Dag:       DefaultDagRules(),
		Epochs:    FakeNetEpochsRules(),
		Economy:   FakeEconomyRules(),
		Blocks: BlocksRules{
			MaxBlockGas:             20500000,
			MaxEmptyBlockSkipPeriod: inter.Timestamp(3 * time.Second),
		},
		Upgrades: Upgrades{
			Berlin: true,
			London: true,
			Llr:    false,
			Sonic:  true,
		},
	}
}

// DefaultEconomyRules returns mainnet economy
func DefaultEconomyRules() EconomyRules {
	rules := EconomyRules{
		BlockMissedSlack: 50,
		Gas:              DefaultGasRules(),
		MinGasPrice:      big.NewInt(1e9),
		ShortGasPower:    DefaultShortGasPowerRules(),
		LongGasPower:     DefaulLongGasPowerRules(),
	}
	return rules
}

// FakeEconomyRules returns fakenet economy
func FakeEconomyRules() EconomyRules {
	cfg := DefaultEconomyRules()
	cfg.ShortGasPower = FakeShortGasPowerRules()
	cfg.LongGasPower = FakeLongGasPowerRules()
	return cfg
}

func DefaultDagRules() DagRules {
	return DagRules{
		MaxParents:     10,
		MaxFreeParents: 3,
		MaxExtraData:   128,
	}
}

func DefaultEpochsRules() EpochsRules {
	return EpochsRules{
		MaxEpochGas:      1500000000,
		MaxEpochDuration: inter.Timestamp(4 * time.Hour),
	}
}

func DefaultGasRules() GasRules {
	return GasRules{
		MaxEventGas:          10000000 + DefaultEventGas,
		EventGas:             DefaultEventGas,
		ParentGas:            2400,
		ExtraDataGas:         25,
		BlockVotesBaseGas:    1024,
		BlockVoteGas:         512,
		EpochVoteGas:         1536,
		MisbehaviourProofGas: 71536,
	}
}

func FakeNetEpochsRules() EpochsRules {
	cfg := DefaultEpochsRules()
	cfg.MaxEpochGas /= 5
	cfg.MaxEpochDuration = inter.Timestamp(10 * time.Minute)
	return cfg
}

// DefaulLongGasPowerRules is long-window config
func DefaulLongGasPowerRules() GasPowerRules {
	return GasPowerRules{
		AllocPerSec:        100 * DefaultEventGas,
		MaxAllocPeriod:     inter.Timestamp(60 * time.Minute),
		StartupAllocPeriod: inter.Timestamp(5 * time.Second),
		MinStartupGas:      DefaultEventGas * 20,
	}
}

// DefaultShortGasPowerRules is short-window config
func DefaultShortGasPowerRules() GasPowerRules {
	// 2x faster allocation rate, 6x lower max accumulated gas power
	cfg := DefaulLongGasPowerRules()
	cfg.AllocPerSec *= 2
	cfg.StartupAllocPeriod /= 2
	cfg.MaxAllocPeriod /= 2 * 6
	return cfg
}

// FakeLongGasPowerRules is fake long-window config
func FakeLongGasPowerRules() GasPowerRules {
	config := DefaulLongGasPowerRules()
	config.AllocPerSec *= 1000
	return config
}

// FakeShortGasPowerRules is fake short-window config
func FakeShortGasPowerRules() GasPowerRules {
	config := DefaultShortGasPowerRules()
	config.AllocPerSec *= 1000
	return config
}

func (r Rules) Copy() Rules {
	cp := r
	cp.Economy.MinGasPrice = new(big.Int).Set(r.Economy.MinGasPrice)
	return cp
}

func (r Rules) String() string {
	b, _ := json.Marshal(&r)
	return string(b)
}
