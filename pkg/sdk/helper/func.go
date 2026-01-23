package helper

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"reflect"
	"store/pkg/sdk/helper/mathz"
	"strings"
	"time"
)

func Pointer[T any](v T) *T {
	return &v
}

func Or_[T any](items ...*T) *T {

	rsp := items[len(items)-1]

	for _, item := range items {
		if !IsNil(item) {
			return item
		}
	}

	return rsp
}

func Choose[T any](b bool, trueValue T, falseValue T) T {
	if b {
		return trueValue
	} else {
		return falseValue
	}
}

func CreateUUID() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	tmp_b := base64.URLEncoding.EncodeToString(b)
	h := md5.New()
	h.Write([]byte(tmp_b))
	return strings.ToLower(hex.EncodeToString(h.Sum(nil)))
}

func Or(params ...interface{}) interface{} {
	if len(params) == 0 {
		return nil
	}
	rsp := params[len(params)-1]
	for _, v := range params {

		switch t := v.(type) {
		case time.Duration:
			if t > 0 {
				return t
			}
		case string:
			if t != "" {
				return t
			}
		case []string:
			if len(t) > 0 {
				return t
			}
		case []int:
			if len(t) > 0 {
				return t
			}
		case []int8:
			if len(t) > 0 {
				return t
			}
		case []int32:
			if len(t) > 0 {
				return t
			}
		case []int64:
			if len(t) > 0 {
				return t
			}
		case []float32:
			if len(t) > 0 {
				return t
			}
		case []float64:
			if len(t) > 0 {
				return t
			}
		}
	}

	return rsp
}

func OrStringSlice(params ...[]string) []string {
	if len(params) == 0 {
		return nil
	}
	rsp := params[len(params)-1]
	for _, v := range params {
		if len(v) != 0 {
			return v
		}
	}
	return rsp
}

func OrString(params ...string) string {
	if len(params) == 0 {
		return ""
	}
	rsp := params[len(params)-1]
	for _, v := range params {
		if v != "" {
			return v
		}
	}
	return rsp
}

func OrNumber[T mathz.Number](params ...T) T {
	if len(params) == 0 {
		return 0
	}
	rsp := params[len(params)-1]
	for _, v := range params {
		if v > 0 {
			return v
		}
	}
	return rsp
}

func OrInt64(params ...int64) int64 {
	if len(params) == 0 {
		return 0
	}
	rsp := params[len(params)-1]
	for _, v := range params {
		if v > 0 {
			return v
		}
	}
	return rsp
}

func CreateMap(params ...interface{}) map[string]interface{} {
	rsp := make(map[string]interface{})

	for index := 0; index < len(params)-1; index = index + 2 {
		if index%2 == 0 {
			rsp[params[index].(string)] = params[index+1]
		}
	}
	return rsp
}

func Select[T any](expect bool, trueValue T, falseValue T) T {
	if expect {
		return trueValue
	}
	return falseValue
}

//
//func SelectString(expect bool, trueValue string, falseValue string) string {
//	return Select[string](expect, trueValue, falseValue)
//}
//
//func SelectInt64(expect bool, trueValue int64, falseValue int64) int64 {
//	return Select[int64](expect, trueValue, falseValue)
//}
//
//func SelectInt(expect bool, trueValue int, falseValue int) int {
//	return Select(expect, trueValue, falseValue).(int)
//}

func Equals(a interface{}, b interface{}) bool {
	aType := reflect.TypeOf(a).String()
	bType := reflect.TypeOf(b).String()

	if aType != bType {
		return false
	}

	rsp := a == b
	return rsp
}

func NotEquals(a interface{}, b interface{}) bool {
	return !Equals(a, b)
}

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	vi := reflect.ValueOf(i)
	return vi.Kind() == reflect.Ptr && vi.IsNil()
}

func AnyEmpty(values ...interface{}) bool {
	for _, v := range values {
		switch t := v.(type) {
		case string:
			if t == "" {
				return true
			}
		case []string:
			if len(t) == 0 {
				return true
			}
		case interface{}:
			if IsNil(t) {
				return true
			}
		}
	}

	return false
}
