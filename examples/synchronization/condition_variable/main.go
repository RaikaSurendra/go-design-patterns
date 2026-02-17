// Package main demonstrates the Condition Variable pattern.
//
// A condition variable lets goroutines wait for a condition to become true
// instead of busy-looping. Go provides sync.Cond for this purpose.
package main

import (
	"fmt"
	"sync"
	"time"
)

type BlockingQueue struct {
	mu    sync.Mutex
	cond  *sync.Cond
	items []int
}

func NewBlockingQueue() *BlockingQueue {
	q := &BlockingQueue{}
	q.cond = sync.NewCond(&q.mu)
	return q
}

func (q *BlockingQueue) Enqueue(item int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, item)
	q.cond.Signal()
}

func (q *BlockingQueue) Dequeue() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	for len(q.items) == 0 {
		q.cond.Wait()
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

func main() {
	q := NewBlockingQueue()

	// Producer
	go func() {
		for i := 1; i <= 5; i++ {
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("Produced: %d\n", i)
			q.Enqueue(i)
		}
	}()

	// Consumer — blocks until items arrive
	for i := 0; i < 5; i++ {
		item := q.Dequeue()
		fmt.Printf("Consumed: %d\n", item)
	}
}
