package ports

import (
	"fmt"
	"net"
	"sync"
)

type Allocator struct {
	mu        sync.Mutex
	allocated map[int]bool
}

func NewAllocator() *Allocator {
	return &Allocator{
		allocated: make(map[int]bool),
	}
}

func (a *Allocator) Allocate() (int, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, fmt.Errorf("cannot allocate port: %v", err)
	}
	defer ln.Close()

	addr := ln.Addr().(*net.TCPAddr)
	port := addr.Port

	// There's a race between closing the listener and the service binding
	// to the same port. For MVP this is documented but not solved.
	a.allocated[port] = true

	return port, nil
}

func (a *Allocator) Release(port int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.allocated, port)
}
