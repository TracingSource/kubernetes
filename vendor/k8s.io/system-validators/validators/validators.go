package system

import (
	"fmt"
)

// Validator is the interface for all validators.
type Validator interface {
	// Name is the name of the validator.
	Name() string
	// Validate is the validate function.
	Validate(SysSpec) ([]error, []error)
}

// Reporter is the interface for the reporters for the validators.
type Reporter interface {
	// Report reports the results of the system verification
	Report(string, string, ValidationResultType) error
}

// Validate uses validators to validate the system and returns a warning or error.
func Validate(spec SysSpec, validators []Validator) ([]error, []error) {
	var errs []error
	var warns []error

	for _, v := range validators {
		fmt.Printf("Validating %s...\n", v.Name())
		warn, err := v.Validate(spec)
		if len(err) != 0 {
			errs = append(errs, err...)
		}
		if len(warn) != 0 {
			warns = append(warns, warn...)
		}
	}
	return warns, errs
}

// ValidateSpec uses all default validators to validate the system and writes to stdout.
func ValidateSpec(spec SysSpec, runtime string) ([]error, []error) {
	// OS-level validators.
	var osValidators = []Validator{
		&OSValidator{Reporter: DefaultReporter},
		&KernelValidator{Reporter: DefaultReporter},
		&CgroupsValidator{Reporter: DefaultReporter},
		&packageValidator{reporter: DefaultReporter},
	}
	// Docker-specific validators.
	var dockerValidators = []Validator{
		&DockerValidator{Reporter: DefaultReporter},
	}

	validators := osValidators
	switch runtime {
	case "docker":
		validators = append(validators, dockerValidators...)
	}
	return Validate(spec, validators)
}
