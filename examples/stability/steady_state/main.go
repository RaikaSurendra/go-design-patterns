// Package main demonstrates the Steady-State pattern.
//
// For every resource that accumulates, there must be a mechanism to
// recycle it. This example shows a cache with automatic TTL-based purging.
package main

import (
	"fmt"
	"sync"
	"time"
)

type CacheEntry struct {
	Value     string
	CreatedAt time.Time
}

type PurgingCache struct {
	mu      sync.RWMutex
	entries map[string]CacheEntry
	maxAge  time.Duration
	stop    chan struct{}
}

func NewPurgingCache(maxAge, purgeInterval time.Duration) *PurgingCache {
	c := &PurgingCache{
		entries: make(map[string]CacheEntry),
		maxAge:  maxAge,
		stop:    make(chan struct{}),
	}

	// Background purger keeps the cache in steady state
	go func() {
		ticker := time.NewTicker(purgeInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.purge()
			case <-c.stop:
				return
			}
		}
	}()

	return c
}

func (c *PurgingCache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = CacheEntry{Value: value, CreatedAt: time.Now()}
}

func (c *PurgingCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[key]
	if !ok {
		return "", false
	}
	return e.Value, true
}

func (c *PurgingCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}

func (c *PurgingCache) purge() {
	c.mu.Lock()
	defer c.mu.Unlock()
	cutoff := time.Now().Add(-c.maxAge)
	purged := 0
	for k, e := range c.entries {
		if e.CreatedAt.Before(cutoff) {
			delete(c.entries, k)
			purged++
		}
	}
	if purged > 0 {
		fmt.Printf("  [purger] removed %d stale entries\n", purged)
	}
}

func (c *PurgingCache) Stop() {
	close(c.stop)
}

func main() {
	// Cache with 300ms TTL, purging every 200ms
	cache := NewPurgingCache(300*time.Millisecond, 200*time.Millisecond)
	defer cache.Stop()

	// Add entries
	cache.Set("a", "alpha")
	cache.Set("b", "bravo")
	cache.Set("c", "charlie")
	fmt.Printf("Initial size: %d\n", cache.Size())

	// Wait for entries to expire and get purged
	time.Sleep(600 * time.Millisecond)
	fmt.Printf("After purge: %d\n", cache.Size())

	// Add new entries — cache stays in steady state
	cache.Set("d", "delta")
	fmt.Printf("After adding new entry: %d\n", cache.Size())
}
