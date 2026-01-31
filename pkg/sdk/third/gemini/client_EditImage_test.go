package gemini

import (
	"context"
	"fmt"
	"os"
	"store/confs"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/genai"
)

// TestClient_StyleTransfer_WithGemini3 uses gemini-3-pro-image-preview for style transfer
// This works with Gemini API via proxy (no Vertex AI required)
func TestClient_StyleTransfer_WithGemini3(t *testing.T) {
	ctx := context.Background()

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

	// Use GenerateImage with multiple images for style transfer
	// This uses gemini-3-pro-image-preview which supports native image editing
	fmt.Println("Calling GenerateImage with style transfer prompt...")
	resp, err := client.GenerateImage(ctx, GenerateImageRequest{
		ImageBytes:  [][]byte{img1Bytes, img2Bytes},
		Prompt:      "Create a new image that combines the content of the first image with the artistic style of the second image. Apply the color palette, textures, and visual aesthetic from the second image to the subject in the first image.",
		ImageSize:   "1K",
		AspectRatio: "9:16",
	})
	if err != nil {
		t.Fatalf("GenerateImage failed: %v", err)
	}

	assert.NotEmpty(t, resp)

	// Save result
	outputPath := "test_style_transfer_output.png"
	err = os.WriteFile(outputPath, resp, 0644)
	if err != nil {
		t.Fatalf("Failed to write output image: %v", err)
	}

	fmt.Printf("Style transfer completed and saved to %s\n", outputPath)
}

// TestClient_EditImage_Imagen uses Imagen's EditImage API
// NOTE: This requires Vertex AI backend with GCP Project/Location and OAuth2 credentials
// It does NOT work with API Key authentication
func TestClient_EditImage_Imagen(t *testing.T) {
	t.Skip("Skipping: Imagen EditImage requires Vertex AI backend with GCP OAuth2 credentials, not API Key")

	ctx := context.Background()

	// Initialize factory with Vertex AI config
	// Requires:
	// 1. Set GOOGLE_CLOUD_PROJECT and GOOGLE_CLOUD_LOCATION environment variables
	// 2. Set up Application Default Credentials: gcloud auth application-default login

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

	// Prepare request - Use imagen-3.0-capability-001 for editing with style mode
	req := EditImageRequest{
		Model:  ModelImagen3Capability001,
		Prompt: "A photo of the [1] with the artistic and cinematic style of [2].",
		ReferenceImages: []ReferenceImageConfig{
			{
				Bytes:    img1Bytes,
				MimeType: "image/jpeg",
				Type:     ReferenceImageTypeRaw,
			},
			{
				Bytes:    img2Bytes,
				MimeType: "image/jpeg",
				Type:     ReferenceImageTypeStyle,
			},
		},
		AspectRatio: "9:16",
		EditMode:    genai.EditModeStyle,
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
