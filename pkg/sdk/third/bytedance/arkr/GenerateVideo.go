package arkr

import (
	"fmt"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

type GenerateVideoRequest struct {
	Model          string
	Prompt         string
	StartFrame     string
	EndFrame       string
	ReferenceImage string
	Duration       int64
	GenerateAudio  *bool
	AspectRatio    string
}

func (t *Client) GenerateVideo(ctx volcengine.Context, params GenerateVideoRequest) (string, error) {

	modelName := params.Model
	if modelName == "" {
		modelName = DefaultModel
	}

	var duration = &[]int64{4}[0]
	if params.Duration != 0 {
		duration = &[]int64{params.Duration}[0]
	}

	var aspectRatio = "9:16"
	if params.AspectRatio != "" {
		aspectRatio = params.AspectRatio
	}

	content := []*model.CreateContentGenerationContentItem{
		{
			Type: "text",
			Text: volcengine.String(params.Prompt + fmt.Sprintf("--rt %s --rs 720p --dur %d --cf false", aspectRatio, *duration)),
		},
	}

	if params.StartFrame != "" {
		content = append(content, &model.CreateContentGenerationContentItem{
			Type: "image_url",
			ImageURL: &model.ImageURL{
				URL: params.StartFrame,
			},
			Role: volcengine.String("first_frame"),
		})
	}

	if params.EndFrame != "" {
		content = append(content, &model.CreateContentGenerationContentItem{
			Type: "image_url",
			ImageURL: &model.ImageURL{
				URL: params.EndFrame,
			},
			Role: volcengine.String("last_frame"),
		})
	}

	req := model.CreateContentGenerationTaskRequest{
		Model:           modelName,
		Content:         content,
		ReturnLastFrame: volcengine.Bool(true),
		GenerateAudio:   params.GenerateAudio,
		//Duration:        duration,
	}

	resp, err := t.c.CreateContentGenerationTask(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}
