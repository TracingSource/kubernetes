package plugins

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
)

func TestTranslateAllowedTopologies(t *testing.T) {
	testCases := []struct {
		name            string
		topology        []v1.TopologySelectorTerm
		expectedToplogy []v1.TopologySelectorTerm
		expErr          bool
	}{
		{
			name:     "no translation",
			topology: generateToplogySelectors(GCEPDTopologyKey, []string{"foo", "bar"}),
			expectedToplogy: []v1.TopologySelectorTerm{
				{
					MatchLabelExpressions: []v1.TopologySelectorLabelRequirement{
						{
							Key:    GCEPDTopologyKey,
							Values: []string{"foo", "bar"},
						},
					},
				},
			},
		},
		{
			name: "translate",
			topology: []v1.TopologySelectorTerm{
				{
					MatchLabelExpressions: []v1.TopologySelectorLabelRequirement{
						{
							Key:    "failure-domain.beta.kubernetes.io/zone",
							Values: []string{"foo", "bar"},
						},
					},
				},
			},
			expectedToplogy: []v1.TopologySelectorTerm{
				{
					MatchLabelExpressions: []v1.TopologySelectorLabelRequirement{
						{
							Key:    GCEPDTopologyKey,
							Values: []string{"foo", "bar"},
						},
					},
				},
			},
		},
		{
			name: "combo",
			topology: []v1.TopologySelectorTerm{
				{
					MatchLabelExpressions: []v1.TopologySelectorLabelRequirement{
						{
							Key:    "failure-domain.beta.kubernetes.io/zone",
							Values: []string{"foo", "bar"},
						},
						{
							Key:    GCEPDTopologyKey,
							Values: []string{"boo", "baz"},
						},
					},
				},
			},
			expectedToplogy: []v1.TopologySelectorTerm{
				{
					MatchLabelExpressions: []v1.TopologySelectorLabelRequirement{
						{
							Key:    GCEPDTopologyKey,
							Values: []string{"foo", "bar"},
						},
						{
							Key:    GCEPDTopologyKey,
							Values: []string{"boo", "baz"},
						},
					},
				},
			},
		},
		{
			name: "some other key",
			topology: []v1.TopologySelectorTerm{
				{
					MatchLabelExpressions: []v1.TopologySelectorLabelRequirement{
						{
							Key:    "test",
							Values: []string{"foo", "bar"},
						},
					},
				},
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		t.Logf("Running test: %v", tc.name)
		gotTop, err := translateAllowedTopologies(tc.topology, GCEPDTopologyKey)
		if err != nil && !tc.expErr {
			t.Errorf("Did not expect an error, got: %v", err)
		}
		if err == nil && tc.expErr {
			t.Errorf("Expected an error but did not get one")
		}

		if !reflect.DeepEqual(gotTop, tc.expectedToplogy) {
			t.Errorf("Expected topology: %v, but got: %v", tc.expectedToplogy, gotTop)
		}
	}
}
