# Read-Write Lock Pattern

A read-write lock allows multiple goroutines to hold the lock simultaneously for
read operations, but only one goroutine can hold it for a write operation. This
improves throughput when reads are far more frequent than writes.

Go provides `sync.RWMutex` in the standard library.

## Implementation

```go
package cache

import "sync"

// Cache is a thread-safe key-value store that uses a read-write lock to allow
// concurrent reads while serializing writes.
type Cache struct {
	mu    sync.RWMutex
	store map[string]string
}

func New() *Cache {
	return &Cache{
		store: make(map[string]string),
	}
}

// Get reads a value. Multiple goroutines can call Get concurrently.
func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	val, ok := c.store[key]
	return val, ok
}

// Set writes a value. Only one goroutine can call Set at a time, and it
// blocks all readers until the write completes.
func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store[key] = value
}

// Delete removes a key.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.store, key)
}

// Len returns the number of items in the cache.
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.store)
}
```

## Usage

```go
c := cache.New()

// Writers
go c.Set("language", "Go")
go c.Set("pattern", "RWMutex")

// Concurrent readers — these do not block each other
var wg sync.WaitGroup
for i := 0; i < 100; i++ {
	wg.Add(1)
	go func() {
		defer wg.Done()
		if val, ok := c.Get("language"); ok {
			fmt.Println(val)
		}
	}()
}
wg.Wait()
```

## Rules of Thumb

- Use `RWMutex` only when reads significantly outnumber writes. If the ratio is roughly equal, a plain `Mutex` has less overhead.
- Never call `Lock()` (write) while already holding `RLock()` (read) in the same goroutine — this causes a deadlock.
- For simple key-value caches, `sync.Map` may be more convenient, but `RWMutex` gives you more control and is generally faster for known access patterns.
