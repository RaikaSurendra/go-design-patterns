// Package main demonstrates the Read-Write Lock pattern.
//
// RWMutex allows multiple concurrent readers but only one writer,
// improving throughput when reads far outnumber writes.
package main

import (
	"fmt"
	"sync"
)

type Cache struct {
	mu    sync.RWMutex
	store map[string]string
}

func NewCache() *Cache {
	return &Cache{store: make(map[string]string)}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.store[key]
	return val, ok
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = value
}

func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.store)
}

func main() {
	c := NewCache()

	// Write some data
	c.Set("language", "Go")
	c.Set("pattern", "RWMutex")
	c.Set("version", "1.24")

	// Concurrent readers
	var wg sync.WaitGroup
	keys := []string{"language", "pattern", "version", "missing"}

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := keys[id%len(keys)]
			if val, ok := c.Get(key); ok {
				fmt.Printf("Reader %d: %s = %s\n", id, key, val)
			} else {
				fmt.Printf("Reader %d: %s not found\n", id, key)
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("Cache size: %d\n", c.Len())
}
