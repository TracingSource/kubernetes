package log

import (
	"fmt"

	"k8s.io/klog"
)

const logFmt = "kubelet config controller: %s"

// Errorf shim that inserts "kubelet config controller" at the beginning of the log message,
// while still reporting the call site of the logging function.
func Errorf(format string, args ...interface{}) {
	var s string
	if len(args) > 0 {
		s = fmt.Sprintf(format, args...)
	} else {
		s = format
	}
	klog.ErrorDepth(1, fmt.Sprintf(logFmt, s))
}

// Infof shim that inserts "kubelet config controller" at the beginning of the log message,
// while still reporting the call site of the logging function.
func Infof(format string, args ...interface{}) {
	var s string
	if len(args) > 0 {
		s = fmt.Sprintf(format, args...)
	} else {
		s = format
	}
	klog.InfoDepth(1, fmt.Sprintf(logFmt, s))
}
