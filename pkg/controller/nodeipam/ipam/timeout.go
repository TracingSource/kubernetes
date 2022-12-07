package ipam

import (
	"time"
)

// Timeout manages the resync loop timing for a given node sync operation. The
// timeout changes depending on whether or not there was an error reported for
// the operation. Consecutive errors will result in exponential backoff to a
// maxBackoff timeout.
type Timeout struct {
	// Resync is the default timeout duration when there are no errors.
	Resync time.Duration
	// MaxBackoff is the maximum timeout when in a error backoff state.
	MaxBackoff time.Duration
	// InitialRetry is the initial retry interval when an error is reported.
	InitialRetry time.Duration

	// errs is the count of consecutive errors that have occurred.
	errs int
	// current is the current backoff timeout.
	current time.Duration
}

// Update the timeout with the current error state.
func (b *Timeout) Update(ok bool) {
	if ok {
		b.errs = 0
		b.current = b.Resync
		return
	}

	b.errs++
	if b.errs == 1 {
		b.current = b.InitialRetry
		return
	}

	b.current *= 2
	if b.current >= b.MaxBackoff {
		b.current = b.MaxBackoff
	}
}

// Next returns the next operation timeout given the disposition of err.
func (b *Timeout) Next() time.Duration {
	if b.errs == 0 {
		return b.Resync
	}
	return b.current
}
