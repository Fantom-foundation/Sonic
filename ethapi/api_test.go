package ethapi

import (
	"context"
	cc "github.com/Fantom-foundation/Carmen/go/common"
	"github.com/Fantom-foundation/Carmen/go/common/amount"
	"github.com/Fantom-foundation/Carmen/go/common/immutable"
	"github.com/Fantom-foundation/Carmen/go/common/witness"
	"github.com/Fantom-foundation/go-opera/inter/state"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"math/big"
	"testing"

	"github.com/Fantom-foundation/go-opera/evmcore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
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

func TestAPI_GetProof(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Input address and keys for witness proof
	addr := cc.Address{1}
	keys := []string{"0x1"}
	hexKeys := []common.Hash{common.HexToHash("0x1")}

	// Return data
	codeHash := cc.Hash{2}
	storageHash := cc.Hash{3}
	balance := amount.New(4)
	nonce := cc.ToNonce(5)
	headerRoot := common.Hash{6}
	storageElements := []immutable.Bytes{immutable.NewBytes([]byte("stElement"))}
	value := cc.Value{7}
	storageProof := StorageResult{
		Key:   hexKeys[0].Hex(),
		Value: (*hexutil.Big)(new(big.Int).SetBytes(value[:])),
		Proof: toHexSlice(storageElements),
	}
	accountElements := []immutable.Bytes{immutable.NewBytes([]byte("accElement"))}

	// Mocks
	mockBackend := NewMockBackend(ctrl)
	mockState := state.NewMockStateDB(ctrl)
	mockProof := witness.NewMockProof(ctrl)
	mockHeader := &evmcore.EvmHeader{Root: headerRoot}

	blkNr := rpc.BlockNumberOrHashWithNumber(rpc.LatestBlockNumber)

	mockBackend.EXPECT().StateAndHeaderByNumberOrHash(gomock.Any(), blkNr).Return(mockState, mockHeader, nil)
	mockState.EXPECT().GetProof(common.Address(addr), hexKeys).Return(mockProof, nil)
	mockProof.EXPECT().GetState(cc.Hash(headerRoot), addr, cc.Key(hexKeys[0])).Return(value, true, nil)
	mockProof.EXPECT().GetStorageElements(cc.Hash(headerRoot), addr, cc.Key(hexKeys[0])).Return(storageElements, true)
	mockProof.EXPECT().GetAccountElements(cc.Hash(headerRoot), addr).Return(accountElements, storageHash, true)
	mockProof.EXPECT().GetCodeHash(cc.Hash(headerRoot), addr).Return(codeHash, true, nil)
	mockProof.EXPECT().GetBalance(cc.Hash(headerRoot), addr).Return(balance, true, nil)
	mockProof.EXPECT().GetNonce(cc.Hash(headerRoot), addr).Return(nonce, true, nil)
	mockState.EXPECT().Error().Return(nil)
	mockState.EXPECT().Release()

	api := NewPublicBlockChainAPI(mockBackend)

	accountProof, err := api.GetProof(context.Background(), common.Address(addr), keys, blkNr)
	require.NoError(t, err, "failed to get account")

	u256Balance := balance.Uint256()
	require.Equal(t, common.Address(addr), accountProof.Address)
	require.Equal(t, toHexSlice(accountElements), accountProof.AccountProof)
	require.Equal(t, (*hexutil.U256)(&u256Balance), accountProof.Balance)
	require.Equal(t, common.Hash(codeHash), accountProof.CodeHash)
	require.Equal(t, hexutil.Uint64(nonce.ToUint64()), accountProof.Nonce)
	require.Equal(t, common.Hash(storageHash), accountProof.StorageHash)
	require.Equal(t, []StorageResult{storageProof}, accountProof.StorageProof)
}

func TestAPI_GetAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	addr := common.Address{1}
	codeHash := common.Hash{2}
	storageRoot := cc.Hash{3}
	balance := uint256.NewInt(4)
	nonce := uint64(5)
	headerRoot := common.Hash{123}

	mockBackend := NewMockBackend(ctrl)
	mockState := state.NewMockStateDB(ctrl)
	mockProof := witness.NewMockProof(ctrl)
	mockHeader := &evmcore.EvmHeader{Root: headerRoot}

	blkNr := rpc.BlockNumberOrHashWithNumber(rpc.LatestBlockNumber)

	mockBackend.EXPECT().StateAndHeaderByNumberOrHash(gomock.Any(), blkNr).Return(mockState, mockHeader, nil)
	mockState.EXPECT().GetProof(addr, nil).Return(mockProof, nil)
	mockProof.EXPECT().GetAccountElements(cc.Hash(headerRoot), cc.Address(addr)).Return(nil, storageRoot, true)
	mockState.EXPECT().GetCodeHash(addr).Return(codeHash)
	mockState.EXPECT().GetBalance(addr).Return(balance)
	mockState.EXPECT().GetNonce(addr).Return(nonce)
	mockState.EXPECT().Error().Return(nil)
	mockState.EXPECT().Release()

	api := NewPublicBlockChainAPI(mockBackend)

	account, err := api.GetAccount(context.Background(), addr, blkNr)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if codeHash.Cmp(account.CodeHash) != 0 {
		t.Errorf("unexpected code hash, got: %s want %s", account.CodeHash, codeHash)
	}

	if common.Hash(storageRoot).Cmp(account.StorageRoot) != 0 {
		t.Errorf("unexpected storage root, got: %s want %s", account.StorageRoot, storageRoot)
	}

	if balance.Cmp((*uint256.Int)(account.Balance)) != 0 {
		t.Errorf("unexpected balance, got: %s want %s", account.Balance, balance)
	}

	if balance.Cmp((*uint256.Int)(account.Balance)) != 0 {
		t.Errorf("unexpected balance, got: %s want %s", account.Balance, balance)
	}

	if nonce != uint64(account.Nonce) {
		t.Errorf("unexpected nonce, got: %d want %d", account.Nonce, nonce)
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
