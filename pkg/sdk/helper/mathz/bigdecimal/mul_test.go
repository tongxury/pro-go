package bigdecimal

import (
	"fmt"
	"testing"
)

func TestMul(t *testing.T) {

	a := New(0.001, 9).MulFloat64(0.01)

	fmt.Println(a.Float64())
}
