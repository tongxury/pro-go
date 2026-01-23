package helper

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitSlice(t *testing.T) {

	src := []int{1, 2, 3, 4, 4, 5, 6, 8, 9, 9, 4}

	r := SubSlice(src, 3)

	fmt.Println(r)
}

func TestSliceContainsAny(t *testing.T) {

	assert.True(t, SliceContainsAny([]string{"a", "b", "c"}, "a", "d"))
	assert.False(t, SliceContainsAny([]string{"a", "b", "c"}, "d", "e"))
}

func TestContains(t *testing.T) {

	assert.True(t, Contains([]string{"a", "b", "c"}, "a"))
	assert.False(t, Contains([]string{"a", "b", "c"}, "d"))
}
