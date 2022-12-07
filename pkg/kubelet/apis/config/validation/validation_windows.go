// +build windows

package validation

import (
	"fmt"

	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	kubeletconfig "k8s.io/kubernetes/pkg/kubelet/apis/config"
)

// validateKubeletOSConfiguration validates os specific kubelet configuration and returns an error if it is invalid.
func validateKubeletOSConfiguration(kc *kubeletconfig.KubeletConfiguration) error {
	message := "invalid configuration: %v (%v) %v is not supported on Windows"
	allErrors := []error{}

	if kc.CgroupsPerQOS {
		allErrors = append(allErrors, fmt.Errorf(message, "CgroupsPerQOS", "--cgroups-per-qos", kc.CgroupsPerQOS))
	}

	if len(kc.EnforceNodeAllocatable) > 0 {
		allErrors = append(allErrors, fmt.Errorf(message, "EnforceNodeAllocatable", "--enforce-node-allocatable", kc.EnforceNodeAllocatable))
	}

	return utilerrors.NewAggregate(allErrors)
}
