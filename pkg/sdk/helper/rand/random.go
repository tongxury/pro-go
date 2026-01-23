package rand

import (
	"math/rand"
	"time"
)

type Item[T any] struct {
	Weight int
	Value  T
}

func Range(start int, end int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(end-start+1) + start
}

func Weighted[T any](items ...Item[T]) T {

	// 计算总权重
	totalWeight := 0
	for _, item := range items {
		totalWeight += item.Weight
	}

	v := Range(0, totalWeight)

	// 根据随机权重选择
	for _, item := range items {
		if v < item.Weight {
			return item.Value
		}
		v -= item.Weight
	}

	return items[len(items)-1.0].Value
}
