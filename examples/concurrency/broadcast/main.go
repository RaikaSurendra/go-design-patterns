// Package main demonstrates the Broadcast pattern.
//
// Broadcast delivers the same message to all registered listeners
// simultaneously, unlike fan-out which distributes different items.
package main

import (
	"fmt"
	"sync"
)

type Broadcaster[T any] struct {
	mu        sync.RWMutex
	listeners []chan T
}

func NewBroadcaster[T any]() *Broadcaster[T] {
	return &Broadcaster[T]{}
}

func (b *Broadcaster[T]) Register() <-chan T {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan T, 1)
	b.listeners = append(b.listeners, ch)
	return ch
}

func (b *Broadcaster[T]) Send(msg T) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, ch := range b.listeners {
		ch <- msg
	}
}

func (b *Broadcaster[T]) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, ch := range b.listeners {
		close(ch)
	}
	b.listeners = nil
}

func main() {
	bc := NewBroadcaster[string]()

	ch1 := bc.Register()
	ch2 := bc.Register()
	ch3 := bc.Register()

	bc.Send("hello everyone")

	fmt.Printf("Listener 1: %s\n", <-ch1)
	fmt.Printf("Listener 2: %s\n", <-ch2)
	fmt.Printf("Listener 3: %s\n", <-ch3)

	bc.Send("second message")

	fmt.Printf("Listener 1: %s\n", <-ch1)
	fmt.Printf("Listener 2: %s\n", <-ch2)
	fmt.Printf("Listener 3: %s\n", <-ch3)

	bc.Close()
	fmt.Println("Broadcaster closed")
}
