package util

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestGetNonZeroRequest(t *testing.T) {
	tests := []struct {
		name           string
		requests       v1.ResourceList
		expectedCPU    int64
		expectedMemory int64
	}{
		{
			"cpu_and_memory_not_found",
			v1.ResourceList{},
			DefaultMilliCPURequest,
			DefaultMemoryRequest,
		},
		{
			"only_cpu_exist",
			v1.ResourceList{
				v1.ResourceCPU: resource.MustParse("200m"),
			},
			200,
			DefaultMemoryRequest,
		},
		{
			"only_memory_exist",
			v1.ResourceList{
				v1.ResourceMemory: resource.MustParse("400Mi"),
			},
			DefaultMilliCPURequest,
			400 * 1024 * 1024,
		},
		{
			"cpu_memory_exist",
			v1.ResourceList{
				v1.ResourceCPU:    resource.MustParse("200m"),
				v1.ResourceMemory: resource.MustParse("400Mi"),
			},
			200,
			400 * 1024 * 1024,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			realCPU, realMemory := GetNonzeroRequests(&test.requests)
			assert.EqualValuesf(t, test.expectedCPU, realCPU, "Failed to test: %s", test.name)
			assert.EqualValuesf(t, test.expectedMemory, realMemory, "Failed to test: %s", test.name)
		})
	}
}

func TestGetLeastRequestResource(t *testing.T) {
	tests := []struct {
		name             string
		requests         v1.ResourceList
		resource         v1.ResourceName
		expectedQuantity int64
	}{
		{
			"extended_resource_not_found",
			v1.ResourceList{},
			v1.ResourceName("intel.com/foo"),
			0,
		},
		{
			"extended_resource_found",
			v1.ResourceList{
				v1.ResourceName("intel.com/foo"): resource.MustParse("4"),
			},
			v1.ResourceName("intel.com/foo"),
			4,
		},
		{
			"cpu_not_found",
			v1.ResourceList{},
			v1.ResourceCPU,
			DefaultMilliCPURequest,
		},
		{
			"memory_not_found",
			v1.ResourceList{},
			v1.ResourceMemory,
			DefaultMemoryRequest,
		},
		{
			"cpu_exist",
			v1.ResourceList{
				v1.ResourceCPU: resource.MustParse("200m"),
			},
			v1.ResourceCPU,
			200,
		},
		{
			"memory_exist",
			v1.ResourceList{
				v1.ResourceMemory: resource.MustParse("400Mi"),
			},
			v1.ResourceMemory,
			400 * 1024 * 1024,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			realQuantity := GetNonzeroRequestForResource(test.resource, &test.requests)
			assert.EqualValuesf(t, test.expectedQuantity, realQuantity, "Failed to test: %s", test.name)
		})
	}
}
