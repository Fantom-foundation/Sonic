package evmcore

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/Fantom-foundation/go-opera/inter/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
	"go.uber.org/mock/gomock"
)

func TestApplyTransaction_InternalTransactionsSkipBaseFeeCharges(t *testing.T) {
	for _, internal := range []bool{true, false} {
		t.Run("internal="+fmt.Sprint(internal), func(t *testing.T) {
			ctxt := gomock.NewController(t)
			state := state.NewMockStateDB(ctxt)

			any := gomock.Any()
			state.EXPECT().GetBalance(any).Return(uint256.NewInt(0))
			state.EXPECT().SubBalance(any, any, any)
			if !internal {
				state.EXPECT().GetNonce(any)
				state.EXPECT().GetCodeHash(any)
			}

			evm := vm.NewEVM(vm.BlockContext{}, vm.TxContext{}, state, &params.ChainConfig{}, vm.Config{})
			gp := new(core.GasPool).AddGas(1000000)

			applyTransaction(&core.Message{
				SkipAccountChecks: internal,
				GasPrice:          big.NewInt(0),
				Value:             big.NewInt(0),
			}, nil, gp, state, nil, common.Hash{}, nil, nil, evm, nil)

			if want, got := internal, evm.Config.NoBaseFee; want != got {
				t.Fatalf("want %v, got %v", want, got)
			}
		})
	}
}
