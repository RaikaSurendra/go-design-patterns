// Package main demonstrates the Bulkhead pattern.
//
// Bulkhead isolates resources into pools so that if one downstream
// service fails, it doesn't consume all resources and cascade.
package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrBulkheadFull = errors.New("bulkhead capacity full")

type Bulkhead struct {
	name    string
	sem     chan struct{}
	timeout time.Duration
}

func NewBulkhead(name string, capacity int, timeout time.Duration) *Bulkhead {
	return &Bulkhead{
		name:    name,
		sem:     make(chan struct{}, capacity),
		timeout: timeout,
	}
}

func (b *Bulkhead) Execute(fn func() error) error {
	select {
	case b.sem <- struct{}{}:
		defer func() { <-b.sem }()
		return fn()
	case <-time.After(b.timeout):
		return fmt.Errorf("%s: %w", b.name, ErrBulkheadFull)
	}
}

func main() {
	orderPool := NewBulkhead("orders", 3, 500*time.Millisecond)
	paymentPool := NewBulkhead("payments", 2, 500*time.Millisecond)

	var wg sync.WaitGroup

	// Simulate 5 order requests (pool of 3)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			err := orderPool.Execute(func() error {
				time.Sleep(200 * time.Millisecond)
				fmt.Printf("  Order %d processed\n", id)
				return nil
			})
			if err != nil {
				fmt.Printf("  Order %d rejected: %v\n", id, err)
			}
		}(i)
	}

	// Simulate 4 payment requests (pool of 2)
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			err := paymentPool.Execute(func() error {
				time.Sleep(200 * time.Millisecond)
				fmt.Printf("  Payment %d processed\n", id)
				return nil
			})
			if err != nil {
				fmt.Printf("  Payment %d rejected: %v\n", id, err)
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("Pools are isolated — payment failures don't affect orders")
}
