package gemini_test

import (
	"context"
	"fmt"
	"os"
	"store/confs"
	"store/pkg/sdk/third/gemini"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/genai"
)

func TestClient_GenerateImage_Retouch(t *testing.T) {
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
	// usage of an existing image in the directory
	imagePath := "截屏2026-02-12 12.08.13.png"
	imageBytes, err := os.ReadFile(imagePath)
	if err != nil {
		t.Fatalf("Failed to read local image %s: %v", imagePath, err)
	}

	// Use GenerateContent (Multimodal) instead of EditImage
	// EditImage without a mask often preserves the original image significantly (inpainting behavior).
	// To perform a semantic transformation (retouching/style transfer), we use GenerateContent with image input and image output.

	prompt := "帮我改成app store 介绍图"
	model := gemini.ModelGenmini3FlashPreview // gemini-2.0-flash-001 supports multimodal input and output

	fmt.Printf("Calling native GenerateContent with model: %s, prompt: %s\n", model, prompt)

	// Construct Parts: Image + Text Prompt
	parts := []*genai.Part{
		genai.NewPartFromBytes(imageBytes, "image/png"),
		genai.NewPartFromText(prompt),
	}

	config := &genai.GenerateContentConfig{
		ResponseModalities: []string{"TEXT", "IMAGE"},
		// You can add generation config parameters here if needed (e.g., candidate count)
	}

	// Native call
	resp, err := client.C().Models.GenerateContent(ctx, model, []*genai.Content{{Role: "user", Parts: parts}}, config)

	if err != nil {
		t.Fatalf("Native GenerateContent Failed: %v", err)
	}

	// Extract Image from Response
	var outputData []byte
	for _, candidate := range resp.Candidates {
		if candidate.Content == nil {
			continue
		}
		for _, part := range candidate.Content.Parts {
			if part.InlineData != nil && strings.HasPrefix(part.InlineData.MIMEType, "image/") {
				outputData = part.InlineData.Data
				break
			}
		}
		if outputData != nil {
			break
		}
	}

	if outputData == nil {
		t.Fatal("No image generated in response")
	}

	// Save output
	outputPath := "retouch_output.png"
	err = os.WriteFile(outputPath, outputData, 0644)
	assert.NoError(t, err)

	fmt.Printf("Image retouching (generation) successful! Saved to %s\n", outputPath)
}
