package tests

import (
	"math/big"
	"testing"

	accessCost "github.com/Fantom-foundation/go-opera/tests/contracts/access_cost"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestAddressAccess(t *testing.T) {
	someAccountAddress := common.Address{1}

	net, err := StartIntegrationTestNet(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to start the fake network: %v", err)
	}
	defer net.Stop()

	contract, receipt, err := DeployContract(net, accessCost.DeployAccessCost)
	checkTxExecution(t, receipt, err)

	// Execute function on an address, cold access
	receipt, err = net.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return contract.TouchAddress(opts, someAccountAddress)
	})
	checkTxExecution(t, receipt, err)
	txColdAccess, err := contract.ParseLogCost(*receipt.Logs[0])
	if err != nil {
		t.Fatalf("Failed to parse log: %v", err)
	}
	_, viewColdAccess, err := contract.GetAddressAccessCost(nil, someAccountAddress)
	if err != nil {
		t.Fatalf("Failed to get address access cost: %v", err)
	}

	t.Run("coinbase yields zero address", func(t *testing.T) {
		coinBaseAddress, err := contract.GetCoinBaseAddress(nil)
		if err != nil {
			t.Fatalf("Failed to get coinbase address: %v", err)
		}

		if want, got := (common.Address{}), coinBaseAddress; want != got {
			t.Errorf("Expected coinbase address %v, got %v", want, got)
		}
	})

	t.Run("tx access is warm", func(t *testing.T) {
		tests := map[string]func(*bind.TransactOpts) (*types.Transaction, error){
			"coinbase": contract.TouchCoinBase,
			"origin":   contract.TouchOrigin,
			"access list": func(ops *bind.TransactOpts) (*types.Transaction, error) {
				ops.GasPrice = nil // < transactions with gas price cannot have access list
				ops.GasFeeCap = big.NewInt(1e12)
				ops.GasTipCap = big.NewInt(1000)
				ops.AccessList = types.AccessList{
					{Address: someAccountAddress},
				}
				return contract.TouchAddress(ops, someAccountAddress)
			},
		}

		for name, access := range tests {
			t.Run(name, func(t *testing.T) {
				receipt, err = net.Apply(access)
				checkTxExecution(t, receipt, err)
				warmAccess, err := contract.ParseLogCost(*receipt.Logs[0])
				if err != nil {
					t.Fatalf("Failed to parse log: %v", err)
				}

				// Difference must be the extra cost of a cold access
				diff := new(big.Int).Sub(txColdAccess.Cost, warmAccess.Cost)
				if want, got := big.NewInt(2500), diff; want.Cmp(got) != 0 {
					t.Errorf("Expected cost difference %v, got %v", want, got)
				}
			})
		}
	})

	t.Run("archive access is warm", func(t *testing.T) {

		tests := map[string]func() (*big.Int, error){
			"origin": func() (*big.Int, error) {
				originAddr, err := contract.GetOrigin(nil)
				if err != nil {
					return nil, err
				}
				_, cost, err := contract.GetAddressAccessCost(nil, originAddr)
				return cost, err
			},
			"coinbase": func() (*big.Int, error) {
				coinbaseAddr, err := contract.GetCoinBaseAddress(nil)
				if err != nil {
					return nil, err
				}
				_, cost, err := contract.GetAddressAccessCost(nil, coinbaseAddr)
				return cost, err
			},
		}

		for name, access := range tests {
			t.Run(name, func(t *testing.T) {

				cost, err := access()
				if err != nil {
					t.Fatalf("Failed to get address access cost: %v", err)
				}
				diff := new(big.Int).Sub(viewColdAccess, cost)
				if want, got := big.NewInt(2500), diff; want.Cmp(got) != 0 {
					t.Errorf("Expected cost difference %v, got %v", want, got)
				}
			})
		}
	})
}

////////////////////////////////////////////////////////////////////////////////
// helpers

func checkTxExecution(t *testing.T, receipt *types.Receipt, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Failed to execute transaction; %v", err)
	}
	if want, got := types.ReceiptStatusSuccessful, receipt.Status; want != got {
		t.Errorf("Expected status %v, got %v", want, got)
	}
}
