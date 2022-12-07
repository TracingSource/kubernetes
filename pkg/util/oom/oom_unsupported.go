// +build !linux

package oom

import (
	"errors"
)

var unsupportedErr = errors.New("setting OOM scores is unsupported in this build")

func NewOOMAdjuster() *OOMAdjuster {
	return &OOMAdjuster{
		ApplyOOMScoreAdj:          unsupportedApplyOOMScoreAdj,
		ApplyOOMScoreAdjContainer: unsupportedApplyOOMScoreAdjContainer,
	}
}

func unsupportedApplyOOMScoreAdj(pid int, oomScoreAdj int) error {
	return unsupportedErr
}

func unsupportedApplyOOMScoreAdjContainer(cgroupName string, oomScoreAdj, maxTries int) error {
	return unsupportedErr
}
