package geminiai

import (
	"bytes"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-resty/resty/v2"
	"google.golang.org/genai"
	"time"
)

func (t *Client) UploadBlob(ctx context.Context, content []byte, mimeType string) (string, error) {

	fileBody := bytes.NewReader(content)

	t1 := time.Now()
	log.Debugw("start to upload file to genai time", t1)

	ctx = context.Background()

	opts := genai.UploadFileOptions{DisplayName: "", MIMEType: mimeType}
	response, err := t.c.UploadFile(ctx, "", fileBody, &opts)
	if err != nil {
		log.Errorw("upload file error", err)
		return "", err
	}

	// 校验
	for response.State == genai.FileStateProcessing {
		time.Sleep(1 * time.Second)
		// Fetch the file from the API again.
		response, err = t.c.GetFile(ctx, response.Name)
		if err != nil {
			log.Fatal(err)
			continue
		}

		log.Debugw("GenaiClient.GetFile response", response)
	}

	return response.URI, nil
}

func (t *Client) UploadFile(ctx context.Context, url, mimeType string) (string, error) {
	result, err := resty.New().R().Get(url)
	if err != nil {
		return "", err
	}

	return t.UploadBlob(ctx, result.Body(), mimeType)
}
