package timed

import (
	"store/pkg/sdk/conv"
	"time"
)

var (
	YYYY_MM_DD_HH_MM_SS     = "2006-01-02 15:04:05"
	YYYY_MM_DD              = "2006-01-02"
	YYYYMMDDHH              = "2006010215"
	YYYY_MM_DD_T_HH_MM_SS_Z = "2006-01-02T15:04:05Z"

	LocAsiaShanghai, _ = time.LoadLocation("Asia/Shanghai")
)

func AsTSSecond(ts int64) int64 {

	s := conv.String(ts)
	if len(s) > 10 {
		s = s[:10]
	}

	return conv.Int64(s)
}

func LastDate(date string) string {

	t, _ := time.Parse(YYYY_MM_DD, date)

	return t.AddDate(0, 0, -1).Format(YYYY_MM_DD)
}

func DateTs(date string) int64 {
	t, _ := time.Parse(YYYY_MM_DD, date)
	return t.Unix()
}

func DateEndTs(date time.Time) int64 {
	year, month, day := date.Date()
	d := time.Date(year, month, day, 23, 59, 59, 0, time.Local)
	return d.Unix()
}

func DateStartTs(date time.Time) int64 {
	year, month, day := date.Date()
	d := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	return d.Unix()
}

func DateStart(t time.Time) time.Time {
	year, month, day := t.Date()
	d := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return d
}

func DateEnd(t time.Time) time.Time {
	year, month, day := t.Date()
	d := time.Date(year, month, day, 23, 59, 59, 0, t.Location())
	return d
}
