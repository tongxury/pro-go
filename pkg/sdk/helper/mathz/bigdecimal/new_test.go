package bigdecimal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {

	cases := []struct {
		value     float64
		decimals  int
		expected1 float64
		expected2 string
	}{
		{
			value:     1,
			decimals:  2,
			expected1: 1,
			expected2: "100",
		},
	}

	for _, c := range cases {
		x := New(c.value, c.decimals)

		assert.Equal(t, c.expected1, x.Float64())
		assert.Equal(t, c.expected2, x.Raw())
	}
}
