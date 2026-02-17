// Package main demonstrates the Fail-Fast pattern.
//
// Fail-fast checks preconditions before executing expensive work.
// If requirements aren't met, it returns immediately with a clear error.
package main

import (
	"errors"
	"fmt"
)

type Checker func() error
type Handler func() error

var (
	ErrDBDown    = errors.New("database is unavailable")
	ErrCacheDown = errors.New("cache is unavailable")
)

func FailFast(handler Handler, checks ...Checker) error {
	for _, check := range checks {
		if err := check(); err != nil {
			return err
		}
	}
	return handler()
}

func main() {
	dbAlive := true
	cacheAlive := true

	checkDB := func() error {
		if !dbAlive {
			return ErrDBDown
		}
		return nil
	}
	checkCache := func() error {
		if !cacheAlive {
			return ErrCacheDown
		}
		return nil
	}
	processOrder := func() error {
		fmt.Println("  Processing order (expensive work)")
		return nil
	}

	// All checks pass
	fmt.Println("Scenario 1: Everything healthy")
	err := FailFast(processOrder, checkDB, checkCache)
	fmt.Printf("  Result: %v\n", err)

	// DB is down — fail fast, skip expensive work
	fmt.Println("\nScenario 2: DB is down")
	dbAlive = false
	err = FailFast(processOrder, checkDB, checkCache)
	fmt.Printf("  Result: %v\n", err)

	// DB recovered, cache is down
	fmt.Println("\nScenario 3: Cache is down")
	dbAlive = true
	cacheAlive = false
	err = FailFast(processOrder, checkDB, checkCache)
	fmt.Printf("  Result: %v\n", err)
}
