package node

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntSet(t *testing.T) {
	i := newIntSet()

	assert.False(t, i.has(1))
	assert.False(t, i.has(2))
	assert.False(t, i.has(3))
	assert.False(t, i.has(4))

	i.startNewGeneration()
	i.mark(1)
	i.mark(2)
	i.sweep()

	assert.True(t, i.has(1))
	assert.True(t, i.has(2))
	assert.False(t, i.has(3))
	assert.False(t, i.has(4))

	i.startNewGeneration()
	i.mark(2)
	i.mark(3)
	i.sweep()

	assert.False(t, i.has(1))
	assert.True(t, i.has(2))
	assert.True(t, i.has(3))
	assert.False(t, i.has(4))

	i.startNewGeneration()
	i.mark(3)
	i.mark(4)
	i.sweep()

	assert.False(t, i.has(1))
	assert.False(t, i.has(2))
	assert.True(t, i.has(3))
	assert.True(t, i.has(4))
}
