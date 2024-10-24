package gossip

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/exp/rand"
	"math"
	"reflect"
	"testing"
)

func TestTxScrambler_AnalyseEntryList_RemovesDuplicateTransactions(t *testing.T) {
	entries := []*scramblerEntry{
		{hash: common.Hash{1}},
		{hash: common.Hash{2}},
		{hash: common.Hash{3}},
		{hash: common.Hash{2}},
		{hash: common.Hash{1}},
	}

	shuffleEntries(entries)
	result, _, _ := analyseEntryList(entries)
	if len(result) == 0 {
		t.Fatal("analyseEntryList returned empty list")
	}

	seen := map[common.Hash]struct{}{}
	for _, entry := range result {
		if _, seen := seen[entry.hash]; seen {
			t.Fatalf("duplicate hash %v", entry.hash)
		}
		seen[entry.hash] = struct{}{}
	}
}

func TestTxScrambler_UnifyEntries_SaltCreationIsDeterministic(t *testing.T) {
	entries := []*scramblerEntry{
		{hash: common.Hash{1}},
		{hash: common.Hash{2}},
		{hash: common.Hash{3}},
		{hash: common.Hash{2}},
		{hash: common.Hash{1}},
	}

	_, wantedSalt, _ := analyseEntryList(entries)
	for range 10 {
		shuffleEntries(entries)
		_, gotSalt, _ := analyseEntryList(entries)
		if gotSalt != wantedSalt {
			t.Fatal("incorrect salt - salt creation is not deterministic")
		}
	}

}

func TestTxScrambler_AnalyseEntryList_ReportsDuplicateAddresses(t *testing.T) {
	tests := []struct {
		name         string
		input        []*scramblerEntry
		hasDuplicate bool
	}{
		{
			name: "has duplicate address",
			input: []*scramblerEntry{
				{sender: common.Address{1}},
				{sender: common.Address{3}},
				{sender: common.Address{2}},
				{sender: common.Address{3}},
			},
			hasDuplicate: true,
		},
		{
			name: "has no duplicate address",
			input: []*scramblerEntry{
				{sender: common.Address{1}},
				{sender: common.Address{2}},
				{sender: common.Address{3}},
			},
			hasDuplicate: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, _, hasDuplicateAddresses := analyseEntryList(test.input)
			if hasDuplicateAddresses != hasDuplicateAddresses {
				t.Error("wrongly reported duplicate address")
			}
		})
	}

}

func TestTxScrambler_ScrambleTransactions_ScrambleIsDeterministic(t *testing.T) {
	res1 := []*scramblerEntry{
		{hash: common.Hash{1}},
		{hash: common.Hash{2}},
		{hash: common.Hash{3}},
		{hash: common.Hash{2}},
		{hash: common.Hash{1}},
	}

	res2 := deepCopyEntries(res1)

	for i := 0; i < 10; i++ {
		salt := createRandomSalt()
		scrambleTransactions(res1, salt)
		for j := 0; j < 10; j++ {
			shuffleEntries(res2)
			scrambleTransactions(res2, salt)
			if !reflect.DeepEqual(res1, res2) {
				t.Error("scramble is not deterministic")
			}
		}
	}
}

func TestTxScrambler_SortTransactionsWithSameSender_SortsByNonce(t *testing.T) {
	entries := []*scramblerEntry{
		{
			hash:   common.Hash{1},
			sender: common.Address{1},
			nonce:  1,
		},
		{
			hash:   common.Hash{2},
			sender: common.Address{2},
			nonce:  1,
		},
		{
			hash:   common.Hash{3},
			sender: common.Address{3},
			nonce:  1,
		},
		{
			hash:   common.Hash{4},
			sender: common.Address{2},
			nonce:  2,
		},
		{
			hash:   common.Hash{5},
			sender: common.Address{1},
			nonce:  2,
		},
	}

	sortTransactionsWithSameSender(entries)
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

func TestTxScrambler_SortTransactionsWithSameSender_SortsByGasIfNonceIsSame(t *testing.T) {
	entries := []*scramblerEntry{
		{
			hash:   common.Hash{1},
			sender: common.Address{1},
			nonce:  1,
			gas:    1,
		},
		{
			hash:   common.Hash{2},
			sender: common.Address{1},
			nonce:  1,
			gas:    2,
		},
		{
			hash:   common.Hash{3},
			sender: common.Address{2},
			nonce:  1,
			gas:    3,
		},
		{
			hash:   common.Hash{4},
			sender: common.Address{2},
			nonce:  1,
			gas:    4,
		},
	}

	shuffleEntries(entries)
	sortTransactionsWithSameSender(entries)
	for i := 0; i < len(entries); i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[i].sender == entries[j].sender {
				if entries[i].nonce > entries[j].nonce {
					t.Errorf("incorrect nonce order %d must be before %d", entries[j].nonce, entries[i].nonce)
				}
				if entries[i].nonce == entries[j].nonce && entries[i].gas < entries[j].gas {
					t.Errorf("incorrect gas order %d must be before %d", entries[i].gas, entries[j].gas)
				}
			}
		}
	}
}

func TestTxScrambler_GetExecutionOrder_SortIsDeterministic_IdenticalData(t *testing.T) {
	tests := []struct {
		name    string
		entries []*scramblerEntry
	}{
		{
			name: "identical hashes",
			entries: []*scramblerEntry{
				{
					hash:   common.Hash{1},
					sender: common.Address{1},
					nonce:  1,
				},
				{
					hash:   common.Hash{1},
					sender: common.Address{1},
					nonce:  1,
				},
				{
					hash:   common.Hash{1},
					sender: common.Address{1},
					nonce:  1,
				},
			},
		},
		{
			name: "identical addresses",
			entries: []*scramblerEntry{
				{
					hash:   common.Hash{1},
					sender: common.Address{1},
					nonce:  1,
				},
				{
					hash:   common.Hash{2},
					sender: common.Address{1},
					nonce:  2,
				},
				{
					hash:   common.Hash{3},
					sender: common.Address{1},
					nonce:  3,
				},
			},
		},
		{
			name: "identical addresses and nonces",
			entries: []*scramblerEntry{
				{
					hash:   common.Hash{1},
					sender: common.Address{1},
					nonce:  1,
					gas:    1,
				},
				{
					hash:   common.Hash{2},
					sender: common.Address{1},
					nonce:  1,
					gas:    2,
				},
				{
					hash:   common.Hash{3},
					sender: common.Address{1},
					nonce:  1,
					gas:    3,
				},
			},
		},
		{
			name: "identical addresses, nonces and gas",
			entries: []*scramblerEntry{
				{
					hash:   common.Hash{1},
					sender: common.Address{1},
					nonce:  1,
					gas:    1,
				},
				{
					hash:   common.Hash{2},
					sender: common.Address{1},
					nonce:  1,
					gas:    1,
				},
				{
					hash:   common.Hash{3},
					sender: common.Address{1},
					nonce:  1,
					gas:    1,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res1 := test.entries
			res2 := deepCopyEntries(res1)
			// shuffle one array
			shuffleEntries(res2)

			res1 = getExecutionOrder(res1)
			res2 = getExecutionOrder(res2)
			if !reflect.DeepEqual(res1, res2) {
				t.Error("slices have different order - algorithm is not deterministic")
			}
		})
	}
}

func TestTxScrambler_GetExecutionOrder_SortIsDeterministic_RepeatedData(t *testing.T) {
	tests := []struct {
		name    string
		entries []*scramblerEntry
	}{
		{
			name: "repeated hashes",
			entries: []*scramblerEntry{
				{
					hash:   common.Hash{1},
					sender: common.Address{1},
					nonce:  1,
					gas:    1,
				},
				{
					hash:   common.Hash{2},
					sender: common.Address{2},
					nonce:  2,
					gas:    2,
				},
				{
					hash:   common.Hash{3},
					sender: common.Address{3},
					nonce:  3,
					gas:    3,
				},
				{
					hash:   common.Hash{2},
					sender: common.Address{2},
					nonce:  2,
					gas:    2,
				},
				{
					hash:   common.Hash{1},
					sender: common.Address{1},
					nonce:  1,
					gas:    1,
				},
			},
		},
		{
			name: "repeated addresses",
			entries: []*scramblerEntry{
				{
					hash:   common.Hash{1},
					sender: common.Address{1},
					nonce:  1,
				},
				{
					hash:   common.Hash{2},
					sender: common.Address{2},
					nonce:  2,
				},
				{
					hash:   common.Hash{3},
					sender: common.Address{3},
					nonce:  3,
				},
				{
					hash:   common.Hash{4},
					sender: common.Address{2},
					nonce:  4,
				},
				{
					hash:   common.Hash{5},
					sender: common.Address{1},
					nonce:  5,
				},
			},
		},
		{
			name: "repeated addresses and nonces",
			entries: []*scramblerEntry{
				{
					hash:   common.Hash{1},
					sender: common.Address{1},
					nonce:  1,
					gas:    1,
				},
				{
					hash:   common.Hash{2},
					sender: common.Address{2},
					nonce:  2,
					gas:    2,
				},
				{
					hash:   common.Hash{3},
					sender: common.Address{3},
					nonce:  3,
					gas:    3,
				},
				{
					hash:   common.Hash{4},
					sender: common.Address{2},
					nonce:  2,
					gas:    4,
				},
				{
					hash:   common.Hash{5},
					sender: common.Address{1},
					nonce:  1,
					gas:    5,
				},
			},
		},
		{
			name: "repeated addresses, nonces and gas",
			entries: []*scramblerEntry{
				{
					hash:   common.Hash{1},
					sender: common.Address{1},
					nonce:  1,
					gas:    1,
				},
				{
					hash:   common.Hash{2},
					sender: common.Address{2},
					nonce:  2,
					gas:    2,
				},
				{
					hash:   common.Hash{3},
					sender: common.Address{3},
					nonce:  3,
					gas:    3,
				},
				{
					hash:   common.Hash{4},
					sender: common.Address{2},
					nonce:  2,
					gas:    4,
				},
				{
					hash:   common.Hash{5},
					sender: common.Address{1},
					nonce:  1,
					gas:    5,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res1 := test.entries
			res2 := deepCopyEntries(res1)
			// shuffle one array
			shuffleEntries(res2)

			res1 = getExecutionOrder(res1)
			res2 = getExecutionOrder(res2)
			if !reflect.DeepEqual(res1, res2) {
				t.Error("slices have different order - algorithm is not deterministic")
			}
		})
	}
}

func TestTxScrambler_GetExecutionOrder_SortRemovesDuplicateHashes(t *testing.T) {
	entries := []*scramblerEntry{
		{hash: common.Hash{1}},
		{hash: common.Hash{2}},
		{hash: common.Hash{3}},
		{hash: common.Hash{2}},
		{hash: common.Hash{1}},
	}
	shuffleEntries(entries)
	entries = getExecutionOrder(entries)

	checkDuplicateHashes(t, entries)
}

func TestTxScrambler_GetExecutionOrder_SortsSameSenderByNonceAndGas(t *testing.T) {
	entries := []*scramblerEntry{
		{
			hash:   common.Hash{1},
			sender: common.Address{1},
			nonce:  1,
		},
		{
			hash:   common.Hash{2},
			sender: common.Address{1},
			nonce:  2,
		},
		{
			hash:   common.Hash{3},
			sender: common.Address{1},
			nonce:  3,
			gas:    1,
		},
		{
			hash:   common.Hash{4},
			sender: common.Address{1},
			nonce:  3,
			gas:    2,
		},
		{
			hash:   common.Hash{5},
			sender: common.Address{1},
			nonce:  4,
		},
	}
	shuffleEntries(entries)
	entries = getExecutionOrder(entries)

	for i := 0; i < len(entries); i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[i].sender == entries[j].sender {
				if entries[i].nonce > entries[j].nonce {
					t.Errorf("incorrect nonce order %d must be before %d", entries[j].nonce, entries[i].nonce)
				}
				if entries[i].nonce == entries[j].nonce && entries[i].gas < entries[j].gas {
					t.Errorf("incorrect gas order %d must be before %d", entries[j].gas, entries[i].gas)
				}
			}
		}
	}
}

func TestTxScrambler_GetExecutionOrder_RandomInput(t *testing.T) {
	// this tests these input sizes:
	// 1, 4, 16, 64, 256, 1024
	for i := int64(1); i <= 1024; i = i * 4 {
		input := createRandomScramblerTestInput(i)
		cpy := deepCopyEntries(input)
		shuffleEntries(cpy)
		input = getExecutionOrder(input)
		cpy = getExecutionOrder(input)
		if !reflect.DeepEqual(input, cpy) {
			t.Error("slices have different order - algorithm is not deterministic")
		}
	}

}

func BenchmarkTxScrambler(b *testing.B) {
	for size := int64(10); size < 100_000; size *= 10 {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			for i := 1; i <= b.N; i++ {
				b.StopTimer()
				entries := createRandomScramblerTestInput(size)
				b.StartTimer()
				getExecutionOrder(entries)
			}
		})
	}
}

func createRandomScramblerTestInput(size int64) []*scramblerEntry {
	var entries []*scramblerEntry
	for i := int64(0); i < size; i++ {
		entries = append(entries, &scramblerEntry{
			hash:   common.Hash{byte(rand.Intn(1000 - 1))},
			sender: common.Address{byte(rand.Intn(100 - 1))},
			nonce:  uint64(rand.Intn(10 - 1)),
			gas:    uint64(rand.Intn(10 - 1)),
		})
	}

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
	seenHashes := make(map[common.Hash]struct{})
	for _, entry := range entries {
		if _, found := seenHashes[entry.hash]; found {
			t.Fatalf("found duplicate hash in entries: %s", entry.hash)
		}
		seenHashes[entry.hash] = struct{}{}
	}
}

func createRandomSalt() [32]byte {
	var salt = [32]byte{}
	for i := 0; i < 32; i++ {
		salt[i] = byte(rand.Intn(math.MaxUint8))
	}
	return salt
}
