# Type-Safe Container <span class="gp-difficulty gp-difficulty--easy">Easy</span>

Go generics (1.18+) allow creating type-safe container data structures that
work with any type — no `interface{}` casts, no code generation. The compiler
enforces type safety at compile time while keeping the implementation reusable.

## Implementation

```go
package container

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

func (s *Stack[T]) Peek() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

func (s *Stack[T]) Len() int {
	return len(s.items)
}

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

func (q *Queue[T]) Len() int {
	return len(q.items)
}
```

## Usage

```go
// Type-safe stack of integers — no casting needed.
s := &container.Stack[int]{}
s.Push(1)
s.Push(2)
s.Push(3)
val, _ := s.Pop() // val is int (not interface{}), val == 3

// Type-safe queue of strings.
q := &container.Queue[string]{}
q.Enqueue("first")
q.Enqueue("second")
msg, _ := q.Dequeue() // msg is string, msg == "first"
```

## Rules of Thumb

- Use generics when the logic is truly type-independent (containers, algorithms). Don't use generics just to avoid writing two similar functions.
- Prefer `any` constraint for containers; use specific constraints (`comparable`, `constraints.Ordered`) when the logic requires equality or ordering.
- The `var zero T` pattern returns the zero value of a generic type — this is idiomatic for "not found" returns.
