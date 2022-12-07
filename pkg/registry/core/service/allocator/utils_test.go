package allocator

import (
	"math/big"
	"testing"
)

func TestBitCount(t *testing.T) {
	for i, c := range bitCounts {
		actual := 0
		for j := 0; j < 8; j++ {
			if ((1 << uint(j)) & i) != 0 {
				actual++
			}
		}
		if actual != int(c) {
			t.Errorf("%d should have %d bits but recorded as %d", i, actual, c)
		}
	}
}

func TestCountBits(t *testing.T) {
	tests := []struct {
		n        *big.Int
		expected int
	}{
		{n: big.NewInt(int64(0)), expected: 0},
		{n: big.NewInt(int64(0xffffffffff)), expected: 40},
	}
	for _, test := range tests {
		actual := countBits(test.n)
		if test.expected != actual {
			t.Errorf("%s should have %d bits but recorded as %d", test.n, test.expected, actual)
		}
	}
}
