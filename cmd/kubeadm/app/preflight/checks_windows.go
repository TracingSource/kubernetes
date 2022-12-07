// +build windows

package preflight

import (
	"os/user"

	"github.com/pkg/errors"
)

// The "Well-known SID" of Administrator group
// https://support.microsoft.com/en-us/help/243330/well-known-security-identifiers-in-windows-operating-systems
const administratorSID = "S-1-5-32-544"

// Check validates if a user has elevated (administrator) privileges.
func (ipuc IsPrivilegedUserCheck) Check() (warnings, errorList []error) {
	currUser, err := user.Current()
	if err != nil {
		return nil, []error{errors.New("cannot get current user")}
	}

	groupIds, err := currUser.GroupIds()
	if err != nil {
		return nil, []error{errors.New("cannot get group IDs for current user")}
	}

	for _, sid := range groupIds {
		if sid == administratorSID {
			return nil, nil
		}
	}

	return nil, []error{errors.New("user is not running as administrator")}
}

// Check validates if Docker is setup to use systemd as the cgroup driver.
// No-op for Windows.
func (idsc IsDockerSystemdCheck) Check() (warnings, errorList []error) {
	return nil, nil
}
