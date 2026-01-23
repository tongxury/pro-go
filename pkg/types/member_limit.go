package types

import helpers "store/pkg/sdk/helper"

type Limit struct {
	Month int64  `json:"month,omitempty"`
	Day   int64  `json:"day,omitempty"`
	Total int64  `json:"total,omitempty"`
	Bonus int64  `json:"bonus,omitempty"`
	Desc  string `json:"desc,omitempty"`
}

func (t *Limit) Value() int64 {
	if t == nil {
		return 0
	}
	return helpers.OrInt64(t.Total, t.Month, t.Day)
}

func (t *Limit) Check(dayUsed, monthUsed, totalUsed int64) bool {
	return t.CheckDayLimit(dayUsed) && t.CheckMonthLimit(monthUsed) && t.CheckTotalLimit(totalUsed)
}

func (t *Limit) CheckDayLimit(used int64) bool {
	if t.Day < 0 {
		return false
	} else if t.Day == 0 {
		return true
	} else {
		return used <= t.Day
	}
}

func (t *Limit) CheckMonthLimit(used int64) bool {
	if t.Month < 0 {
		return false
	} else if t.Month == 0 {
		return true
	} else {
		return used <= t.Month
	}
}

func (t *Limit) CheckTotalLimit(used int64) bool {
	if t.Total < 0 {
		return false
	} else if t.Total == 0 {
		return true
	} else {
		return used <= t.Total
	}
}
