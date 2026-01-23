package gemini

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/genai"
)

type AnalyzeVideoRequest struct {
	VideoBytes []byte
	Prompt     string
	Model      string
	Config     *genai.GenerateContentConfig
}

func (t *Client) AnalyzeVideo(ctx context.Context, req AnalyzeVideoRequest) (string, error) {

	parts := []*genai.Part{
		genai.NewPartFromBytes(req.VideoBytes, "video/mp4"),
		{
			Text: req.Prompt,
		},
	}

	model := req.Model
	if model == "" {
		model = DefaultModel
	}

	result, err := t.c.Models.GenerateContent(ctx, model, []*genai.Content{{Role: "user", Parts: parts}}, req.Config)

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
