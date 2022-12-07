package gce

import "testing"

func TestIsPortsSubset(t *testing.T) {
	tc := map[string]struct {
		required  []string
		coverage  []string
		expectErr bool
	}{
		"Single port coverage": {
			required: []string{"tcp/50"},
			coverage: []string{"tcp/50", "tcp/60", "tcp/70"},
		},
		"Port range coverage": {
			required: []string{"tcp/50"},
			coverage: []string{"tcp/20-30", "tcp/45-60"},
		},
		"Multiple Port range coverage": {
			required: []string{"tcp/50", "tcp/29", "tcp/46"},
			coverage: []string{"tcp/20-30", "tcp/45-60"},
		},
		"Not covered": {
			required:  []string{"tcp/50"},
			coverage:  []string{"udp/50", "tcp/49", "tcp/51-60"},
			expectErr: true,
		},
	}

	for name, c := range tc {
		t.Run(name, func(t *testing.T) {
			gotErr := isPortsSubset(c.required, c.coverage)
			if c.expectErr != (gotErr != nil) {
				t.Errorf("isPortsSubset(%v, %v) = %v, wanted err? %v", c.required, c.coverage, gotErr, c.expectErr)
			}
		})
	}
}
