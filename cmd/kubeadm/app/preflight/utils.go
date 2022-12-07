package preflight

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"k8s.io/apimachinery/pkg/util/version"
	utilsexec "k8s.io/utils/exec"
)

// GetKubeletVersion is helper function that returns version of kubelet available in $PATH
func GetKubeletVersion(execer utilsexec.Interface) (*version.Version, error) {
	kubeletVersionRegex := regexp.MustCompile(`^\s*Kubernetes v((0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)([-0-9a-zA-Z_\.+]*)?)\s*$`)

	command := execer.Command("kubelet", "--version")
	out, err := command.Output()
	if err != nil {
		return nil, errors.Wrap(err, "cannot execute 'kubelet --version'")
	}

	cleanOutput := strings.TrimSpace(string(out))
	subs := kubeletVersionRegex.FindAllStringSubmatch(cleanOutput, -1)
	if len(subs) != 1 || len(subs[0]) < 2 {
		return nil, errors.Errorf("Unable to parse output from Kubelet: %q", cleanOutput)
	}
	return version.ParseSemantic(subs[0][1])
}
