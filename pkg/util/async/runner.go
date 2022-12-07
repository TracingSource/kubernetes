package async

import (
	"sync"
)

// Runner is an abstraction to make it easy to start and stop groups of things that can be
// described by a single function which waits on a channel close to exit.
type Runner struct {
	lock      sync.Mutex
	loopFuncs []func(stop chan struct{})
	stop      *chan struct{}
}

// NewRunner makes a runner for the given function(s). The function(s) should loop until
// the channel is closed.
func NewRunner(f ...func(stop chan struct{})) *Runner {
	return &Runner{loopFuncs: f}
}

// Start begins running.
func (r *Runner) Start() {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.stop == nil {
		c := make(chan struct{})
		r.stop = &c
		for i := range r.loopFuncs {
			go r.loopFuncs[i](*r.stop)
		}
	}
}

// Stop stops running.
func (r *Runner) Stop() {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.stop != nil {
		close(*r.stop)
		r.stop = nil
	}
}
