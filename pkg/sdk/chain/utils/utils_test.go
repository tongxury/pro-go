package chainutils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShortenAddress(t *testing.T) {

	r := ShortenAddress("6QE4iDjAow9c4pDAaLWXSav7hoj4KzHSkT5j56ag7UBn")
	fmt.Println(r)
}

func TestFormatDecimal(t *testing.T) {

	assert.Equal(t, "0.{12}6622", FormatDecimal(0.0000000000006622))
	assert.Equal(t, "0.{7}83481", FormatDecimal(0.00000008348155850669375))
	assert.Equal(t, "0.{7}834815", FormatDecimal(0.00000008348155850669375, FormatDecimalOptions{MaxSignificantCount: 6}))
	assert.Equal(t, "0.{5}66", FormatDecimal(0.0000066))
	assert.Equal(t, "0.{5}66", FormatDecimal(0.000006600))

	fmt.Println(FormatDecimal(10000))
}
