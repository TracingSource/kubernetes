package system

import (
	"strings"
)

// DeprecatedMightBeMasterNode returns true if given node is a registered master.
// This code must not be updated to use node role labels, since node role labels
// may not change behavior of the system.
// DEPRECATED: use a label selector provided by test initialization.
func DeprecatedMightBeMasterNode(nodeName string) bool {
	// We are trying to capture "master(-...)?$" regexp.
	// However, using regexp.MatchString() results even in more than 35%
	// of all space allocations in ControllerManager spent in this function.
	// That's why we are trying to be a bit smarter.
	if strings.HasSuffix(nodeName, "master") {
		return true
	}
	if len(nodeName) >= 10 {
		return strings.HasSuffix(nodeName[:len(nodeName)-3], "master-")
	}
	return false
}
