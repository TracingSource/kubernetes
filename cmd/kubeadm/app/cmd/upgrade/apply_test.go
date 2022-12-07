package upgrade

import (
	"testing"
)

func TestSessionIsInteractive(t *testing.T) {
	var tcases = []struct {
		name     string
		flags    *applyFlags
		expected bool
	}{
		{
			name: "Explicitly non-interactive",
			flags: &applyFlags{
				nonInteractiveMode: true,
			},
			expected: false,
		},
		{
			name: "Implicitly non-interactive since --dryRun is used",
			flags: &applyFlags{
				dryRun: true,
			},
			expected: false,
		},
		{
			name: "Implicitly non-interactive since --force is used",
			flags: &applyFlags{
				force: true,
			},
			expected: false,
		},
		{
			name:     "Interactive session",
			flags:    &applyFlags{},
			expected: true,
		},
	}
	for _, tt := range tcases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.flags.sessionIsInteractive() != tt.expected {
				t.Error("unexpected result")
			}
		})
	}
}
