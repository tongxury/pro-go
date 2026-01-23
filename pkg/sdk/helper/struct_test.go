package helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestStruct struct {
	V1 string
	V2 string
}

func TestMergeStructs(t *testing.T) {

	var v1 = &TestStruct{V1: "1", V2: "2"}

	var v2 = &TestStruct{V1: "11", V2: "22"}
	var v3 = &TestStruct{V1: "11", V2: "222"}
	var v4 = &TestStruct{V1: "11"}

	r1 := MergeStructs[TestStruct](v1, nil, v2)
	assert.Equal(t, "11", r1.V1)
	assert.Equal(t, "22", r1.V2)

	r2 := MergeStructs[TestStruct](v1, v2, v3)
	assert.Equal(t, "11", r2.V1)
	assert.Equal(t, "222", r2.V2)

	r3 := MergeStructs[TestStruct](v1, v3, v4)
	assert.Equal(t, "11", r3.V1)
	assert.Equal(t, "222", r3.V2)
}
