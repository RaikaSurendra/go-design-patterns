// Package main demonstrates the Type-Safe Container pattern using generics.
//
// Go generics (1.18+) allow creating type-safe data structures that
// work with any type — no interface{} casts needed.
package main

import "fmt"

// Stack is a generic LIFO container.
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, true
}

func (s *Stack[T]) Len() int { return len(s.items) }

// Queue is a generic FIFO container.
type Queue[T any] struct {
	items []T
}

func (q *Queue[T]) Enqueue(item T) {
	q.items = append(q.items, item)
}

func (q *Queue[T]) Dequeue() (T, bool) {
	if len(q.items) == 0 {
		var zero T
		return zero, false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

func (q *Queue[T]) Len() int { return len(q.items) }

func main() {
	// Type-safe int stack
	fmt.Println("--- Stack[int] ---")
	s := &Stack[int]{}
	s.Push(10)
	s.Push(20)
	s.Push(30)

	for s.Len() > 0 {
		val, _ := s.Pop()
		fmt.Printf("  popped: %d\n", val)
	}

	// Type-safe string queue
	fmt.Println("\n--- Queue[string] ---")
	q := &Queue[string]{}
	q.Enqueue("first")
	q.Enqueue("second")
	q.Enqueue("third")

	for q.Len() > 0 {
		val, _ := q.Dequeue()
		fmt.Printf("  dequeued: %s\n", val)
	}
}
