package main

import (
	"fmt"
	"sync"
	"time"
)

// package singleton

type cacheEntry struct {
	value  string
	expiry *time.Time
}

func (e cacheEntry) isExpired() bool {
	return e.expiry != nil && time.Now().After(*e.expiry)
}

type CacheManager struct {
	mu    sync.Mutex
	cache map[string]cacheEntry
}

var cacheManager = &CacheManager{
	cache: make(map[string]cacheEntry),
}

func (c *CacheManager) Put(key, value string, ttlSeconds int64) {
	var expiry *time.Time

	if ttlSeconds > 0 {
		t := time.Now().Add(time.Duration(ttlSeconds) * time.Second)
		expiry = &t
	}

	c.mu.Lock()
	c.cache[key] = cacheEntry{value: value, expiry: expiry}
	c.mu.Unlock()
}

func (c *CacheManager) Get(key string) string {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.cache[key]
	if !ok {
		return ""
	}

	if entry.isExpired() {
		delete(c.cache, key)
		return ""
	}

	return entry.value
}

func (c *CacheManager) Remove(key string) {
	c.mu.Lock()
	delete(c.cache, key)
	c.mu.Unlock()
}

func (c *CacheManager) Size() int {
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, e := range c.cache {
		if e.expiry != nil && now.After(*e.expiry) {
			delete(c.cache, k)
		}
	}
	return len(c.cache)
}

func main() {
	// Both references point to the same CacheManager instance
	cache1 := cacheManager
	cache2 := cacheManager

	fmt.Println("Same instance?", cache1 == cache2) // true

	// Component A caches data
	cache1.Put("user:42", "{name: 'Alice'}", 5) // 5 second TTL
	cache1.Put("config:theme", "dark", 0)       // no expiry

	// Component B reads from the same cache
	fmt.Println("user:42 =", cache2.Get("user:42"))
	fmt.Println("config:theme =", cache2.Get("config:theme"))
	fmt.Println("Cache size:", cache2.Size())
}
