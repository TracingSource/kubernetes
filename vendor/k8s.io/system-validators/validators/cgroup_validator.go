package system

import (
	"bufio"
	"os"
	"strings"

	"github.com/pkg/errors"
)

var _ Validator = &CgroupsValidator{}

// CgroupsValidator validates cgroup configuration.
type CgroupsValidator struct {
	Reporter Reporter
}

// Name is part of the system.Validator interface.
func (c *CgroupsValidator) Name() string {
	return "cgroups"
}

const (
	cgroupsConfigPrefix = "CGROUPS_"
)

// Validate is part of the system.Validator interface.
func (c *CgroupsValidator) Validate(spec SysSpec) ([]error, []error) {
	subsystems, err := c.getCgroupSubsystems()
	if err != nil {
		return nil, []error{errors.Wrap(err, "failed to get cgroup subsystems")}
	}
	if err = c.validateCgroupSubsystems(spec.Cgroups, subsystems); err != nil {
		return nil, []error{err}
	}
	return nil, nil
}

func (c *CgroupsValidator) validateCgroupSubsystems(cgroupSpec, subsystems []string) error {
	missing := []string{}
	for _, cgroup := range cgroupSpec {
		found := false
		for _, subsystem := range subsystems {
			if cgroup == subsystem {
				found = true
				break
			}
		}
		item := cgroupsConfigPrefix + strings.ToUpper(cgroup)
		if found {
			c.Reporter.Report(item, "enabled", good)
		} else {
			c.Reporter.Report(item, "missing", bad)
			missing = append(missing, cgroup)
		}
	}
	if len(missing) > 0 {
		return errors.Errorf("missing cgroups: %s", strings.Join(missing, " "))
	}
	return nil

}

func (c *CgroupsValidator) getCgroupSubsystems() ([]string, error) {
	f, err := os.Open("/proc/cgroups")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	subsystems := []string{}
	s := bufio.NewScanner(f)
	for s.Scan() {
		if err := s.Err(); err != nil {
			return nil, err
		}
		text := s.Text()
		if text[0] != '#' {
			parts := strings.Fields(text)
			if len(parts) >= 4 && parts[3] != "0" {
				subsystems = append(subsystems, parts[0])
			}
		}
	}
	return subsystems, nil
}
