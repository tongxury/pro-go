package helper

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	//rsp := Combine1([]int{1, 2, 3, 4, 5, 6}, 2)
	rsp := Combine([]string{"a", "b", "c", "d", "e", "f"}, 2)

	fmt.Println(rsp)
}
