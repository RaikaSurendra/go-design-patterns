# Composite Pattern

The composite pattern composes objects into tree structures to represent
part-whole hierarchies. It allows clients to treat individual objects and
compositions of objects uniformly through a common interface.

## Implementation

```go
package composite

import "fmt"

// Component is the common interface for leaf and composite nodes.
type Component interface {
	Search(keyword string)
}

// File is a leaf node.
type File struct {
	Name string
}

func (f *File) Search(keyword string) {
	fmt.Printf("Searching for '%s' in file: %s\n", keyword, f.Name)
}

// Folder is a composite node that can contain other components.
type Folder struct {
	Name       string
	Components []Component
}

func (f *Folder) Search(keyword string) {
	fmt.Printf("Searching for '%s' in folder: %s\n", keyword, f.Name)
	for _, c := range f.Components {
		c.Search(keyword)
	}
}

func (f *Folder) Add(c Component) {
	f.Components = append(f.Components, c)
}
```

## Usage

```go
file1 := &composite.File{Name: "main.go"}
file2 := &composite.File{Name: "utils.go"}
file3 := &composite.File{Name: "readme.md"}

src := &composite.Folder{Name: "src"}
src.Add(file1)
src.Add(file2)

root := &composite.Folder{Name: "project"}
root.Add(src)
root.Add(file3)

// Treats files and folders uniformly.
root.Search("pattern")
// Searching for 'pattern' in folder: project
// Searching for 'pattern' in folder: src
// Searching for 'pattern' in file: main.go
// Searching for 'pattern' in file: utils.go
// Searching for 'pattern' in file: readme.md
```

## Rules of Thumb

- Use composite when you want clients to ignore the difference between compositions of objects and individual objects.
- The composite pattern trades the type-safety of individual leaf operations for uniformity — you can call any operation on any node.
- Composite and decorator have similar structure diagrams but different intents: composite groups objects, decorator adds behavior.
