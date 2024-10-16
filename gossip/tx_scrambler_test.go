package gossip

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/exp/rand"
	"testing"
)

func TestGetExecutionOrder_SortIsDeterministic(t *testing.T) {
	const size = 10
	tests := []struct {
		name string
		sort scramblerSortFunc
	}{
		{
			"builtInSort",
			builtInSort,
		},
		{
			"quickSort",
			quickSort,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			entries := createCertainScramblerTestInput(size)
			cpy := make([]*scramblerEntry, len(entries))
			// make a deep copy
			for i, e := range entries {
				copied := *e
				cpy[i] = &copied
			}
			// sort two same arrays
			entries = getExecutionOrder(entries, test.sort)
			cpy = getExecutionOrder(cpy, test.sort)
			for index, _ := range entries {
				// first occurrence of changed order means algorithm is not deterministic
				if *entries[index] != *cpy[index] {
					t.Fatal("slices have different order - algorithm is not deterministic")
				}
			}
		})
	}
}

func TestGetExecutionOrder_SortRemovesDuplicateHashes(t *testing.T) {
	const size = 10
	tests := []struct {
		name string
		sort scramblerSortFunc
	}{
		{
			"builtInSort",
			builtInSort,
		},
		{
			"quickSort",
			quickSort,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			entries := createCertainScramblerTestInput(size)
			cpy := make([]*scramblerEntry, len(entries))
			// make a deep copy
			for i, e := range entries {
				copied := *e
				cpy[i] = &copied
			}
			// and append it
			entries = append(entries, cpy...)
			entries = getExecutionOrder(entries, test.sort)

			if got, want := len(entries), size; got != want {
				t.Errorf("duplicate entries does not seem to be removed, got size: %v, want size: %v", got, want)
			}
		})
	}
}

func createCertainScramblerTestInput(size int) []*scramblerEntry {
	var entries []*scramblerEntry
	for i := 0; i < size; i++ {
		entries = append(entries, &scramblerEntry{
			hash:   common.Hash{byte(i)},
			sender: common.Address{byte(i)},
			nonce:  uint64(rand.Intn(10 - 1)),
		})
	}

	return entries
}

func TestGetExecutionOrder_SortsSameSenderByNonce(t *testing.T) {
	sample := []*scramblerEntry{
		{
			hash:   common.Hash{1},
			sender: common.Address{1},
			nonce:  4,
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
			sender: common.Address{1},
			nonce:  2,
		},
		{
			hash:   common.Hash{5},
			sender: common.Address{1},
			nonce:  3,
		},
		{
			hash:   common.Hash{6},
			sender: common.Address{2},
			nonce:  2,
		},
		{
			hash:   common.Hash{7},
			sender: common.Address{1},
			nonce:  1,
		},
	}
	entries := getExecutionOrder(sample, builtInSort)

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

func BenchmarkSortAlgorithms_BuiltIn(b *testing.B) {
	const (
		numLoops   = 4
		multiplier = 10
	)
	size := 10

	for _ = range numLoops {
		b.Run(fmt.Sprintf("builtIn_%d", size), func(b *testing.B) {
			for i := 1; i <= b.N; i++ {
				b.StopTimer()
				entries := createRandomScramblerTestInput(size)
				b.StartTimer()
				getExecutionOrder(entries, builtInSort)
			}
		})
		size = size * multiplier
	}
}

func BenchmarkSortAlgorithms_QuickSort(b *testing.B) {
	const (
		numLoops   = 4
		multiplier = 10
	)
	size := 10

	for _ = range numLoops {
		b.Run(fmt.Sprintf("quickSort_%d", size), func(b *testing.B) {
			for i := 1; i <= b.N; i++ {
				b.StopTimer()
				entries := createRandomScramblerTestInput(size)
				b.StartTimer()
				getExecutionOrder(entries, quickSort)
			}
		})
		size = size * multiplier
	}
}

func createRandomScramblerTestInput(size int) []*scramblerEntry {
	var entries []*scramblerEntry
	for i := 0; i < size; i++ {
		entries = append(entries, &scramblerEntry{
			hash:   common.Hash{byte(rand.Intn(100 - 1))},
			sender: common.Address{byte(rand.Intn(10 - 1))},
			nonce:  uint64(rand.Intn(10 - 1)),
		})
	}

	return entries
}
