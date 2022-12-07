// +build linux

package kubelet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIPTablesMark(t *testing.T) {
	tests := []struct {
		bit    int
		expect string
	}{
		{
			14,
			"0x00004000/0x00004000",
		},
		{
			15,
			"0x00008000/0x00008000",
		},
	}
	for _, tc := range tests {
		res := getIPTablesMark(tc.bit)
		assert.Equal(t, tc.expect, res, "input %d", tc.bit)
	}
}
