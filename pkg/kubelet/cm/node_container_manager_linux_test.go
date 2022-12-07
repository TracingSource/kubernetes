// +build linux

package cm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	evictionapi "k8s.io/kubernetes/pkg/kubelet/eviction/api"
)

func TestNodeAllocatableReservationForScheduling(t *testing.T) {
	memoryEvictionThreshold := resource.MustParse("100Mi")
	cpuMemCases := []struct {
		kubeReserved   v1.ResourceList
		systemReserved v1.ResourceList
		expected       v1.ResourceList
		capacity       v1.ResourceList
		hardThreshold  evictionapi.ThresholdValue
	}{
		{
			kubeReserved:   getResourceList("100m", "100Mi"),
			systemReserved: getResourceList("50m", "50Mi"),
			capacity:       getResourceList("10", "10Gi"),
			expected:       getResourceList("150m", "150Mi"),
		},
		{
			kubeReserved:   getResourceList("100m", "100Mi"),
			systemReserved: getResourceList("50m", "50Mi"),
			hardThreshold: evictionapi.ThresholdValue{
				Quantity: &memoryEvictionThreshold,
			},
			capacity: getResourceList("10", "10Gi"),
			expected: getResourceList("150m", "250Mi"),
		},
		{
			kubeReserved:   getResourceList("100m", "100Mi"),
			systemReserved: getResourceList("50m", "50Mi"),
			capacity:       getResourceList("10", "10Gi"),
			hardThreshold: evictionapi.ThresholdValue{
				Percentage: 0.05,
			},
			expected: getResourceList("150m", "694157320"),
		},

		{
			kubeReserved:   v1.ResourceList{},
			systemReserved: v1.ResourceList{},
			capacity:       getResourceList("10", "10Gi"),
			expected:       getResourceList("", ""),
		},
		{
			kubeReserved:   getResourceList("", "100Mi"),
			systemReserved: getResourceList("50m", "50Mi"),
			capacity:       getResourceList("10", "10Gi"),
			expected:       getResourceList("50m", "150Mi"),
		},

		{
			kubeReserved:   getResourceList("50m", "100Mi"),
			systemReserved: getResourceList("", "50Mi"),
			capacity:       getResourceList("10", "10Gi"),
			expected:       getResourceList("50m", "150Mi"),
		},
		{
			kubeReserved:   getResourceList("", "100Mi"),
			systemReserved: getResourceList("", "50Mi"),
			capacity:       getResourceList("10", ""),
			expected:       getResourceList("", "150Mi"),
		},
	}
	for idx, tc := range cpuMemCases {
		nc := NodeConfig{
			NodeAllocatableConfig: NodeAllocatableConfig{
				KubeReserved:   tc.kubeReserved,
				SystemReserved: tc.systemReserved,
				HardEvictionThresholds: []evictionapi.Threshold{
					{
						Signal:   evictionapi.SignalMemoryAvailable,
						Operator: evictionapi.OpLessThan,
						Value:    tc.hardThreshold,
					},
				},
			},
		}
		cm := &containerManagerImpl{
			NodeConfig: nc,
			capacity:   tc.capacity,
		}
		for k, v := range cm.GetNodeAllocatableReservation() {
			expected, exists := tc.expected[k]
			assert.True(t, exists, "test case %d expected resource %q", idx+1, k)
			assert.Equal(t, expected.MilliValue(), v.MilliValue(), "test case %d failed for resource %q", idx+1, k)
		}
	}

	ephemeralStorageEvictionThreshold := resource.MustParse("100Mi")
	ephemeralStorageTestCases := []struct {
		kubeReserved  v1.ResourceList
		expected      v1.ResourceList
		capacity      v1.ResourceList
		hardThreshold evictionapi.ThresholdValue
	}{
		{
			kubeReserved: getEphemeralStorageResourceList("100Mi"),
			capacity:     getEphemeralStorageResourceList("10Gi"),
			expected:     getEphemeralStorageResourceList("100Mi"),
		},
		{
			kubeReserved: getEphemeralStorageResourceList("100Mi"),
			hardThreshold: evictionapi.ThresholdValue{
				Quantity: &ephemeralStorageEvictionThreshold,
			},
			capacity: getEphemeralStorageResourceList("10Gi"),
			expected: getEphemeralStorageResourceList("200Mi"),
		},
		{
			kubeReserved: getEphemeralStorageResourceList("150Mi"),
			capacity:     getEphemeralStorageResourceList("10Gi"),
			hardThreshold: evictionapi.ThresholdValue{
				Percentage: 0.05,
			},
			expected: getEphemeralStorageResourceList("694157320"),
		},

		{
			kubeReserved: v1.ResourceList{},
			capacity:     getEphemeralStorageResourceList("10Gi"),
			expected:     getEphemeralStorageResourceList(""),
		},
	}
	for idx, tc := range ephemeralStorageTestCases {
		nc := NodeConfig{
			NodeAllocatableConfig: NodeAllocatableConfig{
				KubeReserved: tc.kubeReserved,
				HardEvictionThresholds: []evictionapi.Threshold{
					{
						Signal:   evictionapi.SignalNodeFsAvailable,
						Operator: evictionapi.OpLessThan,
						Value:    tc.hardThreshold,
					},
				},
			},
		}
		cm := &containerManagerImpl{
			NodeConfig: nc,
			capacity:   tc.capacity,
		}
		for k, v := range cm.GetNodeAllocatableReservation() {
			expected, exists := tc.expected[k]
			assert.True(t, exists, "test case %d expected resource %q", idx+1, k)
			assert.Equal(t, expected.MilliValue(), v.MilliValue(), "test case %d failed for resource %q", idx+1, k)
		}
	}
}

func TestNodeAllocatableForEnforcement(t *testing.T) {
	memoryEvictionThreshold := resource.MustParse("100Mi")
	testCases := []struct {
		kubeReserved   v1.ResourceList
		systemReserved v1.ResourceList
		capacity       v1.ResourceList
		expected       v1.ResourceList
		hardThreshold  evictionapi.ThresholdValue
	}{
		{
			kubeReserved:   getResourceList("100m", "100Mi"),
			systemReserved: getResourceList("50m", "50Mi"),
			capacity:       getResourceList("10", "10Gi"),
			expected:       getResourceList("9850m", "10090Mi"),
		},
		{
			kubeReserved:   getResourceList("100m", "100Mi"),
			systemReserved: getResourceList("50m", "50Mi"),
			hardThreshold: evictionapi.ThresholdValue{
				Quantity: &memoryEvictionThreshold,
			},
			capacity: getResourceList("10", "10Gi"),
			expected: getResourceList("9850m", "10090Mi"),
		},
		{
			kubeReserved:   getResourceList("100m", "100Mi"),
			systemReserved: getResourceList("50m", "50Mi"),
			hardThreshold: evictionapi.ThresholdValue{
				Percentage: 0.05,
			},
			capacity: getResourceList("10", "10Gi"),
			expected: getResourceList("9850m", "10090Mi"),
		},

		{
			kubeReserved:   v1.ResourceList{},
			systemReserved: v1.ResourceList{},
			capacity:       getResourceList("10", "10Gi"),
			expected:       getResourceList("10", "10Gi"),
		},
		{
			kubeReserved:   getResourceList("", "100Mi"),
			systemReserved: getResourceList("50m", "50Mi"),
			capacity:       getResourceList("10", "10Gi"),
			expected:       getResourceList("9950m", "10090Mi"),
		},

		{
			kubeReserved:   getResourceList("50m", "100Mi"),
			systemReserved: getResourceList("", "50Mi"),
			capacity:       getResourceList("10", "10Gi"),
			expected:       getResourceList("9950m", "10090Mi"),
		},
		{
			kubeReserved:   getResourceList("", "100Mi"),
			systemReserved: getResourceList("", "50Mi"),
			capacity:       getResourceList("10", ""),
			expected:       getResourceList("10", ""),
		},
	}
	for idx, tc := range testCases {
		nc := NodeConfig{
			NodeAllocatableConfig: NodeAllocatableConfig{
				KubeReserved:   tc.kubeReserved,
				SystemReserved: tc.systemReserved,
				HardEvictionThresholds: []evictionapi.Threshold{
					{
						Signal:   evictionapi.SignalMemoryAvailable,
						Operator: evictionapi.OpLessThan,
						Value:    tc.hardThreshold,
					},
				},
			},
		}
		cm := &containerManagerImpl{
			NodeConfig: nc,
			capacity:   tc.capacity,
		}
		for k, v := range cm.getNodeAllocatableAbsolute() {
			expected, exists := tc.expected[k]
			assert.True(t, exists)
			assert.Equal(t, expected.MilliValue(), v.MilliValue(), "test case %d failed for resource %q", idx+1, k)
		}
	}
}

func TestNodeAllocatableInputValidation(t *testing.T) {
	memoryEvictionThreshold := resource.MustParse("100Mi")
	highMemoryEvictionThreshold := resource.MustParse("2Gi")
	cpuMemTestCases := []struct {
		kubeReserved         v1.ResourceList
		systemReserved       v1.ResourceList
		capacity             v1.ResourceList
		hardThreshold        evictionapi.ThresholdValue
		invalidConfiguration bool
	}{
		{
			kubeReserved:   getResourceList("100m", "100Mi"),
			systemReserved: getResourceList("50m", "50Mi"),
			capacity:       getResourceList("10", "10Gi"),
		},
		{
			kubeReserved:   getResourceList("100m", "100Mi"),
			systemReserved: getResourceList("50m", "50Mi"),
			hardThreshold: evictionapi.ThresholdValue{
				Quantity: &memoryEvictionThreshold,
			},
			capacity: getResourceList("10", "10Gi"),
		},
		{
			kubeReserved:   getResourceList("100m", "100Mi"),
			systemReserved: getResourceList("50m", "50Mi"),
			hardThreshold: evictionapi.ThresholdValue{
				Percentage: 0.05,
			},
			capacity: getResourceList("10", "10Gi"),
		},
		{
			kubeReserved:   v1.ResourceList{},
			systemReserved: v1.ResourceList{},
			capacity:       getResourceList("10", "10Gi"),
		},
		{
			kubeReserved:   getResourceList("", "100Mi"),
			systemReserved: getResourceList("50m", "50Mi"),
			capacity:       getResourceList("10", "10Gi"),
		},
		{
			kubeReserved:   getResourceList("50m", "100Mi"),
			systemReserved: getResourceList("", "50Mi"),
			capacity:       getResourceList("10", "10Gi"),
		},
		{
			kubeReserved:   getResourceList("", "100Mi"),
			systemReserved: getResourceList("", "50Mi"),
			capacity:       getResourceList("10", ""),
		},
		{
			kubeReserved:   getResourceList("5", "10Gi"),
			systemReserved: getResourceList("5", "10Gi"),
			hardThreshold: evictionapi.ThresholdValue{
				Quantity: &highMemoryEvictionThreshold,
			},
			capacity:             getResourceList("10", "11Gi"),
			invalidConfiguration: true,
		},
	}
	for _, tc := range cpuMemTestCases {
		nc := NodeConfig{
			NodeAllocatableConfig: NodeAllocatableConfig{
				KubeReserved:   tc.kubeReserved,
				SystemReserved: tc.systemReserved,
				HardEvictionThresholds: []evictionapi.Threshold{
					{
						Signal:   evictionapi.SignalMemoryAvailable,
						Operator: evictionapi.OpLessThan,
						Value:    tc.hardThreshold,
					},
				},
			},
		}
		cm := &containerManagerImpl{
			NodeConfig: nc,
			capacity:   tc.capacity,
		}
		err := cm.validateNodeAllocatable()
		if err == nil && tc.invalidConfiguration {
			t.Fatalf("Expected invalid node allocatable configuration")
		} else if err != nil && !tc.invalidConfiguration {
			t.Fatalf("Expected valid node allocatable configuration: %v", err)
		}
	}

	ephemeralStorageEvictionThreshold := resource.MustParse("100Mi")
	ephemeralStorageTestCases := []struct {
		kubeReserved         v1.ResourceList
		capacity             v1.ResourceList
		hardThreshold        evictionapi.ThresholdValue
		invalidConfiguration bool
	}{
		{
			kubeReserved: getEphemeralStorageResourceList("100Mi"),
			capacity:     getEphemeralStorageResourceList("500Mi"),
		},
		{
			kubeReserved: getEphemeralStorageResourceList("20Gi"),
			hardThreshold: evictionapi.ThresholdValue{
				Quantity: &ephemeralStorageEvictionThreshold,
			},
			capacity:             getEphemeralStorageResourceList("20Gi"),
			invalidConfiguration: true,
		},
	}
	for _, tc := range ephemeralStorageTestCases {
		nc := NodeConfig{
			NodeAllocatableConfig: NodeAllocatableConfig{
				KubeReserved: tc.kubeReserved,
				HardEvictionThresholds: []evictionapi.Threshold{
					{
						Signal:   evictionapi.SignalNodeFsAvailable,
						Operator: evictionapi.OpLessThan,
						Value:    tc.hardThreshold,
					},
				},
			},
		}
		cm := &containerManagerImpl{
			NodeConfig: nc,
			capacity:   tc.capacity,
		}
		err := cm.validateNodeAllocatable()
		if err == nil && tc.invalidConfiguration {
			t.Fatalf("Expected invalid node allocatable configuration")
		} else if err != nil && !tc.invalidConfiguration {
			t.Fatalf("Expected valid node allocatable configuration: %v", err)
		}
	}
}

// getEphemeralStorageResourceList returns a ResourceList with the
// specified ephemeral storage resource values
func getEphemeralStorageResourceList(storage string) v1.ResourceList {
	res := v1.ResourceList{}
	if storage != "" {
		res[v1.ResourceEphemeralStorage] = resource.MustParse(storage)
	}
	return res
}
