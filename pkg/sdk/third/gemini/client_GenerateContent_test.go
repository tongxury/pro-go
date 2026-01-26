package gemini

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/helper/imagez"
	"store/pkg/sdk/third/douyin"
	"store/pkg/sdk/third/openaiz"
	"testing"
	"time"

	"store/confs"

	"github.com/openai/openai-go/v3"
	"google.golang.org/genai"
)

func TestClient_GenerateContent2(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	// 1. 读取本地图片：已修复，使用 os.ReadFile 加载字节流
	imagePath := "微信图片_20260126185147_30_3.jpg"
	imageBytes, err := os.ReadFile(imagePath)
	if err != nil {
		t.Fatalf("Failed to read local image %s: %v", imagePath, err)
	}

	imagePart := genai.NewPartFromBytes(imageBytes, "image/jpeg")

	parts := []*genai.Part{
		imagePart,
		NewTextPart(
			`
请帮我判断这张照片是不是被美化过的，给出美化前的最可能的样子
	`,
		),
	}

	c := NewGenaiFactory(&FactoryConfig{
		Configs: []*Config{
			{Proxy: "http://proxy:strOngPAssWOrd@45.78.194.147:6060", Key: confs.AQKey},
		},
	})
	content, err := c.Get().GenerateContent(ctx, GenerateContentRequest{
		Parts: parts,
	})
	if err != nil {
		t.Fatalf("failed to generate content: %v", err)
	}

	// 3. 生成 Gemini 参考图
	blob3, err := c.Get().GenerateBlob(ctx, GenerateContentRequest{
		Model: ModelGenmini3ProImagePreview,
		Parts: parts,
		Config: &genai.GenerateContentConfig{
			ImageConfig: &genai.ImageConfig{
				AspectRatio: "9:16",
				ImageSize:   "2K",
			},
		},
	})
	if err != nil {
		t.Fatalf("Error generating image from Gemini: %v", err)
	}

	fmt.Println(content)

	// 将生成的 Gemini 参考图保存到本地供查看
	refImagePath := "gemini_reference.jpg"
	if err := os.WriteFile(refImagePath, blob3.Data, 0644); err != nil {
		t.Logf("Warning: failed to save Gemini reference image: %v", err)
	} else {
		t.Logf("Gemini reference image saved to %s", refImagePath)
	}

}

func TestClient_GenerateContent(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	dy := douyin.NewClient()

	c := NewGenaiFactory(&FactoryConfig{
		Configs: []*Config{
			{Proxy: "http://proxy:strOngPAssWOrd@45.78.194.147:6060", Key: confs.AQKey},
		},
	})

	metadata, err := dy.GetCommodityMetadata(ctx, "https://v.douyin.com/fLWcoNGHL5w/")
	if err != nil {
		t.Logf("Warning: failed to get commodity metadata: %v", err)
	}

	parts := []*genai.Part{}

	// 2. 将抖音图片加入 Prompt
	if metadata != nil {
		for i, x := range metadata.Images {
			if i > 5 {
				break
			}
			y, err := NewImagePart(x)
			if err != nil {
				t.Logf("Skipping remote image %d: %v", i, err)
				continue
			}
			parts = append(parts, y)
		}
	}

	// 生成sora2提示词
	parts = append(parts, NewTextPart(
		`
		结合商品的最佳卖点,生成一段能够提供给sora2生成视频的脚本，视频长度12s。要有节奏感，居家写实风，人物动作语气神态要真实自然，视觉冲击力强，针对抖音电商受众，中国人主角。
	`,
	))

	//为了确保sora2生成时的人物IP高度固定，请在脚本中强制嵌入以下【四大一致性支柱】描述：
	//
	//1. 【物理锚点 (Identity Anchors)】：
	//- 发型：vibrant silver wolf-cut hairstyle (银色狼尾发)
	//- 纹身：a tiny delicate blue butterfly tattoo on the side of the neck (颈部蓝色蝴蝶纹身)
	//- 配饰：a pair of oversized liquid gold irregular earrings (夸张的流金不规则耳环)
	//
	//2. 【动作锚点 (Behavioral IP)】：
	//- 固定习惯动作：人物在每一个场景中必须有一次“伸手轻轻调整左耳环并自信微笑 (adjusting her left earring with a slight confident smile)”的动作。
	//
	//3. 【环境锚点 (Environmental Anchor - 关键固定项)】：
	//- 固定拍摄背景：在每一个场景中，背景必须严格固定为“高奢极简风格的法式阳台，阳光明媚，有柔和的金色阴影，白色大理石地面，背后是模糊的巴黎铁塔远景 (a luxury minimalist French balcony, sun-drenched with golden hour light, white marble floors, with a blurred Eiffel Tower in the distant background)”。
	//
	//要求：以上四点必须作为脚本生成的“不变量”硬性嵌入，确保Sora2生成结果具备极高的账号辨识度和品牌高级感。

	content, err := c.Get().GenerateContent(ctx, GenerateContentRequest{
		Parts: parts,
	})
	if err != nil {
		t.Fatalf("failed to generate content: %v", err)
	}

	// 为了防止 content 中包含 % 导致 Sprintf 解析错误，我们使用字符串拼接或直接构造 Parts
	parts = append(parts, NewTextPart(
		fmt.Sprintf(`
			结合我的提示词生成一张能够提供给 sora2 生成视频使用的参考图, 加上居家风的背景，不要有人物漏出。
		`,
		)))

	// 3. 生成 Gemini 参考图
	blob3, err := c.Get().GenerateBlob(ctx, GenerateContentRequest{
		Model: ModelGenmini3ProImagePreview,
		Parts: parts,
		Config: &genai.GenerateContentConfig{
			ImageConfig: &genai.ImageConfig{
				AspectRatio: "9:16",
				ImageSize:   "2K",
			},
		},
	})
	if err != nil {
		t.Fatalf("Error generating image from Gemini: %v", err)
	}

	// 将生成的 Gemini 参考图保存到本地供查看
	refImagePath := "gemini_reference.jpg"
	if err := os.WriteFile(refImagePath, blob3.Data, 0644); err != nil {
		t.Logf("Warning: failed to save Gemini reference image: %v", err)
	} else {
		t.Logf("Gemini reference image saved to %s", refImagePath)
	}

	// 4. 初始化 OpenAI/Sora 客户端
	openaiClient := openaiz.NewClient(openaiz.Config{
		AppKey:  confs.OpenAIKeys[0],
		BaseUrl: "http://45.78.194.147:6000",
	})

	// 缩放图片以符合 Sora 要求
	resizedImage, err := imagez.ResizeKeepRatio(blob3.Data, 720, 1280)
	if err != nil {
		t.Fatalf("Failed to resize image: %v", err)
	}

	inputReference := openai.File(bytes.NewReader(resizedImage), helper.CreateUUID()+".jpg", "image/jpeg")
	//5. 提交 Sora 视频任务
	video, err := openaiClient.Videos().New(ctx, openai.VideoNewParams{
		Model:          openai.VideoModelSora2,
		Prompt:         content,
		InputReference: inputReference,
		Seconds:        "12",
	})

	if err != nil {
		t.Fatalf("Error creating Sora video task: %v", err)
	}

	t.Logf("Video creation started, ID: %s", video.ID)

	// 6. 安全轮询视频状态
	maxRetries := 90 // 约 15 分钟 (10s 间隔)
	for i := 0; i < maxRetries; i++ {
		v, err := openaiClient.Videos().Get(ctx, video.ID)
		if err != nil {
			t.Logf("Attempt %d: Error fetching video status: %v", i+1, err)
			time.Sleep(10 * time.Second)
			continue
		}

		t.Logf("Status: %s", v.Status)
		if v.Status == "completed" {
			// 下载并写入本地保存
			fileName := fmt.Sprintf("sora_output_%s.mp4", v.ID)
			t.Logf("Video completed! Downloading to %s", fileName)

			resp, err := openaiClient.Videos().DownloadContent(ctx, v.ID, openai.VideoDownloadContentParams{
				Variant: "video",
			})
			if err != nil {
				t.Fatalf("Failed to download video content: %v", err)
			}
			defer resp.Body.Close()

			f, err := os.Create(fileName)
			if err != nil {
				t.Fatalf("Failed to create video file: %v", err)
			}
			defer f.Close()

			if _, err := io.Copy(f, resp.Body); err != nil {
				t.Fatalf("Failed to save video content: %v", err)
			}

			t.Logf("Successfully saved Sora video to %s", fileName)
			return
		}

		if v.Status == "failed" {
			t.Fatalf("Sora video generation failed: %v", conv.S2J(v))
		}

		time.Sleep(10 * time.Second)
	}

	t.Errorf("Timeout: Sora video generation did not complete within limit")
}

func TestClient_GetVideo(t *testing.T) {
	videoID := "video_69739c1a2eb08190ad9a9c9e9cafb7b70e9d4a5f9bd05409" // 替换为您想查询的视频 ID
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	openaiClient := openaiz.NewClient(openaiz.Config{
		AppKey:  confs.OpenAIKeys[0],
		BaseUrl: "http://45.78.194.147:6000",
	})

	t.Logf("Querying status for video: %s", videoID)
	v, err := openaiClient.Videos().Get(ctx, videoID)
	if err != nil {
		t.Fatalf("Error getting video status: %v", err)
	}

	t.Logf("Current Status: %s", v.Status)
	if v.Status == "failed" {
		t.Fatalf("Video generation failed: %v", conv.S2J(v))
	}

	if v.Status == "completed" {
		fileName := fmt.Sprintf("sora_download_%s.mp4", v.ID)
		t.Logf("Video is ready! Downloading to %s", fileName)

		resp, err := openaiClient.Videos().DownloadContent(ctx, v.ID, openai.VideoDownloadContentParams{
			Variant: "video",
		})
		if err != nil {
			t.Fatalf("Failed to download video content: %v", err)
		}
		defer resp.Body.Close()

		f, err := os.Create(fileName)
		if err != nil {
			t.Fatalf("Failed to create video file: %v", err)
		}
		defer f.Close()

		if _, err := io.Copy(f, resp.Body); err != nil {
			t.Fatalf("Failed to save video content: %v", err)
		}

		t.Logf("Successfully downloaded video to %s", fileName)
	} else {
		t.Logf("Video is still in progress (Status: %s). Please run this test again later.", v.Status)
	}
}

// video_69739c1a2eb08190ad9a9c9e9cafb7b70e9d4a5f9bd05409
