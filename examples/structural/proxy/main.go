// Package main demonstrates the Proxy pattern.
//
// Proxy provides an object that controls access to another object,
// intercepting all calls for authorization, logging, or lazy init.
package main

import "fmt"

// ITerminal is the common interface for real and proxy terminals.
type ITerminal interface {
	Execute(cmd string) (string, error)
}

// GopherTerminal is the real subject.
type GopherTerminal struct {
	User string
}

func (gt *GopherTerminal) Execute(cmd string) (string, error) {
	prefix := fmt.Sprintf("%s@go_term$:", gt.User)
	switch cmd {
	case "say_hi":
		return fmt.Sprintf("%s Hi %s", prefix, gt.User), nil
	case "man":
		return fmt.Sprintf("%s Visit 'https://golang.org/doc/' for Golang documentation", prefix), nil
	default:
		return "", fmt.Errorf("%s Unknown command: %s", prefix, cmd)
	}
}

// Terminal is the proxy that validates users before delegating.
type Terminal struct {
	currentUser    string
	gopherTerminal *GopherTerminal
}

func NewTerminal(user string) (*Terminal, error) {
	if user == "" {
		return nil, fmt.Errorf("user can't be empty")
	}
	if err := authorizeUser(user); err != nil {
		return nil, fmt.Errorf("you (%s) are not allowed to use terminal", user)
	}
	return &Terminal{currentUser: user}, nil
}

func (t *Terminal) Execute(command string) (string, error) {
	t.gopherTerminal = &GopherTerminal{User: t.currentUser}
	fmt.Printf("PROXY: Intercepted execution of user (%s), command (%s)\n", t.currentUser, command)
	return t.gopherTerminal.Execute(command)
}

func authorizeUser(user string) error {
	allowed := map[string]bool{"gopher": true, "admin": true}
	if !allowed[user] {
		return fmt.Errorf("user %s not in allow list", user)
	}
	return nil
}

func main() {
	t, err := NewTerminal("gopher")
	if err != nil {
		panic(err)
	}

	resp, err := t.Execute("say_hi")
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
	}
	fmt.Println(resp)

	fmt.Println()

	resp, err = t.Execute("man")
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
	}
	fmt.Println(resp)

	fmt.Println()

	// Unauthorized user
	_, err = NewTerminal("hacker")
	if err != nil {
		fmt.Printf("Access denied: %s\n", err)
	}
}
