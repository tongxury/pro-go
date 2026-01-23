package timed

import (
	"fmt"
	"testing"
	"unicode/utf8"
)

func TestName(t *testing.T) {
	fmt.Println(utf8.RuneCountInString("正念练习"))
}
