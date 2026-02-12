package service

import (
	"context"
	"fmt"
	projpb "store/api/proj"
	"store/pkg/sdk/helper/stringz"
	"store/pkg/sdk/third/bytedance/volcengine"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

func (t *WorkFlowService) CreateRemixTask(ctx context.Context, req *projpb.CreateRemixTaskRequest) (*projpb.RemixTask, error) {

	if len(req.Items) == 0 {
		return nil, fmt.Errorf("items is empty")
	}

	var groups []volcengine.Group

	for _, item := range req.Items {

		material, err := t.data.Volcengine.CreateUrlMaterial(ctx, volcengine.CreateMaterialParams{
			MediaFirstCategory: "video",
			Title:              "title",
			MaterialUrl:        item.Url,
			Wait:               true,
		})
		if err != nil {
			log.Errorw("CreateUrlMaterial err", err, "url", item.Url)
			return nil, err
		}

		var tts []volcengine.TtsConfig
		var timeStrategy int

		if item.Subtitle != "" {
			parts := strings.Split(item.Subtitle, "，")
			var texts []volcengine.TextConfig

			for _, xx := range parts {
				// Split long sentences
				for _, xxx := range stringz.SplitToSlice(xx, 20) {
					texts = append(texts, volcengine.TextConfig{
						Text:         xxx,
						FontSize:     "15",
						InnerPadding: 0.5,
						Bold:         1,
						Alignment:    "center",
						Y:            0.8,
					})
				}
			}

			tts = []volcengine.TtsConfig{
				{
					AddAudio:    1,
					AddSubtitle: 1,
					Subtitle:    texts,
					// Default settings, can be parameterized if needed
					Speed: 150,
				},
			}
			timeStrategy = 1
		}

		groups = append(groups, volcengine.Group{
			MuseList: []volcengine.Muse{
				{
					Type:      "video",
					MuseId:    material.MediaId,
					StartTime: 0,
					EndTime:   0,
				},
			},
			TtsConfig: tts,
			TimeStrategy: volcengine.TimeStrategy{
				Type: timeStrategy,
			},
			Type:   2,
			Volume: 0,
		})
	}

	params := volcengine.CreateMixCutTaskAsyncParams{
		Name:             "混剪",
		Group:            groups,
		Height:           1920,
		Width:            1080,
		TransitionConfig: &volcengine.TransitionConfig{
			//Type: 1,
		},
		FilterConfig: &volcengine.FilterConfig{
			Type: 0,
		},
	}

	cutTask, err := t.data.Volcengine.CreateMixCutTask(ctx, params)
	if err != nil {
		log.Errorw("create cut task err", err)
		return nil, err
	}

	return &projpb.RemixTask{
		TaskId: cutTask.TaskKey,
		Status: "running", // Initial status
	}, nil
}

func (t *WorkFlowService) GetRemixTask(ctx context.Context, req *projpb.GetRemixTaskRequest) (*projpb.RemixTask, error) {

	result, err := t.data.Volcengine.QueryMixCutTaskResult(ctx, volcengine.QueryMixCutTaskResultParams{
		TaskKey: req.TaskId,
	})
	if err != nil {
		return nil, err
	}

	status := "running"
	url := ""

	// Status 200 means success in Volcengine API
	if result.Data.Task.Status == 200 {
		status = "completed"
		if len(result.Data.Task.VideoList) > 0 {
			url = result.Data.Task.VideoList[0].DownloadUrl
		}
	} else if result.Data.Task.Status == 201 { // Assuming 201 or similar code for failed. Need to check status codes.
		// For now map 200 as completed, others running/failed.
		// If status is not documented, we can just return string representation of status code
		status = fmt.Sprintf("%d", result.Data.Task.Status)
	} else {
		// Just return the status code as string if not 200
		status = fmt.Sprintf("%d", result.Data.Task.Status)
	}

	return &projpb.RemixTask{
		TaskId: req.TaskId,
		Status: status,
		Url:    url,
	}, nil
}
