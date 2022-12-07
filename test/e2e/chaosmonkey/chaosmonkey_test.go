package chaosmonkey

import (
	"sync/atomic"
	"testing"
)

func TestDoWithPanic(t *testing.T) {
	var counter int64
	cm := New(func() {})
	tests := []Test{
		// No panic
		func(sem *Semaphore) {
			defer atomic.AddInt64(&counter, 1)
			sem.Ready()
		},
		// Panic after sem.Ready()
		func(sem *Semaphore) {
			defer atomic.AddInt64(&counter, 1)
			sem.Ready()
			panic("Panic after calling sem.Ready()")
		},
		// Panic before sem.Ready()
		func(sem *Semaphore) {
			defer atomic.AddInt64(&counter, 1)
			panic("Panic before calling sem.Ready()")
		},
	}
	for _, test := range tests {
		cm.Register(test)
	}
	cm.Do()
	// Check that all funcs in tests were called.
	if int(counter) != len(tests) {
		t.Errorf("Expected counter to be %v, but it was %v", len(tests), counter)
	}
}
