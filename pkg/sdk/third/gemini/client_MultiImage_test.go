package gemini_test

import (
	"context"
	"fmt"
	"os"
	"store/confs"
	"store/pkg/sdk/third/gemini"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient_GenerateImage_MultiImage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	credentialsJSON := confs.VertexAiSecret
	if credentialsJSON == "" {
		t.Skip("Skipping test: VertexAiSecret is empty")
	}

	factory := gemini.NewGenaiFactory(&gemini.FactoryConfig{
		Configs: []*gemini.Config{
			{
				Project:         "yuzhi-483807",
				Location:        "us-central1",
				APIVersion:      "v1",
				CredentialsJSON: credentialsJSON,
				Proxy:           "http://proxy:strOngPAssWOrd@45.78.194.147:6060",
			},
		},
	})

	client := factory.Get()
	if client == nil {
		t.Fatal("Failed to create client from factory")
	}

	// Read local image
	imagePath := "微信图片_20260130174959_1320_402.jpg"
	imageBytes, err := os.ReadFile(imagePath)
	if err != nil {
		t.Fatalf("Failed to read local image %s: %v", imagePath, err)
	}

	// Test multi-image: using the same image twice as a sample
	resp, err := client.GenerateImage(ctx, gemini.GenerateImageRequest{

		ImageBytes:  [][]byte{imageBytes, imageBytes},
		Prompt:      "Combine elements from these two images into a new high quality commercial background for fashion video.",
		AspectRatio: "9:16",
	})

	if err != nil {
		t.Fatalf("GenerateImage Failed: %v", err)
	}

	assert.NotEmpty(t, resp)

	// Save output
	outputPath := "multi_image_output.png"
	err = os.WriteFile(outputPath, resp, 0644)
	assert.NoError(t, err)

	fmt.Printf("Multi-image generation successful! Saved to %s\n", outputPath)

}

func TestClient_GenerateImage_MultiImageV2(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	credentialsJSON := confs.VertexAiSecret
	if credentialsJSON == "" {
		t.Skip("Skipping test: VertexAiSecret is empty")
	}

	factory := gemini.NewGenaiFactory(&gemini.FactoryConfig{
		Configs: []*gemini.Config{
			{
				Project:         "yuzhi-483807",
				Location:        "us-central1",
				APIVersion:      "v1",
				CredentialsJSON: credentialsJSON,
				Proxy:           "http://proxy:strOngPAssWOrd@45.78.194.147:6060",
			},
		},
	})

	client := factory.Get()
	if client == nil {
		t.Fatal("Failed to create client from factory")
	}

	// Read local image
	imagePath := "微信图片_20260130174959_1320_402.jpg"
	imageBytes, err := os.ReadFile(imagePath)
	if err != nil {
		t.Fatalf("Failed to read local image %s: %v", imagePath, err)
	}

	// Test multi-image: using the same image twice as a sample
	resp, err := client.GenerateImageV2(ctx, gemini.GenerateImageRequest{

		ImageBytes:  [][]byte{imageBytes, imageBytes},
		Prompt:      "Combine elements from these two images into a new high quality commercial background for fashion video.",
		AspectRatio: "9:16",
	})

	if err != nil {
		t.Fatalf("GenerateImage Failed: %v", err)
	}

	assert.NotEmpty(t, resp)

	// Save output
	outputPath := "multi_image_output.png"
	err = os.WriteFile(outputPath, resp, 0644)
	assert.NoError(t, err)

	fmt.Printf("Multi-image generation successful! Saved to %s\n", outputPath)

}
