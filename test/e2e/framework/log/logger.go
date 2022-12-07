// Package log will be removed after switching to use core framework log.
// Do not make further changes here!
package log

import (
	"fmt"
	"time"

	"github.com/onsi/ginkgo"

	"k8s.io/kubernetes/test/e2e/framework/ginkgowrapper"
)

func nowStamp() string {
	return time.Now().Format(time.StampMilli)
}

func log(level string, format string, args ...interface{}) {
	fmt.Fprintf(ginkgo.GinkgoWriter, nowStamp()+": "+level+": "+format+"\n", args...)
}

// Logf logs the info.
func Logf(format string, args ...interface{}) {
	log("INFO", format, args...)
}

// Failf logs the fail info.
func Failf(format string, args ...interface{}) {
	FailfWithOffset(1, format, args...)
}

// FailfWithOffset calls "Fail" and logs the error at "offset" levels above its caller
// (for example, for call chain f -> g -> FailfWithOffset(1, ...) error would be logged for "f").
func FailfWithOffset(offset int, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log("FAIL", msg)
	ginkgowrapper.Fail(nowStamp()+": "+msg, 1+offset)
}

// Fail is a replacement for ginkgo.Fail which logs the problem as it occurs
// and then calls ginkgowrapper.Fail.
func Fail(msg string, callerSkip ...int) {
	skip := 1
	if len(callerSkip) > 0 {
		skip += callerSkip[0]
	}
	log("FAIL", msg)
	ginkgowrapper.Fail(nowStamp()+": "+msg, skip)
}
