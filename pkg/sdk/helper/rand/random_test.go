package rand

import (
	"fmt"
	"testing"
)

func TestWeighted(t *testing.T) {
	items := []Item[string]{
		{Weight: 1, Value: "1"},
		{Weight: 3, Value: "3"},
		{Weight: 6, Value: "6"},
	}

	selections := make(map[string]int)
	iterations := 500000

	// 进行多次随机选择并记录结果
	for i := 0; i < iterations; i++ {
		selected := Weighted(items...)
		selections[selected]++
	}

	// 输出每个项被选择的次数
	for item, count := range selections {
		fmt.Printf("%s: %d\n", item, count)
	}
}
