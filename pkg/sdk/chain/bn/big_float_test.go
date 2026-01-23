package bn

import (
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestBigFloat_Val(t *testing.T) {

	cases := []struct {
		bg   *BigFloat
		want float64
	}{
		{bg: MustFloat(big.NewFloat(35749377858), 6), want: 35749.377858},
		{bg: MustFloat(big.NewFloat(100), 6), want: 0.0001},
		{bg: MustFloat(big.NewFloat(0), 0), want: 0},
		{bg: MustFloat(big.NewFloat(0), 10), want: 0},
		{bg: MustFloat(big.NewFloat(10), 0), want: 10},
		{bg: MustFloat(big.NewFloat(100000000), 6), want: 100},
	}

	for _, c := range cases {
		assert.Equal(t, c.want, c.bg.Val())
	}

}

func TestBigFloat_New(t *testing.T) {

	cases := []struct {
		bg   *BigFloat
		want float64
	}{
		{bg: MustFloatFromFloat(0.00014350000299145904, 9), want: 0.00014350000299145904},
		{bg: MustFloat("35749377858", 6), want: 35749.377858},
		{bg: MustFloat(100, 6), want: 0.0001},
		{bg: MustFloat(uint64(100), 6), want: 0.0001},
		{bg: MustFloat(int64(100), 6), want: 0.0001},
		{bg: MustFloat(int8(100), 6), want: 0.0001},
	}

	for _, c := range cases {
		assert.Equal(t, c.want, c.bg.Val())
	}

}
