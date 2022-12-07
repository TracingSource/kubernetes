package clientset_test

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/metrics/pkg/client/clientset/versioned/fake"
)

// TestFakeList is a basic sanity check that makes sure the fake Clientset is working properly.
func TestFakeList(t *testing.T) {
	client := fake.NewSimpleClientset()
	if _, err := client.MetricsV1alpha1().PodMetricses("").List(metav1.ListOptions{}); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if _, err := client.MetricsV1alpha1().NodeMetricses().List(metav1.ListOptions{}); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if _, err := client.MetricsV1beta1().PodMetricses("").List(metav1.ListOptions{}); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if _, err := client.MetricsV1beta1().NodeMetricses().List(metav1.ListOptions{}); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
