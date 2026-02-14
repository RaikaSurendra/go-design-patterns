# Flyweight Pattern <span class="gp-difficulty gp-difficulty--medium">Medium</span>

The flyweight pattern minimizes memory usage by sharing as much data as possible
with similar objects. It separates object state into **intrinsic** (shared,
immutable) and **extrinsic** (unique, context-dependent) parts. Shared intrinsic
state is stored once and referenced by many objects.

## Implementation

```go
package flyweight

import "fmt"

// TreeType holds the intrinsic (shared) state.
type TreeType struct {
	Name    string
	Color   string
	Texture string
}

func (t *TreeType) Draw(x, y int) string {
	return fmt.Sprintf("Drawing '%s' tree (%s) at (%d, %d)", t.Name, t.Color, x, y)
}

// TreeFactory caches and reuses TreeType instances.
type TreeFactory struct {
	types map[string]*TreeType
}

func NewTreeFactory() *TreeFactory {
	return &TreeFactory{types: make(map[string]*TreeType)}
}

func (f *TreeFactory) GetTreeType(name, color, texture string) *TreeType {
	key := name + "_" + color + "_" + texture
	if t, ok := f.types[key]; ok {
		return t
	}

	t := &TreeType{Name: name, Color: color, Texture: texture}
	f.types[key] = t
	return t
}

// Tree holds the extrinsic (unique) state plus a reference to the shared type.
type Tree struct {
	X, Y     int
	TreeType *TreeType
}

func (t *Tree) Draw() string {
	return t.TreeType.Draw(t.X, t.Y)
}
```

## Usage

```go
factory := flyweight.NewTreeFactory()

// Thousands of trees share only a few TreeType instances.
trees := []flyweight.Tree{
	{X: 1, Y: 2, TreeType: factory.GetTreeType("Oak", "green", "rough")},
	{X: 5, Y: 3, TreeType: factory.GetTreeType("Oak", "green", "rough")},    // reuses same TreeType
	{X: 8, Y: 1, TreeType: factory.GetTreeType("Pine", "dark green", "smooth")},
	{X: 3, Y: 7, TreeType: factory.GetTreeType("Oak", "green", "rough")},    // reuses same TreeType
}

for _, t := range trees {
	fmt.Println(t.Draw())
}
// Drawing 'Oak' tree (green) at (1, 2)
// Drawing 'Oak' tree (green) at (5, 3)
// Drawing 'Pine' tree (dark green) at (8, 1)
// Drawing 'Oak' tree (green) at (3, 7)

// Only 2 TreeType objects created despite 4 Tree instances.
```

## Rules of Thumb

- Use flyweight when the application creates a large number of objects that share most of their state.
- The shared (intrinsic) state must be immutable — if it changes, all references see the change.
- Go's `string` type is already a flyweight: identical string literals share the same underlying memory.
