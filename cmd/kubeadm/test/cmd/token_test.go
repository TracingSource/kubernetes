package kubeadm

import (
	"flag"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"testing"
)

const (
	TokenExpectedRegex = "^\\S{6}\\.\\S{16}\n$"
)

var kubeadmPathFlag = flag.String("kubeadm-path", filepath.Join(os.Getenv("KUBE_ROOT"), "cluster/kubeadm.sh"), "Location of kubeadm")

func getKubeadmPath() string {
	kubeadmPath := *kubeadmPathFlag // TEST_SRCDIR is provided by Bazel.
	if srcDir := os.Getenv("TEST_SRCDIR"); srcDir != "" {
		kubeadmPath = path.Join(srcDir, os.Getenv("TEST_WORKSPACE"), kubeadmPath)
	}

	return kubeadmPath
}

var kubeadmCmdSkip = flag.Bool("kubeadm-cmd-skip", false, "Skip kubeadm cmd tests")

func TestCmdTokenGenerate(t *testing.T) {
	kubeadmPath := getKubeadmPath()
	if *kubeadmCmdSkip {
		t.Log("kubeadm cmd tests being skipped")
		t.Skip()
	}
	stdout, _, _, err := RunCmd(kubeadmPath, "token", "generate")
	if err != nil {
		t.Fatalf("'kubeadm token generate' exited uncleanly: %v", err)
	}

	matched, err := regexp.MatchString(TokenExpectedRegex, stdout)
	if err != nil {
		t.Fatalf("encountered an error while trying to match 'kubeadm token generate' stdout: %v", err)
	}
	if !matched {
		t.Errorf("'kubeadm token generate' stdout did not match expected regex; wanted: [%q], got: [%s]", TokenExpectedRegex, stdout)
	}
}

func TestCmdTokenGenerateTypoError(t *testing.T) {
	/*
		Since we expect users to do things like this:

			$ TOKEN=$(kubeadm token generate)

		we want to make sure that if they have a typo in their command, we exit
		with a non-zero status code after showing the command's usage, so that
		the usage itself isn't captured as a token without the user noticing.
	*/
	if *kubeadmCmdSkip {
		t.Log("kubeadm cmd tests being skipped")
		t.Skip()
	}

	kubeadmPath := getKubeadmPath()
	_, _, _, err := RunCmd(kubeadmPath, "token", "genorate") // subtle typo
	if err == nil {
		t.Error("'kubeadm token genorate' (a deliberate typo) exited without an error when we expected non-zero exit status")
	}
}
func TestCmdTokenDelete(t *testing.T) {
	if *kubeadmCmdSkip {
		t.Log("kubeadm cmd tests being skipped")
		t.Skip()
	}

	var tests = []struct {
		name     string
		args     string
		expected bool
	}{
		{"no token provided", "", false},
		{"invalid token", "foobar", false},
	}

	kubeadmPath := getKubeadmPath()
	for _, rt := range tests {
		t.Run(rt.name, func(t *testing.T) {
			_, _, _, actual := RunCmd(kubeadmPath, "token", "delete", rt.args)
			if (actual == nil) != rt.expected {
				t.Errorf(
					"failed CmdTokenDelete running 'kubeadm token %s' with an error: %v\n\texpected: %t\n\t  actual: %t",
					rt.args,
					actual,
					rt.expected,
					(actual == nil),
				)
			}
			kubeadmReset()
		})
	}
}
