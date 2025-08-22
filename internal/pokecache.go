package internal

import (
	"sync"
	"time"
)

type CacheEntry struct {
	CreatedAt time.Time
	Val       []byte
}

type Cache struct {
	Entries map[string]CacheEntry
	Mu      sync.Mutex
}

func (c *Cache) Add(key string, value []byte) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	c.Entries[key] = CacheEntry{
		CreatedAt: time.Now(),
		Val:       value,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	if value, ok := c.Entries[key]; ok {
		return value.Val, true
	} else {
		return nil, false
	}
}

func (c *Cache) reapLoop(interval time.Duration) {

	ticker := time.NewTicker(interval)

	for {
		// wait for the ticker signal after interval
		<-ticker.C

		// remove data older than interval
		now := time.Now()
		var keyToDelete []string

		// lock before access/delete
		c.Mu.Lock()

		for key, entry := range c.Entries {
			if entry.CreatedAt.Add(interval).Before(now) {
				keyToDelete = append(keyToDelete, key)
			}
		}

		for _, key := range keyToDelete {
			delete(c.Entries, key)
		}

		// unlock after all delete
		c.Mu.Unlock()
	}
}

func NewCache(interval time.Duration) *Cache {
	newCache := Cache{
		Entries: make(map[string]CacheEntry),
		Mu:      sync.Mutex{},
	}

	go newCache.reapLoop(interval)

	return &newCache
}
