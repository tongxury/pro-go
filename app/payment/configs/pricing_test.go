package configs

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {

	a := GetPlanById("com_veogo_l2_month")

	fmt.Println(a)
}
