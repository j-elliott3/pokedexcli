package pokecache

import (
	"time"
	"sync"
)

type Cache struct {
	mu 			sync.Mutex
	interval 	time.Duration
	contents 	map[string]cacheEntry
}

type cacheEntry struct {
	createdAt 	time.Time
	val 		[]byte
}

func NewCache(interval time.Duration) *Cache{
	newCache := &Cache{
		interval: 	interval,
		contents: 	make(map[string]cacheEntry),
	}
	go newCache.reapLoop()
	return newCache
}

func (cache *Cache)Add(key string, val []byte) {
	cache.mu.Lock()
	newEntry := cacheEntry{
		createdAt: time.Now(),
		val: val,
	}
	cache.contents[key] = newEntry
	cache.mu.Unlock()
}

func (cache *Cache)Get(key string) ([]byte, bool) {
	cache.mu.Lock()
	entry, ok := cache.contents[key]
	if !ok {
		return []byte{}, false
	}
	defer cache.mu.Unlock()
	return entry.val, true
}

func (cache *Cache)reapLoop() {
	ticker := time.NewTicker(cache.interval)
	for range ticker.C {
		cache.mu.Lock()
		for key, entry := range cache.contents {
			if time.Since(entry.createdAt) > cache.interval {
				delete(cache.contents, key)
			}
		}
		cache.mu.Unlock()
	}
}