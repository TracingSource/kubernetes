package scheduling

import (
	"fmt"
	"strings"

	"k8s.io/api/core/v1"
)

func scheduleSuccessEvent(ns, podName, nodeName string) func(*v1.Event) bool {
	return func(e *v1.Event) bool {
		return e.Type == v1.EventTypeNormal &&
			e.Reason == "Scheduled" &&
			strings.HasPrefix(e.Name, podName) &&
			strings.Contains(e.Message, fmt.Sprintf("Successfully assigned %v/%v to %v", ns, podName, nodeName))
	}
}

func scheduleFailureEvent(podName string) func(*v1.Event) bool {
	return func(e *v1.Event) bool {
		return strings.HasPrefix(e.Name, podName) &&
			e.Type == "Warning" &&
			e.Reason == "FailedScheduling"
	}
}
