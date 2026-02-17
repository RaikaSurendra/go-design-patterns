// Package main demonstrates the Handshaking pattern.
//
// Handshaking proactively asks a service if it can accept work before
// sending the actual request. Unlike circuit breaker (reactive),
// handshaking is cooperative.
package main

import (
	"errors"
	"fmt"
	"sync/atomic"
)

var ErrServiceAtCapacity = errors.New("service is at capacity")

type Request struct {
	Payload string
}

type Response struct {
	Result string
}

type Service interface {
	IsReady() bool
	Do(req Request) (Response, error)
}

type TrackedService struct {
	active   int64
	capacity int64
}

func NewTrackedService(capacity int64) *TrackedService {
	return &TrackedService{capacity: capacity}
}

func (s *TrackedService) IsReady() bool {
	return atomic.LoadInt64(&s.active) < s.capacity
}

func (s *TrackedService) Do(req Request) (Response, error) {
	atomic.AddInt64(&s.active, 1)
	defer atomic.AddInt64(&s.active, -1)
	return Response{Result: "processed: " + req.Payload}, nil
}

func Call(svc Service, req Request) (Response, error) {
	if !svc.IsReady() {
		return Response{}, ErrServiceAtCapacity
	}
	return svc.Do(req)
}

func main() {
	svc := NewTrackedService(2)

	requests := []Request{
		{Payload: "order-1"},
		{Payload: "order-2"},
		{Payload: "order-3"},
	}

	// Simulate: first two succeed, third checks capacity
	// (In practice this would be concurrent; simplified for clarity)
	atomic.StoreInt64(&svc.active, 0)
	for _, req := range requests[:2] {
		resp, err := Call(svc, req)
		if err != nil {
			fmt.Printf("%-10s -> %v\n", req.Payload, err)
		} else {
			fmt.Printf("%-10s -> %s\n", req.Payload, resp.Result)
		}
	}

	// Simulate service at capacity
	atomic.StoreInt64(&svc.active, 2)
	resp, err := Call(svc, requests[2])
	if errors.Is(err, ErrServiceAtCapacity) {
		fmt.Printf("%-10s -> rejected: %v\n", requests[2].Payload, err)
	} else {
		fmt.Printf("%-10s -> %s\n", requests[2].Payload, resp.Result)
	}
}
