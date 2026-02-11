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

type VideoReplication3_VideoSegmentsGenerationJob struct {
	data *data.Data
}

func (t VideoReplication3_VideoSegmentsGenerationJob) Initialize(ctx context.Context, options Options) error {
	return nil
}

func (t VideoReplication3_VideoSegmentsGenerationJob) GetName() string {
	return "videoSegmentsGenerationJob"
}

func (t VideoReplication3_VideoSegmentsGenerationJob) Execute(ctx context.Context, jobState *projpb.Job, wfState *projpb.Workflow) (status *ExecuteResult, err error) {

	dataBus := GetDataBus(wfState)

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "VideoSegmentsGenerationJob.Execute",
		"workflowId ", wfState.XId,
		"jobState.Name ", jobState.Name,
		"jobState.Index ", jobState.Index,
		"jobState.Status ", jobState.Status,
	))

	segments := dataBus.SegmentScript.GetSegments()
	if len(segments) == 0 {
		return nil, errors.New("segments data not found")
	}

	keyFrames := dataBus.KeyFrames
	if len(keyFrames.Frames) == 0 {
		return nil, errors.New("keyFrames data not found")
	}

	jobData := dataBus.VideoGenerations
	if len(jobData) == 0 {

		var vgs []*projpb.VideoGeneration

		//for i := 0; i < len(keyFrames.GetFrames())-1; i++ {
		//
		//	vgs = append(vgs, &projpb.VideoGeneration{
		//		Url:        "",
		//		TaskId:     "",
		//		Status:     "running",
		//		CoverUrl:   "",
		//		FirstFrame: keyFrames.Frames[i].Url,
		//		LastFrame:  keyFrames.Frames[i+1].Url,
		//		Prompt:     "",
		//	})
		//}

		for i := 0; i < len(keyFrames.GetFrames()); i = i + 2 {

			first := keyFrames.GetFrames()[i].Url

			var last string

			if len(keyFrames.GetFrames())-1 >= i+1 {
				last = keyFrames.GetFrames()[i+1].Url
			}

			seg := segments[i/2]

			//

			vgs = append(vgs, &projpb.VideoGeneration{
				Url:         "",
				TaskId:      "",
				Status:      "running",
				CoverUrl:    "",
				FirstFrame:  first,
				LastFrame:   last,
				AspectRatio: wfState.GetDataBus().GetSettings().GetAspectRatio(),
				// tmp 用ai 根据这生成新的
				TmpPrompt: fmt.Sprintf("%s\n%s\n\n%s\n\n%s\n\n%s",
					seg.GetCoreAction(),
					seg.GetElementTransformation(),
					seg.GetVisualChange(),
					//seg.GetDescription(),
					seg.GetSceneStyle(),
					seg.GetScript(),
					//strings.Join(seg.GetTypedTags().Tags(), " "),
					//"依据以上分镜描述， 结合【图1的整体风格和场景】+【图2中的商品】,生成一张视频首帧图片",
					//helper.Select(isFirstFrame,
					//	,
					//"依据以上分镜描述， 结合【图1的整体风格和场景】,推演视频尾帧图片。要符合物理规律",
					//),
				),
			})
		}

		t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId,
			mgz.Op().Set(fmt.Sprintf("jobs.%d.dataBus.videoGenerations", jobState.Index), vgs),
		)

		return nil, nil
	}

	// 已完成
	runnings := helper.Filter(dataBus.VideoGenerations, func(param *projpb.VideoGeneration) bool {
		return param.Prompt == ""
	})

	if len(runnings) == 0 {
		return &ExecuteResult{
			Status: ExecuteStatusCompleted,
		}, nil
	}

	prompt, err := t.data.Mongo.Settings.GetPrompt(ctx, "jimeng_video_generation_by_first_last_frames")
	if err != nil {
		return nil, err
	}

	// 生成运镜描述
	if len(helper.Filter(dataBus.VideoGenerations, func(param *projpb.VideoGeneration) bool {
		return param.Prompt == ""
	})) > 0 {

		for i, x := range dataBus.VideoGenerations {
			if x.Prompt != "" {
				continue
			}

			dataBus.VideoGenerations[i].Prompt, err = t.data.GenaiFactory.Get().GenerateContentV2(ctx, gemini.GenerateContentRequestV2{
				ImageUrls: []string{x.FirstFrame, x.LastFrame},
				Prompt: fmt.Sprintf(`

生成提示词的逻辑:
%s

===

参考一下信息：

- 分镜信息: %s
`, prompt,
					x.TmpPrompt,
				),
			})
			if err != nil {
				logger.Errorw("GenerateContentV2 err", err)
				//return nil, err
				continue
			}
		}

		// 发起后立即更新 避免重复
		t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId,
			mgz.Op().Set(fmt.Sprintf("jobs.%d.dataBus.videoGenerations", jobState.Index), dataBus.VideoGenerations),
		)

		return nil, nil
	}

	//for i := range dataBus.VideoGenerations {
	//
	//	x := dataBus.VideoGenerations[i]
	//
	//	if x.Url != "" {
	//		continue
	//	}
	//
	//	if x.TaskId == "" {
	//
	//		taskId, err := t.data.Arkr.GenerateVideo(ctx, arkr.GenerateVideoRequest{
	//			Prompt:        x.Prompt,
	//			StartFrame:    x.FirstFrame,
	//			EndFrame:      x.LastFrame,
	//			GenerateAudio: helper.Pointer(false),
	//			Duration:      x.Duration,
	//			AspectRatio:   x.AspectRatio,
	//		})
	//		if err != nil {
	//			logger.Errorw("GenerateVideo err", err)
	//			//return nil, err
	//			continue
	//		}
	//
	//		// 发起后立即更新 避免重复
	//		t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId,
	//			mgz.Op().SetListItem(fmt.Sprintf("jobs.%d.dataBus.videoGenerations", jobState.Index), i, "taskId", taskId),
	//		)
	//
	//	} else {
	//
	//		task, err := t.data.Arkr.GetTask(ctx, x.TaskId)
	//		if err != nil {
	//			logger.Errorw("GetTask err", err)
	//			return nil, err
	//		}
	//
	//		logger.Debugw("GetTask task ", "ing", "id", task.ID, "status", task.Status)
	//
	//		if task.Status == "succeeded" {
	//
	//			video, err := t.data.TOS.PutVideo(ctx, task.Content.VideoURL)
	//			if err != nil {
	//				return nil, err
	//			}
	//
	//			image, err := t.data.TOS.PutImage(ctx, task.Content.LastFrameURL)
	//			if err != nil {
	//
	//				return nil, err
	//			}
	//
	//			t.data.Mongo.Workflow.UpdateByIDIfExists(ctx,
	//				wfState.XId,
	//				mgz.Op().
	//					SetListItem(fmt.Sprintf("jobs.%d.dataBus.videoGenerations", jobState.Index), i, "url", video).
	//					SetListItem(fmt.Sprintf("jobs.%d.dataBus.videoGenerations", jobState.Index), i, "duration", *task.Duration).
	//					SetListItem(fmt.Sprintf("jobs.%d.dataBus.videoGenerations", jobState.Index), i, "status", "completed").
	//					SetListItem(fmt.Sprintf("jobs.%d.dataBus.videoGenerations", jobState.Index), i, "lastFrame", image),
	//			)
	//		} else if task.Status == "failed" {
	//			t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId,
	//				mgz.Op().
	//					SetListItem(fmt.Sprintf("jobs.%d.dataBus.videoGenerations", jobState.Index), i, "status", "failed").
	//					SetListItem(fmt.Sprintf("jobs.%d.dataBus.videoGenerations", jobState.Index), i, "error", task.Error.Message),
	//			)
	//		}
	//	}
	//}

	return nil, nil
}
