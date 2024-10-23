package tests

import (
	"testing"

	"github.com/Fantom-foundation/go-opera/tests/contracts/coinbase"
	"github.com/ethereum/go-ethereum/common"
)

func TestCoinBase_CoinbaseYieldsZero(t *testing.T) {

	net, err := StartIntegrationTestNet(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to start the fake network: %v", err)
	}
	defer net.Stop()

	contract, _, err := DeployContract(net, coinbase.DeployCoinbase)
	if err != nil {
		t.Fatalf("failed to deploy contract; %v", err)
	}

	for i := 0; i < 10; i++ {
		// execution out of transaction yields 0 address
		cbValue, err := contract.GetCoinbase(nil)
		if err != nil {
			t.Fatalf("failed to get coinbase value; %v", err)
		}
		if want, got := cbValue, (common.Address{}); want != got {
			t.Fatalf("unexpected coinbase value; expected %v, got %v", want, got)
		}

		// execution in transaction yields 0 address
		receipt, err := net.Apply(contract.LogCoinbase)
		if err != nil {
			t.Fatalf("failed to log coinbase; %v", err)
		}
		if len(receipt.Logs) != 1 {
			t.Fatalf("unexpected number of logs; expected 1, got %d", len(receipt.Logs))
		}
		if want, got := (common.Address{}), common.BytesToAddress(receipt.Logs[0].Data); want != got {
			t.Fatalf("unexpected coinbase address; expected %v, got %v", want, got)
		}
	}
}
