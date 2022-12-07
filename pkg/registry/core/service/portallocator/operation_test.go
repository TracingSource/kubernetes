package portallocator

import (
	"testing"

	"k8s.io/apimachinery/pkg/util/net"
)

// TestDryRunAllocate tests the Allocate function in dry run mode
func TestDryRunAllocate(t *testing.T) {
	pr, err := net.ParsePortRange("10000-10200")
	if err != nil {
		t.Fatal(err)
	}

	// Allocate some ports before calling
	previouslyAllocated := []int{10000, 10010, 10020}
	r, err := NewPortAllocator(*pr)
	if err != nil {
		t.Fatal(err)
	}
	for _, port := range previouslyAllocated {
		_ = r.Allocate(port)
	}
	freeAtStart := r.Free()

	// Do some allocations with a dry run operation
	toAllocate := []int{
		10000,
		10030,
		10030,
		10040,
	}
	expectedErrors := []error{
		ErrAllocated,
		nil,
		ErrAllocated,
		nil,
	}
	op := StartOperation(r, true)
	for i, port := range toAllocate {
		err := op.Allocate(port)
		if err != expectedErrors[i] {
			t.Errorf("%v: expected error %v but got %v", i, expectedErrors[i], err)
		}
	}

	// Make sure no port allocations were actually made by the dry run
	freeAtEnd := r.Free()
	if freeAtStart != freeAtEnd {
		t.Errorf("expected %v free ports but got %v", freeAtStart, freeAtEnd)
	}
}

// TestDryRunAllocateNext tests the AllocateNext function in dry run mode
func TestDryRunAllocateNext(t *testing.T) {
	pr, err := net.ParsePortRange("10000-10200")
	if err != nil {
		t.Fatal(err)
	}

	// Allocate some ports before calling
	previouslyAllocated := []int{10000, 10010, 10020}
	r, err := NewPortAllocator(*pr)
	if err != nil {
		t.Fatal(err)
	}
	for _, port := range previouslyAllocated {
		_ = r.Allocate(port)
	}
	freeAtStart := r.Free()

	// AllocateNext without a previously unused dry run operation
	op := StartOperation(r, true)
	port, err := op.AllocateNext()
	if port == 0 {
		t.Errorf("expected non zero port but got: %v", port)
	}
	if err != nil {
		t.Errorf("expected no error but got: %v", err)
	}

	// Try to allocate the returned port using the same operation
	if e, a := ErrAllocated, op.Allocate(port); e != a {
		t.Errorf("expected %v but got: %v", e, a)
	}

	// AllocateNext with a previously used dry run operation
	op = StartOperation(r, true)
	_ = op.Allocate(12345)
	port, err = op.AllocateNext()
	if port == 0 {
		t.Errorf("expected non zero port but got: %v", port)
	}
	if port == 12345 {
		t.Errorf("expected port not to be 12345 but got %v", port)
	}
	if err != nil {
		t.Errorf("expected no error but got: %v", err)
	}

	// Make sure no port allocations were actually made by the dry run
	freeAtEnd := r.Free()
	if freeAtStart != freeAtEnd {
		t.Errorf("expected %v free ports but got %v", freeAtStart, freeAtEnd)
	}
}
