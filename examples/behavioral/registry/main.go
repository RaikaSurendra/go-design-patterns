// Package main demonstrates the Registry pattern.
//
// Registry provides a central lookup table where implementations are
// registered by name and retrieved when needed.
package main

import (
	"fmt"
	"sync"
)

// Service is the common interface for registered services.
type Service interface {
	Name() string
	Execute() string
}

// Registry is a thread-safe service locator.
type Registry struct {
	mu       sync.RWMutex
	services map[string]Service
}

func NewRegistry() *Registry {
	return &Registry{services: make(map[string]Service)}
}

func (r *Registry) Register(svc Service) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.services[svc.Name()] = svc
}

func (r *Registry) Lookup(name string) (Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	svc, ok := r.services[name]
	if !ok {
		return nil, fmt.Errorf("service %q not found", name)
	}
	return svc, nil
}

// --- Concrete services ---

type EmailService struct{}

func (e *EmailService) Name() string    { return "email" }
func (e *EmailService) Execute() string { return "sending email" }

type SMSService struct{}

func (s *SMSService) Name() string    { return "sms" }
func (s *SMSService) Execute() string { return "sending sms" }

func main() {
	r := NewRegistry()
	r.Register(&EmailService{})
	r.Register(&SMSService{})

	for _, name := range []string{"email", "sms", "push"} {
		svc, err := r.Lookup(name)
		if err != nil {
			fmt.Printf("%-8s -> %s\n", name, err)
			continue
		}
		fmt.Printf("%-8s -> %s\n", name, svc.Execute())
	}
}
