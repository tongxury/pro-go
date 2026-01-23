package stringz

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {

	for i := 0; i < 10; i++ {
		code := "RN_" + RandString(8)
		fmt.Println(fmt.Sprintf("https://knowee.ai/discount?promo=%s, %s", code, code))
	}

	for i := 0; i < 10; i++ {
		code := "WC_" + RandString(8)
		fmt.Println(fmt.Sprintf("https://knowee.ai/discount?promo=%s, %s", code, code))
	}
}
