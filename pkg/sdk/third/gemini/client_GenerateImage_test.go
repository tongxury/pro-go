package gemini

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GenerateImage_WithLocalFile(t *testing.T) {
	ctx := context.Background()

	// 初始化客户端（请确保你的配置文件或环境变量中有正确的 Key 和 Proxy）
	// 这里参考你之前的测试配置
	c := NewGenaiFactory(&FactoryConfig{
		Configs: []*Config{
			{Proxy: "http://proxy:strOngPAssWOrd@45.78.194.147:6060", Key: "AQ.Ab8RN6Kj4bK5fzaq8rNj1xHc1vtcWSKgQ4h-F7SwItVNbk8RzQ"},
		},
	})

	// 读取本地图片文件
	//imagePath := "截屏2026-01-18 21.46.17.png"
	imagePath := "截屏2026-01-18 21.47.23.png"
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		t.Fatalf("failed to read test image: %v", err)
	}

	imagePath1 := "generated_outputa1.png"
	imageData1, err := os.ReadFile(imagePath1)
	if err != nil {
		t.Fatalf("failed to read test image: %v", err)
	}

	// 调用生成图片接口
	resp, err := c.Get().GenerateImage(ctx, GenerateImageRequest{
		ImageBytes:  [][]byte{imageData, imageData1},
		Prompt:      "根据【图1】，帮我构想 在ipad 上的截图，输出和【图2】风格类似的图",
		ImageSize:   "1K",
		AspectRatio: "4:3",
	})

	if err != nil {
		fmt.Printf("Error generating image: %v\n", err)
		return
	}

	assert.NoError(t, err)
	assert.NotEmpty(t, resp)

	// 将生成的图片保存到本地以便查看效果
	outputPath := "generated_output.png"
	err = os.WriteFile(outputPath, resp, 0644)
	assert.NoError(t, err)

	fmt.Printf("Image generated successfully and saved to %s\n", outputPath)
}
