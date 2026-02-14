# Memento Pattern <span class="gp-difficulty gp-difficulty--medium">Medium</span>

The memento pattern captures and externalizes an object's internal state so that
the object can be restored to this state later, without violating encapsulation.
It is commonly used for implementing undo functionality.

## Implementation

```go
package memento

// Memento stores a snapshot of the editor's state.
type Memento struct {
	content string
}

// Editor is the originator whose state we want to save and restore.
type Editor struct {
	content string
}

func (e *Editor) Type(text string) {
	e.content += text
}

func (e *Editor) Content() string {
	return e.content
}

func (e *Editor) Save() *Memento {
	return &Memento{content: e.content}
}

func (e *Editor) Restore(m *Memento) {
	e.content = m.content
}

// History is the caretaker that stores mementos.
type History struct {
	snapshots []*Memento
}

func (h *History) Push(m *Memento) {
	h.snapshots = append(h.snapshots, m)
}

func (h *History) Pop() *Memento {
	if len(h.snapshots) == 0 {
		return nil
	}
	last := h.snapshots[len(h.snapshots)-1]
	h.snapshots = h.snapshots[:len(h.snapshots)-1]
	return last
}
```

## Usage

```go
editor := &memento.Editor{}
history := &memento.History{}

editor.Type("Hello, ")
history.Push(editor.Save())

editor.Type("World!")
history.Push(editor.Save())

editor.Type(" Extra text.")
fmt.Println(editor.Content())
// Hello, World! Extra text.

// Undo last change
editor.Restore(history.Pop())
fmt.Println(editor.Content())
// Hello, World!

// Undo again
editor.Restore(history.Pop())
fmt.Println(editor.Content())
// Hello,
```

## Rules of Thumb

- The caretaker (history) should never inspect or modify the memento's contents — it is an opaque token.
- Use memento when a direct interface to obtain the object's state would expose implementation details.
- Be mindful of memory usage: storing frequent snapshots of large objects can be costly. Consider diffing or compressing snapshots if needed.
