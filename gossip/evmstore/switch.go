package evmstore

import (
	cc "github.com/Fantom-foundation/Carmen/go/common"
	"github.com/Fantom-foundation/Carmen/go/evmstore"
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
)

type Backend interface {
	// SetTxPosition stores transaction block and position.
	SetTxPosition(txid common.Hash, position TxPosition) error

	// GetTxPosition returns stored transaction block and position.
	GetTxPosition(txid common.Hash) (*TxPosition, error)

	// SetTx stores non-event transaction.
	SetTx(txid common.Hash, tx *types.Transaction) error

	// GetTx returns stored non-event transaction.
	GetTx(txid common.Hash) (*types.Transaction, error)

	// SetRawReceipts stores raw transaction receipts for one block.
	SetRawReceipts(n idx.Block, receipts []byte) error

	// GetRawReceipts loads raw transaction receipts.
	GetRawReceipts(n idx.Block) ([]byte, error)

	// Flush the store
	Flush() error

	io.Closer
}

type defaultBackend struct {
	s *Store
}

func (s defaultBackend) SetTxPosition(txid common.Hash, position TxPosition) error {
	s.s.rlp.Set(s.s.table.TxPositions, txid.Bytes(), &position)
	return nil
}

func (s defaultBackend) GetTxPosition(txid common.Hash) (*TxPosition, error) {
	txPosition, _ := s.s.rlp.Get(s.s.table.TxPositions, txid.Bytes(), &TxPosition{}).(*TxPosition)
	return txPosition, nil
}

func (s defaultBackend) SetTx(txid common.Hash, tx *types.Transaction) error {
	s.s.rlp.Set(s.s.table.Txs, txid.Bytes(), tx)
	return nil
}

func (s defaultBackend) GetTx(txid common.Hash) (*types.Transaction, error) {
	tx, _ := s.s.rlp.Get(s.s.table.Txs, txid.Bytes(), &types.Transaction{}).(*types.Transaction)
	return tx, nil
}

func (s defaultBackend) SetRawReceipts(n idx.Block, receipts []byte) error {
	return s.s.table.Receipts.Put(n.Bytes(), receipts)
}

func (s defaultBackend) GetRawReceipts(n idx.Block) ([]byte, error) {
	return s.s.table.Receipts.Get(n.Bytes())
}

func (s defaultBackend) Flush() error {
	return nil
}

func (s defaultBackend) Close() error {
	return nil
}

type carmenBackend struct {
	c evmstore.EvmStore
}

func (s carmenBackend) SetTxPosition(txid common.Hash, position TxPosition) error {
	return s.c.SetTxPosition(cc.Hash(txid), evmstore.TxPosition{
		Block:       uint64(position.Block),
		Event:       cc.Hash(position.Event),
		EventOffset: position.EventOffset,
		BlockOffset: position.BlockOffset,
	})
}

func (s carmenBackend) GetTxPosition(txid common.Hash) (*TxPosition, error) {
	txPosition, err := s.c.GetTxPosition(cc.Hash(txid))
	if err != nil {
		return nil, err
	}
	return &TxPosition{
		Block:       idx.Block(txPosition.Block),
		Event:       hash.Event(txPosition.Event),
		EventOffset: txPosition.EventOffset,
		BlockOffset: txPosition.BlockOffset,
	}, nil
}

func (s carmenBackend) SetTx(txid common.Hash, tx *types.Transaction) error {
	buf, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return err
	}
	return s.c.SetTx(cc.Hash(txid), buf)
}

func (s carmenBackend) GetTx(txid common.Hash) (*types.Transaction, error) {
	buf, err := s.c.GetTx(cc.Hash(txid))
	var tx = &types.Transaction{}
	err = rlp.DecodeBytes(buf, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (s carmenBackend) SetRawReceipts(n idx.Block, receipts []byte) error {
	return s.c.SetRawReceipts(uint64(n), receipts)
}

func (s carmenBackend) GetRawReceipts(n idx.Block) ([]byte, error) {
	return s.c.GetRawReceipts(uint64(n))
}

func (s carmenBackend) Flush() error {
	return s.c.Flush()
}

func (s carmenBackend) Close() error {
	return s.c.Close()
}
