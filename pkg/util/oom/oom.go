package oom

// This is a struct instead of an interface to allow injection of process ID listers and
// applying OOM score in tests.
// TODO: make this an interface, and inject a mock ioutil struct for testing.
type OOMAdjuster struct {
	pidLister                 func(cgroupName string) ([]int, error)
	ApplyOOMScoreAdj          func(pid int, oomScoreAdj int) error
	ApplyOOMScoreAdjContainer func(cgroupName string, oomScoreAdj, maxTries int) error
}
