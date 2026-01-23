package mathz

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMin(t *testing.T) {

	assert.True(t, Min(1, 1) == 1)
	assert.True(t, Min(1, 2) == 1)
	assert.True(t, Min(1, 5) == 1)
}

func TestToFixed(t *testing.T) {

	//fmt.Println(ToFixed4(0.342345))
	//fmt.Println(ToFixed4(20000.0000000))
	//values := []float64{
	//	0.000333335,
	//	0.33333,
	//	0.00012345,
	//	0.000123456,
	//	0.100005,
	//	1.234567,
	//	0.00000012345,
	//}
	//
	//for _, value := range values {
	//	fmt.Printf("Original: %.10f => Formatted: %s\n", value, FormatWithFiveSignificantDecimals(value))
	//}
}

func TestCmn(t *testing.T) {

	fmt.Println(Cmn(8, 4))
}
