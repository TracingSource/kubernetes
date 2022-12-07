package migration

import (
	"errors"
	"reflect"
	"testing"

	"k8s.io/kubernetes/pkg/scheduler/algorithm/predicates"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
)

func TestPredicateResultToFrameworkStatus(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		reasons    []predicates.PredicateFailureReason
		wantStatus *framework.Status
	}{
		{
			name: "Success",
		},
		{
			name:       "Error",
			err:        errors.New("Failed with error"),
			wantStatus: framework.NewStatus(framework.Error, "Failed with error"),
		},
		{
			name:       "Error with reason",
			err:        errors.New("Failed with error"),
			reasons:    []predicates.PredicateFailureReason{predicates.ErrDiskConflict},
			wantStatus: framework.NewStatus(framework.Error, "Failed with error"),
		},
		{
			name:       "Unschedulable",
			reasons:    []predicates.PredicateFailureReason{predicates.ErrExistingPodsAntiAffinityRulesNotMatch},
			wantStatus: framework.NewStatus(framework.Unschedulable, "node(s) didn't satisfy existing pods anti-affinity rules"),
		},
		{
			name:       "Unschedulable and Unresolvable",
			reasons:    []predicates.PredicateFailureReason{predicates.ErrDiskConflict, predicates.ErrNodeSelectorNotMatch},
			wantStatus: framework.NewStatus(framework.UnschedulableAndUnresolvable, "node(s) didn't match node selector"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStatus := PredicateResultToFrameworkStatus(tt.reasons, tt.err)
			if !reflect.DeepEqual(tt.wantStatus, gotStatus) {
				t.Errorf("Got status %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}
}
