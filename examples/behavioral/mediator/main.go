// Package main demonstrates the Mediator pattern.
//
// Mediator defines an object that encapsulates how a set of objects interact,
// promoting loose coupling by centralizing communication.
package main

import "fmt"

// Mediator defines the communication interface.
type Mediator interface {
	Notify(sender string, event string)
}

// --- Concrete components ---

type AuthService struct {
	mediator Mediator
}

func (a *AuthService) Login(user string) string {
	msg := fmt.Sprintf("%s logged in", user)
	a.mediator.Notify("auth", "login:"+user)
	return msg
}

type Logger struct {
	mediator Mediator
}

func (l *Logger) Log(message string) {
	fmt.Printf("  LOG: %s\n", message)
}

type Notifier struct {
	mediator Mediator
}

func (n *Notifier) Send(message string) {
	fmt.Printf("  NOTIFY: %s\n", message)
}

// --- Concrete mediator ---

type AppMediator struct {
	Auth     *AuthService
	Logger   *Logger
	Notifier *Notifier
}

func NewAppMediator() *AppMediator {
	m := &AppMediator{
		Auth:     &AuthService{},
		Logger:   &Logger{},
		Notifier: &Notifier{},
	}
	m.Auth.mediator = m
	m.Logger.mediator = m
	m.Notifier.mediator = m
	return m
}

func (m *AppMediator) Notify(sender string, event string) {
	if sender == "auth" {
		m.Logger.Log("user login event: " + event)
		m.Notifier.Send("welcome back!")
	}
}

func main() {
	app := NewAppMediator()

	fmt.Println(app.Auth.Login("alice"))
	fmt.Println()
	fmt.Println(app.Auth.Login("bob"))
}
