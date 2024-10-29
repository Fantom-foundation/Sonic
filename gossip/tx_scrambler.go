package gossip

import (
	"bytes"
	"cmp"
	"github.com/ethereum/go-ethereum/common"
	"slices"
)

type scramblerEntry struct {
	hash   common.Hash
	sender common.Address
	nonce  uint64
	gas    uint64
}

// getExecutionOrder first removes any entries with duplicate hashes, then sorts the list by XORed hashes.
// Furthermore, if there are entries with same sender, these entries are sorted by their nonce (lower comes first).
func getExecutionOrder(entries []*scramblerEntry) []*scramblerEntry {
	uniqueList, salt, hasDuplicateAddresses := analyseEntryList(entries)
	scrambleTransactions(uniqueList, salt)

	// do we need to sort more?
	if hasDuplicateAddresses {
		sortTransactionsWithSameSender(uniqueList)
	}

	return uniqueList
}

// sortTransactionsWithSameSender finds any duplicate senders and sorts their transactions by nonce ascending.
func sortTransactionsWithSameSender(entries []*scramblerEntry) {
	senderNonceOrder := slices.Clone(entries)
	// sort copied slice so that it has all txs from same address together + sorted by nonce ascending
	slices.SortFunc(senderNonceOrder, func(a, b *scramblerEntry) int {
		cmp := a.sender.Cmp(b.sender)
		if cmp != 0 {
			return cmp
		}
		// if addresses are same, sort by nonce
		res := compareAsc(a.nonce, b.nonce)
		if res != 0 {
			return res
		}
		// if nonce is same, sort by gas
		res = compareDesc(a.gas, b.gas)
		if res != 0 {
			return res
		}
		// if both nonce and gas are equal, sort by hash
		// note: at this point, hashes can never be same - duplicates are removed
		return a.hash.Cmp(b.hash)
	})

	// find the first entry for each sender in the senderNonceOrder
	senderIndex := make(map[common.Address]int)
	for idx, entry := range senderNonceOrder {
		if _, found := senderIndex[entry.sender]; !found {
			senderIndex[entry.sender] = idx
		}
	}

	// replace already scrambled entries so that they are sorted by nonce
	for idx := range entries {
		sender := entries[idx].sender
		entries[idx] = senderNonceOrder[senderIndex[sender]]
		senderIndex[sender]++
	}

	return
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

	seenHashes := make(map[common.Hash]struct{})
	seenAddresses := make(map[common.Address]struct{})
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
		seenHashes[entry.hash] = struct{}{}
		seenAddresses[entry.sender] = struct{}{}
	}

	return uniqueList, salt, hasDuplicateAddresses
}

func xorBytes32(a, b [32]byte) (dst [32]byte) {
	for i := 0; i < 32; i++ {
		dst[i] = a[i] ^ b[i]
	}
	return
}

func compareAsc[T cmp.Ordered](a T, b T) int {
	if a > b {
		return 1
	}
	if a < b {
		return -1
	}
	return 0

}

func compareDesc[T cmp.Ordered](a T, b T) int {
	if a > b {
		return -1
	}
	if a < b {
		return 1
	}
	return 0

}
