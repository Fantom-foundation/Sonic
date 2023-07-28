package makefakegenesis

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Fantom-foundation/go-opera/integration/makegenesis"
	"github.com/Fantom-foundation/go-opera/inter/drivertype"
	"github.com/Fantom-foundation/go-opera/inter/iblockproc"
	"github.com/Fantom-foundation/go-opera/inter/ier"
	"github.com/Fantom-foundation/go-opera/opera"
	"github.com/Fantom-foundation/go-opera/opera/genesis"
	"github.com/Fantom-foundation/go-opera/opera/genesisstore"
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/Fantom-foundation/lachesis-base/inter/pos"
	"github.com/Fantom-foundation/lachesis-base/kvdb/memorydb"
	"github.com/Fantom-foundation/lachesis-base/lachesis"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"os"
)

type GenesisJson struct {
	Rules    NetworkRules
	Accounts []Account     `json:",omitempty"`
	Txs      []Transaction `json:",omitempty"`
}

type NetworkRules struct {
	NetworkName         string
	NetworkID           hexutil.Uint64
	MaxBlockGas         *uint64 `json:",omitempty"`
	MaxEpochGas         *uint64 `json:",omitempty"`
	MaxEventGas         *uint64 `json:",omitempty"`
	LongGasAllocPerSec  *uint64 `json:",omitempty"`
	ShortGasAllocPerSec *uint64 `json:",omitempty"`
	OverrideMinGasPrice *uint64 `json:",omitempty"`
}

type Account struct {
	Name    string
	Address common.Address
	Balance *big.Int                    `json:",omitempty"`
	Code    VariableLenCode             `json:",omitempty"`
	Storage map[common.Hash]common.Hash `json:",omitempty"`
}

type Transaction struct {
	Name string
	To   common.Address
	Data VariableLenCode `json:",omitempty"`
}

func LoadGenesisJson(filename string) (*GenesisJson, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read genesis json file; %v", err)
	}
	var decoded GenesisJson
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal genesis json file; %v", err)
	}
	return &decoded, nil
}

func ApplyGenesisJson(json *GenesisJson) (*genesisstore.Store, error) {
	builder := makegenesis.NewGenesisBuilder(memorydb.NewProducer(""))

	for _, acc := range json.Accounts {
		if acc.Balance != nil {
			builder.AddBalance(acc.Address, acc.Balance)
		}
		if acc.Code != nil {
			builder.SetCode(acc.Address, acc.Code)
		}
		if acc.Storage != nil {
			for key, val := range acc.Storage {
				builder.SetStorage(acc.Address, key, val)
			}
		}
	}

	rules := opera.FakeNetRules()
	rules.Name = json.Rules.NetworkName
	rules.NetworkID = uint64(json.Rules.NetworkID)
	if json.Rules.MaxBlockGas != nil {
		rules.Blocks.MaxBlockGas = *json.Rules.MaxBlockGas
	}
	if json.Rules.MaxEventGas != nil {
		rules.Economy.Gas.MaxEventGas = *json.Rules.MaxEventGas
	}
	if json.Rules.MaxEpochGas != nil {
		rules.Epochs.MaxEpochGas = *json.Rules.MaxEpochGas
	}
	if json.Rules.ShortGasAllocPerSec != nil {
		rules.Economy.ShortGasPower.AllocPerSec = *json.Rules.ShortGasAllocPerSec
	}
	if json.Rules.LongGasAllocPerSec != nil {
		rules.Economy.LongGasPower.AllocPerSec = *json.Rules.LongGasAllocPerSec
	}
	if json.Rules.OverrideMinGasPrice != nil {
		rules.Economy.OverrideMinGasPrice = big.NewInt(int64(*json.Rules.OverrideMinGasPrice))
	}

	builder.SetCurrentEpoch(ier.LlrIdxFullEpochRecord{
		LlrFullEpochRecord: ier.LlrFullEpochRecord{
			BlockState: iblockproc.BlockState{
				LastBlock: iblockproc.BlockCtx{
					Idx:     0,
					Time:    FakeGenesisTime,
					Atropos: hash.Event{},
				},
				FinalizedStateRoot:    hash.Hash{},
				EpochGas:              0,
				EpochCheaters:         lachesis.Cheaters{},
				CheatersWritten:       0,
				ValidatorStates:       make([]iblockproc.ValidatorBlockState, 0),
				NextValidatorProfiles: make(map[idx.ValidatorID]drivertype.Validator),
				DirtyRules:            nil,
				AdvanceEpochs:         0,
			},
			EpochState: iblockproc.EpochState{
				Epoch:             1,
				EpochStart:        FakeGenesisTime,
				PrevEpochStart:    FakeGenesisTime - 1,
				EpochStateRoot:    hash.Zero,
				Validators:        pos.NewBuilder().Build(),
				ValidatorStates:   make([]iblockproc.ValidatorEpochState, 0),
				ValidatorProfiles: make(map[idx.ValidatorID]drivertype.Validator),
				Rules:             rules,
			},
		},
		Idx: 1,
	})

	blockProc := makegenesis.DefaultBlockProc()
	buildTx := txBuilder()
	genesisTxs := make(types.Transactions, 0, len(json.Txs))
	for _, tx := range json.Txs {
		genesisTxs = append(genesisTxs, buildTx(tx.Data, tx.To))
	}
	err := builder.ExecuteGenesisTxs(blockProc, genesisTxs)
	if err != nil {
		return nil, fmt.Errorf("failed to execute json genesis txs; %v", err)
	}

	return builder.Build(genesis.Header{
		GenesisID:   builder.CurrentHash(),
		NetworkID:   uint64(json.Rules.NetworkID),
		NetworkName: json.Rules.NetworkName,
	}), nil
}

type VariableLenCode []byte

func (c *VariableLenCode) MarshalJSON() ([]byte, error) {
	out := make([]byte, hex.EncodedLen(len(*c))+4)
	out[0], out[1], out[2] = '"', '0', 'x'
	hex.Encode(out[3:], *c)
	out[len(*c)-1] = '"'
	return out, nil
}

func (c *VariableLenCode) UnmarshalJSON(data []byte) error {
	if !bytes.HasPrefix(data, []byte(`"`)) || !bytes.HasSuffix(data, []byte(`"`)) {
		return fmt.Errorf("code must be in a string")
	}
	data = bytes.Trim(data, "\"")
	data = bytes.TrimPrefix(data, []byte("0x"))
	decoded := make([]byte, hex.DecodedLen(len(data)))
	_, err := hex.Decode(decoded, data)
	if err != nil {
		return err
	}
	*c = decoded
	return nil
}
