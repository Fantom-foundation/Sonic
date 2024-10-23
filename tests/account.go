package tests

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Account struct {
	PrivateKey *ecdsa.PrivateKey
}

func NewAccount() *Account {
	key, _ := crypto.GenerateKey()
	return &Account{
		PrivateKey: key,
	}
}

func (a *Account) Address() common.Address {
	return crypto.PubkeyToAddress(a.PrivateKey.PublicKey)
}
