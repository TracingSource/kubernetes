package cm

import (
	"k8s.io/api/core/v1"
)

func NewFakeInternalContainerLifecycle() *fakeInternalContainerLifecycle {
	return &fakeInternalContainerLifecycle{}
}

type fakeInternalContainerLifecycle struct{}

func (f *fakeInternalContainerLifecycle) PreStartContainer(pod *v1.Pod, container *v1.Container, containerID string) error {
	return nil
}

func (f *fakeInternalContainerLifecycle) PreStopContainer(containerID string) error {
	return nil
}

func (f *fakeInternalContainerLifecycle) PostStopContainer(containerID string) error {
	return nil
}
