package helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type SomeStruct struct{}

func TestOr(t *testing.T) {

	//assert.Equal(t, 1, Or(0, int64(0), d, 1))
}

func TestIsNil(t *testing.T) {
	var a SomeStruct
	assert.Equal(t, true, IsNil(a))
}
