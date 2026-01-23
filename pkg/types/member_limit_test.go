package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestName(t *testing.T) {

	cases := []struct {
		limit     *Limit
		dayUsed   int64
		monthUsed int64
		totalUsed int64
		expect    bool
	}{
		{limit: &Limit{Day: 2}, dayUsed: 1, monthUsed: 1, totalUsed: 1, expect: true},
		{limit: &Limit{Day: 2}, dayUsed: 1, monthUsed: 2, totalUsed: 2, expect: true},
		{limit: &Limit{Day: 2}, dayUsed: 2, monthUsed: 3, totalUsed: 3, expect: false},
		{limit: &Limit{Day: 2}, dayUsed: 2, monthUsed: 3, totalUsed: 3, expect: false},

		{limit: &Limit{Month: 3}, dayUsed: 1, monthUsed: 2, totalUsed: 2, expect: true},
		{limit: &Limit{Month: 3}, dayUsed: 1, monthUsed: 3, totalUsed: 3, expect: false},
		{limit: &Limit{Month: 3}, dayUsed: 1, monthUsed: 2, totalUsed: 2, expect: true},

		{limit: &Limit{Total: 3}, dayUsed: 1, monthUsed: 2, totalUsed: 2, expect: true},
		{limit: &Limit{Total: 3}, dayUsed: 1, monthUsed: 1, totalUsed: 3, expect: false},
		{limit: &Limit{Total: 3}, dayUsed: 1, monthUsed: 1, totalUsed: 4, expect: false},

		{limit: &Limit{Day: 2, Month: 5}, dayUsed: 1, monthUsed: 5, totalUsed: 5, expect: false},
		{limit: &Limit{Day: 2, Month: 5}, dayUsed: 2, monthUsed: 3, totalUsed: 5, expect: false},
		{limit: &Limit{Day: 2, Month: 5}, dayUsed: 1, monthUsed: 3, totalUsed: 5, expect: true},
	}

	for i, x := range cases {
		assert.Equal(t, x.expect, x.limit.Check(x.dayUsed, x.monthUsed, x.totalUsed), i)
	}

}
