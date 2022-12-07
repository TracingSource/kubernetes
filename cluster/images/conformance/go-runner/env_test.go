package main

import (
	"os"
	"testing"
)

func TestEnv(t *testing.T) {
	testCases := []struct {
		desc    string
		preHook func()
		env     Getenver
		expect  map[string]string
	}{
		{
			desc: "OS env",
			env:  &osEnv{},
			preHook: func() {
				os.Setenv("key1", "1")
			},
			expect: map[string]string{"key1": "1"},
		}, {
			desc: "OS env falls defaults to empty",
			env:  &osEnv{},
			preHook: func() {
				os.Unsetenv("key1")
			},
			expect: map[string]string{"key1": ""},
		}, {
			desc: "First choice of env respected",
			env: &defaultEnver{
				firstChoice: &explicitEnv{
					vals: map[string]string{
						"key1": "1",
					},
				},
				defaults: map[string]string{
					"key1": "default1",
					"key2": "default2",
				},
			},
			expect: map[string]string{
				"key1": "1",
				"key2": "default2",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			for k, expectVal := range tc.expect {
				if tc.preHook != nil {
					tc.preHook()
				}
				val := tc.env.Getenv(k)
				if val != expectVal {
					t.Errorf("Expected %q but got %q", expectVal, val)
				}
			}
		})
	}
}
