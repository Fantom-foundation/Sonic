package tests

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestGasPrices_EvolutionFollowsGasPriceModel(t *testing.T) {

	net, err := StartIntegrationTestNet(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to start the fake network: %v", err)
	}
	defer net.Stop()

	// Produce a few blocks on the network.
	for range 10 {
		_, err := net.EndowAccount(common.Address{42}, 100)
		if err != nil {
			t.Fatalf("failed to endow account; %v", err)
		}
	}

	client, err := net.GetClient()
	if err != nil {
		t.Fatalf("failed to get client; %v", err)
	}
	defer client.Close()

	lastBlock, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		t.Fatalf("failed to get block header; %v", err)
	}

	for i := int64(0); i < int64(lastBlock.NumberU64()); i++ {
		block, err := client.BlockByNumber(context.Background(), big.NewInt(i))
		if err != nil {
			t.Fatalf("failed to get block header; %v", err)
		}
		fmt.Printf("Block %d - %v\n", i, block.BaseFee())
	}
	// TODO: check the gas price evolution
	//t.Fail()
}


// TODO:
//  - add a test checking the accuracy of the gas price suggestions
//  - test that transactions are charged the base fee, not their maximum price
