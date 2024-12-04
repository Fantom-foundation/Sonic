package eventid

import (
	"sync"

	"github.com/Fantom-foundation/lachesis-base/ltypes"
)

type Cache struct {
	ids     map[ltypes.EventHash]bool
	mu      sync.RWMutex
	maxSize int
	epoch   ltypes.EpochID
}

func NewCache(maxSize int) *Cache {
	return &Cache{
		maxSize: maxSize,
	}
}

func (c *Cache) Reset(epoch ltypes.EpochID) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ids = make(map[ltypes.EventHash]bool)
	c.epoch = epoch
}

func (c *Cache) Has(id ltypes.EventHash) (has bool, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.ids == nil {
		return false, false
	}
	if c.epoch != id.Epoch() {
		return false, false
	}
	return c.ids[id], true
}

func (c *Cache) Add(id ltypes.EventHash) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.ids == nil {
		return false
	}
	if c.epoch != id.Epoch() {
		return false
	}
	if len(c.ids) >= c.maxSize {
		c.ids = nil
		return false
	}
	c.ids[id] = true
	return true
}

func (c *Cache) Remove(id ltypes.EventHash) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.ids == nil {
		return
	}
	delete(c.ids, id)
}
