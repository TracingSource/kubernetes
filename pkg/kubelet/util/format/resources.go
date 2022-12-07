package format

import (
	"fmt"
	"sort"
	"strings"

	"k8s.io/api/core/v1"
)

// ResourceList returns a string representation of a resource list in a human readable format.
func ResourceList(resources v1.ResourceList) string {
	resourceStrings := make([]string, 0, len(resources))
	for key, value := range resources {
		resourceStrings = append(resourceStrings, fmt.Sprintf("%v=%v", key, value.String()))
	}
	// sort the results for consistent log output
	sort.Strings(resourceStrings)
	return strings.Join(resourceStrings, ",")
}
