// Package main demonstrates the Semaphore pattern.
//
// A semaphore limits concurrent access to a finite number of resources
// using a buffered channel as the token pool.
package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrNoTickets      = errors.New("semaphore: could not acquire semaphore")
	ErrIllegalRelease = errors.New("semaphore: can't release without acquiring first")
)

type Semaphore struct {
	sem     chan struct{}
	timeout time.Duration
}

func NewSemaphore(tickets int, timeout time.Duration) *Semaphore {
	return &Semaphore{
		sem:     make(chan struct{}, tickets),
		timeout: timeout,
	}
}

func (s *Semaphore) Acquire() error {
	select {
	case s.sem <- struct{}{}:
		return nil
	case <-time.After(s.timeout):
		return ErrNoTickets
	}
}

func (s *Semaphore) Release() error {
	select {
	case <-s.sem:
		return nil
	case <-time.After(s.timeout):
		return ErrIllegalRelease
	}
}

func main() {
	sem := NewSemaphore(3, 2*time.Second)

	var wg sync.WaitGroup
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if err := sem.Acquire(); err != nil {
				fmt.Printf("Worker %d: %s\n", id, err)
				return
			}
			fmt.Printf("Worker %d: acquired (working...)\n", id)
			time.Sleep(200 * time.Millisecond)
			sem.Release()
			fmt.Printf("Worker %d: released\n", id)
		}(i)
	}

	wg.Wait()
	fmt.Println("All workers done")
}
