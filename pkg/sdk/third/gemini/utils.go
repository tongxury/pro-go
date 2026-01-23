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
