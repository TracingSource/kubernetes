package custommetrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/api/core/v1"
)

func TestGetCAdvisorCustomMetricsDefinitionPath(t *testing.T) {

	regularContainer := &v1.Container{
		Name: "test_container",
	}

	cmContainer := &v1.Container{
		Name: "test_container",
		VolumeMounts: []v1.VolumeMount{
			{
				Name:      "cm",
				MountPath: CustomMetricsDefinitionDir,
			},
		},
	}
	path, err := GetCAdvisorCustomMetricsDefinitionPath(regularContainer)
	assert.Nil(t, path)
	assert.NoError(t, err)

	path, err = GetCAdvisorCustomMetricsDefinitionPath(cmContainer)
	assert.NotEmpty(t, *path)
	assert.NoError(t, err)
}
