package timed

import (
	"fmt"
	"math"
	"strings"
	"time"
)

/**
* @desc 时间转换函数
* @param atime string 要转换的时间戳（秒）
* @return string
 */
func SmartTime(atime int64) string {

	var byTime = []int64{365 * 24 * 60 * 60, 24 * 60 * 60, 60 * 60, 60, 1}
	var unit = []string{"年", "天", "小时", "分钟", "秒钟"}
	now := time.Now().Unix()

	delta := now - atime

	if delta <= 0 {
		return "刚刚"
	}

	if delta > 5*24*60*60 {
		return time.Unix(atime, 0).Format(time.DateOnly)
	}

	var parts []string

	for i := 0; i < len(byTime); i++ {
		if delta < byTime[i] {
			continue
		}
		var temp = math.Floor(float64(delta / byTime[i]))

		delta = delta % byTime[i]
		if temp > 0 {
			parts = append(parts, fmt.Sprintf("%.0f%s", temp, unit[i]))
		}
	}

	if len(parts) > 1 {
		parts = parts[:1]
	}

	return strings.Join(parts, "") + "前"
}
