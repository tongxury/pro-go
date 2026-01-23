package conv

import (
	"fmt"
	"strconv"
	"strings"
)

func Str(value any) string {

	switch v := value.(type) {
	case string:
		return v
	default:
		return fmt.Sprintf("%v", value)
	}
}

// Duplicated
// case 目前遇到一个加一个
func String(e interface{}) string {
	var res string
	switch v := e.(type) {
	case string:
		res = v
	case float32:
		res = strconv.FormatFloat(Float64(v), 'f', 6, 64)
	case float64:
		res = strconv.FormatFloat(v, 'f', 6, 64)
	case int:
		res = strconv.Itoa(v)
	case int8:
		res = strconv.Itoa(Int(v))
	case int32, int64:
		res = strconv.Itoa(Int(v))
	//case int64:
	//	res = strconv.Itoa(Int(v))
	case uint64:
		res = strconv.Itoa(Int(v))
	default:
		value, ok := e.(string)
		if ok {
			res = value
		}
	}

	return res
}

// snake string, XxYy to xx_yy , XxYY to xx_yy
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

// camel string, xx_yy to XxYy
func CamelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}
