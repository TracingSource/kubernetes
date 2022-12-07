// +build darwin

package preflight

// This is a MacOS stub

// Check validates if Docker is setup to use systemd as the cgroup driver.
// No-op for Darwin (MacOS).
func (idsc IsDockerSystemdCheck) Check() (warnings, errorList []error) {
	return nil, nil
}
