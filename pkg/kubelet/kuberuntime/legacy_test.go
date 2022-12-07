package kuberuntime

import (
	"fmt"
	"math/rand"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func TestLogSymLink(t *testing.T) {
	as := assert.New(t)
	containerLogsDir := "/foo/bar"
	podFullName := randStringBytes(128)
	containerName := randStringBytes(70)
	dockerID := randStringBytes(80)
	// The file name cannot exceed 255 characters. Since .log suffix is required, the prefix cannot exceed 251 characters.
	expectedPath := path.Join(containerLogsDir, fmt.Sprintf("%s_%s-%s", podFullName, containerName, dockerID)[:251]+".log")
	as.Equal(expectedPath, logSymlink(containerLogsDir, podFullName, containerName, dockerID))
}

func TestLegacyLogSymLink(t *testing.T) {
	as := assert.New(t)
	containerID := randStringBytes(80)
	containerName := randStringBytes(70)
	podName := randStringBytes(128)
	podNamespace := randStringBytes(10)
	// The file name cannot exceed 255 characters. Since .log suffix is required, the prefix cannot exceed 251 characters.
	expectedPath := path.Join(legacyContainerLogsDir, fmt.Sprintf("%s_%s_%s-%s", podName, podNamespace, containerName, containerID)[:251]+".log")
	as.Equal(expectedPath, legacyLogSymlink(containerID, containerName, podName, podNamespace))
}
