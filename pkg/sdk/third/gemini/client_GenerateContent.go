package gemini

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper"

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

	var parts []*genai.Part

	imageSize := helper.OrString(req.ImageSize, "1K")
	aspectRatio := helper.OrString(req.AspectRatio, "9:16")

	for i := range req.Images {
		p, err := NewImagePart(req.Images[i])
		if err != nil {
			return nil, err
		}

		parts = append(parts, p)
	}

	for i := range req.ImageBytes {
		parts = append(parts, genai.NewPartFromBytes(req.ImageBytes[i], "image/jpeg"))
	}

	for i := range req.Videos {
		//gsFile, err := t.UploadBlob(ctx, req.Videos[i], "video/mp4")
		//if err != nil {
		//	return nil, err
		//}
		//parts = append(parts, &genai.Part{
		//	FileData: &genai.FileData{
		//		MIMEType: "video/mp4",
		//		FileURI:  gsFile,
		//	},
		//})

		parts = append(parts, genai.NewPartFromBytes(req.Videos[i], "video/mp4"))
		//if err != nil {, "video/mp4"))
	}

	parts = append(parts, &genai.Part{
		Text: req.Prompt + ", high quality, photorealistic, cinematic lighting, highly detailed",
	})

	blob3, err := t.GenerateBlob(ctx, GenerateContentRequest{
		Model: ModelGenmini3ProImagePreview,
		Parts: parts,
		//Config: config,
		Config: &genai.GenerateContentConfig{
			ImageConfig: &genai.ImageConfig{
				AspectRatio: aspectRatio,
				ImageSize:   imageSize,
				//OutputMIMEType:           "",
				//OutputCompressionQuality: nil,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	//b64Image := base64.StdEncoding.EncodeToString(blob3.Data)

	return blob3.Data, nil

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
