package tests

import (
	"context"
	"errors"
	"github.com/Fantom-foundation/go-opera/tests/contracts/prevrandao"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"math/big"
	"testing"
)

func TestPrevRandao(t *testing.T) {
	tests := []struct {
		name string
		f    func(net *IntegrationTestNet, contract *prevrandao.Prevrandao)
	}{
		{
			name: "is_correctly_computed",
			f: func(net *IntegrationTestNet, contract *prevrandao.Prevrandao) {
				// Collect the current PrevRandao from the head state.
				receipt, err := net.Apply(contract.LogCurrentPrevRandao)
				if err != nil {
					t.Fatalf("failed to log current prevrandao; %v", err)
				}

				if len(receipt.Logs) != 1 {
					t.Fatalf("unexpected number of logs; expected 1, got %d", len(receipt.Logs))
				}

				entry, err := contract.ParseCurrentPrevRandao(*receipt.Logs[0])
				if err != nil {
					t.Fatalf("failed to parse log; %v", err)
				}
				fromLog := entry.Prevrandao

				client, err := net.GetClient()
				if err != nil {
					t.Fatalf("failed to get client; %v", err)
				}
				defer client.Close()

				block, err := client.BlockByNumber(context.Background(), receipt.BlockNumber)
				if err != nil {
					t.Fatalf("failed to get block header; %v", err)
				}
				fromLatestBlock := block.MixDigest().Big() // MixDigest == MixHash == PrevRandao
				if block.Difficulty().Uint64() != 0 {
					t.Errorf("incorrect header difficulty got: %d, want: %d", block.Difficulty().Uint64(), 0)
				}

				// Collect the prevrandao from the archive.
				fromArchive, err := contract.GetPrevRandao(&bind.CallOpts{BlockNumber: receipt.BlockNumber})
				if err != nil {
					t.Fatalf("failed to get prevrandao from archive; %v", err)
				}

				if fromLog.Sign() < 1 {
					t.Fatalf("invalid prevrandao from log; %v", fromLog)
				}

				if fromLog.Cmp(fromLatestBlock) != 0 {
					t.Errorf("prevrandao mismatch; from log %v, from block %v", fromLog, fromLatestBlock)
				}
				if fromLog.Cmp(fromArchive) != 0 {
					t.Errorf("prevrandao mismatch; from log %v, from archive %v", fromLog, fromArchive)
				}
			},
		},
		{
			name: "is_different_for_each_block",
			f: func(net *IntegrationTestNet, contract *prevrandao.Prevrandao) {
				client, err := net.GetClient()
				if err != nil {
					t.Fatalf("failed to get client; %v", err)
				}
				defer client.Close()

				compared := make(map[int64]*big.Int)
				for i := int64(0); i < 100; i++ {
					blk := big.NewInt(i)
					block, err := client.BlockByNumber(context.Background(), blk)
					if err != nil {
						// all blocks have been checked
						if errors.Is(err, ethereum.NotFound) {
							return
						}
						t.Fatalf("failed to get block header; %v", err)
					}
					if block.Difficulty().Uint64() != 0 {
						t.Errorf("incorrect difficulty for block %d got: %d, want: %d", i, block.Difficulty().Uint64(), 0)
					}

					currentPrevrandao := block.MixDigest().Big() // MixDigest == MixHash == PrevRandao
					for cmpBlk, cmp := range compared {
						if currentPrevrandao.Cmp(cmp) == 0 {
							t.Errorf("found same prevrandao for blocks %d and %d: %s, %s", cmpBlk, i, cmp, currentPrevrandao)
						}
					}
					compared[i] = currentPrevrandao
				}
			},
		},
	}

	net, err := StartIntegrationTestNet(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to start the fake network: %v", err)
	}
	defer net.Stop()
	// Deploy the contract.
	contract, _, err := DeployContract(net, prevrandao.DeployPrevrandao)
	if err != nil {
		t.Fatalf("failed to deploy contract; %v", err)
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.f(net, contract)
		})
	}

}
