# Command Pattern

The command pattern encapsulates a request as an object, allowing you to
parameterize clients with different requests, queue or log requests, and
support undoable operations. It decouples the invoker (who triggers the action)
from the receiver (who performs it).

## Implementation

```go
package command

// Command is the interface all commands implement.
type Command interface {
	Execute() string
	Undo() string
}

// Receiver is the object that performs the actual work.
type Light struct {
	IsOn bool
}

// --- Concrete commands ---

type TurnOnCommand struct {
	Light *Light
}

func (c *TurnOnCommand) Execute() string {
	c.Light.IsOn = true
	return "light turned on"
}

func (c *TurnOnCommand) Undo() string {
	c.Light.IsOn = false
	return "light turned off (undo)"
}

type TurnOffCommand struct {
	Light *Light
}

func (c *TurnOffCommand) Execute() string {
	c.Light.IsOn = false
	return "light turned off"
}

func (c *TurnOffCommand) Undo() string {
	c.Light.IsOn = true
	return "light turned on (undo)"
}

// Invoker stores and executes commands.
type RemoteControl struct {
	history []Command
}

func (r *RemoteControl) Press(cmd Command) string {
	r.history = append(r.history, cmd)
	return cmd.Execute()
}

func (r *RemoteControl) UndoLast() string {
	if len(r.history) == 0 {
		return "nothing to undo"
	}
	last := r.history[len(r.history)-1]
	r.history = r.history[:len(r.history)-1]
	return last.Undo()
}
```

## Usage

```go
light := &command.Light{}
remote := &command.RemoteControl{}

on := &command.TurnOnCommand{Light: light}
off := &command.TurnOffCommand{Light: light}

fmt.Println(remote.Press(on))   // light turned on
fmt.Println(remote.Press(off))  // light turned off
fmt.Println(remote.UndoLast())  // light turned on (undo)
```

## Rules of Thumb

- Use command when you need undo/redo, request queuing, or transaction logging.
- Commands can be serialized and sent over the network for distributed execution.
- In Go, simple commands can be represented as `func()` closures rather than full interface implementations.
