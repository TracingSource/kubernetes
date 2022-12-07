package util

import (
	"strings"

	"github.com/pkg/errors"

	utilsexec "k8s.io/utils/exec"
)

const (
	// CgroupDriverSystemd holds the systemd driver type
	CgroupDriverSystemd = "systemd"
	// CgroupDriverCgroupfs holds the cgroupfs driver type
	CgroupDriverCgroupfs = "cgroupfs"
)

// TODO: add support for detecting the cgroup driver for CRI other than
// Docker. Currently only Docker driver detection is supported:
// Discussion:
//     https://github.com/kubernetes/kubeadm/issues/844

// GetCgroupDriverDocker runs 'docker info -f "{{.CgroupDriver}}"' to obtain the docker cgroup driver
func GetCgroupDriverDocker(execer utilsexec.Interface) (string, error) {
	driver, err := callDockerInfo(execer)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(driver, "\n"), nil
}

func callDockerInfo(execer utilsexec.Interface) (string, error) {
	out, err := execer.Command("docker", "info", "-f", "{{.CgroupDriver}}").Output()
	if err != nil {
		return "", errors.Wrap(err, "cannot execute 'docker info -f {{.CgroupDriver}}'")
	}
	return string(out), nil
}
