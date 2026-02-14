# Steady-State Pattern

The steady-state pattern states that for every service that accumulates a
resource, some other mechanism must recycle that resource. Without active
cleanup, unbounded growth of logs, caches, temporary files, or connections will
eventually exhaust the system and cause failures.

The goal is to keep the system in a stable, predictable operating range without
human intervention.

## Implementation

Below is a generic `Purger` that periodically cleans up accumulated resources to
maintain a steady state.

```go
package steadystate

import (
	"log"
	"time"
)

// Resource represents an accumulating resource that can report its size and
// be purged.
type Resource interface {
	// Size returns the current amount of accumulated resources.
	Size() int64

	// Purge removes resources that are older than the given threshold.
	Purge(olderThan time.Duration) (purged int64, err error)
}

// Purger periodically checks a resource and purges entries that exceed the
// maximum age, keeping the system in a steady state.
type Purger struct {
	resource Resource
	maxAge   time.Duration
	interval time.Duration
	stop     chan struct{}
}

// NewPurger creates a purger that checks the resource at the given interval
// and removes entries older than maxAge.
func NewPurger(r Resource, maxAge, interval time.Duration) *Purger {
	return &Purger{
		resource: r,
		maxAge:   maxAge,
		interval: interval,
		stop:     make(chan struct{}),
	}
}

// Start begins the periodic purge loop in a background goroutine.
func (p *Purger) Start() {
	ticker := time.NewTicker(p.interval)

	go func() {
		for {
			select {
			case <-ticker.C:
				before := p.resource.Size()
				purged, err := p.resource.Purge(p.maxAge)
				if err != nil {
					log.Printf("purge error: %v", err)
					continue
				}
				log.Printf("purged %d items (before: %d, after: %d)",
					purged, before, before-purged)
			case <-p.stop:
				ticker.Stop()
				return
			}
		}
	}()
}

// Stop terminates the purge loop.
func (p *Purger) Stop() {
	close(p.stop)
}
```

## Usage

```go
// LogDir implements the steadystate.Resource interface for a log directory.
type LogDir struct {
	path string
}

func (d *LogDir) Size() int64 {
	entries, _ := os.ReadDir(d.path)
	return int64(len(entries))
}

func (d *LogDir) Purge(olderThan time.Duration) (int64, error) {
	entries, err := os.ReadDir(d.path)
	if err != nil {
		return 0, err
	}

	var purged int64
	cutoff := time.Now().Add(-olderThan)
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue
		}
		if info.ModTime().Before(cutoff) {
			os.Remove(filepath.Join(d.path, e.Name()))
			purged++
		}
	}
	return purged, nil
}

// Purge log files older than 7 days, checking every hour.
purger := steadystate.NewPurger(&LogDir{path: "/var/log/myapp"}, 7*24*time.Hour, 1*time.Hour)
purger.Start()
defer purger.Stop()
```

## Rules of Thumb

- Every accumulating resource (logs, temp files, cache entries, sessions) must have a corresponding cleanup mechanism.
- Prefer time-based purging over size-based purging — it is simpler and more predictable.
- Monitor the resource size over time. If it trends upward despite purging, the purge interval or threshold needs adjustment.
- Run purgers as background goroutines with graceful shutdown support to avoid data loss.
