package gossip

import (
	"bytes"
	"github.com/ethereum/go-ethereum/common"
	"slices"
)

type scramblerEntry struct {
	hash   common.Hash
	sender common.Address
	nonce  uint64
}

// getExecutionOrder first removes any entries with duplicate hashes, then sorts the list by XORed hashes.
// Furthermore, if there are entries with same sender, these entries are sorted by their nonce (lower comes first).
func getExecutionOrder(entries []*scramblerEntry) []*scramblerEntry {
	uniqueList, salt, hasDuplicateAddresses := analyseEntryList(entries)
	scrambleTransactions(uniqueList, salt)

	// no need to sort more
	if !hasDuplicateAddresses {
		return uniqueList
	}

	return sortTransactionsByNonce(uniqueList)
}

// sortTransactionsByNonce finds any duplicate senders and sorts their transactions by nonce ascending.
func sortTransactionsByNonce(entries []*scramblerEntry) []*scramblerEntry {
	senderNonceOrder := deepCopyEntries(entries)
	// sort copied slice so that it has all txs from same address together + sorted by nonce ascending
	slices.SortFunc(senderNonceOrder, func(a, b *scramblerEntry) int {
		cmp := a.sender.Cmp(b.sender)
		if cmp != 0 {
			return cmp
		}
		if a.nonce > b.nonce {
			return 1
		}
		// todo resolve same nonce and same address - add gas comparsion and only allow bigger gas
		return -1
	})

	// find the first entry for each sender in the senderNonceOrder
	senderIndex := make(map[common.Address]int)
	for idx, entry := range senderNonceOrder {
		if _, ok := senderIndex[entry.sender]; !ok {
			senderIndex[entry.sender] = idx
		}
	}

	// replace already scrambled entries so that they are sorted by nonce
	for idx := 0; idx < len(entries); idx++ {
		sender := entries[idx].sender
		entries[idx] = senderNonceOrder[senderIndex[sender]]
		senderIndex[sender]++
	}

	return entries
}

// scrambleTransactions scrambles transactions by comparing its XORed hashes with salt
func scrambleTransactions(list []*scramblerEntry, salt [32]byte) {
	var aX, bX [32]byte
	slices.SortFunc(list, func(a, b *scramblerEntry) int {
		aX = xorBytes32(a.hash, salt)
		bX = xorBytes32(b.hash, salt)
		return bytes.Compare(aX[:], bX[:])
	})
}

// analyseEntryList removes any transactions with duplicate hashes and creates the XOR salt from the unique tx list.
// Furthermore, it returns whether given list of entries contains duplicate addresses.
func analyseEntryList(entries []*scramblerEntry) ([]*scramblerEntry, [32]byte, bool) {
	var (
		salt                  [32]byte
		hasDuplicateAddresses bool
	)

	seenHashes := make(map[common.Hash]bool)
	seenAddresses := make(map[common.Address]bool)
	uniqueList := make([]*scramblerEntry, 0, len(entries))
	for _, entry := range entries {
		// skip any duplicate hashes
		if _, ok := seenHashes[entry.hash]; ok {
			continue
		}
		// mark whether we have duplicate addresses
		if _, ok := seenAddresses[entry.sender]; ok {
			hasDuplicateAddresses = true
		}

		salt = xorBytes32(salt, entry.hash)
		uniqueList = append(uniqueList, entry)
		seenHashes[entry.hash] = true
		seenAddresses[entry.sender] = true
	}

	return uniqueList, salt, hasDuplicateAddresses
}

func xorBytes32(a, b [32]byte) (dst [32]byte) {
	for i := 0; i < 32; i++ {
		dst[i] = a[i] ^ b[i]
	}
	return
}

// deepCopyEntries returns deep copy of entries.
func deepCopyEntries(entries []*scramblerEntry) []*scramblerEntry {
	cpy := make([]*scramblerEntry, len(entries))
	// make a deep copy
	for i, e := range entries {
		copied := *e
		cpy[i] = &copied
	}
	return cpy
}
