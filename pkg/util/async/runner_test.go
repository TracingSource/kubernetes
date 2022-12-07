package async

import (
	"fmt"
	"sync"
	"testing"
)

func TestRunner(t *testing.T) {
	var (
		lock   sync.Mutex
		events []string
		funcs  []func(chan struct{})
	)
	done := make(chan struct{}, 20)
	for i := 0; i < 10; i++ {
		iCopy := i
		funcs = append(funcs, func(c chan struct{}) {
			lock.Lock()
			events = append(events, fmt.Sprintf("%v starting\n", iCopy))
			lock.Unlock()
			<-c
			lock.Lock()
			events = append(events, fmt.Sprintf("%v stopping\n", iCopy))
			lock.Unlock()
			done <- struct{}{}
		})
	}

	r := NewRunner(funcs...)
	r.Start()
	r.Stop()
	for i := 0; i < 10; i++ {
		<-done
	}
	if len(events) != 20 {
		t.Errorf("expected 20 events, but got:\n%v\n", events)
	}
}
