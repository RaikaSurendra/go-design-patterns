# State Pattern <span class="gp-difficulty gp-difficulty--medium">Medium</span>

The state pattern allows an object to alter its behavior when its internal state
changes. The object appears to change its type. Each state is represented as a
separate type implementing a common interface, and the context delegates
behavior to the current state object.

## Implementation

```go
package vending

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

func New() *Machine {
	m := &Machine{}
	m.idle = &IdleState{}
	m.coined = &CoinedState{}
	m.sold = &SoldState{}
	m.current = m.idle
	return m
}

func (m *Machine) SetState(s State)  { m.current = s }
func (m *Machine) InsertCoin()       { m.current.InsertCoin(m) }
func (m *Machine) SelectItem()       { m.current.SelectItem(m) }
func (m *Machine) Dispense()         { m.current.Dispense(m) }

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
```

## Usage

```go
m := vending.New()

m.SelectItem()  // Insert coin first
m.InsertCoin()  // Coin inserted
m.InsertCoin()  // Coin already inserted
m.SelectItem()  // Item selected
m.Dispense()    // Item dispensed
m.Dispense()    // Insert coin first
```

## Rules of Thumb

- Use state when an object's behavior depends on its state and it must change behavior at runtime based on that state.
- State eliminates large conditional statements (`switch` / `if-else` chains) that select behavior based on the current state.
- State objects can be shared across contexts if they carry no instance-specific data (they become flyweights).
