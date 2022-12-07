package selinux

import (
	corev1 "k8s.io/api/core/v1"
	policy "k8s.io/api/policy/v1beta1"
	"testing"
)

func TestRunAsAnyOptions(t *testing.T) {
	_, err := NewRunAsAny(nil)
	if err != nil {
		t.Fatalf("unexpected error initializing NewRunAsAny %v", err)
	}
	_, err = NewRunAsAny(&policy.SELinuxStrategyOptions{})
	if err != nil {
		t.Errorf("unexpected error initializing NewRunAsAny %v", err)
	}
}

func TestRunAsAnyGenerate(t *testing.T) {
	s, err := NewRunAsAny(&policy.SELinuxStrategyOptions{})
	if err != nil {
		t.Fatalf("unexpected error initializing NewRunAsAny %v", err)
	}
	uid, err := s.Generate(nil, nil)
	if uid != nil {
		t.Errorf("expected nil uid but got %v", *uid)
	}
	if err != nil {
		t.Errorf("unexpected error generating uid %v", err)
	}
}

func TestRunAsAnyValidate(t *testing.T) {
	s, err := NewRunAsAny(&policy.SELinuxStrategyOptions{
		SELinuxOptions: &corev1.SELinuxOptions{
			Level: "foo",
		},
	},
	)
	if err != nil {
		t.Fatalf("unexpected error initializing NewRunAsAny %v", err)
	}
	errs := s.Validate(nil, nil, nil, nil)
	if len(errs) != 0 {
		t.Errorf("unexpected errors validating with ")
	}
	s, err = NewRunAsAny(&policy.SELinuxStrategyOptions{})
	if err != nil {
		t.Fatalf("unexpected error initializing NewRunAsAny %v", err)
	}
	errs = s.Validate(nil, nil, nil, nil)
	if len(errs) != 0 {
		t.Errorf("unexpected errors validating %v", errs)
	}
}
