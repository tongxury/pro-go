package gemini

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"google.golang.org/genai"
)

const (
	DefaultModel                 = ModelGenmini3FlashPreview
	ModelGenmini3ProPreview      = "gemini-3-pro-preview"
	ModelGenmini3ProImagePreview = "gemini-3-pro-image-preview"
	ModelGenmini3FlashPreview    = "gemini-3-flash-preview"

	ModelImagen40FastGenerate001 = "imagen-4.0-fast-generate-001"
	ModelImagen3Capability001    = "imagen-3.0-capability-001"
	//	imagen-4.0-generate-001
	//imagen-4.0-fast-generate-001
	//imagen-4.0-ultra-generate-001
	//imagen-3.0-generate-002
	//imagen-3.0-generate-001
	//imagen-3.0-fast-generate-001
	//imagen-3.0-capability-001
)

func NewVideoPart(url string) (*genai.Part, error) {

	v, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	all, err := io.ReadAll(v.Body)
	if err != nil {
		return nil, err
	}

	return genai.NewPartFromBytes(all, "video/mp4"), nil
}

func NewTextPart(text string) *genai.Part {
	return &genai.Part{
		Text: text,
	}
}

func NewImageParts(images []string) ([]*genai.Part, error) {

	var parts []*genai.Part

	for i := range images {
		p, err := NewImagePart(images[i])
		if err != nil {
			return nil, err
		}

		parts = append(parts, p)
	}

	return parts, nil
}

func NewMediaPart(url, mimeType string) (*genai.Part, error) {
	v, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	all, err := io.ReadAll(v.Body)
	if err != nil {
		return nil, err
	}

	return genai.NewPartFromBytes(all, mimeType), nil
}

func NewImagePart(url string) (*genai.Part, error) {

	v, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	all, err := io.ReadAll(v.Body)
	if err != nil {
		return nil, err
	}

	return genai.NewPartFromBytes(all, "image/jpeg"), nil
}

func ResponseToString(resp *genai.GenerateContentResponse) string {
	var b strings.Builder
	for i, x := range resp.Candidates {
		if len(resp.Candidates) > 1 {
			_, _ = fmt.Fprintf(&b, "%d:", i+1)
		}
		b.WriteString(ContentToString(x.Content))
	}
	return b.String()
}

func ContentToString(c *genai.Content) string {
	var b strings.Builder
	if c == nil || c.Parts == nil {
		return ""
	}
	for _, part := range c.Parts {
		if part.Text != "" {
			b.WriteString(part.Text)
		}
	}
	return b.String()
}
