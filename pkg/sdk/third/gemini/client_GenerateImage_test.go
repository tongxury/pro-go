package gemini

import (
	"context"
	"fmt"
	"os"
	"store/confs"
	"testing"

	"github.com/stretchr/testify/assert"
)

//// 1. 读取本地图片：已修复，使用 os.ReadFile 加载字节流
//imagePath := "WechatIMG30.jpg"
//imageBytes, err := os.ReadFile(imagePath)
//if err != nil {
//t.Fatalf("Failed to read local image %s: %v", imagePath, err)
//}
//
//imagePart := genai.NewPartFromBytes(imageBytes, "image/jpeg")

func TestClient_GenerateImage_WithLocalFile(t *testing.T) {
	ctx := context.Background()

	// 初始化客户端（请确保你的配置文件或环境变量中有正确的 Key 和 Proxy）
	// 这里参考你之前的测试配置
	c := NewGenaiFactory(&FactoryConfig{
		Configs: []*Config{
			{Proxy: "http://proxy:strOngPAssWOrd@45.78.194.147:6060", Key: confs.AQKey},
		},
	})

	// 调用生成图片接口
	resp, err := c.Get().GenerateImage(ctx, GenerateImageRequest{
		//ImageBytes:  [][]byte{imageData, imageData1},
		Prompt:      "请帮我生成一张能够用来做女装带货视频的通用视频背景图，要高级，有冲击力",
		ImageSize:   "1K",
		AspectRatio: "9:16",
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
