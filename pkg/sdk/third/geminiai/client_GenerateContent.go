package geminiai

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/genai"
)

type GenerateContentRequest struct {
	Model  string
	Config *genai.GenerationConfig
	Parts  []genai.Part
}

func (t *Client) GenerateContent(ctx context.Context, req GenerateContentRequest) (string, error) {

	gm := t.c.GenerativeModel(req.Model)

	if req.Config != nil {
		gm.GenerationConfig = *req.Config
	}

	response, err := gm.GenerateContent(ctx, req.Parts...)

	if err != nil {
		log.Errorw("GenerateContent err", err, "response", response)
		return "", err
	}

	text := ResponseToString(response)

	return text, nil
}

type GenerateImageRequest struct {
	Images []string
	Prompt string
}

func (t *Client) GenerateImage(ctx context.Context, req GenerateImageRequest) ([]byte, error) {

	var parts []genai.Part

	for i := range req.Images {
		gsFile, err := t.UploadFile(ctx, req.Images[i], "image/jpeg")
		if err != nil {

		}
		parts = append(parts, genai.FileData{
			MIMEType: "image/jpeg",
			URI:      gsFile,
		})
	}

	parts = append(parts, genai.Text(req.Prompt))

	blob3, err := t.GenerateBlob(ctx, GenerateContentRequest{
		Model: "gemini-2.5-flash-image-preview",
		Parts: parts,
	})
	if err != nil {
		return nil, err
	}

	//b64Image := base64.StdEncoding.EncodeToString(blob3.Data)

	return blob3.Data, nil

}

func (t *Client) GenerateBlob(ctx context.Context, req GenerateContentRequest) (*genai.Blob, error) {

	for {
		blob, err := t.generateBlob(ctx, req)
		if err != nil {
			return nil, err
		}

		if len(blob.Data) == 0 {
			continue
		}

		return blob, nil
	}
}

func (t *Client) generateBlob(ctx context.Context, req GenerateContentRequest) (*genai.Blob, error) {

	gm := t.c.GenerativeModel(req.Model)

	if req.Config != nil {
		gm.GenerationConfig = *req.Config
	}

	response, err := gm.GenerateContent(ctx, req.Parts...)

	if err != nil {
		log.Errorw("GenerateContent err", err)
		return nil, err
	}

	for _, x := range response.Candidates {
		for _, xx := range x.Content.Parts {
			switch y := xx.(type) {
			case genai.Text:
				fmt.Println(string(y))
			case genai.Blob:
				return &y, nil
			}
		}
	}

	return nil, errors.New("blob content is empty")
}
