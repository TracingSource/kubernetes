// +build linux

package preflight

import (
	"github.com/pkg/errors"
	"k8s.io/kubernetes/cmd/kubeadm/app/util"
	"k8s.io/utils/exec"
)

// Check validates if Docker is setup to use systemd as the cgroup driver.
func (idsc IsDockerSystemdCheck) Check() (warnings, errorList []error) {
	driver, err := util.GetCgroupDriverDocker(exec.New())
	if err != nil {
		return nil, []error{err}
	}
	if driver != util.CgroupDriverSystemd {
		err = errors.Errorf("detected %q as the Docker cgroup driver. "+
			"The recommended driver is %q. "+
			"Please follow the guide at https://kubernetes.io/docs/setup/cri/",
			driver,
			util.CgroupDriverSystemd)
		return []error{err}, nil
	}
	return nil, nil
}
