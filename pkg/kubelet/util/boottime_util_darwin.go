// +build darwin

package util

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

// GetBootTime returns the time at which the machine was started, truncated to the nearest second
func GetBootTime() (time.Time, error) {
	output, err := unix.SysctlRaw("kern.boottime")
	if err != nil {
		return time.Time{}, err
	}
	var timeval syscall.Timeval
	if len(output) != int(unsafe.Sizeof(timeval)) {
		return time.Time{}, fmt.Errorf("unexpected output when calling syscall kern.bootime.  Expected len(output) to be %v, but got %v",
			int(unsafe.Sizeof(timeval)), len(output))
	}
	timeval = *(*syscall.Timeval)(unsafe.Pointer(&output[0]))
	sec, nsec := timeval.Unix()
	return time.Unix(sec, nsec).Truncate(time.Second), nil
}
