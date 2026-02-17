// Package main demonstrates the N-Barrier pattern.
//
// A barrier blocks N goroutines until all have reached the barrier point,
// then releases them simultaneously. Useful for phased computations.
package main

import (
	"fmt"
	"sync"
	"time"
)

type Barrier struct {
	n     int
	count int
	mu    sync.Mutex
	cond  *sync.Cond
}

func NewBarrier(n int) *Barrier {
	b := &Barrier{n: n}
	b.cond = sync.NewCond(&b.mu)
	return b
}

func (b *Barrier) Wait() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.count++
	if b.count == b.n {
		b.count = 0
		b.cond.Broadcast()
		return
	}
	b.cond.Wait()
}

func main() {
	const workers = 4
	b := NewBarrier(workers)

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Phase 1
			time.Sleep(time.Duration(id*50) * time.Millisecond)
			fmt.Printf("Worker %d: phase 1 done\n", id)
			b.Wait()

			// Phase 2 — all workers start together
			fmt.Printf("Worker %d: phase 2 done\n", id)
			b.Wait()

			fmt.Printf("Worker %d: finished\n", id)
		}(i)
	}

	wg.Wait()
	fmt.Println("All phases complete")
}
