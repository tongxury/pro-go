package biz

import (
	"context"
	"errors"
	"fmt"
	projpb "store/api/proj"
	"store/app/proj-pro/internal/data"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/third/gemini"

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

	script := dataBus.GetSegmentScript().GetScript()
	if script == "" {
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
我提供给你的脚本: %s
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

	log.Debug("start 1")

	aspectRatio := helper.OrString(wfState.GetDataBus().GetSettings().GetAspectRatio(), "9:16")
	blob, err := t.data.GenaiFactory.Get().GenerateImage(ctx, gemini.GenerateImageRequest{
		Images: keyFrame.Refs,
		//Videos: [][]byte{seg.Content},
		Prompt:      keyFrame.Prompt,
		AspectRatio: aspectRatio,
		ImageSize:   "2K",
		//Count: 8,
	})

	log.Debug("start 2")

	if err != nil {
		logger.Errorw("create image fail err", err)
		return nil, err
	}

	imageUrl, err := t.data.TOS.PutImageBytes(ctx, blob)
	if err != nil {
		logger.Errorw("put image bytes error", err)
		return nil, err
	}
	log.Debug("start 3")

	keyFrame.Url = imageUrl
	keyFrame.Status = ExecuteStatusReviewing

	_, err = t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId, mgz.Op().
		Set(fmt.Sprintf("jobs.%d.dataBus.keyFrames", jobState.Index), &projpb.KeyFrames{
			Frames: []*projpb.KeyFrames_Frame{keyFrame},
		}))
	if err != nil {
		logger.Errorw("update segment keyframe fail", "err", err)
		return nil, err
	}

	return &ExecuteResult{
		Status: ExecuteStatusCompleted,
	}, nil
}
