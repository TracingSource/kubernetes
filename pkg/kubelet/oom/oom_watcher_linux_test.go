package oom

import (
	"testing"

	"github.com/stretchr/testify/assert"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
)

// TestBasic verifies that the OOMWatch works without error.
func TestBasic(t *testing.T) {
	fakeRecorder := &record.FakeRecorder{}
	node := &v1.ObjectReference{}
	oomWatcher := NewWatcher(fakeRecorder)
	assert.NoError(t, oomWatcher.Start(node))

	// TODO: Improve this test once cadvisor exports events.EventChannel as an interface
	// and thereby allow using a mock version of cadvisor.
}
