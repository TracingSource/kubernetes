// +build linux

package rlimit

import (
	"golang.org/x/sys/unix"
)

func RlimitNumFiles(maxOpenFiles uint64) {
	unix.Setrlimit(unix.RLIMIT_NOFILE, &unix.Rlimit{Max: maxOpenFiles, Cur: maxOpenFiles})
}
