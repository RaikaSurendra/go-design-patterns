// Package main demonstrates the Singleton pattern.
//
// Singleton restricts instantiation of a type to a single object,
// using sync.Once for thread-safe lazy initialization.
package main

import (
	"fmt"
	"sync"
)

type registry struct {
	data map[string]string
}

var (
	once     sync.Once
	instance *registry
)

func GetRegistry() *registry {
	once.Do(func() {
		fmt.Println("Initializing singleton registry...")
		instance = &registry{data: make(map[string]string)}
	})
	return instance
}

func main() {
	// First call initializes the singleton.
	r1 := GetRegistry()
	r1.data["service.url"] = "https://api.example.com"
	r1.data["service.timeout"] = "30s"

	// Second call returns the same instance.
	r2 := GetRegistry()

	fmt.Printf("Same instance: %v\n", r1 == r2)
	fmt.Printf("Value from r2: service.url = %s\n", r2.data["service.url"])
	fmt.Printf("Value from r2: service.timeout = %s\n", r2.data["service.timeout"])
}
