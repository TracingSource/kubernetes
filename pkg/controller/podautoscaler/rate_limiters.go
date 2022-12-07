package podautoscaler

import (
	"time"

	"k8s.io/client-go/util/workqueue"
)

// FixedItemIntervalRateLimiter limits items to a fixed-rate interval
type FixedItemIntervalRateLimiter struct {
	interval time.Duration
}

var _ workqueue.RateLimiter = &FixedItemIntervalRateLimiter{}

func NewFixedItemIntervalRateLimiter(interval time.Duration) workqueue.RateLimiter {
	return &FixedItemIntervalRateLimiter{
		interval: interval,
	}
}

func (r *FixedItemIntervalRateLimiter) When(item interface{}) time.Duration {
	return r.interval
}

func (r *FixedItemIntervalRateLimiter) NumRequeues(item interface{}) int {
	return 1
}

func (r *FixedItemIntervalRateLimiter) Forget(item interface{}) {
}

// NewDefaultHPARateLimiter creates a rate limiter which limits overall (as per the
// default controller rate limiter), as well as per the resync interval
func NewDefaultHPARateLimiter(interval time.Duration) workqueue.RateLimiter {
	return NewFixedItemIntervalRateLimiter(interval)
}
