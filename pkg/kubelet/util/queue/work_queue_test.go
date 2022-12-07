package queue

import (
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/apimachinery/pkg/util/sets"
)

func newTestBasicWorkQueue() (*basicWorkQueue, *clock.FakeClock) {
	fakeClock := clock.NewFakeClock(time.Now())
	wq := &basicWorkQueue{
		clock: fakeClock,
		queue: make(map[types.UID]time.Time),
	}
	return wq, fakeClock
}

func compareResults(t *testing.T, expected, actual []types.UID) {
	expectedSet := sets.NewString()
	for _, u := range expected {
		expectedSet.Insert(string(u))
	}
	actualSet := sets.NewString()
	for _, u := range actual {
		actualSet.Insert(string(u))
	}
	if !expectedSet.Equal(actualSet) {
		t.Errorf("Expected %#v, got %#v", expectedSet.List(), actualSet.List())
	}
}

func TestGetWork(t *testing.T) {
	q, clock := newTestBasicWorkQueue()
	q.Enqueue(types.UID("foo1"), -1*time.Minute)
	q.Enqueue(types.UID("foo2"), -1*time.Minute)
	q.Enqueue(types.UID("foo3"), 1*time.Minute)
	q.Enqueue(types.UID("foo4"), 1*time.Minute)
	expected := []types.UID{types.UID("foo1"), types.UID("foo2")}
	compareResults(t, expected, q.GetWork())
	compareResults(t, []types.UID{}, q.GetWork())
	// Dial the time to 1 hour ahead.
	clock.Step(time.Hour)
	expected = []types.UID{types.UID("foo3"), types.UID("foo4")}
	compareResults(t, expected, q.GetWork())
	compareResults(t, []types.UID{}, q.GetWork())
}
