package jsonz

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type Demo struct {
	V1 string
	V2 float64
}

func TestAPI_MapToStruct(t *testing.T) {

	api := New()

	src := map[string]interface{}{
		"V1": "v1",
		"V2": 1.0,
	}

	var d Demo

	err := api.MapToStruct(src, &d)

	assert.Nil(t, err)

	assert.Equal(t, d.V1, "v1")
	assert.Equal(t, d.V2, 1.0)

}
