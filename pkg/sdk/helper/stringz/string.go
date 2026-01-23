package stringz

import (
	"strconv"
	"strings"
)

func HasAnySuffix(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.HasSuffix(s, sub) {
			return true
		}
	}
	return false
}

func IsFloat(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func IsInt(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64) // 使用基数 10
	return err == nil
}

func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func SplitToSlice(str string, length int) []string {
	if length <= 0 {
		return []string{str}
	}

	var result []string
	runes := []rune(str)

	for i := 0; i < len(runes); i += length {
		end := i + length
		if end > len(runes) {
			end = len(runes)
		}
		result = append(result, string(runes[i:end]))
	}
	return result
}
