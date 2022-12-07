package kubeadm

import "testing"

func TestCmdCompletion(t *testing.T) {
	kubeadmPath := getKubeadmPath()

	if *kubeadmCmdSkip {
		t.Log("kubeadm cmd tests being skipped")
		t.Skip()
	}

	var tests = []struct {
		name     string
		args     string
		expected bool
	}{
		{"shell not expected", "", false},
		{"unsupported shell type", "foo", false},
	}

	for _, rt := range tests {
		t.Run(rt.name, func(t *testing.T) {
			_, _, _, actual := RunCmd(kubeadmPath, "completion", rt.args)
			if (actual == nil) != rt.expected {
				t.Errorf(
					"failed CmdCompletion running 'kubeadm completion %s' with an error: %v\n\texpected: %t\n\t  actual: %t",
					rt.args,
					actual,
					rt.expected,
					(actual == nil),
				)
			}
		})
	}
}
