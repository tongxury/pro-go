package videoz

import (
	"fmt"
	"os"
	"testing"
)

// 使用示例
func TestGetFrames(t *testing.T) {

	file, err := os.ReadFile("34f25d2ce7a68dd2ec22cddafbe9aef0.mp4")
	if err != nil {
		return
	}

	frames, err := GetFrame(file, 1.0)
	if err != nil {
		return
	}

	fmt.Println(frames)
}
