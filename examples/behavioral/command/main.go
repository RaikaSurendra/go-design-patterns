// Package main demonstrates the Command pattern.
//
// Command encapsulates a request as an object, supporting undo,
// queuing, and decoupling invoker from receiver.
package main

import "fmt"

// Command is the interface all commands implement.
type Command interface {
	Execute() string
	Undo() string
}

// Light is the receiver that performs actual work.
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

// RemoteControl is the invoker that stores and executes commands.
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

func main() {
	light := &Light{}
	remote := &RemoteControl{}

	on := &TurnOnCommand{Light: light}
	off := &TurnOffCommand{Light: light}

	fmt.Println(remote.Press(on))
	fmt.Println(remote.Press(off))
	fmt.Println(remote.UndoLast())
	fmt.Printf("Light is on: %v\n", light.IsOn)
}
