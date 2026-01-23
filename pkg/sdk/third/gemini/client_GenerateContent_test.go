package gemini

import (
	"context"
	"fmt"
	"testing"

	"store/confs"

	"google.golang.org/genai"
)

func TestClient_GenerateContent(t *testing.T) {

	ctx := context.Background()

	c := NewGenaiFactory(&FactoryConfig{
		Configs: []*Config{
			{Proxy: "http://proxy:strOngPAssWOrd@45.78.194.147:6060", Key: confs.AQKey},
		},
	})

	content, err := c.Get().GenerateContent(ctx, GenerateContentRequest{
		Model: "gemini-3-pro-preview",
		Config: &genai.GenerateContentConfig{
			ResponseMIMEType: "application/json",
			ResponseSchema: &genai.Schema{
				Type:     genai.TypeObject,
				Required: []string{"subtitle"},
				Properties: map[string]*genai.Schema{
					"subtitle": {
						Type:        genai.TypeString,
						Description: `口播文案`,
					},
				},
			},
		},
		Parts: []*genai.Part{
			{
				Text: "给我一段口播文案",
			},
			{
				Text: fmt.Sprintf(`
帮我生成新的口播文案，需要结合分镜头信息和新的商品信息`),
			},
		},
	})

	fmt.Println(content, err)
}
