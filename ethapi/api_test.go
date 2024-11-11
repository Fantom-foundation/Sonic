package ethapi

import (
	"context"
	"math/big"
	"testing"

	"github.com/Fantom-foundation/go-opera/evmcore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
	gomock "go.uber.org/mock/gomock"
)

func TestGetBlockReceipts(t *testing.T) {

	tests := []struct {
		name  string
		block rpc.BlockNumberOrHash
	}{
		{
			name:  "number",
			block: rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(42)),
		},
		{
			name:  "latest",
			block: rpc.BlockNumberOrHashWithNumber(rpc.LatestBlockNumber),
		},
		{
			name:  "pending",
			block: rpc.BlockNumberOrHashWithNumber(rpc.PendingBlockNumber),
		},
		{
			name:  "hash",
			block: rpc.BlockNumberOrHashWithHash(common.Hash{42}, false),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receipts, err := testGetBlockReceipts(t, tt.block)
			if err != nil {
				t.Fatal(err)
			}

			if len(receipts) != 1 {
				t.Fatalf("expected 1 receipt, got %d", len(receipts))
			}
		})
	}
}

func testGetBlockReceipts(t *testing.T, blockParam rpc.BlockNumberOrHash) ([]map[string]interface{}, error) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockObj := NewMockBackend(ctrl)

	header, transaction, receipts, err := getTestData()
	if err != nil {
		return nil, err
	}

	if blockParam.BlockNumber != nil {
		mockObj.EXPECT().HeaderByNumber(gomock.Any(), *blockParam.BlockNumber).Return(header, nil)
	}

	if blockParam.BlockHash != nil {
		mockObj.EXPECT().HeaderByHash(gomock.Any(), *blockParam.BlockHash).Return(header, nil)
	}

	mockObj.EXPECT().GetReceiptsByNumber(gomock.Any(), gomock.Any()).Return(receipts, nil)
	mockObj.EXPECT().GetTransaction(gomock.Any(), transaction.Hash()).Return(transaction, uint64(0), uint64(0), nil)
	mockObj.EXPECT().ChainConfig().Return(&params.ChainConfig{}).AnyTimes()

	api := NewPublicTransactionPoolAPI(
		mockObj,
		&AddrLocker{},
	)

	receiptsRes, err := api.GetBlockReceipts(context.Background(), blockParam)
	if err != nil {
		return nil, err
	}

	return receiptsRes, nil
}

func getTestData() (*evmcore.EvmHeader, *types.Transaction, types.Receipts, error) {

	key, err := crypto.GenerateKey()
	if err != nil {
		return nil, nil, nil, err
	}

	address := crypto.PubkeyToAddress(key.PublicKey)
	chainId := big.NewInt(1)

	transaction, err := types.SignTx(types.NewTx(&types.AccessListTx{
		ChainID:  chainId,
		Gas:      21000,
		GasPrice: big.NewInt(1),
		To:       &address,
		Nonce:    0,
	}), types.NewLondonSigner(chainId), key)
	if err != nil {
		return nil, nil, nil, err
	}

	header := &evmcore.EvmHeader{
		Number: big.NewInt(1),
	}

	receipt := types.Receipt{
		Status:  1,
		TxHash:  transaction.Hash(),
		GasUsed: 0,
	}

	receipts := types.Receipts{
		&receipt,
	}
	return header, transaction, receipts, nil
}
