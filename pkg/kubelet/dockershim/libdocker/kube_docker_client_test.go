package libdocker

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsContainerNotFoundError(t *testing.T) {
	// Expected error message from docker.
	containerNotFoundError := fmt.Errorf("Error response from daemon: No such container: 96e914f31579e44fe49b239266385330a9b2125abeb9254badd9fca74580c95a")
	otherError := fmt.Errorf("Error response from daemon: Other errors")

	assert.True(t, IsContainerNotFoundError(containerNotFoundError))
	assert.False(t, IsContainerNotFoundError(otherError))
}
