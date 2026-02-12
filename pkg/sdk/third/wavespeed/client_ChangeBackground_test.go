package wavespeed

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"store/pkg/sdk/helper/filez"
	"testing"
	"time"
)

func TestClient_Gemini3ProImage_ChangeBackground(t *testing.T) {
	c := NewClient()
	ctx := context.Background()

	// 1. Read image file
	imgData, err := os.ReadFile("tpl.png")
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", "", err)
	}

	// 2. Convert to Base64 Data URI
	// "data:image/png;base64,..."
	tplDataURI := filez.BytesToURI(imgData)

	// List of images to process
	imageFiles := []string{
		//"截屏2026-02-12 12.08.13.png",
		//"截屏2026-02-12 12.11.23.png",
		"截屏2026-02-12 12.12.00.png",
	}

	for _, filename := range imageFiles {
		t.Run(filename, func(t *testing.T) {

			// 1. Read image file
			imgData, err := os.ReadFile(filename)
			if err != nil {
				t.Fatalf("Failed to read file %s: %v", filename, err)
			}

			// 2. Convert to Base64 Data URI
			// "data:image/png;base64,..."
			dataURI := filez.BytesToURI(imgData)

			fmt.Printf("[%s] Submitting for background replacement...\n", filename)

			// 3. Call API
			// Using "Change background to..." prompt
			req := Gemini3ProImageRequest{
				Prompt:             "参考图1的背景和布局，帮我将主体图改成图2,文案也改一下",
				Images:             []string{tplDataURI, dataURI},
				AspectRatio:        "3:4", // Matching source aspect ratio roughly
				OutputFormat:       "png",
				EnableSyncMode:     false,
				EnableBase64Output: false, // We want URL
			}

			resp, err := c.Gemini3ProImage(ctx, req)
			if err != nil {
				t.Fatalf("[%s] API Request failed: %v", filename, err)
			}

			if resp.Code != 200 {
				t.Fatalf("[%s] API Error: %s", filename, resp.Message)
			}

			taskID := resp.Data.Id
			if taskID == "" {
				t.Fatalf("[%s] No Task ID returned", filename)
			}

			fmt.Printf("[%s] Task ID: %s. Polling...\n", filename, taskID)

			// 4. Poll for results
			// Poll every 2 seconds, up to 30 times (60s timeout)
			timeout := time.After(60 * time.Second)
			ticker := time.NewTicker(2 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-timeout:
					t.Fatalf("[%s] Timed out polling for result", filename)
				case <-ticker.C:
					res, err := c.GetResult(ctx, taskID)
					if err != nil {
						t.Logf("[%s] Poll error: %v", filename, err)
						continue
					}

					status := res.Data.Status
					fmt.Printf("[%s] Status: %s\n", filename, status)

					if status == "completed" {
						if len(res.Data.Outputs) > 0 {
							url := res.Data.Outputs[0]
							fmt.Printf("[%s] SUCCESS! Output URL: %s\n", filename, url)

							// Download and save the image
							outputFilename := "output_" + filename
							err := downloadFile(url, outputFilename)
							if err != nil {
								t.Fatalf("[%s] Failed to download image: %v", filename, err)
							}
							fmt.Printf("[%s] Saved to %s\n", filename, outputFilename)

						} else {
							t.Errorf("[%s] Succeeded but no outputs!", filename)
						}
						return // Done with this image
					} else if status == "failed" {
						t.Fatalf("[%s] Task failed: %s", filename, res.Data.Error)
					}
				}
			}
		})
	}
}

func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
