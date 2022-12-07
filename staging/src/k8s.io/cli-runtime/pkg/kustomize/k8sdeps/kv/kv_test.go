package kv

import (
	"reflect"
	"testing"
)

func TestKeyValuesFromLines(t *testing.T) {
	tests := []struct {
		desc          string
		content       string
		expectedPairs []Pair
		expectedErr   bool
	}{
		{
			desc: "valid kv content parse",
			content: `
		k1=v1
		k2=v2
		`,
			expectedPairs: []Pair{
				{Key: "k1", Value: "v1"},
				{Key: "k2", Value: "v2"},
			},
			expectedErr: false,
		},
		{
			desc: "content with comments",
			content: `
		k1=v1
		#k2=v2
		`,
			expectedPairs: []Pair{
				{Key: "k1", Value: "v1"},
			},
			expectedErr: false,
		},
		// TODO: add negative testcases
	}

	for _, test := range tests {
		pairs, err := KeyValuesFromLines([]byte(test.content))
		if test.expectedErr && err == nil {
			t.Fatalf("%s should not return error", test.desc)
		}

		if !reflect.DeepEqual(pairs, test.expectedPairs) {
			t.Errorf("%s should succeed, got:%v exptected:%v", test.desc, pairs, test.expectedPairs)
		}

	}
}
