package timed

import (
	"fmt"
	"regexp"
	"strconv"
)

/*
	^(\d+d)?(\d+h)?(\d+m)?$
 	1d1h1m、10d、5h、30m 以及任何组合形式，比如 1d2h 或 1h30m 等等。

	time.ParseDuration 没有对d的支持
*/

func ParseDuration(input string) (int64, error) {
	if input == "" {
		return 0, nil
	}

	regex := regexp.MustCompile(`^(?:(\d+d)|(\d+h)|(\d+m))+$`)
	matches := regex.FindStringSubmatch(input)

	if len(matches) == 4 {

		var totalDuration int64

		if days := matches[1]; days != "" {
			d, _ := strconv.Atoi(days[0 : len(days)-1])
			totalDuration += int64(d) * 24 * 60 * 60
		}
		if hours := matches[2]; hours != "" {
			h, _ := strconv.Atoi(hours[0 : len(hours)-1])
			totalDuration += int64(h) * 60 * 60
		}
		if minutes := matches[3]; minutes != "" {
			m, _ := strconv.Atoi(minutes[0 : len(minutes)-1])
			totalDuration += int64(m) * 60
		}

		return totalDuration, nil
	}

	return 0, fmt.Errorf("invalid duration: %s", input)
}
