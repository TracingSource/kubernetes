package system

import (
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

var _ Validator = &OSValidator{}

// OSValidator validates OS.
type OSValidator struct {
	Reporter Reporter
}

// Name is part of the system.Validator interface.
func (o *OSValidator) Name() string {
	return "os"
}

// Validate is part of the system.Validator interface.
func (o *OSValidator) Validate(spec SysSpec) ([]error, []error) {
	os, err := exec.Command("uname").CombinedOutput()
	if err != nil {
		return nil, []error{errors.Wrap(err, "failed to get os name")}
	}
	if err = o.validateOS(strings.TrimSpace(string(os)), spec.OS); err != nil {
		return nil, []error{err}
	}
	return nil, nil
}

func (o *OSValidator) validateOS(os, specOS string) error {
	if os != specOS {
		o.Reporter.Report("OS", os, bad)
		return errors.Errorf("unsupported operating system: %s", os)
	}
	o.Reporter.Report("OS", os, good)
	return nil
}
