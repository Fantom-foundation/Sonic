package ethapi

import (
	"context"
	"fmt"

	"github.com/Fantom-foundation/Carmen/go/common/amount"
	"github.com/Fantom-foundation/go-opera/gossip/evmstore"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// totalSupplyProvider is the part of the Backend needed by the sonic API.
type totalSupplyProvider interface {
	GetTotalSupply() (*evmstore.TotalSupplySnapshot, error)
}

// PublicSonicAPI provides Sonic specific RPC methods in the `sonic` namespace.
type PublicSonicAPI struct {
	b totalSupplyProvider
}

// NewPublicSonicAPI creates a new sonic API.
func NewPublicSonicAPI(b Backend) *PublicSonicAPI {
	return &PublicSonicAPI{b}
}

// TotalSupplyResult is the result of the `sonic_totalSupply` RPC method.
// All values are valid for the block identified by BlockNumber.
type TotalSupplyResult struct {
	// BlockNumber is the block the values are valid for.
	BlockNumber hexutil.Uint64 `json:"blockNumber"`
	// TotalSupply is the sum of the balances of all accounts.
	TotalSupply *hexutil.Big `json:"totalSupply"`
	// TotalSupplyWithoutBurned is TotalSupply reduced by the balances of the
	// zero address and the 0x000000000000000000000000000000000000dEaD address.
	TotalSupplyWithoutBurned *hexutil.Big `json:"totalSupplyWithoutBurned"`
	// BurnedZero is the balance of the zero address.
	BurnedZero *hexutil.Big `json:"burnedZero"`
	// BurnedDead is the balance of the 0x000000000000000000000000000000000000dEaD address.
	BurnedDead *hexutil.Big `json:"burnedDead"`
}

// TotalSupply implements the `sonic_totalSupply` RPC method. It returns the
// total supply of the native token together with the total supply reduced by
// the balances of the burn addresses (the zero address and 0x...dEaD), and the
// number of the block the values are valid for.
//
// The values are captured during block processing, so they are consistent:
// the balances of the burn addresses reported by `eth_getBalance` for the
// returned block equal the amounts subtracted here.
func (s *PublicSonicAPI) TotalSupply(_ context.Context) (*TotalSupplyResult, error) {
	snapshot, err := s.b.GetTotalSupply()
	if err != nil {
		return nil, err
	}
	burned, overflow := amount.AddOverflow(snapshot.ZeroBalance, snapshot.DeadBalance)
	if overflow {
		return nil, fmt.Errorf("balance of burn addresses overflows")
	}
	withoutBurned, underflow := amount.SubUnderflow(snapshot.Total, burned)
	if underflow {
		return nil, fmt.Errorf("balance of burn addresses exceeds the total supply")
	}
	return &TotalSupplyResult{
		BlockNumber:              hexutil.Uint64(snapshot.Block),
		TotalSupply:              (*hexutil.Big)(snapshot.Total.ToBig()),
		TotalSupplyWithoutBurned: (*hexutil.Big)(withoutBurned.ToBig()),
		BurnedZero:               (*hexutil.Big)(snapshot.ZeroBalance.ToBig()),
		BurnedDead:               (*hexutil.Big)(snapshot.DeadBalance.ToBig()),
	}, nil
}
