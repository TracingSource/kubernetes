package format

import (
	"testing"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestResourceList(t *testing.T) {
	resourceList := v1.ResourceList{}
	resourceList[v1.ResourceCPU] = resource.MustParse("100m")
	resourceList[v1.ResourceMemory] = resource.MustParse("5Gi")
	actual := ResourceList(resourceList)
	expected := "cpu=100m,memory=5Gi"
	if actual != expected {
		t.Errorf("Unexpected result, actual: %v, expected: %v", actual, expected)
	}
}
