package upgrade

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/pkg/errors"
	"k8s.io/kubernetes/cmd/kubeadm/app/constants"
)

func createTestRunDiffFile(contents []byte) (string, error) {
	file, err := ioutil.TempFile("", "kubeadm-upgrade-diff-config-*.yaml")
	if err != nil {
		return "", errors.Wrap(err, "failed to create temporary test file")
	}
	if _, err := file.Write([]byte(contents)); err != nil {
		return "", errors.Wrap(err, "failed to write to temporary test file")
	}
	if err := file.Close(); err != nil {
		return "", errors.Wrap(err, "failed to close temporary test file")
	}
	return file.Name(), nil
}

func TestRunDiff(t *testing.T) {
	currentVersion := "v" + constants.CurrentKubernetesVersion.String()

	// create a temporary file with valid ClusterConfiguration
	testUpgradeDiffConfigContents := []byte("apiVersion: kubeadm.k8s.io/v1beta2\n" +
		"kind: ClusterConfiguration\n" +
		"kubernetesVersion: " + currentVersion)
	testUpgradeDiffConfig, err := createTestRunDiffFile(testUpgradeDiffConfigContents)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testUpgradeDiffConfig)

	// create a temporary manifest file with dummy contents
	testUpgradeDiffManifestContents := []byte("some-contents")
	testUpgradeDiffManifest, err := createTestRunDiffFile(testUpgradeDiffManifestContents)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testUpgradeDiffManifest)

	flags := &diffFlags{
		cfgPath: "",
		out:     ioutil.Discard,
	}

	// TODO: Add test cases for empty cfgPath, it should automatically fetch cfg from cluster
	testCases := []struct {
		name            string
		args            []string
		setManifestPath bool
		manifestPath    string
		cfgPath         string
		expectedError   bool
	}{
		{
			name:            "valid: run diff on valid manifest path",
			cfgPath:         testUpgradeDiffConfig,
			setManifestPath: true,
			manifestPath:    testUpgradeDiffManifest,
			expectedError:   false,
		},
		{
			name:          "invalid: missing config file",
			cfgPath:       "missing-path-to-a-config",
			expectedError: true,
		},
		{
			name:            "invalid: valid config but empty manifest path",
			cfgPath:         testUpgradeDiffConfig,
			setManifestPath: true,
			manifestPath:    "",
			expectedError:   true,
		},
		{
			name:            "invalid: valid config but bad manifest path",
			cfgPath:         testUpgradeDiffConfig,
			setManifestPath: true,
			manifestPath:    "bad-path",
			expectedError:   true,
		},
		{
			name:          "invalid: badly formatted version as argument",
			cfgPath:       testUpgradeDiffConfig,
			args:          []string{"bad-version"},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			flags.cfgPath = tc.cfgPath
			if tc.setManifestPath {
				flags.apiServerManifestPath = tc.manifestPath
				flags.controllerManagerManifestPath = tc.manifestPath
				flags.schedulerManifestPath = tc.manifestPath
			}
			if err := runDiff(flags, tc.args); (err != nil) != tc.expectedError {
				t.Fatalf("expected error: %v, saw: %v, error: %v", tc.expectedError, (err != nil), err)
			}
		})
	}
}
