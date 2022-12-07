package queue

import (
	"testing"
	"time"

	ktypes "k8s.io/apimachinery/pkg/types"
)

func TestBackoffPod(t *testing.T) {
	bpm := NewPodBackoffMap(1*time.Second, 10*time.Second)

	tests := []struct {
		podID            ktypes.NamespacedName
		expectedDuration time.Duration
		advanceClock     time.Duration
	}{
		{
			podID:            ktypes.NamespacedName{Namespace: "default", Name: "foo"},
			expectedDuration: 1 * time.Second,
		},
		{
			podID:            ktypes.NamespacedName{Namespace: "default", Name: "foo"},
			expectedDuration: 2 * time.Second,
		},
		{
			podID:            ktypes.NamespacedName{Namespace: "default", Name: "foo"},
			expectedDuration: 4 * time.Second,
		},
		{
			podID:            ktypes.NamespacedName{Namespace: "default", Name: "foo"},
			expectedDuration: 8 * time.Second,
		},
		{
			podID:            ktypes.NamespacedName{Namespace: "default", Name: "foo"},
			expectedDuration: 10 * time.Second,
		},
		{
			podID:            ktypes.NamespacedName{Namespace: "default", Name: "foo"},
			expectedDuration: 10 * time.Second,
		},
		{
			podID:            ktypes.NamespacedName{Namespace: "default", Name: "bar"},
			expectedDuration: 1 * time.Second,
		},
	}

	for _, test := range tests {
		// Backoff the pod
		bpm.BackoffPod(test.podID)
		// Get backoff duration for the pod
		duration := bpm.calculateBackoffDuration(test.podID)

		if duration != test.expectedDuration {
			t.Errorf("expected: %s, got %s for pod %s", test.expectedDuration.String(), duration.String(), test.podID)
		}
	}
}

func TestClearPodBackoff(t *testing.T) {
	bpm := NewPodBackoffMap(1*time.Second, 60*time.Second)
	// Clear backoff on an not existed pod
	bpm.clearPodBackoff(ktypes.NamespacedName{Namespace: "ns", Name: "not-existed"})
	// Backoff twice for pod foo
	podID := ktypes.NamespacedName{Namespace: "ns", Name: "foo"}
	bpm.BackoffPod(podID)
	bpm.BackoffPod(podID)
	if duration := bpm.calculateBackoffDuration(podID); duration != 2*time.Second {
		t.Errorf("Expected backoff of 1s for pod %s, got %s", podID, duration.String())
	}
	// Clear backoff for pod foo
	bpm.clearPodBackoff(podID)
	// Backoff once for pod foo
	bpm.BackoffPod(podID)
	if duration := bpm.calculateBackoffDuration(podID); duration != 1*time.Second {
		t.Errorf("Expected backoff of 1s for pod %s, got %s", podID, duration.String())
	}
}
