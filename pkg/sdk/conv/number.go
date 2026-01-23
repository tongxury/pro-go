package conv

import (
	"encoding/json"
	"strconv"
)

func Uint64(e interface{}, defaultValue ...uint64) uint64 {
	return uint64(Int64(e))
}

func Float64(e interface{}) float64 {

	var res float64
	switch v := e.(type) {
	case float32:
		res = float64(v)
		break
	case float64:
		res = v
		break
	case int:
		res = float64(v)
		break
	case int16:
		res = float64(v)
		break
	case int32:
		res = float64(v)
		break
	case int64:
		res = float64(v)
		break
	case string:
		if e.(string) != "" {
			result, err := strconv.ParseFloat(e.(string), 64)
			if err != nil {
			} else {
				res = result
			}
		}
	}

	return res
}

func Int64(num interface{}, defaultValue ...int64) int64 {

	var rsp int64
	var err error

	switch t := num.(type) {
	case string:
		rsp, err = strconv.ParseInt(t, 10, 64)
	case int:
		rsp = int64(t)
	case int8:
		rsp = int64(t)
	case int16:
		rsp = int64(t)
	case int32:
		rsp = int64(t)
	case int64:
		rsp = t
	case uint64:
		rsp = int64(t)
	case uint8:
		rsp = int64(t)
	case uint32:
		rsp = int64(t)
	case float32:
		rsp = int64(t)
	case float64:
		rsp = int64(t)
	default:
	}
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
	}
	return rsp
}

func Int(e interface{}) int {
	var res int
	switch v := e.(type) {
	case int:
		res = v
		break
	case float32:
		res = int(v)
		break
	case float64:
		res = int(v)
		break
	case int8:
		res = int(v)
		break
	case int16:
		res = int(v)
		break
	case int32:
		res = int(v)
		break
	case int64:
		res = int(v)
		break
	case uint64:
		res = int(v)
		break
	case string:
		res, _ = strconv.Atoi(e.(string))
		break
	case json.Number:
		i, err := e.(json.Number).Int64()
		if err == nil {
			res = int(i)
		}
		break
	}
	return res
}
