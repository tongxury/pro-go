package gemini

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/genai"
)

type GenerateImagesRequest struct {
	Model       string
	Prompt      string
	ImageSize   string
	AspectRatio string
	Count       int
	Images      []string
	ImageBytes  [][]byte
}

func (t *Client) GenerateImages(ctx context.Context, req GenerateImagesRequest) (*genai.GenerateImagesResponse, error) {

	model := req.Model
	if model == "" {
		model = ModelGenmini3ProImagePreview
	}

	var parts []*genai.Part

	for i := range req.Images {
		p, err := NewImagePart(req.Images[i])
		if err != nil {
			return nil, err
		}
		parts = append(parts, p)
	}

	for i := range req.ImageBytes {
		parts = append(parts, genai.NewPartFromBytes(req.ImageBytes[i], "image/png"))
	}

	if req.Prompt != "" {
		parts = append(parts, &genai.Part{
			Text: req.Prompt,
		})
	}

	config := &genai.GenerateImagesConfig{
		AspectRatio:    req.AspectRatio,
		ImageSize:      req.ImageSize,
		NumberOfImages: int32(req.Count),
	}

	return t.c.Models.GenerateImages(ctx, model, req.Prompt, config)
}

type EditImageRequest struct {
	Model           string
	Prompt          string
	ReferenceImages []ReferenceImageConfig
	AspectRatio     string
	EditMode        genai.EditMode
}

type ReferenceImageConfig struct {
	URL      string
	Bytes    []byte
	MimeType string
	Type     ReferenceImageType
}

type ReferenceImageType string

const (
	ReferenceImageTypeRaw     ReferenceImageType = "raw"
	ReferenceImageTypeStyle   ReferenceImageType = "style"
	ReferenceImageTypeSubject ReferenceImageType = "subject"
	ReferenceImageTypeContent ReferenceImageType = "content"
)

func (t *Client) EditImage(ctx context.Context, req EditImageRequest) (*genai.EditImageResponse, error) {
	model := req.Model
	if model == "" {
		model = ModelGenmini3ProImagePreview
	}

	var refImages []genai.ReferenceImage
	for i, ref := range req.ReferenceImages {
		var imgBytes []byte
		var err error
		mimeType := ref.MimeType
		if mimeType == "" {
			mimeType = "image/jpeg"
		}

		if len(ref.Bytes) > 0 {
			imgBytes = ref.Bytes
		} else if ref.URL != "" {
			var resp *http.Response
			resp, err = http.Get(ref.URL)
			if err != nil {
				return nil, fmt.Errorf("failed to download image %d: %w", i, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("failed to download image %d, status: %d", i, resp.StatusCode)
			}

			imgBytes, err = io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to read body for image %d: %w", i, err)
			}
		}

		if len(imgBytes) == 0 {
			return nil, fmt.Errorf("image %d is empty", i)
		}

		img := &genai.Image{
			ImageBytes: imgBytes,
			MIMEType:   mimeType,
		}

		var refImg genai.ReferenceImage
		// Use 1-indexed IDs to match prompt references like [1], [2]...
		refID := int32(i + 1)

		switch ref.Type {
		case ReferenceImageTypeRaw:
			refImg = genai.NewRawReferenceImage(img, refID)
		case ReferenceImageTypeStyle:
			refImg = genai.NewStyleReferenceImage(img, refID, nil)
		case ReferenceImageTypeSubject:
			refImg = genai.NewSubjectReferenceImage(img, refID, nil)
		case ReferenceImageTypeContent:
			refImg = genai.NewContentReferenceImage(img, refID)
		default:
			// Default to raw for first image, subject for others
			if i == 0 {
				refImg = genai.NewRawReferenceImage(img, refID)
			} else {
				refImg = genai.NewSubjectReferenceImage(img, refID, nil)
			}
		}
		refImages = append(refImages, refImg)
	}

	config := &genai.EditImageConfig{
		AspectRatio: req.AspectRatio,
		EditMode:    req.EditMode,
	}

	return t.c.Models.EditImage(ctx, model, req.Prompt, refImages, config)
}

func (t *Client) EditImageToBytes(ctx context.Context, req EditImageRequest) ([]byte, error) {
	resp, err := t.EditImage(ctx, req)
	if err != nil {
		return nil, err
	}

	if len(resp.GeneratedImages) > 0 {
		return resp.GeneratedImages[0].Image.ImageBytes, nil
	}

	return nil, fmt.Errorf("no image generated")
}
