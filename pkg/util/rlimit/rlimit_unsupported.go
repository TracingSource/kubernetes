// +build !linux

package rlimit

import (
	"errors"
)

func RlimitNumFiles(maxOpenFiles uint64) error {
	return errors.New("SetRLimit unsupported in this platform")
}
