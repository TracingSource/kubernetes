package protobuf

import "testing"

func TestProtoSafePackage(t *testing.T) {
	tests := []struct {
		pkg      string
		expected string
	}{
		{
			pkg:      "foo",
			expected: "foo",
		},
		{
			pkg:      "foo/bar",
			expected: "foo.bar",
		},
		{
			pkg:      "foo/bar/baz",
			expected: "foo.bar.baz",
		},
		{
			pkg:      "foo/bar-baz/x/y-z/q",
			expected: "foo.bar_baz.x.y_z.q",
		},
	}

	for _, test := range tests {
		actual := protoSafePackage(test.pkg)
		if e, a := test.expected, actual; e != a {
			t.Errorf("%s: expected %s, got %s", test.pkg, e, a)
		}
	}
}
