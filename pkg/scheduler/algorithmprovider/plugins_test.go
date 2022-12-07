package algorithmprovider

import (
	"fmt"
	"testing"

	"k8s.io/kubernetes/pkg/scheduler"
)

var (
	algorithmProviderNames = []string{
		scheduler.DefaultProvider,
	}
)

func TestDefaultConfigExists(t *testing.T) {
	p, err := scheduler.GetAlgorithmProvider(scheduler.DefaultProvider)
	if err != nil {
		t.Errorf("error retrieving default provider: %v", err)
	}
	if p == nil {
		t.Error("algorithm provider config should not be nil")
	}
	if len(p.FitPredicateKeys) == 0 {
		t.Error("default algorithm provider shouldn't have 0 fit predicates")
	}
}

func TestAlgorithmProviders(t *testing.T) {
	for _, pn := range algorithmProviderNames {
		t.Run(pn, func(t *testing.T) {
			p, err := scheduler.GetAlgorithmProvider(pn)
			if err != nil {
				t.Fatalf("error retrieving provider: %v", err)
			}
			if len(p.PriorityFunctionKeys) == 0 {
				t.Errorf("algorithm provider shouldn't have 0 priority functions")
			}
			for _, pf := range p.PriorityFunctionKeys.List() {
				t.Run(fmt.Sprintf("priorityfunction/%s", pf), func(t *testing.T) {
					if !scheduler.IsPriorityFunctionRegistered(pf) {
						t.Errorf("priority function is not registered but is used in the algorithm provider")
					}
				})
			}
			for _, fp := range p.FitPredicateKeys.List() {
				t.Run(fmt.Sprintf("fitpredicate/%s", fp), func(t *testing.T) {
					if !scheduler.IsFitPredicateRegistered(fp) {
						t.Errorf("fit predicate is not registered but is used in the algorithm provider")
					}
				})
			}
		})
	}
}

func TestApplyFeatureGates(t *testing.T) {
	for _, pn := range algorithmProviderNames {
		t.Run(pn, func(t *testing.T) {
			p, err := scheduler.GetAlgorithmProvider(pn)
			if err != nil {
				t.Fatalf("Error retrieving provider: %v", err)
			}

			if !p.FitPredicateKeys.Has("PodToleratesNodeTaints") {
				t.Fatalf("Failed to find predicate: 'PodToleratesNodeTaints'")
			}
		})
	}

	defer ApplyFeatureGates()()

	for _, pn := range algorithmProviderNames {
		t.Run(pn, func(t *testing.T) {
			p, err := scheduler.GetAlgorithmProvider(pn)
			if err != nil {
				t.Fatalf("Error retrieving '%s' provider: %v", pn, err)
			}

			if !p.FitPredicateKeys.Has("PodToleratesNodeTaints") {
				t.Fatalf("Failed to find predicate: 'PodToleratesNodeTaints'")
			}
		})
	}
}
