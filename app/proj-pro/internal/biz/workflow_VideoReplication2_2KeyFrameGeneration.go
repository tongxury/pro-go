package biz

import (
	"context"
	"errors"
	"fmt"
	projpb "store/api/proj"
	"store/app/proj-pro/internal/data"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/third/wavespeed"

	"github.com/go-kratos/kratos/v2/log"
)

type VideoReplication2_KeyFrameGenerationJob struct {
	data *data.Data
}

func (t VideoReplication2_KeyFrameGenerationJob) Initialize(ctx context.Context, options Options) error {
	return nil
}

func (t VideoReplication2_KeyFrameGenerationJob) GetName() string {
	return "keyFramesGenerationJob"
}

func (t VideoReplication2_KeyFrameGenerationJob) Execute(ctx context.Context, jobState *projpb.Job, wfState *projpb.Workflow) (status *ExecuteResult, err error) {

	dataBus := GetDataBus(wfState)

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "VideoReplication2_KeyFrameGenerationJob.Execute",
		"workflowId ", wfState.XId,
		"jobState.Name ", jobState.Name,
		"jobState.Index ", jobState.Index,
		"jobState.Status ", jobState.Status,
	))

	segment := dataBus.Segment
	if segment == nil {
		return nil, errors.New("segment not found")
	}

	commodity := dataBus.Commodity
	if commodity == nil {
		return nil, errors.New("commodity not found")
	}

	script := dataBus.GetSegmentScript().GetSegments()
	if len(script) == 0 {
		return nil, errors.New("script not found")
	}

	if len(dataBus.KeyFrames.GetFrames()) > 0 && dataBus.KeyFrames.GetFrames()[0].Url != "" {
		return &ExecuteResult{
			Status: ExecuteStatusCompleted,
		}, nil
	}

	// 初始化
	var keyFrame *projpb.KeyFrames_Frame
	if len(dataBus.KeyFrames.GetFrames()) == 0 {

		prompt, err := t.data.Mongo.Settings.GetPrompt(ctx, "sora2_grid_images")
		if err != nil {
			return nil, err
		}

		keyFrame = &projpb.KeyFrames_Frame{
			Status: ExecuteStatusRunning,
			Refs: []string{
				segment.HighlightFrames[0].Url,
				helper.SliceElement[string](
					helper.Mapping(dataBus.Commodity.GetMedias(),
						func(x *projpb.Media) string {
							return x.Url
						}),
					0, false),
			},
			Prompt: fmt.Sprintf(`
	%s
===
我提供给你的脚本: %v
`, prompt.Content,
				script,
			),
		}

		_, err = t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId, mgz.Op().
			Set(fmt.Sprintf("jobs.%d.dataBus.keyFrames", jobState.Index), &projpb.KeyFrames{
				Frames: []*projpb.KeyFrames_Frame{keyFrame},
			}))
		if err != nil {
			logger.Errorw("update segment keyframe fail", "err", err)
			return nil, err
		}

	} else {
		keyFrame = dataBus.KeyFrames.GetFrames()[0]
	}

	if keyFrame.Url != "" {
		return &ExecuteResult{
			Status: ExecuteStatusCompleted,
		}, nil
	}

	if keyFrame.TaskId != "" {
		result, err := t.data.Wavespeed.GetResult(ctx, keyFrame.TaskId)
		if err != nil {
			logger.Errorw("wavespeed.GetResult fail", "err", err)
			return nil, err
		}

		if len(result.Data.Outputs) > 0 {
			keyFrame.Url = result.Data.Outputs[0]
			keyFrame.Status = ExecuteStatusReviewing

			_, err = t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId, mgz.Op().
				Set(fmt.Sprintf("jobs.%d.dataBus.keyFrames.frames.0", jobState.Index), keyFrame))
			if err != nil {
				logger.Errorw("update keyframe fail", "err", err)
				return nil, err
			}

			return &ExecuteResult{
				Status: ExecuteStatusCompleted,
			}, nil
		}

		if result.Data.Status == "failed" {
			keyFrame.TaskId = ""
			_, err = t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId, mgz.Op().
				Set(fmt.Sprintf("jobs.%d.dataBus.keyFrames.frames.0", jobState.Index), keyFrame))
			if err != nil {
				logger.Errorw("update keyframe fail", "err", err)
				return nil, err
			}
			return nil, nil
		}

		return nil, nil // Running
	}

	aspectRatio := helper.OrString(wfState.GetDataBus().GetSettings().GetAspectRatio(), "9:16")
	res, err := t.data.Wavespeed.Gemini3ProImage(ctx, wavespeed.Gemini3ProImageRequest{
		Prompt:      keyFrame.Prompt,
		Images:      keyFrame.Refs,
		AspectRatio: aspectRatio,
		Resolution:  "2k",
	})

	if err != nil {
		logger.Errorw("wavespeed.Gemini3ProImage fail", "err", err)
		return nil, err
	}

	keyFrame.TaskId = res.Data.Id
	_, err = t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId, mgz.Op().
		Set(fmt.Sprintf("jobs.%d.dataBus.keyFrames.frames.0", jobState.Index), keyFrame))
	if err != nil {
		logger.Errorw("update keyframe taskId fail", "err", err)
		return nil, err
	}

	return nil, nil
}
