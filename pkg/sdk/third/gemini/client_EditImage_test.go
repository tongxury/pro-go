package gemini

import (
	"context"
	"fmt"
	"os"
	"store/confs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_EditImage_WithLocalFiles(t *testing.T) {
	ctx := context.Background()

	// Initialize factory
	factory := NewGenaiFactory(&FactoryConfig{
		Configs: []*Config{
			{
				Proxy: "http://proxy:strOngPAssWOrd@45.78.194.147:6060",
				Key:   confs.AQKey,
			},
		},
	})
	client := factory.Get()

	// Load local images
	img1Path := "微信图片_20260117200030_66_59.jpg"
	img2Path := "微信图片_20260130174959_1320_402.jpg"

	img1Bytes, err := os.ReadFile(img1Path)
	if err != nil {
		t.Fatalf("Failed to read local image %s: %v", img1Path, err)
	}

	img2Bytes, err := os.ReadFile(img2Path)
	if err != nil {
		t.Fatalf("Failed to read local image %s: %v", img2Path, err)
	}

	// Prepare request
	req := EditImageRequest{
		Model:  ModelImagen40FastGenerate001,
		Prompt: "A photo of the [1] with the artistic and cinematic style of [2].",
		ReferenceImages: []ReferenceImageConfig{
			{
				Bytes:    img1Bytes,
				MimeType: "image/jpeg",
				Type:     ReferenceImageTypeContent,
			},
			{
				Bytes:    img2Bytes,
				MimeType: "image/jpeg",
				Type:     ReferenceImageTypeStyle,
			},
		},
		AspectRatio: "9:16",
	}

	// Call EditImageToBytes
	fmt.Println("Calling EditImageToBytes...")
	resp, err := client.EditImageToBytes(ctx, req)
	if err != nil {
		t.Fatalf("EditImageToBytes failed: %v", err)
	}

	assert.NotEmpty(t, resp)

	// Save result
	outputPath := "test_edit_output.png"
	err = os.WriteFile(outputPath, resp, 0644)
	if err != nil {
		t.Fatalf("Failed to write output image: %v", err)
	}

	fmt.Printf("Image edited successfully and saved to %s\n", outputPath)
}
