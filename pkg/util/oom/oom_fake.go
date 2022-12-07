package oom

type FakeOOMAdjuster struct{}

func NewFakeOOMAdjuster() *OOMAdjuster {
	return &OOMAdjuster{
		pidLister:                 func(cgroupName string) ([]int, error) { return make([]int, 0), nil },
		ApplyOOMScoreAdj:          fakeApplyOOMScoreAdj,
		ApplyOOMScoreAdjContainer: fakeApplyOOMScoreAdjContainer,
	}
}

func fakeApplyOOMScoreAdj(pid int, oomScoreAdj int) error {
	return nil
}

func fakeApplyOOMScoreAdjContainer(cgroupName string, oomScoreAdj, maxTries int) error {
	return nil
}
