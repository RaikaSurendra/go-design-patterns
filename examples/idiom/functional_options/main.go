// Package main demonstrates the Functional Options pattern.
//
// Functional options provide a clean API for configuring objects with
// sensible defaults and optional overrides.
package main

import "fmt"

type Server struct {
	host    string
	port    int
	timeout int
	maxConn int
}

type Option func(*Server)

func WithHost(host string) Option {
	return func(s *Server) { s.host = host }
}

func WithPort(port int) Option {
	return func(s *Server) { s.port = port }
}

func WithTimeout(timeout int) Option {
	return func(s *Server) { s.timeout = timeout }
}

func WithMaxConn(maxConn int) Option {
	return func(s *Server) { s.maxConn = maxConn }
}

func NewServer(opts ...Option) *Server {
	s := &Server{
		host:    "localhost",
		port:    8080,
		timeout: 30,
		maxConn: 100,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Server) String() string {
	return fmt.Sprintf("Server{host:%s, port:%d, timeout:%ds, maxConn:%d}",
		s.host, s.port, s.timeout, s.maxConn)
}

func main() {
	// Default configuration
	s1 := NewServer()
	fmt.Println("Default:", s1)

	// Custom configuration
	s2 := NewServer(
		WithHost("0.0.0.0"),
		WithPort(9090),
		WithTimeout(60),
	)
	fmt.Println("Custom: ", s2)

	// Minimal override
	s3 := NewServer(WithPort(3000))
	fmt.Println("Minimal:", s3)
}
