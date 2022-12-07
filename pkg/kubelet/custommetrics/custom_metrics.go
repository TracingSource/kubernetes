// Package custommetrics contains support for instrumenting cAdvisor to gather custom metrics from pods.
package custommetrics

import (
	"path"

	"k8s.io/api/core/v1"
)

const (
	// CustomMetricsDefinitionContainerFile is the file in container that stores Custom Metrics definition
	CustomMetricsDefinitionContainerFile = "definition.json"

	// CustomMetricsDefinitionDir is the dir where Custom Metrics definition is stored
	CustomMetricsDefinitionDir = "/etc/custom-metrics"
)

// GetCAdvisorCustomMetricsDefinitionPath returns a path to a cAdvisor-specific custom metrics configuration.
// Alpha implementation.
func GetCAdvisorCustomMetricsDefinitionPath(container *v1.Container) (*string, error) {
	// Assumes that the container has Custom Metrics enabled if it has "/etc/custom-metrics" directory
	// mounted as a volume. Custom Metrics definition is expected to be in "definition.json".
	if container.VolumeMounts != nil {
		for _, volumeMount := range container.VolumeMounts {
			if path.Clean(volumeMount.MountPath) == path.Clean(CustomMetricsDefinitionDir) {
				// TODO: add definition file validation.
				definitionPath := path.Clean(path.Join(volumeMount.MountPath, CustomMetricsDefinitionContainerFile))
				return &definitionPath, nil
			}
		}
	}
	// No Custom Metrics definition available.
	return nil, nil
}
