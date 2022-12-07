package clientv3

import (
	"math/rand"
	"reflect"
	"runtime"
	"strings"
	"time"
)

// jitterUp adds random jitter to the duration.
//
// This adds or subtracts time from the duration within a given jitter fraction.
// For example for 10s and jitter 0.1, it will return a time within [9s, 11s])
//
// Reference: https://godoc.org/github.com/grpc-ecosystem/go-grpc-middleware/util/backoffutils
func jitterUp(duration time.Duration, jitter float64) time.Duration {
	multiplier := jitter * (rand.Float64()*2 - 1)
	return time.Duration(float64(duration) * (1 + multiplier))
}

// Check if the provided function is being called in the op options.
func isOpFuncCalled(op string, opts []OpOption) bool {
	for _, opt := range opts {
		v := reflect.ValueOf(opt)
		if v.Kind() == reflect.Func {
			if opFunc := runtime.FuncForPC(v.Pointer()); opFunc != nil {
				if strings.Contains(opFunc.Name(), op) {
					return true
				}
			}
		}
	}
	return false
}
