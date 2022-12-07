package printers

import (
	"testing"
)

func TestIllegalPackageSourceChecker(t *testing.T) {
	disallowedPrefixes := []string{
		"foo/bar",
		"k8s.io/foo/bar/vendor/k8s.io/baz/buz",
		"bar/foo/baz",
	}

	testCases := []struct {
		name            string
		pkgPath         string
		shouldBeAllowed bool
	}{
		{
			name:            "package path beginning with forbidden prefix is rejected",
			pkgPath:         "foo/bar/baz/buz",
			shouldBeAllowed: false,
		},
		{
			name:            "package path not fully matching forbidden prefix is allowed",
			pkgPath:         "bar/foo",
			shouldBeAllowed: true,
		},
		{
			name:            "package path containing forbidden prefix (not as prefix) is allowed",
			pkgPath:         "k8s.io/bar/foo/baz/etc",
			shouldBeAllowed: true,
		},
	}

	checker := &illegalPackageSourceChecker{disallowedPrefixes}

	for _, tc := range testCases {
		if checker.IsForbidden(tc.pkgPath) {
			if tc.shouldBeAllowed {
				t.Fatalf("expected package path %q to have been allowed", tc.pkgPath)
			}
			continue
		}

		if !tc.shouldBeAllowed {
			t.Fatalf("expected package path %q to have been rejected", tc.pkgPath)
		}
	}
}
