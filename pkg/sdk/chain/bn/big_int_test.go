package bn

import (
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestBigInt_Val(t *testing.T) {

	cases := []struct {
		bg   *BigInt
		want float64
	}{
		{bg: MustInt(big.NewInt(35749377858), 6), want: 35749.377858},
		{bg: MustInt(big.NewInt(100), 6), want: 0.0001},
		{bg: MustInt(big.NewInt(0), 0), want: 0},
		{bg: MustInt(big.NewInt(0), 10), want: 0},
		{bg: MustInt(big.NewInt(10), 0), want: 10},
		{bg: MustInt(big.NewInt(100000000), 6), want: 100},
	}

	for _, c := range cases {
		assert.Equal(t, c.want, c.bg.Val())
	}

}

func TestBigInt_New(t *testing.T) {

	cases := []struct {
		bg   *BigInt
		want float64
	}{
		{bg: MustInt("35749377858", 6), want: 35749.377858},
		{bg: MustInt(100, 6), want: 0.0001},
		{bg: MustInt(uint64(100), 6), want: 0.0001},
		{bg: MustInt(int64(100), 6), want: 0.0001},
		{bg: MustInt(int8(100), 6), want: 0.0001},
	}

	for _, c := range cases {
		assert.Equal(t, c.want, c.bg.Val())
	}

}

func TestBigInt_Raw(t *testing.T) {

	cases := []struct {
		bg   *BigInt
		want string
	}{
		//{bg: MustIntFromInt(35749.377858, 6), want: "35749377858"},
	}

	for _, c := range cases {
		assert.Equal(t, c.want, c.bg.RawStr())
	}

}
