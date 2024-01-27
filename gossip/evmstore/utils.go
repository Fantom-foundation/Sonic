package evmstore

import (
	"bytes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
)

func IsMptKey(key []byte) bool {
	return len(key) == common.HashLength ||
		(bytes.HasPrefix(key, rawdb.CodePrefix) && len(key) == len(rawdb.CodePrefix)+common.HashLength)
}

func IsPreimageKey(key []byte) bool {
	preimagePrefix := []byte("secure-key-")
	return bytes.HasPrefix(key, preimagePrefix) && len(key) == (len(preimagePrefix)+common.HashLength)
}
