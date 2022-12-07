package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetContainerName(t *testing.T) {
	var cases = []struct {
		labels        map[string]string
		containerName string
	}{
		{
			labels: map[string]string{
				"io.kubernetes.container.name": "c1",
			},
			containerName: "c1",
		},
		{
			labels: map[string]string{
				"io.kubernetes.container.name": "c2",
			},
			containerName: "c2",
		},
	}
	for _, data := range cases {
		containerName := GetContainerName(data.labels)
		assert.Equal(t, data.containerName, containerName)
	}
}

func TestGetPodName(t *testing.T) {
	var cases = []struct {
		labels  map[string]string
		podName string
	}{
		{
			labels: map[string]string{
				"io.kubernetes.pod.name": "p1",
			},
			podName: "p1",
		},
		{
			labels: map[string]string{
				"io.kubernetes.pod.name": "p2",
			},
			podName: "p2",
		},
	}
	for _, data := range cases {
		podName := GetPodName(data.labels)
		assert.Equal(t, data.podName, podName)
	}
}

func TestGetPodUID(t *testing.T) {
	var cases = []struct {
		labels map[string]string
		podUID string
	}{
		{
			labels: map[string]string{
				"io.kubernetes.pod.uid": "uid1",
			},
			podUID: "uid1",
		},
		{
			labels: map[string]string{
				"io.kubernetes.pod.uid": "uid2",
			},
			podUID: "uid2",
		},
	}
	for _, data := range cases {
		podUID := GetPodUID(data.labels)
		assert.Equal(t, data.podUID, podUID)
	}
}

func TestGetPodNamespace(t *testing.T) {
	var cases = []struct {
		labels       map[string]string
		podNamespace string
	}{
		{
			labels: map[string]string{
				"io.kubernetes.pod.namespace": "ns1",
			},
			podNamespace: "ns1",
		},
		{
			labels: map[string]string{
				"io.kubernetes.pod.namespace": "ns2",
			},
			podNamespace: "ns2",
		},
	}
	for _, data := range cases {
		podNamespace := GetPodNamespace(data.labels)
		assert.Equal(t, data.podNamespace, podNamespace)
	}
}
