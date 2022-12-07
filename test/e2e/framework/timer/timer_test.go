package timer

import (
	"testing"
	"time"

	"github.com/onsi/gomega"

	"k8s.io/kubernetes/test/e2e/framework"
)

var currentTime time.Time

func init() {
	setCurrentTimeSinceEpoch(0)
	now = func() time.Time { return currentTime }
}

func setCurrentTimeSinceEpoch(duration time.Duration) {
	currentTime = time.Unix(0, duration.Nanoseconds())
}

func testUsageWithDefer(timer *TestPhaseTimer) {
	defer timer.StartPhase(33, "two").End()
	setCurrentTimeSinceEpoch(6*time.Second + 500*time.Millisecond)
}

func TestTimer(t *testing.T) {
	gomega.RegisterTestingT(t)

	timer := NewTestPhaseTimer()
	setCurrentTimeSinceEpoch(1 * time.Second)
	phaseOne := timer.StartPhase(1, "one")
	setCurrentTimeSinceEpoch(3 * time.Second)
	testUsageWithDefer(timer)

	gomega.Expect(timer.PrintJSON()).To(gomega.MatchJSON(`{
		"version": "v1",
		"dataItems": [
			{
				"data": {
					"001-one": 5.5,
					"033-two": 3.5
				},
				"unit": "s",
				"labels": {
					"test": "phases",
					"ended": "false"
				}
			}
		]
	}`))
	framework.ExpectEqual(timer.PrintHumanReadable(), `Phase 001-one: 5.5s so far
Phase 033-two: 3.5s
`)

	setCurrentTimeSinceEpoch(7*time.Second + 500*time.Millisecond)
	phaseOne.End()

	gomega.Expect(timer.PrintJSON()).To(gomega.MatchJSON(`{
		"version": "v1",
		"dataItems": [
			{
				"data": {
					"001-one": 6.5,
					"033-two": 3.5
				},
				"unit": "s",
				"labels": {
					"test": "phases"
				}
			}
		]
	}`))
	framework.ExpectEqual(timer.PrintHumanReadable(), `Phase 001-one: 6.5s
Phase 033-two: 3.5s
`)
}
