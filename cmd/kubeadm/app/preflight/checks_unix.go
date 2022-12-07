// +build !windows

package preflight

import (
	"os"

	"github.com/pkg/errors"
)

// Check validates if an user has elevated (root) privileges.
func (ipuc IsPrivilegedUserCheck) Check() (warnings, errorList []error) {
	if os.Getuid() != 0 {
		return nil, []error{errors.New("user is not running as root")}
	}

	return nil, nil
}
