package gemini

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/genai"
)

type GenerateContentRequest struct {
	Model  string
	Config *genai.GenerateContentConfig
	Parts  []*genai.Part
}

type GenerateContentRequestV2 struct {
	Model     string
	Config    *genai.GenerateContentConfig
	ImageUrls []string
	Prompt    string
}

func (t *Client) GenerateContentV2(ctx context.Context, req GenerateContentRequestV2) (string, error) {

	var parts []*genai.Part

	for i := range req.ImageUrls {
		//gsFile, err := t.UploadFile(ctx, req.ImageUrls[i], "image/jpeg")
		//if err != nil {
		//	return "", err
		//}

		p, err := NewImagePart(req.ImageUrls[i])
		if err != nil {
			return "", err
		}

		parts = append(parts, p)
	}

	parts = append(parts, &genai.Part{Text: req.Prompt})

	return t.GenerateContent(ctx, GenerateContentRequest{
		Model:  req.Model,
		Config: req.Config,
		Parts:  parts,
	})
}

func (t *Client) GenerateContent(ctx context.Context, req GenerateContentRequest) (string, error) {

	//parts := []*genai.Part{
	//	{Text: "What's this image about?"},
	//	{InlineData: &genai.Blob{Data: imageBytes, MIMEType: "image/jpeg"}},
	//}
	//

	model := req.Model
	if model == "" {
		model = DefaultModel
	}

	result, err := t.c.Models.GenerateContent(ctx, model, []*genai.Content{{Role: "user", Parts: req.Parts}}, req.Config)

	if err != nil {
		log.Errorw("GenerateContent err", err, "response", result)
		return "", err
	}

	for _, x := range result.Candidates {

		if x.Content == nil {
			continue
		}

		for _, xx := range x.Content.Parts {
			if xx.Text != "" {
				return xx.Text, nil
			}
		}

	}

	return "", errors.New("no content found")
}

func (t *Client) GenerateContentStream(ctx context.Context, req GenerateContentRequest) iter.Seq2[*genai.GenerateContentResponse, error] {
	model := req.Model
	if model == "" {
		model = DefaultModel
	}

	return t.c.Models.GenerateContentStream(ctx, model, []*genai.Content{{Role: "user", Parts: req.Parts}}, req.Config)
}

type GenerateImageRequest struct {
	Model       string
	Images      []string
	ImageBytes  [][]byte
	Videos      [][]byte
	Prompt      string
	ImageSize   string
	AspectRatio string
	//Count  int
}

func (t *Client) GenerateImage(ctx context.Context, req GenerateImageRequest) ([]byte, error) {
	var loop int
	for {
		fmt.Printf("GenerateImage attempt %d\n", loop+1)
		res, err := t.generateImage(ctx, req)
		if err == nil && len(res) > 0 {
			return res, nil
		}

		loop++
		if loop >= 3 {
			if err != nil {
				return nil, err
			}
			return nil, errors.New("GenerateImage failed: empty result after retries")
		}

		fmt.Printf("GenerateImage failed (err: %v), retrying in 1s...\n", err)
		time.Sleep(1 * time.Second)
	}
}

func (t *Client) GenerateImageV2(ctx context.Context, req GenerateImageRequest) ([]byte, error) {
	var loop int
	for {
		fmt.Printf("GenerateImageV2 attempt %d\n", loop+1)
		res, err := t.generateImageV2(ctx, req)
		if err == nil && len(res) > 0 {
			return res, nil
		}

		loop++
		if loop >= 3 {
			if err != nil {
				return nil, err
			}
			return nil, errors.New("GenerateImageV2 failed: empty result after retries")
		}

		fmt.Printf("GenerateImageV2 failed (err: %v), retrying in 1s...\n", err)
		time.Sleep(1 * time.Second)
	}
}

func (t *Client) generateImage(ctx context.Context, req GenerateImageRequest) ([]byte, error) {
	imageSize := helper.OrString(req.ImageSize, "1K")
	aspectRatio := helper.OrString(req.AspectRatio, "9:16")

	model := req.Model
	if model == "" {
		if len(req.ImageBytes) > 0 || len(req.Images) > 0 {
			model = ModelGenmini3CapabilityPreview
		} else {
			model = ModelGenmini3ProImagePreview
		}
	}

	prompt := req.Prompt
	if prompt == "" {
		prompt = "high quality, photorealistic, cinematic lighting, highly detailed"
	}

	// Case 1: Text-to-Image (No input images)
	if len(req.ImageBytes) == 0 && len(req.Images) == 0 {
		resp, err := t.c.Models.GenerateImages(ctx, model, prompt, &genai.GenerateImagesConfig{
			AspectRatio:      aspectRatio,
			ImageSize:        imageSize,
			PersonGeneration: genai.PersonGenerationAllowAll,
		})
		if err != nil {
			return nil, fmt.Errorf("GenerateImages failed: %w", err)
		}
		if len(resp.GeneratedImages) == 0 || resp.GeneratedImages[0].Image == nil {
			return nil, errors.New("no image generated")
		}
		return resp.GeneratedImages[0].Image.ImageBytes, nil
	}

	// Case 2: Image Editing or Fusion
	var refs []genai.ReferenceImage

	// Collect all images from URL and Bytes
	var allImages [][]byte
	for _, url := range req.Images {
		p, err := NewImagePart(url)
		if err != nil {
			return nil, err
		}
		if p.InlineData != nil {
			allImages = append(allImages, p.InlineData.Data)
		}
	}
	allImages = append(allImages, req.ImageBytes...)

	if len(allImages) == 0 {
		return nil, errors.New("no images provided for editing")
	}

	// Construct ReferenceImages
	// First image is Raw (Base)
	refs = append(refs, genai.NewRawReferenceImage(&genai.Image{
		ImageBytes: allImages[0],
		MIMEType:   "image/jpeg",
	}, 1))

	// Subsequent images are Content (referenced as [2], [3], etc.)
	for i := 1; i < len(allImages); i++ {
		refs = append(refs, genai.NewContentReferenceImage(&genai.Image{
			ImageBytes: allImages[i],
			MIMEType:   "image/jpeg",
		}, int32(i+1)))

		// Ensure the prompt mentions the subjects if not already there
		// This is a bit of a hack to make multi-image fusion more likely to work
		if !strings.Contains(prompt, fmt.Sprintf("[%d]", i+1)) {
			prompt += fmt.Sprintf(" and refer to image [%d]", i+1)
		}
	}

	resp, err := t.c.Models.EditImage(ctx, model, prompt, refs, &genai.EditImageConfig{
		AspectRatio:      aspectRatio,
		PersonGeneration: genai.PersonGenerationAllowAll,
	})
	if err != nil {
		return nil, fmt.Errorf("EditImage failed: %w", err)
	}

	if len(resp.GeneratedImages) == 0 || resp.GeneratedImages[0].Image == nil {
		return nil, errors.New("no edited image generated")
	}

	return resp.GeneratedImages[0].Image.ImageBytes, nil
}

func (t *Client) generateImageV2(ctx context.Context, req GenerateImageRequest) ([]byte, error) {
	model := req.Model
	if model == "" {
		model = "gemini-3-pro-image-preview"
	}

	var parts []*genai.Part

	// Add images
	for _, url := range req.Images {
		p, err := NewImagePart(url)
		if err != nil {
			return nil, err
		}
		parts = append(parts, p)
	}
	for _, b := range req.ImageBytes {
		parts = append(parts, genai.NewPartFromBytes(b, "image/jpeg"))
	}

	// Add text prompt
	prompt := req.Prompt
	if prompt == "" {
		prompt = "high quality, photorealistic, cinematic lighting, highly detailed"
	}
	parts = append(parts, genai.NewPartFromText(prompt))

	config := &genai.GenerateContentConfig{
		ResponseModalities: []string{"TEXT", "IMAGE"},
	}

	resp, err := t.c.Models.GenerateContent(ctx, model, []*genai.Content{{Role: "user", Parts: parts}}, config)
	if err != nil {
		return nil, fmt.Errorf("GenerateContent for image failed: %w", err)
	}

	for _, candidate := range resp.Candidates {
		if candidate.Content == nil {
			continue
		}
		for _, part := range candidate.Content.Parts {
			if part.InlineData != nil && strings.HasPrefix(part.InlineData.MIMEType, "image/") {
				return part.InlineData.Data, nil
			}
		}
	}

	return nil, errors.New("no image part found in multimodal response")
}

func (t *Client) GenerateBlob(ctx context.Context, req GenerateContentRequest) (*genai.Blob, error) {

	var loop int

	for {

		fmt.Println("GenerateBlob ing", loop)

		blob, err := t.generateBlob(ctx, req)
		if err != nil {
			return nil, err
		}

		fmt.Println("GenerateBlob done", loop)

		if blob == nil || len(blob.Data) == 0 {
			loop += 1
			fmt.Println("GenerateBlob got empty blob:  try " + conv.Str(loop))

			if loop >= 2 {
				return nil, errors.New("too many blob upload")
			}

			continue
		}

		return blob, nil
	}
}

func (t *Client) generateBlob(ctx context.Context, req GenerateContentRequest) (*genai.Blob, error) {

	response, err := t.c.Models.GenerateContent(ctx, req.Model, []*genai.Content{{Role: "user", Parts: req.Parts}}, req.Config)

	if err != nil {
		return nil, err
	}

	for _, x := range response.Candidates {
		for _, xx := range x.Content.Parts {
			if xx.InlineData != nil && len(xx.InlineData.Data) > 0 {
				return xx.InlineData, nil
			}
		}
	}

	// 尝试提取拒绝原因
	var refusal string
	for _, x := range response.Candidates {
		if x.Content != nil {
			for _, xx := range x.Content.Parts {
				if xx.Text != "" {
					refusal += xx.Text + " "
				}
			}
		}
	}

	return nil, fmt.Errorf("%s", refusal)
}
