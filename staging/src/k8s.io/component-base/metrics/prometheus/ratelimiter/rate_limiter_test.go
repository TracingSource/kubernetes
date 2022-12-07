package ratelimiter

import (
	"strings"
	"testing"

	"k8s.io/client-go/util/flowcontrol"
)

func TestRegisterMetricAndTrackRateLimiterUsage(t *testing.T) {
	testCases := []struct {
		ownerName   string
		rateLimiter flowcontrol.RateLimiter
		err         string
	}{
		{
			ownerName:   "owner_name",
			rateLimiter: flowcontrol.NewTokenBucketRateLimiter(1, 1),
			err:         "",
		},
		{
			ownerName:   "owner_name",
			rateLimiter: flowcontrol.NewTokenBucketRateLimiter(1, 1),
			err:         "already registered",
		},
		{
			ownerName:   "invalid-owner-name",
			rateLimiter: flowcontrol.NewTokenBucketRateLimiter(1, 1),
			err:         "error registering rate limiter usage metric",
		},
	}

	for i, tc := range testCases {
		e := RegisterMetricAndTrackRateLimiterUsage(tc.ownerName, tc.rateLimiter)
		if e != nil {
			if tc.err == "" {
				t.Errorf("[%d] unexpected error: %v", i, e)
			} else if !strings.Contains(e.Error(), tc.err) {
				t.Errorf("[%d] expected an error containing %q: %v", i, tc.err, e)
			}
		}
	}
}
