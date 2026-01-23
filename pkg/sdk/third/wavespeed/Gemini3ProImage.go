package wavespeed

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Gemini3ProImageRequest struct {
	Prompt string   `json:"prompt,omitempty"`
	Images []string `json:"images,omitempty"`
	/*
		1:1, 3:2, 2:3, 3:4, 4:3, 4:5, 5:4, 9:16, 16:9, 21:9
	*/
	AspectRatio string `json:"aspect_ratio,omitempty"`
	// 1k, 2k, 4k
	Resolution string `json:"resolution,omitempty"`
	// png, jpeg
	OutputFormat       string `json:"output_format,omitempty"`
	EnableSyncMode     bool   `json:"enable_sync_mode,omitempty"`
	EnableBase64Output bool   `json:"enable_base64_output,omitempty"`
}

type Gemini3ProImageResponse struct{}

func (t *Client) Gemini3ProImage(ctx context.Context, req Gemini3ProImageRequest) (*Response, error) {

	var response Response
	err := t.invoke(ctx, "/api/v3/google/gemini-3-pro-image/edit", req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (t *Client) GetResult(ctx context.Context, id string) (*Response, error) {

	post, err := resty.New().R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+t.conf.APIKey).
		Get(t.conf.BaseURL + fmt.Sprintf("/api/v3/predictions/%s/result", id))

	if err != nil {
		return nil, err
	}

	//fmt.Println(string(post.Body()))

	var result Response
	err = json.Unmarshal(post.Body(), &result)
	if err != nil {
		return nil, err
	}

	if result.Code != 200 {
		return nil, errors.New(result.Message)
	}

	return &result, nil
}

//curl --location --request GET "https://api.wavespeed.ai/api/v3/predictions/${requestId}/result" \
//--header "Authorization: Bearer ${WAVESPEED_API_KEY}"
