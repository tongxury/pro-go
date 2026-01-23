package helper

import (
	"sort"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
)

func SliceElement[T any](src []T, index int, reverse bool) T {
	var zero T

	// Empty slice check
	if len(src) == 0 {
		return zero
	}

	// Handle reverse (count from end)
	if reverse {
		index = len(src) - 1 - index
	}

	// Bounds check
	if index < 0 || index >= len(src) {
		return zero
	}

	return src[index]
}

func SubSlice[T any](src []T, size int) []T {

	if size == 0 {
		return nil
	}

	if len(src) <= size {
		return src
	}

	return src[:size]
}

func SplitSlice[T any](src []T, partSize int) [][]T {
	l := len(src)
	if l <= partSize {
		return [][]T{src}
	}

	var rsp [][]T
	for i := 0; i < l; i += partSize {

		var current []T
		if i+partSize > l {
			current = src[i:]
		} else {
			current = src[i : i+partSize]
		}

		rsp = append(rsp, current)
	}

	return rsp
}

func ToSet[T comparable](list []T) []T {

	s := mapset.NewSet[T]()
	for _, x := range list {
		s.Add(x)
	}
	return s.ToSlice()
}

func FilterEmpty(list []string) []string {

	var rsp []string

	for _, s := range list {
		if s == "" {
			continue
		}

		rsp = append(rsp, s)
	}

	return rsp
}

func ReverseSlice(slice []string) []string {

	l := len(slice)

	if l == 0 {
		return slice
	}

	var rsp []string
	for i := l - 1; i >= 0; i-- {
		rsp = append(rsp, slice[i])
	}

	return rsp
}

func InSliceIgnoreCase(str string, list []string) bool {

	str = strings.ToLower(str)

	for _, s := range list {
		if str == strings.ToLower(s) {
			return true
		}
	}
	return false
}

func ContainsAll(src []string, targets ...string) bool {

	if len(targets) == 0 {
		return false
	}

	for _, x := range targets {
		if !InSlice(x, src) {
			return false
		}
	}

	return true
}

func ContainsAny(src []string, targets ...string) bool {

	for _, x := range targets {
		if InSlice(x, src) {
			return true
		}
	}
	return false
}

func SliceContainsAny(src []string, targets ...string) bool {

	mp := make(map[string]struct{}, len(targets))
	for _, x := range targets {
		mp[x] = struct{}{}
	}

	for _, x := range src {
		if _, ok := mp[x]; ok {
			return true
		}
	}
	return false
}

func InSlice(str string, list []string) bool {
	for _, s := range list {
		if str == s {
			return true
		}
	}

	return false
}

func NotInSlice(str string, list []string) bool {
	return !InSlice(str, list)
}

func InIntSlice(str int, list []int) bool {
	for _, s := range list {
		if str == s {
			return true
		}
	}

	return false
}

func UnionSlices(arrays ...[]interface{}) []interface{} {

	var rsp []interface{}

	for _, arr := range arrays {

		if arr != nil && len(arr) > 0 {
			for _, v := range arr {
				rsp = append(rsp, v)
			}
		}

	}
	return rsp
}

func UniqueSlice(list []string) []string {
	m := make(map[string]bool)
	var uniqueList []string
	for _, e := range list {
		if !m[e] {
			uniqueList = append(uniqueList, e)
			m[e] = true
		}
	}
	return uniqueList
}

func SliceIntersect(listA, listB []string) []string {
	res := []string{}
	if len(listA) == 0 || len(listB) == 0 {
		return res
	}
	m := make(map[string]bool)
	for _, e := range listA {
		m[e] = true
	}
	for _, e := range listB {
		if m[e] {
			res = append(res, e)
		}
	}
	return res
}

func SliceOverlap(listA, listB []string) bool {
	if len(listA) == 0 || len(listB) == 0 {
		return false
	}
	m := make(map[string]bool)
	for _, e := range listA {
		m[e] = true
	}
	for _, e := range listB {
		if m[e] {
			return true
		}
	}
	return false
}

func SliceElementsEqual(listA, listB []string) bool {
	if len(listA) != len(listB) {
		return false
	}
	if len(listA) == 0 && len(listB) == 0 {
		return true
	}
	m := make(map[string]bool)
	for _, e := range listA {
		m[e] = true
	}
	for _, e := range listB {
		if !m[e] {
			return false
		}
	}
	return true
}

func SliceEqualRegardlessOfOrder(listA, listB []string) bool {
	if len(listA) != len(listB) {
		return false
	}
	a, b := make([]string, len(listA)), make([]string, len(listB))
	copy(a, listA)
	copy(b, listB)
	sort.Strings(a)
	sort.Strings(b)
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func SliceSubstract(listA, listB []string) []string {
	res := []string{}
	m := make(map[string]bool)
	for _, e := range listB {
		m[e] = true
	}
	for _, e := range listA {
		if !m[e] {
			res = append(res, e)
		}
	}
	return res
}

func ConvertMapToStrings(mapLists map[string][]string) []string {
	list := []string{}
	for _, val := range mapLists {
		list = append(list, val...)
	}
	return list
}

// ConvertSliceToSlice2D is used to convert slice to two-dimensional slice.
func SplitSliceToSlice2D[T any](list []T, num int) [][]T {
	count := len(list) / num
	remainder := len(list) % num

	// calculate silce capacity
	capacity := count
	if remainder > 0 {
		capacity++
	}

	var result = make([][]T, capacity)
	var i = 0
	for ; i < count; i++ {
		result[i] = list[i*num : i*num+num]
	}

	if remainder > 0 {
		result[i] = list[i*num:]
	}

	return result
}

func SplitByMultiSepAndTrimSpace(ori string, ss ...string) []string {
	if len(ss) == 0 {
		return []string{strings.TrimSpace(ori)}
	}
	f := ss[0]
	for _, s := range ss[1:] {
		ori = strings.Replace(ori, s, f, -1)
	}
	var res []string
	for _, s := range strings.Split(ori, f) {
		if m := strings.TrimSpace(s); m != "" {
			res = append(res, m)
		}
	}
	return res
}

func Append(arr1 []string, arr2 []string) []string {

	var arr = arr1

	if len(arr2) > 0 {

		for _, v := range arr2 {
			arr = append(arr, v)
		}
	}

	return arr
}

func AsStringSliceIgnoreEmpty(args ...string) []string {
	var rsp []string

	for _, a := range args {
		if a != "" {
			rsp = append(rsp, a)
		}
	}

	return rsp
}

func AsArray(args ...interface{}) []interface{} {
	var rsp []interface{}

	for _, a := range args {
		rsp = append(rsp, a)
	}

	return rsp
}

// 有交集
func Intersect(src []string, target []string) ([]string, []string) {
	if len(src) == 0 || len(target) == 0 {
		return nil, nil
	}

	targetMap := make(map[string]bool, len(src))
	for _, v := range target {
		targetMap[v] = true
	}

	var hits []string
	var unHits []string
	for _, vv := range src {
		if targetMap[vv] {
			hits = append(hits, vv)
		} else {
			unHits = append(unHits, vv)
		}
	}

	return hits, unHits
}

func Contains(src []string, target string) bool {

	if len(src) == 0 {
		return false
	}

	for _, v := range src {
		if v == target {
			return true
		}
	}

	return false
}
