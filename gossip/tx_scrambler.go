package gossip

import (
	"bytes"
	"crypto/sha256"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/bitutil"
	"slices"
)

type scramblerEntry struct {
	hash   common.Hash
	sender common.Address
	nonce  uint64
}

type scramblerSortFunc func(entries []*scramblerEntry, seed []byte) []*scramblerEntry

// getExecutionOrder first removes any entries with duplicate hashes, then sorts the list by XORed hashes.
// Furthermore, if there are entries with same sender, these entries are sorted by their nonce (lower comes first).
func getExecutionOrder(entries []*scramblerEntry, sortFunc scramblerSortFunc) []*scramblerEntry {
	hashList := make([]byte, 0)
	seenAddresses := make(map[common.Address]bool)
	duplicateAddresses := make(map[common.Address]int)

	m := make(map[common.Hash]bool)
	uniqueList := make([]*scramblerEntry, 0)
	for _, entry := range entries {
		// skip any duplicate hashes
		if _, ok := m[entry.hash]; ok {
			continue
		}

		// Remove txs with duplicate hashes using map
		m[entry.hash] = true
		uniqueList = append(uniqueList, entry)
		hashList = append(hashList, entry.hash.Bytes()...)

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
	seed := sha256.Sum256(hashList)
	sorted := sortFunc(uniqueList, seed[:])

	// if no duplicate addresses, return early
	if len(duplicateAddresses) == 0 {
		return sorted
	}

	indexMap := make(map[common.Address][]int)
	for i, entry := range sorted {
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
				a := sorted[idxs[i]].nonce
				b := sorted[idxs[j]].nonce
				// txs with smaller nonce must be executed first, otherwise they will never be executed
				if a < b {
					e = *sorted[idxs[i]]
					*sorted[idxs[i]] = *sorted[idxs[j]]
					*sorted[idxs[j]] = e
				}

			}
		}
	}

	return sorted
}

// builtInSort uses golang built in sort func.
func builtInSort(entries []*scramblerEntry, seed []byte) []*scramblerEntry {
	slices.SortFunc(entries, func(a, b *scramblerEntry) int {
		var (
			aX = make([]byte, 32)
			bX = make([]byte, 32)
		)

		bitutil.XORBytes(aX, a.hash.Bytes()[:], seed[:])
		bitutil.XORBytes(bX, b.hash.Bytes()[:], seed[:])
		return bytes.Compare(aX[:], bX[:])
	})
	return entries
}

// quickSort is a quicksort implementation for comparing XORed entries hashes.
func quickSort(entries []*scramblerEntry, seed []byte) []*scramblerEntry {
	qSort(entries, seed, 0, len(entries)-1)
	return entries
}

func qSort(entries []*scramblerEntry, seed []byte, low, high int) {
	if low < high {
		pivot := partition(entries, seed, low, high)
		qSort(entries, seed, low, pivot-1)
		qSort(entries, seed, pivot+1, high)
	}
}

func partition(entries []*scramblerEntry, seed []byte, low, high int) int {
	pivot := entries[high]
	pivotX := xorWithSeed(pivot.hash.Bytes(), seed)
	i := low - 1

	for j := low; j < high; j++ {
		currentX := xorWithSeed(entries[j].hash.Bytes(), seed)
		if bytes.Compare(currentX, pivotX) < 0 {
			i++
			// wwap entries[i] and entries[j]
			entries[i], entries[j] = entries[j], entries[i]
		}
	}
	// wwap the pivot to its correct position
	entries[i+1], entries[high] = entries[high], entries[i+1]
	return i + 1
}

// xorWithSeed returns XORed hash with seed.
func xorWithSeed(hash, seed []byte) []byte {
	dst := make([]byte, 32)
	for i := range dst {
		dst[i] = hash[i] ^ seed[i]
	}
	return dst
}
