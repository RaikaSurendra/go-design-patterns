// Package main demonstrates the Observer pattern.
//
// Observer allows a type instance to publish events to other instances
// (observers) who wish to be updated when a particular event occurs.
package main

import "fmt"

type Event struct {
	Data string
}

type Observer interface {
	OnNotify(Event)
}

type Notifier interface {
	Register(Observer)
	Deregister(Observer)
	Notify(Event)
}

type eventObserver struct {
	id int
}

func (o *eventObserver) OnNotify(e Event) {
	fmt.Printf("  Observer %d received: %s\n", o.id, e.Data)
}

type eventNotifier struct {
	observers map[Observer]struct{}
}

func newEventNotifier() *eventNotifier {
	return &eventNotifier{observers: map[Observer]struct{}{}}
}

func (n *eventNotifier) Register(o Observer) {
	n.observers[o] = struct{}{}
}

func (n *eventNotifier) Deregister(o Observer) {
	delete(n.observers, o)
}

func (n *eventNotifier) Notify(e Event) {
	for o := range n.observers {
		o.OnNotify(e)
	}
}

func main() {
	notifier := newEventNotifier()

	obs1 := &eventObserver{id: 1}
	obs2 := &eventObserver{id: 2}
	obs3 := &eventObserver{id: 3}

	notifier.Register(obs1)
	notifier.Register(obs2)
	notifier.Register(obs3)

	fmt.Println("Publishing 'user.login':")
	notifier.Notify(Event{Data: "user.login"})

	fmt.Println("\nDeregistering observer 2...")
	notifier.Deregister(obs2)

	fmt.Println("\nPublishing 'user.logout':")
	notifier.Notify(Event{Data: "user.logout"})
}
