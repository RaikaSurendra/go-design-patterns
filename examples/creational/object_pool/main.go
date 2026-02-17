// Package main demonstrates the Object Pool pattern.
//
// Object pool prepares and keeps multiple instances ready for use,
// avoiding expensive initialization on each request.
package main

import "fmt"

type Object struct {
	ID int
}

func (o *Object) Do(task string) {
	fmt.Printf("  Object %d processing: %s\n", o.ID, task)
}

type Pool chan *Object

func NewPool(total int) Pool {
	p := make(Pool, total)
	for i := 0; i < total; i++ {
		p <- &Object{ID: i + 1}
	}
	return p
}

func main() {
	pool := NewPool(3)
	fmt.Printf("Pool created with %d objects\n\n", len(pool))

	tasks := []string{"parse config", "validate input", "transform data", "send response"}

	for _, task := range tasks {
		select {
		case obj := <-pool:
			fmt.Printf("Acquired object from pool (%d remaining):\n", len(pool))
			obj.Do(task)
			pool <- obj
			fmt.Printf("  Returned object %d to pool (%d available)\n\n", obj.ID, len(pool))
		default:
			fmt.Printf("No objects available for: %s (pool exhausted)\n\n", task)
		}
	}
}
