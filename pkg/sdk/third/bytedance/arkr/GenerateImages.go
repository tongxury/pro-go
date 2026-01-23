package arkr

import (
	"context"
	"store/pkg/sdk/helper/wg"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
)

func (t *Client) GenerateImagesInBatch(ctx context.Context, params model.GenerateImagesRequest, count int) ([]string, error) {

	var images []string

	var params_ []model.GenerateImagesRequest
	for i := 0; i < count; i++ {
		params_ = append(params_, params)
	}

	images, errs := wg.WaitGroupResults(ctx, params_, func(ctx context.Context, param model.GenerateImagesRequest) (string, error) {

		imagesResponse, err := t.c.GenerateImages(ctx, param)
		if err != nil {
			return "", err
		}

		return *imagesResponse.Data[0].Url, nil
	})

	if errs != nil {
		return nil, errs[0]
	}

	return images, nil
}
