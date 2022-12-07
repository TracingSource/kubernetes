package app

import (
	"testing"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/stretchr/testify/assert"
)

func TestIsControllerEnabled(t *testing.T) {
	tcs := []struct {
		name                         string
		controllerName               string
		controllers                  []string
		disabledByDefaultControllers []string
		expected                     bool
	}{
		{
			name:                         "on by name",
			controllerName:               "bravo",
			controllers:                  []string{"alpha", "bravo", "-charlie"},
			disabledByDefaultControllers: []string{"delta", "echo"},
			expected:                     true,
		},
		{
			name:                         "off by name",
			controllerName:               "charlie",
			controllers:                  []string{"alpha", "bravo", "-charlie"},
			disabledByDefaultControllers: []string{"delta", "echo"},
			expected:                     false,
		},
		{
			name:                         "on by default",
			controllerName:               "alpha",
			controllers:                  []string{"*"},
			disabledByDefaultControllers: []string{"delta", "echo"},
			expected:                     true,
		},
		{
			name:                         "off by default",
			controllerName:               "delta",
			controllers:                  []string{"*"},
			disabledByDefaultControllers: []string{"delta", "echo"},
			expected:                     false,
		},
		{
			name:                         "off by default implicit, no star",
			controllerName:               "foxtrot",
			controllers:                  []string{"alpha", "bravo", "-charlie"},
			disabledByDefaultControllers: []string{"delta", "echo"},
			expected:                     false,
		},
	}

	for _, tc := range tcs {
		actual := IsControllerEnabled(tc.controllerName, sets.NewString(tc.disabledByDefaultControllers...), tc.controllers)
		assert.Equal(t, tc.expected, actual, "%v: expected %v, got %v", tc.name, tc.expected, actual)
	}

}
