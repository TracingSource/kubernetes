package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidKey(t *testing.T) {
	testcases := []struct {
		key   string
		valid bool
	}{
		{
			"    ",
			false,
		},
		{
			"/foo/bar",
			false,
		},
		{
			".foo",
			false,
		},
		{
			"a78768279290d33d0b82eaea43cb8346f500057cb5bd250e88c97a5585385d66",
			true,
		},
		{
			"a7.87-6_8",
			true,
		},
		{
			"a7.87-677-",
			false,
		},
	}

	for _, tc := range testcases {
		if tc.valid {
			assert.NoError(t, ValidateKey(tc.key))
		} else {
			assert.Error(t, ValidateKey(tc.key))
		}
	}
}
