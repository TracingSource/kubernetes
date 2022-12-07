package utils

import (
	"fmt"

	"k8s.io/api/core/v1"
	podutil "k8s.io/kubernetes/pkg/api/v1/pod"
)

type ContainerFailures struct {
	status   *v1.ContainerStateTerminated
	Restarts int
}

// PodRunningReady checks whether pod p's phase is running and it has a ready
// condition of status true.
func PodRunningReady(p *v1.Pod) (bool, error) {
	// Check the phase is running.
	if p.Status.Phase != v1.PodRunning {
		return false, fmt.Errorf("want pod '%s' on '%s' to be '%v' but was '%v'",
			p.ObjectMeta.Name, p.Spec.NodeName, v1.PodRunning, p.Status.Phase)
	}
	// Check the ready condition is true.
	if !podutil.IsPodReady(p) {
		return false, fmt.Errorf("pod '%s' on '%s' didn't have condition {%v %v}; conditions: %v",
			p.ObjectMeta.Name, p.Spec.NodeName, v1.PodReady, v1.ConditionTrue, p.Status.Conditions)
	}
	return true, nil
}

func PodRunningReadyOrSucceeded(p *v1.Pod) (bool, error) {
	// Check if the phase is succeeded.
	if p.Status.Phase == v1.PodSucceeded {
		return true, nil
	}
	return PodRunningReady(p)
}

// FailedContainers inspects all containers in a pod and returns failure
// information for containers that have failed or been restarted.
// A map is returned where the key is the containerID and the value is a
// struct containing the restart and failure information
func FailedContainers(pod *v1.Pod) map[string]ContainerFailures {
	var state ContainerFailures
	states := make(map[string]ContainerFailures)

	statuses := pod.Status.ContainerStatuses
	if len(statuses) == 0 {
		return nil
	}
	for _, status := range statuses {
		if status.State.Terminated != nil {
			states[status.ContainerID] = ContainerFailures{status: status.State.Terminated}
		} else if status.LastTerminationState.Terminated != nil {
			states[status.ContainerID] = ContainerFailures{status: status.LastTerminationState.Terminated}
		}
		if status.RestartCount > 0 {
			var ok bool
			if state, ok = states[status.ContainerID]; !ok {
				state = ContainerFailures{}
			}
			state.Restarts = int(status.RestartCount)
			states[status.ContainerID] = state
		}
	}

	return states
}

// TerminatedContainers inspects all containers in a pod and returns a map
// of "container name: termination reason", for all currently terminated
// containers.
func TerminatedContainers(pod *v1.Pod) map[string]string {
	states := make(map[string]string)
	statuses := pod.Status.ContainerStatuses
	if len(statuses) == 0 {
		return states
	}
	for _, status := range statuses {
		if status.State.Terminated != nil {
			states[status.Name] = status.State.Terminated.Reason
		}
	}
	return states
}

// PodNotReady checks whether pod p's has a ready condition of status false.
func PodNotReady(p *v1.Pod) (bool, error) {
	// Check the ready condition is false.
	if podutil.IsPodReady(p) {
		return false, fmt.Errorf("pod '%s' on '%s' didn't have condition {%v %v}; conditions: %v",
			p.ObjectMeta.Name, p.Spec.NodeName, v1.PodReady, v1.ConditionFalse, p.Status.Conditions)
	}
	return true, nil
}
