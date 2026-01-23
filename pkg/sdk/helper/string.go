package helper

import (
	"reflect"
	"store/pkg/sdk/conv"
	"strings"
	"unsafe"
)

func SubString(src string, start, length int) string {
	if start < 0 {
		start = 0
	}
	if start > len(src) {
		return ""
	}
	if length < 0 {
		return ""
	}
	end := start + length
	if end > len(src) {
		end = len(src)
	}
	return src[start:end]
}

func Underline(values ...any) string {

	var vals []string
	for _, x := range values {
		vals = append(vals, conv.String(x))
	}

	return strings.Join(vals, "_")
}

func inBytes(src byte, arr ...byte) bool {

	for _, b := range arr {
		if src == b {
			return true
		}
	}

	return false
}

func CleanBytes(src []byte, targets ...byte) []byte {

	var rsp []byte

	for _, b := range src {
		if inBytes(b, targets...) {
			continue
		}

		rsp = append(rsp, b)
	}

	return rsp
}

func StringToBytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func CleanString(src string) string {
	//strings.FieldsFunc(src, unicode.IsSpace)
	src = strings.TrimSpace(src)
	return src
}

func Replace(src string, replaceMap map[string]string) string {

	var s = src

	for k, v := range replaceMap {
		s = strings.ReplaceAll(s, k, v)
	}

	return s
}

func ContainsIgnoreCase(s string, sub string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(sub))
}

func ContainsAnyIgnoreCase(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(strings.ToLower(s), strings.ToLower(sub)) {
			return true
		}
	}

	return false
}

func EqualsIgnoreCase(a string, b string) bool {
	return Equals(strings.ToLower(a), strings.ToLower(b))
}

func EqualsAnyIgnoreCase(a string, targets ...string) bool {
	if len(targets) == 0 {
		return false
	}

	a = strings.ToLower(a)

	for _, x := range targets {
		x = strings.ToLower(x)
		if x == a {
			return true
		}
	}

	return false
}
