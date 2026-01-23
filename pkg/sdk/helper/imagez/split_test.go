package imagez

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

func TestSplit3x3(t *testing.T) {
	// Create a 300x300 image
	width := 300
	height := 300
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with som content to be sure (optional, just testing split logic)
	// Iterate to verify distinctness? Not necessary for basic dimension check.
	// Let's color the center pixel red to check if the middle one catches it.
	img.Set(150, 150, color.RGBA{255, 0, 0, 255})

	// Encode to PNG to get bytes
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("failed to encode test image: %v", err)
	}

	partsData, err := Split3x3(buf.Bytes())
	if err != nil {
		t.Fatalf("Split3x3 failed: %v", err)
	}

	if len(partsData) != 9 {
		t.Errorf("expected 9 parts, got %d", len(partsData))
	}

	expectedWidth := 100
	expectedHeight := 100

	var parts []image.Image
	for _, p := range partsData {
		pi, _, err := image.Decode(bytes.NewReader(p))
		if err != nil {
			t.Fatalf("failed to decode part: %v", err)
		}
		parts = append(parts, pi)
	}

	for i, part := range parts {
		b := part.Bounds()
		if b.Dx() != expectedWidth {
			t.Errorf("part %d width mismatch: got %d, want %d", i, b.Dx(), expectedWidth)
		}
		if b.Dy() != expectedHeight {
			t.Errorf("part %d height mismatch: got %d, want %d", i, b.Dy(), expectedHeight)
		}
	}

	// Check center image (index 4) contains the red pixel
	centerImg := parts[4]
	// Since we re-encoded/decoded, bounds are reset to (0,0) usually
	// The Center of 100x100 is (50,50)

	midX := centerImg.Bounds().Min.X + 50
	midY := centerImg.Bounds().Min.Y + 50

	r, _, _, _ := centerImg.At(midX, midY).RGBA()
	if r < 0x8000 { // Check if red component is significant
		t.Errorf("center image did not contain the center pixel properly, got r=%d", r)
	}
}

func TestSplit3x3_SmallImage(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	var buf bytes.Buffer
	png.Encode(&buf, img)

	_, err := Split3x3(buf.Bytes())
	if err == nil {
		t.Error("expected error for small image, got nil")
	}
}

func TestSplit3x3_EmptyData(t *testing.T) {
	_, err := Split3x3(nil)
	if err == nil {
		t.Error("expected error for empty data, got nil")
	}
}

func TestSplit3x3_RealImage(t *testing.T) {
	// Read the real image file
	data, err := os.ReadFile("021766687487854c7ba22cbcbef0b3ad5a364f1d5ffd4d538dc89_0.jpeg")
	if err != nil {
		t.Skipf("skipping real image test, file not found: %v", err)
		return
	}

	partsData, err := Split3x3(data)
	if err != nil {
		t.Fatalf("Split3x3 failed: %v", err)
	}

	if len(partsData) != 9 {
		t.Errorf("expected 9 parts, got %d", len(partsData))
	}

	// Verify we can decode the parts back to images
	for i, p := range partsData {
		img, _, err := image.Decode(bytes.NewReader(p))
		if err != nil {
			t.Errorf("failed to decode part %d: %v", i, err)
			continue
		}

		// Basic check that dimension is reasonable (w > 0, h > 0)
		if img.Bounds().Dx() == 0 || img.Bounds().Dy() == 0 {
			t.Errorf("part %d has zero dimension", i)
		}

		// Save to file for visual inspection
		filename := fmt.Sprintf("output_%d.jpeg", i)
		if err := os.WriteFile(filename, p, 0644); err != nil {
			t.Errorf("failed to write file %s: %v", filename, err)
		}
	}
}

func TestSplit3x3_OddDimensions(t *testing.T) {
	// 301 x 302 image
	width := 301
	height := 302
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// Set a pixel at the very edge
	img.Set(300, 301, color.RGBA{255, 0, 0, 255})

	// Encode to PNG to get bytes
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("failed to encode test image: %v", err)
	}

	partsData, err := Split3x3(buf.Bytes())
	if err != nil {
		t.Fatalf("Split3x3 failed: %v", err)
	}

	totalArea := 0
	foundRedPixel := false

	for i, p := range partsData {
		pi, _, err := image.Decode(bytes.NewReader(p))
		if err != nil {
			t.Fatalf("failed to decode part %d: %v", i, err)
		}

		b := pi.Bounds()
		totalArea += b.Dx() * b.Dy()

		// The red pixel should be in the last part (index 8)
		if i == 8 {
			// In the last part, the pixel should be at bottom-right
			// Part 8 bounds relative to itself: (0,0) - (w,h)
			// The pixel was at (300, 301) in original.
			// Cell sizes:
			// W: 301/3 = 100. Cols: 0-100, 100-200, 200-301 (width 101)
			// H: 302/3 = 100. Rows: 0-100, 100-200, 200-302 (height 102)
			// So last part should be 101x102.

			if b.Dx() != 101 {
				t.Errorf("expected last part width 101, got %d", b.Dx())
			}
			if b.Dy() != 102 {
				t.Errorf("expected last part height 102, got %d", b.Dy())
			}

			// Local coordinate of the pixel:
			// Origin of part 8 is (200, 200)
			// Pixel (300, 301) -> Local (100, 101)
			pr, _, _, _ := pi.At(100, 101).RGBA()
			if pr > 0 {
				foundRedPixel = true
			}
		}
	}

	if totalArea != width*height {
		t.Errorf("total area mismatch: expected %d, got %d", width*height, totalArea)
	}
	if !foundRedPixel {
		t.Error("failed to find the edge pixel in the split results")
	}
}
