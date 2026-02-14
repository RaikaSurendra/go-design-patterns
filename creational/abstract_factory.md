# Abstract Factory Pattern <span class="gp-difficulty gp-difficulty--medium">Medium</span>

The abstract factory pattern provides an interface for creating families of
related objects without specifying their concrete types. The client code works
with factories and products through abstract interfaces, making it independent
of the actual types being created.

## Implementation

```go
package factory

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

func (b *WindowsButton) Paint() string { return "Windows button" }

type WindowsCheckbox struct{}

func (c *WindowsCheckbox) Check() string { return "Windows checkbox" }

type WindowsFactory struct{}

func (f *WindowsFactory) CreateButton() Button     { return &WindowsButton{} }
func (f *WindowsFactory) CreateCheckbox() Checkbox  { return &WindowsCheckbox{} }

// --- Mac family ---

type MacButton struct{}

func (b *MacButton) Paint() string { return "Mac button" }

type MacCheckbox struct{}

func (c *MacCheckbox) Check() string { return "Mac checkbox" }

type MacFactory struct{}

func (f *MacFactory) CreateButton() Button     { return &MacButton{} }
func (f *MacFactory) CreateCheckbox() Checkbox  { return &MacCheckbox{} }
```

## Usage

```go
func BuildUI(f factory.GUIFactory) {
	button := f.CreateButton()
	checkbox := f.CreateCheckbox()

	fmt.Println(button.Paint())
	fmt.Println(checkbox.Check())
}

// The client code is decoupled from concrete product types.
BuildUI(&factory.WindowsFactory{})
// Windows button
// Windows checkbox

BuildUI(&factory.MacFactory{})
// Mac button
// Mac checkbox
```

## Rules of Thumb

- Use abstract factory when the system needs to be independent of how its products are created and composed.
- Abstract factory is often implemented using factory methods under the hood.
- If you only have one product family, you likely only need the simpler factory method pattern.
