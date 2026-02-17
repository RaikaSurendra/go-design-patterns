// Package main demonstrates the Abstract Factory pattern.
//
// The abstract factory provides an interface for creating families of related
// objects without specifying their concrete types.
package main

import "fmt"

// Button and Checkbox are the abstract product interfaces.
type Button interface {
	Paint() string
}

type Checkbox interface {
	Check() string
}

// GUIFactory is the abstract factory interface.
type GUIFactory interface {
	CreateButton() Button
	CreateCheckbox() Checkbox
}

// --- Windows family ---

type WindowsButton struct{}

func (b *WindowsButton) Paint() string { return "Rendering Windows button" }

type WindowsCheckbox struct{}

func (c *WindowsCheckbox) Check() string { return "Windows checkbox toggled" }

type WindowsFactory struct{}

func (f *WindowsFactory) CreateButton() Button    { return &WindowsButton{} }
func (f *WindowsFactory) CreateCheckbox() Checkbox { return &WindowsCheckbox{} }

// --- Mac family ---

type MacButton struct{}

func (b *MacButton) Paint() string { return "Rendering Mac button" }

type MacCheckbox struct{}

func (c *MacCheckbox) Check() string { return "Mac checkbox toggled" }

type MacFactory struct{}

func (f *MacFactory) CreateButton() Button    { return &MacButton{} }
func (f *MacFactory) CreateCheckbox() Checkbox { return &MacCheckbox{} }

// BuildUI assembles a UI from the given factory, decoupled from concrete types.
func BuildUI(name string, f GUIFactory) {
	button := f.CreateButton()
	checkbox := f.CreateCheckbox()

	fmt.Printf("--- %s UI ---\n", name)
	fmt.Println(button.Paint())
	fmt.Println(checkbox.Check())
}

func main() {
	BuildUI("Windows", &WindowsFactory{})
	fmt.Println()
	BuildUI("Mac", &MacFactory{})
}
