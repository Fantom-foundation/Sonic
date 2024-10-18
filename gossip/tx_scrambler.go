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

type scramblerSortFunc func(entries []*scramblerEntry, salt [32]byte) []*scramblerEntry

// getExecutionOrder first removes any entries with duplicate hashes, then sorts the list by XORed hashes.
// Furthermore, if there are entries with same sender, these entries are sorted by their nonce (lower comes first).
func getExecutionOrder(entries []*scramblerEntry, sortFunc scramblerSortFunc) []*scramblerEntry {
	var salt [32]byte

	seenAddresses := make(map[common.Address]bool)
	duplicateAddresses := make(map[common.Address]int)
	seenHashes := make(map[common.Hash]bool)
	uniqueList := make([]*scramblerEntry, 0, len(entries))
	for _, entry := range entries {
		// skip any duplicate hashes
		if _, ok := seenHashes[entry.hash]; ok {
			continue
		}

		salt = xorBytes32(salt, entry.hash)

		// Remove txs with duplicate hashes using map
		seenHashes[entry.hash] = true
		uniqueList = append(uniqueList, entry)

		if _, ok := seenAddresses[entry.sender]; ok {
			// map number of occurrence
			if _, ok := duplicateAddresses[entry.sender]; ok {
				// if already marked as duplicate, only increment the occurrence
				duplicateAddresses[entry.sender]++
			} else {
				// if not yet marked, mark it with 2 because at this point, this address has been seen second time
				duplicateAddresses[entry.sender] = 2
			}
		} else {
			// mark seen address
			seenAddresses[entry.sender] = true
		}
	}

	entries = sortFunc(uniqueList, salt)

	// if no duplicate addresses, return early
	if len(duplicateAddresses) == 0 {
		return entries
	}

	indexMap := make(map[common.Address][]int)
	for i, entry := range entries {
		// find if address is duplicate
		occurrence, ok := duplicateAddresses[entry.sender]
		if !ok {
			continue
		}
		// if first time found, mark create array and mark first index
		idxs, ok := indexMap[entry.sender]
		if !ok {
			indexMap[entry.sender] = []int{i}
			continue
		}

		// if we have not found all indexes of given address, append and continue
		if len(idxs) != occurrence {
			indexMap[entry.sender] = append(idxs, i)
			continue
		}
	}

	for _, idxs := range indexMap {
		for i := 0; i < len(idxs); i++ {
			for j := 0; j < len(idxs); j++ {
				var e scramblerEntry
				a := entries[idxs[i]].nonce
				b := entries[idxs[j]].nonce
				// txs with smaller nonce must be executed first, otherwise they will never be executed
				if a < b {
					e = *entries[idxs[i]]
					*entries[idxs[i]] = *entries[idxs[j]]
					*entries[idxs[j]] = e
				}

			}
		}
	}

	return entries
}

// builtInSort uses golang built in sort func.
func builtInSort(entries []*scramblerEntry, salt [32]byte) []*scramblerEntry {
	var aX, bX [32]byte
	slices.SortFunc(entries, func(a, b *scramblerEntry) int {
		aX = xorBytes32(a.hash, salt)
		bX = xorBytes32(b.hash, salt)
		return bytes.Compare(aX[:], bX[:])
	})
	return entries
}

func xorBytes32(a, b [32]byte) (dst [32]byte) {
	for i := 0; i < 32; i++ {
		dst[i] = a[i] ^ b[i]
	}
	return
}
