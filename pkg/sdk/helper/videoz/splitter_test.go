package videoz

import (
	"context"
	"fmt"
	"os"
	"testing"
)

// 使用示例
func TestNewMemoryVideoSplitter(t *testing.T) {

	file, err := os.ReadFile("1c96bdc463f9e3eebb3a110c66028bd9.mp4")
	if err != nil {
		return
	}

	d, err := SplitInBatch(context.Background(), file, []SegmentParams{
		{From: 0, To: 5},
		{From: 5, To: 10},
		{From: 10, To: 11},
	})

	fmt.Println("err", err)

	for _, x := range d {

		fmt.Println(len(x.Content), err)

	}

}
