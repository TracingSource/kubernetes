package modes

import "testing"

func TestIsValidAuthorizationMode(t *testing.T) {
	var tests = []struct {
		authzMode string
		expected  bool
	}{
		{"", false},
		{"rBAC", false},        // not supported
		{"falsy value", false}, // not supported
		{"RBAC", true},         // supported
		{"ABAC", true},         // supported
		{"Webhook", true},      // supported
		{"AlwaysAllow", true},  // supported
		{"AlwaysDeny", true},   // supported
	}
	for _, rt := range tests {
		actual := IsValidAuthorizationMode(rt.authzMode)
		if actual != rt.expected {
			t.Errorf(
				"failed ValidAuthorizationMode:\n\texpected: %t\n\t  actual: %t",
				rt.expected,
				actual,
			)
		}
	}
}
