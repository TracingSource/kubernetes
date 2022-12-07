package ipam

import (
	"errors"
	"testing"
	"time"
)

func TestTimeout(t *testing.T) {
	time10s := 10 * time.Second
	time5s := 5 * time.Second
	timeout := &Timeout{
		Resync:       time10s,
		MaxBackoff:   time5s,
		InitialRetry: time.Second,
	}

	for _, testStep := range []struct {
		err  error
		want time.Duration
	}{
		{nil, time10s},
		{nil, time10s},
		{errors.New("x"), time.Second},
		{errors.New("x"), 2 * time.Second},
		{errors.New("x"), 4 * time.Second},
		{errors.New("x"), 5 * time.Second},
		{errors.New("x"), 5 * time.Second},
		{nil, time10s},
		{nil, time10s},
		{errors.New("x"), time.Second},
		{errors.New("x"), 2 * time.Second},
		{nil, time10s},
	} {
		timeout.Update(testStep.err == nil)
		next := timeout.Next()
		if next != testStep.want {
			t.Errorf("timeout.next(%v) = %v, want %v", testStep.err, next, testStep.want)
		}
	}
}
