package ports

import "testing"

func TestAllocate(t *testing.T) {
	a := NewAllocator()
	port, err := a.Allocate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if port <= 0 || port > 65535 {
		t.Errorf("port %d out of range", port)
	}
}

func TestAllocateMultiple(t *testing.T) {
	a := NewAllocator()
	ports := make(map[int]bool)

	for i := 0; i < 5; i++ {
		port, err := a.Allocate()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ports[port] {
			t.Errorf("duplicate port %d", port)
		}
		ports[port] = true
	}
}

func TestRelease(t *testing.T) {
	a := NewAllocator()
	port, err := a.Allocate()
	if err != nil {
		t.Fatal(err)
	}
	a.Release(port)
}
