package dockershim

import (
	"testing"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"

	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

func TestConvertDockerStatusToRuntimeAPIState(t *testing.T) {
	testCases := []struct {
		input    string
		expected runtimeapi.ContainerState
	}{
		{input: "Up 5 hours", expected: runtimeapi.ContainerState_CONTAINER_RUNNING},
		{input: "Exited (0) 2 hours ago", expected: runtimeapi.ContainerState_CONTAINER_EXITED},
		{input: "Created", expected: runtimeapi.ContainerState_CONTAINER_CREATED},
		{input: "Random string", expected: runtimeapi.ContainerState_CONTAINER_UNKNOWN},
	}

	for _, test := range testCases {
		actual := toRuntimeAPIContainerState(test.input)
		assert.Equal(t, test.expected, actual)
	}
}

func TestConvertToPullableImageID(t *testing.T) {
	testCases := []struct {
		id       string
		image    *dockertypes.ImageInspect
		expected string
	}{
		{
			id: "image-1",
			image: &dockertypes.ImageInspect{
				RepoDigests: []string{"digest-1"},
			},
			expected: DockerPullableImageIDPrefix + "digest-1",
		},
		{
			id: "image-2",
			image: &dockertypes.ImageInspect{
				RepoDigests: []string{},
			},
			expected: DockerImageIDPrefix + "image-2",
		},
	}

	for _, test := range testCases {
		actual := toPullableImageID(test.id, test.image)
		assert.Equal(t, test.expected, actual)
	}
}
