package cache

import (
	"sync"
	"time"
)

type Cache[K comparable, V any] interface {
	Get(key K) (V, bool)
	Set(key K, value V)
	Delete(key K)
	Clear()
}

type MemoryCache[K comparable, V any] struct {
	data            map[K]cacheItem[V]
	mutex           sync.RWMutex
	ttl             time.Duration
	maxItems        int
	cleanupInterval time.Duration
}

type cacheItem[V any] struct {
	value      V
	expiration time.Time
	lastAccess time.Time
}

func NewMemoryCache[K comparable, V any](ttl time.Duration, maxItems int) *MemoryCache[K, V] {
	cache := &MemoryCache[K, V]{
		data:            make(map[K]cacheItem[V]),
		ttl:             ttl,
		maxItems:        maxItems,
		cleanupInterval: ttl / 2,
	}

	go cache.cleanup()

	return cache
}

func (c *MemoryCache[K, V]) Get(key K) (V, bool) {
	c.mutex.RLock()
	item, exists := c.data[key]
	defer c.mutex.RUnlock()

	if !exists || time.Now().After(item.expiration) {
		var zero V
		return zero, false
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if item, stillExists := c.data[key]; stillExists {
		item.lastAccess = time.Now()
		c.data[key] = item
	}

	return item.value, true
}

func (c *MemoryCache[K, V]) Set(key K, value V) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if len(c.data) >= c.maxItems {
		c.evictLRU()
	}

	c.data[key] = cacheItem[V]{
		value:      value,
		expiration: time.Now().Add(c.ttl),
		lastAccess: time.Now(),
	}
}

func (c *MemoryCache[K, V]) Delete(key K) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, key)
}

func (c *MemoryCache[K, V]) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = make(map[K]cacheItem[V])
}

func (c *MemoryCache[K, V]) evictLRU() {
	var oldestKey K
	var oldestTime time.Time

	first := true
	for k, item := range c.data {
		if first || item.lastAccess.Before(oldestTime) {
			oldestKey = k
			oldestTime = item.lastAccess
			first = false
		}
	}

	if !first {
		delete(c.data, oldestKey)
	}
}

func (c *MemoryCache[K, V]) cleanup() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		for key, item := range c.data {
			if time.Now().After(item.expiration) {
				delete(c.data, key)
			}
		}

		c.mutex.Unlock()
	}
}
