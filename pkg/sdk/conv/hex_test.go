package conv

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHexToBigInt(t *testing.T) {

	tests := []struct {
		test string
		exp  string
	}{
		//{"0x2df3d3b524cd3b394", int64(52979566816118686000)},
		{"0x2df3d3b524cd3b394", "52979566816118682516"},
	}

	for _, x := range tests {

		got := HexToBigInt(x.test)
		assert.Equal(t, got.String(), x.exp)
	}
}
