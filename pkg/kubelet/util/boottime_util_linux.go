// +build freebsd linux

package util

import (
	"fmt"
	"time"

	"golang.org/x/sys/unix"
)

// GetBootTime returns the time at which the machine was started, truncated to the nearest second
func GetBootTime() (time.Time, error) {
	currentTime := time.Now()
	var info unix.Sysinfo_t
	if err := unix.Sysinfo(&info); err != nil {
		return time.Time{}, fmt.Errorf("error getting system uptime: %s", err)
	}
	return currentTime.Add(-time.Duration(info.Uptime) * time.Second).Truncate(time.Second), nil
}
