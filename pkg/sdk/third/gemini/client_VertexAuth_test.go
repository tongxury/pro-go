package gemini

import (
	"context"
	"fmt"
	"store/confs"
	"testing"
	"time"

	"google.golang.org/genai"
)

func TestClient_VertexAuth(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// 使用 secrets.go 中定义的 VertexAiSecret
	credentialsJSON := confs.VertexAiSecret
	if credentialsJSON == "" {
		t.Skip("Skipping test: VertexAiSecret is empty")
	}

	factory := NewGenaiFactory(&FactoryConfig{
		Configs: []*Config{
			{
				Project:         "yuzhi-483807", // 从 secrets.go 或 configs.go 中提取的 Project ID
				Location:        "us-central1",  // 常用 Vertex AI Location
				APIVersion:      "v1",
				CredentialsJSON: credentialsJSON,                                  // 使用 JSON 字符串凭据
				Proxy:           "http://proxy:strOngPAssWOrd@45.78.194.147:6060", // 如有需要
			},
		},
	})

	client := factory.Get()
	if client == nil {
		t.Fatal("Failed to create client from factory")
	}

	// 简单的测试 Prompt
	prompt := "Hello, are you working?"

	resp, err := client.GenerateContent(ctx, GenerateContentRequest{
		Parts: []*genai.Part{
			NewTextPart(prompt),
		},
	})
	if err != nil {
		t.Fatalf("Vertex AI Auth Failed: %v", err)
	}

	fmt.Printf("Vertex AI Response: %s\n", resp)
	t.Logf("Successfully authenticated and generated content via Vertex AI!")
}
