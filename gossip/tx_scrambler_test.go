package gossip

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/exp/rand"
	"math"
	"math/big"
	"reflect"
	"testing"
)

func TestTxScrambler_AnalyseEntryList_RemovesDuplicateTransactions(t *testing.T) {
	entries := createScramblerTestInputRepeatedHashes(2, 10)
	entries, _, _ = analyseEntryList(entries)

	checkDuplicateHashes(t, entries)
}

func TestTxScrambler_UnifyEntries_SaltCreationIsDeterministic(t *testing.T) {
	entries := createScramblerTestInputRepeatedAddr(2, 10)
	_, wantedSalt, _ := analyseEntryList(entries)
	for _ = range 10 {
		shuffleEntries(entries)
		_, gotSalt, _ := analyseEntryList(entries)
		if gotSalt != wantedSalt {
			t.Fatal("incorrect salt - salt creation is not deterministic")
		}
	}

}

func TestTxScrambler_AnalyseEntryList_ReportsDuplicateAddresses(t *testing.T) {
	entries := createScramblerTestInputRepeatedAddr(2, 10)
	_, _, hasDuplicateAddresses := analyseEntryList(entries)
	if !hasDuplicateAddresses {
		t.Error("entries have duplicate addresses")
	}

	entries = createScramblerTestInputRepeatedHashes(2, 10)
	_, _, hasDuplicateAddresses = analyseEntryList(entries)
	if hasDuplicateAddresses {
		t.Error("entries does not have duplicate addresses")
	}
}

func TestTxScrambler_ScrambleTransactions_ScrambleIsDeterministic(t *testing.T) {
	res1 := createScramblerTestInputRepeatedAddr(2, 10)
	res2 := deepCopyEntries(res1)
	// shuffle one array
	shuffleEntries(res2)

	for i := 0; i < 10; i++ {
		salt := createRandomSalt()
		scrambleTransactions(res1, salt)
		for j := 0; j < 10; j++ {
			scrambleTransactions(res2, salt)
			if !reflect.DeepEqual(res1, res2) {
				t.Error("scramble is not deterministic")
			}
		}
	}

}

func TestTxScrambler_SortTransactionsByNonce_SortsSameSenderByNonce(t *testing.T) {
	entries := createScramblerTestInputRepeatedAddr(2, 10)
	entries = sortTransactionsByNonce(entries)
	for i := 0; i < len(entries); i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[i].sender == entries[j].sender {
				if entries[i].nonce > entries[j].nonce {
					t.Errorf("incorrect nonce order %d must be before %d", entries[j].nonce, entries[i].nonce)
				}
			}
		}
	}
}

func TestTxScrambler_GetExecutionOrder_SortIsDeterministic_IdenticalData(t *testing.T) {
	const numLoops = 5
	tests := []struct {
		name        string
		createInput func(size int64) []*scramblerEntry
	}{
		{
			name:        "identical hashes",
			createInput: createScramblerTestInputOnlySameHashes,
		},
		{
			name:        "identical sender different nonces",
			createInput: createScramblerTestInputOnlySameAddr,
		},
		{
			name:        "identical sender and nonces",
			createInput: createScramblerTestInputOnlySameAddrAndNonce,
		},
	}

	size := int64(3)
	for i := 0; i < numLoops; i++ {
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				res1 := test.createInput(size)
				res2 := deepCopyEntries(res1)
				// shuffle one array
				shuffleEntries(res2)

				res1 = getExecutionOrder(res1)
				res2 = getExecutionOrder(res2)
				if !reflect.DeepEqual(res1, res2) {
					t.Errorf("slices have different order - algorithm is not deterministic; input size: %d", size)
				}
			})
		}
		size = size * 4
	}
}

func TestTxScrambler_GetExecutionOrder_SortIsDeterministic_RepeatedData(t *testing.T) {
	const numLoops = 5
	tests := []struct {
		name        string
		createInput func(size int64, numRepeats int64) []*scramblerEntry
	}{
		{
			name:        "repeated hashes",
			createInput: createScramblerTestInputRepeatedHashes,
		},
		{
			name:        "repeated sender different nonces",
			createInput: createScramblerTestInputRepeatedAddr,
		},
		{
			name:        "repeated sender and nonces",
			createInput: createScramblerTestInputRepeatedAddrAndNonce,
		},
	}

	size := int64(3)
	for i := 0; i < numLoops; i++ {
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				// skip numRepeats 0 and 1
				res1 := test.createInput(size, int64(i+2))
				res2 := deepCopyEntries(res1)
				// shuffle one array
				shuffleEntries(res2)

				res1 = getExecutionOrder(res1)
				res2 = getExecutionOrder(res2)
				if !reflect.DeepEqual(res1, res2) {
					t.Errorf("slices have different order - algorithm is not deterministic; input size: %d", size)
				}
			})
		}
		size = size * 4
	}
}

func TestTxScrambler_GetExecutionOrder_SortRemovesDuplicateHashes(t *testing.T) {
	entries := createScramblerTestInputRepeatedHashes(5, 10)
	shuffleEntries(entries)
	entries = getExecutionOrder(entries)

	checkDuplicateHashes(t, entries)
}

func TestTxScrambler_GetExecutionOrder_SortsSameSenderByNonce(t *testing.T) {
	const numLoops = 5
	for l := 1; l <= numLoops; l++ {
		entries := createScramblerTestInputRepeatedAddr(int64(l*2), int64(10*l))
		entries = getExecutionOrder(entries)

		for i := 0; i < len(entries); i++ {
			for j := i + 1; j < len(entries); j++ {
				if entries[i].sender == entries[j].sender {
					if entries[i].nonce > entries[j].nonce {
						t.Errorf("incorrect nonce order %d must be before %d", entries[j].nonce, entries[i].nonce)
					}
				}
			}
		}
	}
}

func BenchmarkTxScrambler(b *testing.B) {
	const (
		numLoops   = 4
		multiplier = 10
	)
	size := int64(10)

	for _ = range numLoops {
		b.Run(fmt.Sprintf("builtIn_%d", size), func(b *testing.B) {
			for i := 1; i <= b.N; i++ {
				b.StopTimer()
				entries := createRandomScramblerTestInput(size)
				b.StartTimer()
				getExecutionOrder(entries)
			}
		})
		size = size * multiplier
	}
}

// createRandomScramblerTestInput creates a testing createInput with randomized
// hash and address. This means both address and hashes can repeat.
func createRandomScramblerTestInput(size int64) []*scramblerEntry {
	var entries []*scramblerEntry
	for i := int64(0); i < size; i++ {
		r := rand.Intn(10 - 1)
		entries = append(entries, &scramblerEntry{
			hash:   common.Hash{byte(r)},
			sender: common.Address{byte(r)},
			nonce:  uint64(r),
		})
	}

	return entries
}

func createScramblerTestInputOnlySameAddrAndNonce(size int64) []*scramblerEntry {
	return createScramblerTestInputRepeatedAddrAndNonce(size, size)
}

// createScramblerTestInputRepeatedAddrAndNonce creates testing input of given size which
// swaps the sender and nonce only every X entries. This relies on the numRepeats param.
// If numRepeats == size it only creates duplicate inputs.
// If numRepeats == 1 only different inputs are created.
// Note: Entries are randomly shuffled after creation.
func createScramblerTestInputRepeatedAddrAndNonce(numRepeats int64, size int64) []*scramblerEntry {
	var (
		entries []*scramblerEntry
		addr    common.Address
		nonce   uint64
	)
	for i := int64(0); i < size; i++ {
		if i%numRepeats == 0 {
			// note: same hash different addr/nonce can never happen
			addr = common.BigToAddress(big.NewInt(i))
			nonce = uint64(i)
		}
		entries = append(entries, &scramblerEntry{
			hash:   common.BigToHash(big.NewInt(i)),
			sender: addr,
			nonce:  nonce,
		})
	}
	shuffleEntries(entries)

	return entries
}

func createScramblerTestInputOnlySameHashes(size int64) []*scramblerEntry {
	return createScramblerTestInputRepeatedHashes(size, size)
}

// createScramblerTestInputRepeatedHashes creates testing input of given size
// which swaps the hash only every X entries. This relies on the numRepeats param.
// If numRepeats == size it only creates duplicate inputs.
// If numRepeats == 1 only different inputs are created.
// Note: Entries are randomly shuffled after creation.
func createScramblerTestInputRepeatedHashes(numRepeats int64, size int64) []*scramblerEntry {
	var (
		entries []*scramblerEntry
		hash    common.Hash
		addr    common.Address
		nonce   uint64
	)
	for i := int64(0); i < size; i++ {
		if i%numRepeats == 0 {
			// note: same hash different addr/nonce can never happen
			hash = common.BigToHash(big.NewInt(i))
			addr = common.BigToAddress(big.NewInt(i))
			nonce = uint64(i)
		}
		entries = append(entries, &scramblerEntry{
			hash:   hash,
			sender: addr,
			nonce:  nonce,
		})
	}
	shuffleEntries(entries)

	return entries
}

func createScramblerTestInputOnlySameAddr(size int64) []*scramblerEntry {
	return createScramblerTestInputRepeatedAddr(size, size)
}

// createScramblerTestInputRepeatedAddr creates testing input of given size which
// swaps the address only every X entries. This relies on the numRepeats param.
// If numRepeats == size it only creates inputs with same sender.
// If numRepeats == 1 only different inputs are created.
// Note: Entries are randomly shuffled after creation.
func createScramblerTestInputRepeatedAddr(numRepeats int64, size int64) []*scramblerEntry {
	var (
		entries []*scramblerEntry
		addr    common.Address
	)
	for i := int64(0); i < size; i++ {
		if i%numRepeats == 0 {
			addr = common.BigToAddress(big.NewInt(i))
		}
		entries = append(entries, &scramblerEntry{
			hash:   common.BigToHash(big.NewInt(i)),
			sender: addr,
			nonce:  uint64(i),
		})
	}
	shuffleEntries(entries)

	return entries
}

// shuffleEntries shuffles given entries randomly.
func shuffleEntries(entries []*scramblerEntry) {
	rand.Shuffle(len(entries), func(i, j int) {
		entries[i], entries[j] = entries[j], entries[i]
	})
}

// checkDuplicateHashes checks hash of each entry and fails test if duplicate hash is found.
func checkDuplicateHashes(t *testing.T, entries []*scramblerEntry) {
	seenHashes := make(map[common.Hash]bool)
	for _, entry := range entries {
		if _, ok := seenHashes[entry.hash]; ok {
			t.Fatal("found duplicate hash in entries")
		}
		seenHashes[entry.hash] = true
	}
}

func createRandomSalt() [32]byte {
	var salt = [32]byte{}
	for i := 0; i < 32; i++ {
		salt[i] = byte(rand.Intn(math.MaxUint8))
	}
	return salt
}
