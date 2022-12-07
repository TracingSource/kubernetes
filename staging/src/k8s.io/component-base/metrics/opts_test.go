package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultStabilityLevel(t *testing.T) {
	var tests = []struct {
		name        string
		inputValue  StabilityLevel
		expectValue StabilityLevel
		expectPanic bool
	}{
		{
			name:        "empty should take ALPHA by default",
			inputValue:  "",
			expectValue: ALPHA,
			expectPanic: false,
		},
		{
			name:        "ALPHA remain unchanged",
			inputValue:  ALPHA,
			expectValue: ALPHA,
			expectPanic: false,
		},
		{
			name:        "STABLE remain unchanged",
			inputValue:  STABLE,
			expectValue: STABLE,
			expectPanic: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var stability = tc.inputValue

			stability.setDefaults()
			assert.Equalf(t, tc.expectValue, stability, "Got %s, expected: %v ", stability, tc.expectValue)
		})
	}
}
