package mathz

import (
	"fmt"
	"math"
	mathRand "math/rand"
	"strconv"
	"time"
)

type Int interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | UInt
}

type UInt interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Float interface {
	~float64 | ~float32
}

type Number interface {
	Float | Int
}

func RandNumber(start, end int) int {
	r := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
	return r.Intn(end-start+1) + start
}

func Max[T Number](ts ...T) T {

	if len(ts) == 0 {
		return 0
	}

	var rsp T
	for i, t := range ts {
		if i == 0 {
			rsp = t
			continue
		}

		if rsp < t {
			rsp = t
		}
	}

	return rsp
}

func Pow10[V Number, E Int, R Number](value V, e E) R {
	return R(float64(value) * math.Pow10(int(e)))
}

func Min[T Number](ts ...T) T {

	if len(ts) == 0 {
		return 0
	}

	var rsp T
	for i, t := range ts {
		if i == 0 {
			rsp = t
			continue
		}

		if rsp > t {
			rsp = t
		}
	}

	return rsp
}

func Sum[T Number](ts ...T) T {
	var rsp T
	for _, t := range ts {
		rsp += t
	}

	return rsp
}

func Avg[T Float](ts ...T) T {

	if len(ts) == 0 {
		return 0
	}

	sum := Sum(ts...)

	return sum / T(len(ts))
}

func Factorial[T Int](num T) T {
	if num > 0 {
		return num * Factorial[T](num-1)
	} else {
		return 1
	}
}

func Cmn(total, target int) int {

	if target > total {
		return 0
	}

	return Factorial(total) / (Factorial(total-target) * Factorial(target))
}

func ToFixed4(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", value), 64)
	return value
}

func ToFixed3(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", value), 64)
	return value
}

func ToFixed2(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func ToFixed8(value float64) float64 {

	value, _ = strconv.ParseFloat(fmt.Sprintf("%.8f", value), 64)
	return value
}
