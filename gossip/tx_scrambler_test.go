package gossip

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/exp/rand"
	"math/big"
	"reflect"
	"testing"
)

// scramblerTestInputSize determines size of input for tx scrambler tests.
const scramblerTestInputSize = 10

func TestTxScrambler_UnifyEntries_RemovesDuplicateTransactions(t *testing.T) {
	entries := createScramblerTestInputWithRepeatedHashes()
	entries, _, _ = unifyEntries(entries)
	if got, want := len(entries), scramblerTestInputSize/2; got != want {
		t.Errorf("duplicate entries does not seem to be removed, got size: %v, want size: %v", got, want)
	}
}

func TestTxScrambler_UnifyEntries_CorrectlyCreatesSalt(t *testing.T) {
	entries := createScramblerTestInputWithRepeatedAddresses()
	var wantedSalt [32]byte
	for _, entry := range entries {
		wantedSalt = xorBytes32(wantedSalt, entry.hash)
	}
	_, gotSalt, _ := unifyEntries(entries)
	if gotSalt != wantedSalt {
		t.Error("incorrect salt")
	}
}

func TestTxScrambler_UnifyEntries_ReportsDuplicateAddresses(t *testing.T) {
	entries := createScramblerTestInputWithRepeatedAddresses()
	_, _, hasDuplicateAddresses := unifyEntries(entries)
	if !hasDuplicateAddresses {
		t.Error("entries have duplicate addresses")
	}

	entries = createScramblerTestInputWithRepeatedHashes()
	_, _, hasDuplicateAddresses = unifyEntries(entries)
	if hasDuplicateAddresses {
		t.Error("entries does not have duplicate addresses")
	}
}

func TestTxScrambler_ScrambleTransactions_ScrambleIsDeterministic(t *testing.T) {
	res1 := createScramblerTestInputWithRepeatedAddresses()
	res2 := deepCopyEntries(res1)
	// shuffle one array
	shuffleEntries(res2)

	salt := [32]byte{1}
	scrambleTransactions(res1, salt)
	scrambleTransactions(res2, salt)
	if !reflect.DeepEqual(res1, res2) {
		t.Error("scramble is not deterministic")
	}
}

func TestTxScrambler_SortTransactionsByNonce_SortsSameSenderByNonce(t *testing.T) {
	entries := createScramblerTestInputWithRepeatedAddresses()
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

func TestTxScrambler_GetExecutionOrder_SortIsDeterministic(t *testing.T) {
	res1 := createRandomScramblerTestInput()
	res2 := deepCopyEntries(res1)
	// shuffle one array
	shuffleEntries(res2)

	res1 = getExecutionOrder(res1)
	res2 = getExecutionOrder(res2)
	if !reflect.DeepEqual(res1, res2) {
		t.Fatal("slices have different order - algorithm is not deterministic")
	}
}

func TestTxScrambler_GetExecutionOrder_SortRemovesDuplicateHashes(t *testing.T) {
	const size = 10
	entries := createScramblerTestInputWithRepeatedHashes()
	shuffleEntries(entries)

	entries = getExecutionOrder(entries)

	if got, want := len(entries), size/2; got != want {
		t.Errorf("duplicate entries does not seem to be removed, got size: %v, want size: %v", got, want)
	}
}

func TestTxScrambler_GetExecutionOrder_SortsSameSenderByNonce(t *testing.T) {
	entries := createScramblerTestInputWithRepeatedAddresses()
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
				entries := createRandomScramblerTestInput()
				b.StartTimer()
				getExecutionOrder(entries)
			}
		})
		size = size * multiplier
	}
}

// createRandomScramblerTestInput creates a testing input with randomized
// hash and address. This means both address and hashes can repeat.
func createRandomScramblerTestInput() []*scramblerEntry {
	var entries []*scramblerEntry
	for i := int64(0); i < scramblerTestInputSize; i++ {
		r := rand.Intn(10 - 1)
		entries = append(entries, &scramblerEntry{
			hash:   common.Hash{byte(r)},
			sender: common.Address{byte(r)},
			nonce:  uint64(r),
		})
	}

	return entries
}

// createScramblerTestInputWithRepeatedHashes creates testing input which swaps the
// hash only every other entry. creating half of the entries without unique hash.
// Note: Entries are randomly shuffled after creation.
func createScramblerTestInputWithRepeatedHashes() []*scramblerEntry {
	var (
		entries []*scramblerEntry
		hash    common.Hash
	)
	for i := int64(0); i < scramblerTestInputSize; i++ {
		if i%2 == 0 {
			hash = common.BigToHash(big.NewInt(i))
		}
		entries = append(entries, &scramblerEntry{
			hash:   hash,
			sender: common.BigToAddress(big.NewInt(i)),
			nonce:  uint64(i),
		})
	}
	shuffleEntries(entries)

	return entries
}

// createScramblerTestInputWithRepeatedAddresses creates testing input which swaps the
// address only every other entry. Creating half of the entries without unique address.
// Note: Entries are randomly shuffled after creation.
func createScramblerTestInputWithRepeatedAddresses() []*scramblerEntry {
	var (
		entries []*scramblerEntry
		addr    common.Address
	)
	for i := int64(0); i < scramblerTestInputSize; i++ {
		if i%2 == 0 {
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
