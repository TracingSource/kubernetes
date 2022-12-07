package bandwidth

import (
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	api "k8s.io/kubernetes/pkg/apis/core"
)

func TestExtractPodBandwidthResources(t *testing.T) {
	four, _ := resource.ParseQuantity("4M")
	ten, _ := resource.ParseQuantity("10M")
	twenty, _ := resource.ParseQuantity("20M")

	testPod := func(ingress, egress string) *api.Pod {
		pod := &api.Pod{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{}}}
		if len(ingress) != 0 {
			pod.Annotations["kubernetes.io/ingress-bandwidth"] = ingress
		}
		if len(egress) != 0 {
			pod.Annotations["kubernetes.io/egress-bandwidth"] = egress
		}
		return pod
	}

	tests := []struct {
		pod             *api.Pod
		expectedIngress *resource.Quantity
		expectedEgress  *resource.Quantity
		expectError     bool
	}{
		{
			pod: &api.Pod{},
		},
		{
			pod:             testPod("10M", ""),
			expectedIngress: &ten,
		},
		{
			pod:            testPod("", "10M"),
			expectedEgress: &ten,
		},
		{
			pod:             testPod("4M", "20M"),
			expectedIngress: &four,
			expectedEgress:  &twenty,
		},
		{
			pod:         testPod("foo", ""),
			expectError: true,
		},
	}
	for _, test := range tests {
		ingress, egress, err := ExtractPodBandwidthResources(test.pod.Annotations)
		if test.expectError {
			if err == nil {
				t.Errorf("unexpected non-error")
			}
			continue
		}
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}
		if !reflect.DeepEqual(ingress, test.expectedIngress) {
			t.Errorf("expected: %v, saw: %v", ingress, test.expectedIngress)
		}
		if !reflect.DeepEqual(egress, test.expectedEgress) {
			t.Errorf("expected: %v, saw: %v", egress, test.expectedEgress)
		}
	}
}
