package printers

import (
	"strings"
)

var (
	InternalObjectPrinterErr = "a versioned object must be passed to a printer"

	// disallowedPackagePrefixes contains regular expression templates
	// for object package paths that are not allowed by printers.
	disallowedPackagePrefixes = []string{
		"k8s.io/kubernetes/pkg/apis/",
	}
)

var InternalObjectPreventer = &illegalPackageSourceChecker{disallowedPackagePrefixes}

func IsInternalObjectError(err error) bool {
	if err == nil {
		return false
	}

	return err.Error() == InternalObjectPrinterErr
}

// illegalPackageSourceChecker compares a given
// object's package path, and determines if the
// object originates from a disallowed source.
type illegalPackageSourceChecker struct {
	// disallowedPrefixes is a slice of disallowed package path
	// prefixes for a given runtime.Object that we are printing.
	disallowedPrefixes []string
}

func (c *illegalPackageSourceChecker) IsForbidden(pkgPath string) bool {
	for _, forbiddenPrefix := range c.disallowedPrefixes {
		if strings.HasPrefix(pkgPath, forbiddenPrefix) || strings.Contains(pkgPath, "/vendor/"+forbiddenPrefix) {
			return true
		}
	}

	return false
}
