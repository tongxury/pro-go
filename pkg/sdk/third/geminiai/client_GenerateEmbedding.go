package geminiai

import (
	"context"
	"fmt"
	"google.golang.org/genai"
	"sort"
	"strings"
)

func (t *Client) GenerateEmbedding(ctx context.Context, keywords ...string) ([]float32, error) {
	model := t.c.EmbeddingModel("text-embedding-004")

	sort.Strings(keywords)

	p2 := strings.Join(keywords, " ")

	part := genai.Text(p2)
	//var parts []genai.Part
	//
	//for _, x := range keywords {
	//	parts = append(parts, genai.Text(x))
	//}

	resp, err := model.EmbedContent(ctx, part)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	if len(resp.Embedding.Values) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	return resp.Embedding.Values, nil
}

func (t *Client) GenerateEmbeddingV2(ctx context.Context, keyword string) ([]float32, error) {
	model := t.c.EmbeddingModel("text-embedding-004")

	part := genai.Text(keyword)
	//var parts []genai.Part
	//
	//for _, x := range keywords {
	//	parts = append(parts, genai.Text(x))
	//}

	resp, err := model.EmbedContent(ctx, part)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	if len(resp.Embedding.Values) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	return resp.Embedding.Values, nil
}
