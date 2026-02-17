// Package main demonstrates the Monitor pattern.
//
// A monitor combines a mutex with condition variables to protect shared
// state while allowing goroutines to wait for specific conditions.
package main

import (
	"fmt"
	"sync"
)

type BoundedBuffer struct {
	mu       sync.Mutex
	notFull  *sync.Cond
	notEmpty *sync.Cond
	buf      []int
	capacity int
}

func NewBoundedBuffer(capacity int) *BoundedBuffer {
	b := &BoundedBuffer{
		buf:      make([]int, 0, capacity),
		capacity: capacity,
	}
	b.notFull = sync.NewCond(&b.mu)
	b.notEmpty = sync.NewCond(&b.mu)
	return b
}

func (b *BoundedBuffer) Put(item int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for len(b.buf) == b.capacity {
		b.notFull.Wait()
	}
	b.buf = append(b.buf, item)
	b.notEmpty.Signal()
}

func (b *BoundedBuffer) Get() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	for len(b.buf) == 0 {
		b.notEmpty.Wait()
	}
	item := b.buf[0]
	b.buf = b.buf[1:]
	b.notFull.Signal()
	return item
}

func main() {
	buf := NewBoundedBuffer(3)
	var wg sync.WaitGroup

	// 2 producers, each producing 5 items
	for p := 0; p < 2; p++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < 5; i++ {
				val := id*100 + i
				buf.Put(val)
				fmt.Printf("Producer %d put: %d\n", id, val)
			}
		}(p)
	}

	// 1 consumer, consuming all 10 items
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			val := buf.Get()
			fmt.Printf("Consumer got: %d\n", val)
		}
	}()

	wg.Wait()
}
