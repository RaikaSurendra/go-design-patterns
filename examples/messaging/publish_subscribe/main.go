// Package main demonstrates the Publish/Subscribe pattern.
//
// An event bus decouples publishers from subscribers. Publishers send
// messages to topics, and all subscribers to that topic receive them.
package main

import (
	"fmt"
	"sync"
)

type EventBus struct {
	mu          sync.RWMutex
	subscribers map[string][]chan string
}

func NewEventBus() *EventBus {
	return &EventBus{subscribers: make(map[string][]chan string)}
}

func (eb *EventBus) Subscribe(topic string) <-chan string {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	ch := make(chan string, 8)
	eb.subscribers[topic] = append(eb.subscribers[topic], ch)
	return ch
}

func (eb *EventBus) Publish(topic, message string) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	for _, ch := range eb.subscribers[topic] {
		ch <- message
	}
}

func (eb *EventBus) Close(topic string) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	for _, ch := range eb.subscribers[topic] {
		close(ch)
	}
	delete(eb.subscribers, topic)
}

func main() {
	bus := NewEventBus()

	// Two subscribers on "orders" topic
	sub1 := bus.Subscribe("orders")
	sub2 := bus.Subscribe("orders")

	// One subscriber on "logs" topic
	sub3 := bus.Subscribe("logs")

	// Publish
	bus.Publish("orders", "order-123 created")
	bus.Publish("orders", "order-456 created")
	bus.Publish("logs", "system started")

	// Consume
	fmt.Println("Subscriber 1 (orders):", <-sub1)
	fmt.Println("Subscriber 1 (orders):", <-sub1)
	fmt.Println("Subscriber 2 (orders):", <-sub2)
	fmt.Println("Subscriber 2 (orders):", <-sub2)
	fmt.Println("Subscriber 3 (logs):  ", <-sub3)

	bus.Close("orders")
	bus.Close("logs")
}
