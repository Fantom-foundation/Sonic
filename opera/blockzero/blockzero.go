package blockzero

import (
	"github.com/Fantom-foundation/go-opera/evmcore"
	"github.com/Fantom-foundation/go-opera/gossip/gasprice"
	"github.com/Fantom-foundation/go-opera/inter"
	"github.com/Fantom-foundation/go-opera/opera"
	"github.com/ethereum/go-ethereum/common"
)

func GetBlockZero(rules opera.Rules) *inter.Block {
	gasLimit := rules.Blocks.MaxBlockGas
	return inter.NewBlockBuilder().
		WithNumber(0).
		WithTime(evmcore.FakeGenesisTime - 1). // TODO: extend genesis generator to provide time
		WithGasLimit(gasLimit).
		WithStateRoot(common.Hash{}). // TODO: get proper state root from genesis data
		WithBaseFee(gasprice.GetInitialBaseFee(rules.Economy)).
		Build()
}
