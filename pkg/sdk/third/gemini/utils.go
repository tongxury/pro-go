package gemini

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"google.golang.org/genai"
)

const (
	DefaultModel = ModelGenmini3FlashPreview
	//ModelGenmini3ProPreview      = "gemini-1.5-pro-002"
	//ModelGenmini3ProImagePreview = "imagen-4.0-generate-001"
	ModelGenmini3ProImagePreview   = "imagen-3.0-fast-generate-001"
	ModelGenmini3CapabilityPreview = "imagen-3.0-capability-001"
	ModelGenmini4ImagePreview      = "imagen-4.0-generate-001"
	ModelGenmini3FlashPreview      = "gemini-2.0-flash-001"
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
