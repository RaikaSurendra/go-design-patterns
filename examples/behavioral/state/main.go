// Package main demonstrates the State pattern.
//
// State allows an object to alter its behavior when its internal state
// changes. Each state is a separate type implementing a common interface.
package main

import "fmt"

// State defines behavior for a particular state of the vending machine.
type State interface {
	InsertCoin(v *Machine)
	SelectItem(v *Machine)
	Dispense(v *Machine)
}

// Machine is the context that holds the current state.
type Machine struct {
	current State
	idle    State
	coined  State
	sold    State
}

func NewMachine() *Machine {
	m := &Machine{}
	m.idle = &IdleState{}
	m.coined = &CoinedState{}
	m.sold = &SoldState{}
	m.current = m.idle
	return m
}

func (m *Machine) SetState(s State) { m.current = s }
func (m *Machine) InsertCoin()      { m.current.InsertCoin(m) }
func (m *Machine) SelectItem()      { m.current.SelectItem(m) }
func (m *Machine) Dispense()        { m.current.Dispense(m) }

// --- Concrete states ---

type IdleState struct{}

func (s *IdleState) InsertCoin(v *Machine) {
	fmt.Println("Coin inserted")
	v.SetState(v.coined)
}
func (s *IdleState) SelectItem(v *Machine) { fmt.Println("Insert coin first") }
func (s *IdleState) Dispense(v *Machine)   { fmt.Println("Insert coin first") }

type CoinedState struct{}

func (s *CoinedState) InsertCoin(v *Machine) { fmt.Println("Coin already inserted") }
func (s *CoinedState) SelectItem(v *Machine) {
	fmt.Println("Item selected")
	v.SetState(v.sold)
}
func (s *CoinedState) Dispense(v *Machine) { fmt.Println("Select an item first") }

type SoldState struct{}

func (s *SoldState) InsertCoin(v *Machine)  { fmt.Println("Wait, dispensing item") }
func (s *SoldState) SelectItem(v *Machine)  { fmt.Println("Wait, dispensing item") }
func (s *SoldState) Dispense(v *Machine) {
	fmt.Println("Item dispensed")
	v.SetState(v.idle)
}

func main() {
	m := NewMachine()

	fmt.Println("--- Attempt without coin ---")
	m.SelectItem()

	fmt.Println("\n--- Normal purchase ---")
	m.InsertCoin()
	m.InsertCoin() // already inserted
	m.SelectItem()
	m.Dispense()

	fmt.Println("\n--- Second purchase ---")
	m.InsertCoin()
	m.SelectItem()
	m.Dispense()
}
