// +build linux

package dockershim

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/api/core/v1"
)

func TestGetSeccompSecurityOpts(t *testing.T) {
	tests := []struct {
		msg            string
		seccompProfile string
		expectedOpts   []string
	}{{
		msg:            "No security annotations",
		seccompProfile: "",
		expectedOpts:   []string{"seccomp=unconfined"},
	}, {
		msg:            "Seccomp unconfined",
		seccompProfile: "unconfined",
		expectedOpts:   []string{"seccomp=unconfined"},
	}, {
		msg:            "Seccomp default",
		seccompProfile: v1.SeccompProfileRuntimeDefault,
		expectedOpts:   nil,
	}, {
		msg:            "Seccomp deprecated default",
		seccompProfile: v1.DeprecatedSeccompProfileDockerDefault,
		expectedOpts:   nil,
	}}

	for i, test := range tests {
		opts, err := getSeccompSecurityOpts(test.seccompProfile, '=')
		assert.NoError(t, err, "TestCase[%d]: %s", i, test.msg)
		assert.Len(t, opts, len(test.expectedOpts), "TestCase[%d]: %s", i, test.msg)
		for _, opt := range test.expectedOpts {
			assert.Contains(t, opts, opt, "TestCase[%d]: %s", i, test.msg)
		}
	}
}

func TestLoadSeccompLocalhostProfiles(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "seccomp-local-profile-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)
	testProfile := `{"foo": "bar"}`
	err = ioutil.WriteFile(filepath.Join(tmpdir, "test"), []byte(testProfile), 0644)
	require.NoError(t, err)

	tests := []struct {
		msg            string
		seccompProfile string
		expectedOpts   []string
		expectErr      bool
	}{{
		msg:            "Seccomp localhost/test profile should return correct seccomp profiles",
		seccompProfile: "localhost/" + filepath.Join(tmpdir, "test"),
		expectedOpts:   []string{`seccomp={"foo":"bar"}`},
		expectErr:      false,
	}, {
		msg:            "Non-existent profile should return error",
		seccompProfile: "localhost/" + filepath.Join(tmpdir, "fixtures/non-existent"),
		expectedOpts:   nil,
		expectErr:      true,
	}, {
		msg:            "Relative profile path should return error",
		seccompProfile: "localhost/fixtures/test",
		expectedOpts:   nil,
		expectErr:      true,
	}}

	for i, test := range tests {
		opts, err := getSeccompSecurityOpts(test.seccompProfile, '=')
		if test.expectErr {
			assert.Error(t, err, fmt.Sprintf("TestCase[%d]: %s", i, test.msg))
			continue
		}
		assert.NoError(t, err, "TestCase[%d]: %s", i, test.msg)
		assert.Len(t, opts, len(test.expectedOpts), "TestCase[%d]: %s", i, test.msg)
		for _, opt := range test.expectedOpts {
			assert.Contains(t, opts, opt, "TestCase[%d]: %s", i, test.msg)
		}
	}
}
