package lifecycle

import (
	"k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/scheduler/algorithm/predicates"
)

// AdmissionFailureHandlerStub is an AdmissionFailureHandler that does not perform any handling of admission failure.
// It simply passes the failure on.
type AdmissionFailureHandlerStub struct{}

var _ AdmissionFailureHandler = &AdmissionFailureHandlerStub{}

func NewAdmissionFailureHandlerStub() *AdmissionFailureHandlerStub {
	return &AdmissionFailureHandlerStub{}
}

func (n *AdmissionFailureHandlerStub) HandleAdmissionFailure(admitPod *v1.Pod, failureReasons []predicates.PredicateFailureReason) (bool, []predicates.PredicateFailureReason, error) {
	return false, failureReasons, nil
}
