// Package main demonstrates the Dependency Injection pattern.
//
// Dependencies are passed as interfaces via constructors, making code
// testable and decoupled from concrete implementations.
package main

import "fmt"

// Repository is the dependency interface.
type Repository interface {
	Save(id, data string) error
	FindByID(id string) (string, error)
}

// Notifier is another dependency.
type Notifier interface {
	Notify(userID, message string) error
}

// Service has its dependencies injected via the constructor.
type Service struct {
	repo     Repository
	notifier Notifier
}

func NewService(repo Repository, notifier Notifier) *Service {
	return &Service{repo: repo, notifier: notifier}
}

func (s *Service) CreateOrder(id, data string) error {
	if err := s.repo.Save(id, data); err != nil {
		return fmt.Errorf("save: %w", err)
	}
	return s.notifier.Notify(id, "Order created: "+data)
}

// --- Concrete implementations ---

type InMemoryRepo struct {
	store map[string]string
}

func (r *InMemoryRepo) Save(id, data string) error {
	r.store[id] = data
	fmt.Printf("  [repo] saved %s = %s\n", id, data)
	return nil
}

func (r *InMemoryRepo) FindByID(id string) (string, error) {
	data, ok := r.store[id]
	if !ok {
		return "", fmt.Errorf("not found: %s", id)
	}
	return data, nil
}

type ConsoleNotifier struct{}

func (n *ConsoleNotifier) Notify(userID, message string) error {
	fmt.Printf("  [notify] %s: %s\n", userID, message)
	return nil
}

func main() {
	// Wire up with real implementations
	repo := &InMemoryRepo{store: make(map[string]string)}
	notifier := &ConsoleNotifier{}
	svc := NewService(repo, notifier)

	fmt.Println("Creating orders:")
	svc.CreateOrder("order-1", "Widget x3")
	svc.CreateOrder("order-2", "Gadget x1")

	fmt.Println("\nLooking up order:")
	data, err := repo.FindByID("order-1")
	if err != nil {
		fmt.Printf("  error: %v\n", err)
	} else {
		fmt.Printf("  found: %s\n", data)
	}
}
