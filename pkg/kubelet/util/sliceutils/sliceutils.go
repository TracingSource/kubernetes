package sliceutils

import (
	"k8s.io/api/core/v1"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
)

// StringInSlice returns true if s is in list
func StringInSlice(s string, list []string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}

	return false
}

// PodsByCreationTime makes an array of pods sortable by their creation
// timestamps in ascending order.
type PodsByCreationTime []*v1.Pod

func (s PodsByCreationTime) Len() int {
	return len(s)
}

func (s PodsByCreationTime) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s PodsByCreationTime) Less(i, j int) bool {
	return s[i].CreationTimestamp.Before(&s[j].CreationTimestamp)
}

// ByImageSize makes an array of images sortable by their size in descending
// order.
type ByImageSize []kubecontainer.Image

func (a ByImageSize) Less(i, j int) bool {
	if a[i].Size == a[j].Size {
		return a[i].ID > a[j].ID
	}
	return a[i].Size > a[j].Size
}
func (a ByImageSize) Len() int      { return len(a) }
func (a ByImageSize) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
