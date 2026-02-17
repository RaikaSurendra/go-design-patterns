// Package main demonstrates the Strategy pattern.
//
// Strategy enables an algorithm's behavior to be selected at runtime
// by encapsulating algorithms behind a common interface.
package main

import "fmt"

type Operator interface {
	Apply(int, int) int
}

type Operation struct {
	Operator Operator
}

func (o *Operation) Operate(leftValue, rightValue int) int {
	return o.Operator.Apply(leftValue, rightValue)
}

// --- Concrete strategies ---

type Addition struct{}

func (Addition) Apply(lval, rval int) int { return lval + rval }

type Multiplication struct{}

func (Multiplication) Apply(lval, rval int) int { return lval * rval }

type Subtraction struct{}

func (Subtraction) Apply(lval, rval int) int { return lval - rval }

func main() {
	a, b := 10, 3

	strategies := []struct {
		name string
		op   Operator
	}{
		{"Addition", Addition{}},
		{"Subtraction", Subtraction{}},
		{"Multiplication", Multiplication{}},
	}

	for _, s := range strategies {
		op := &Operation{Operator: s.op}
		fmt.Printf("%-15s: %d op %d = %d\n", s.name, a, b, op.Operate(a, b))
	}
}
