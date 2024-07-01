package swapmap

import (
	"errors"
	"github.com/Fantom-foundation/lachesis-base/kvdb"
)

var errClosed = errors.New("swapmap already closed")

type Callbacks[K comparable, V any] struct {
	SizeOf       func(k K, v V) int
	SerializeK   func(k K) []byte
	SerializeV   func(v V) []byte
	DeserializeV func(b []byte) V
	GetSwapDB    func() kvdb.Store
	Fatal        func(err error)
}

// Map spills data to disk when overflown, similarly to OS memory swap
type Map[K comparable, V any] struct {
	data       map[K]V
	memSize    int
	maxMemSize int

	callback Callbacks[K, V]

	swap kvdb.Store
}

// New initializes and returns a new instance of Map.
func New[K comparable, V any](c Callbacks[K, V], maxMemSize int) *Map[K, V] {
	return &Map[K, V]{
		data:       make(map[K]V),
		maxMemSize: maxMemSize,
		callback:   c,
	}
}

// Set adds or updates an element in the map.
func (sm *Map[K, V]) Set(key K, value V) {
	sm.set(key, value)
	sm.mayUnload()
}

func (sm *Map[K, V]) set(key K, value V) {
	sizeBefore := len(sm.data)
	sm.data[key] = value
	if len(sm.data) > sizeBefore {
		sm.memSize += sm.callback.SizeOf(key, value)
	}
}

func (sm *Map[K, V]) AddMap(mm map[K]V) {
	for k, v := range mm {
		sm.set(k, v)
	}
	sm.mayUnload()
}

func (sm *Map[K, V]) AddSlice(kk []K, vv []V) {
	for i := range kk {
		sm.set(kk[i], vv[i])
	}
	sm.mayUnload()
}

func (sm *Map[K, V]) mayUnload() {
	for sm.memSize > sm.maxMemSize {
		if sm.swap == nil {
			sm.swap = sm.callback.GetSwapDB()
		}
		err := sm.unload(kvdb.IdealBatchSize)
		if err != nil {
			sm.callback.Fatal(err)
			return
		}
	}
}

func (sm *Map[K, V]) unload(toUnload int) error {
	batch := sm.swap.NewBatch()
	defer batch.Reset()

	diff := 0
	for key, val := range sm.data {
		err := batch.Put(sm.callback.SerializeK(key), sm.callback.SerializeV(val))
		if err != nil {
			return err
		}

		delete(sm.data, key)
		diff += sm.callback.SizeOf(key, val)

		if batch.ValueSize() >= toUnload {
			break
		}
	}
	if diff <= sm.memSize {
		sm.memSize -= diff
	} else {
		sm.memSize = 0
	}

	err := batch.Write()
	if err != nil {
		return err
	}

	return nil
}

// Get retrieves an element from the map. Returns the value and a boolean indicating if the key was found.
func (sm *Map[K, V]) Get(key K) (V, bool) {
	value, found := sm.data[key]
	if !found && sm.swap != nil {
		valB, err := sm.swap.Get(sm.callback.SerializeK(key))
		if err != nil {
			sm.callback.Fatal(err)
			return value, false
		}
		if valB == nil {
			return value, false
		}
		return sm.callback.DeserializeV(valB), true
	}
	return value, found
}

func (sm *Map[K, V]) Has(key K) bool {
	_, found := sm.data[key]
	if !found && sm.swap != nil {
		ok, err := sm.swap.Has(sm.callback.SerializeK(key))
		if err != nil {
			sm.callback.Fatal(err)
			return false
		}
		return ok
	}
	return found
}

func (sm *Map[K, V]) Close() error {
	if sm.data == nil {
		return errClosed
	}
	sm.data = nil
	if sm.swap == nil {
		return nil
	}
	if err := sm.swap.Close(); err != nil {
		return err
	}
	sm.swap.Drop()
	sm.swap = nil
	return nil
}
