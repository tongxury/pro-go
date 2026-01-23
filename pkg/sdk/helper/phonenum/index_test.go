package phonenum

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValid(t *testing.T) {

	assert.True(t, IsValid("13564345332"))
	assert.True(t, !IsValid("1356434533"))
	assert.True(t, !IsValid("12564345332"))
	assert.True(t, !IsValid("135643453321"))
	assert.True(t, !IsValid("23564345332"))
}
