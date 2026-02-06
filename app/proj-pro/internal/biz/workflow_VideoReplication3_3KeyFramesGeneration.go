package biz

import (
	"context"
	"fmt"
	projpb "store/api/proj"
	"store/app/proj-pro/internal/data"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/helper/wg"
	"store/pkg/sdk/third/gemini"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

type VideoReplication3_KeyFramesGenerationJob struct {
	data *data.Data
}

func (t VideoReplication3_KeyFramesGenerationJob) Initialize(ctx context.Context, options Options) error {
	return nil
}

func (t VideoReplication3_KeyFramesGenerationJob) GetName() string {
	return "keyFramesGenerationJob"
}

func (t VideoReplication3_KeyFramesGenerationJob) Execute(ctx context.Context, jobState *projpb.Job, wfState *projpb.Workflow) (status *ExecuteResult, err error) {

	dataBus := GetDataBus(wfState)

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "KeyFramesGenerationJob.Execute",
		"workflowId ", wfState.XId,
		"jobState.Name ", jobState.Name,
		"jobState.Index ", jobState.Index,
		"jobState.Status ", jobState.Status,
	))

	logger.Debug("keyFrames generation job start")

	doings := helper.Filter(dataBus.KeyFrames.GetFrames(),
		func(x *projpb.KeyFrames_Frame) bool {
			return x.Url == ""
		})

	// 已完成
	if len(dataBus.KeyFrames.GetFrames()) > 0 && len(doings) == 0 {
		return &ExecuteResult{
			Status: ExecuteStatusCompleted,
		}, nil
	}

	frames := dataBus.KeyFrames.GetFrames()
	// 初始化
	if len(frames) == 0 {
		var keyFrames []*projpb.KeyFrames_Frame

		for i := 0; i < len(dataBus.SegmentScript.GetSegments())*2; i++ {
			keyFrames = append(keyFrames, &projpb.KeyFrames_Frame{
				Status: ExecuteStatusWaiting,
			})
		}

		frames = keyFrames
	}

	firstFramePrompt, err := t.data.Mongo.Settings.GetPrompt(ctx, "segment_first_frame_generation")
	if err != nil {
		return nil, err
	}

	lastFramePrompt, err := t.data.Mongo.Settings.GetPrompt(ctx, "segment_last_frame_generation")
	if err != nil {
		return nil, err
	}

	//
	for i, x := range frames {

		if x.Status != ExecuteStatusWaiting {
			continue
		}

		isFirstFrame := i%2 == 0

		segment := dataBus.SegmentScript.GetSegments()[i/2]

		if isFirstFrame {
			frames[i] = &projpb.KeyFrames_Frame{
				Status: ExecuteStatusRunning,
				Refs: []string{
					segment.StartFrame,
					helper.SliceElement[string](
						dataBus.Commodity.GetRefs(),
						0, false),
				},
				Prompt: fmt.Sprintf("%s\n%s\n%s\n\n%s\n\n%s\n\n%s\n\n%s\n\n%s",
					segment.GetIntention(),
					segment.GetCoreAction(),
					segment.GetElementTransformation(),
					segment.GetVisualChange(),
					segment.GetDescription(),
					segment.GetSceneStyle(),
					strings.Join(segment.GetTypedTags().Tags(), " "),
					firstFramePrompt.Content,
					//helper.Select(isFirstFrame,
					//	,
					//	"依据以上分镜描述， 结合【图1的整体风格和场景】+【图2中的商品】,生成一张视频尾帧图片",
					//),
				),
				Category: "first_frame",
			}

		} else {
			lastFrame := frames[i-1]
			if lastFrame.Url != "" {
				frames[i] = &projpb.KeyFrames_Frame{
					Status: ExecuteStatusCompleted,
					Url:    lastFrame.Url,
					Refs: []string{
						lastFrame.Url,
						helper.SliceElement[string](
							dataBus.Commodity.GetRefs(),
							0, false),
					},
					Prompt: fmt.Sprintf("%s\n%s\n\n%s\n\n%s\n\n%s\n\n%s\n\n%s",
						segment.GetCoreAction(),
						segment.GetElementTransformation(),
						segment.GetVisualChange(),
						segment.GetDescription(),
						segment.GetSceneStyle(),
						strings.Join(segment.GetTypedTags().Tags(), " "),
						lastFramePrompt.Content,
						//"依据以上分镜描述， 结合【图1的整体风格和场景】+【图2中的商品】,生成一张视频首帧图片",
						//helper.Select(isFirstFrame,
						//	,
						//"依据以上分镜描述， 结合【图1的整体风格和场景】,推演视频尾帧图片。要符合物理规律, 图片中不要包含文字",
						//),
					),
				}
			}

		}
	}

	_, err = t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId, mgz.Op().
		Set(fmt.Sprintf("jobs.%d.dataBus.keyFrames", jobState.Index), &projpb.KeyFrames{
			Frames: frames,
		}))
	if err != nil {
		return nil, err
	}

	logger.Debugw("settings", wfState.DataBus.GetSettings())

	// Waved AI 版
	//wg.WaitGroupIndexed(ctx, frames, func(ctx context.Context, x *projpb.KeyFrames_Frame, index int) error {
	//	if x.Status != ExecuteStatusRunning {
	//		return nil
	//	}
	//
	//	if x.TaskId != "" {
	//		result, err := t.data.Wavespeed.GetResult(ctx, x.GetTaskId())
	//		if err != nil {
	//			return err
	//		}
	//
	//		logger.Debugw("keyFrames generation job task check result", result)
	//
	//		if len(result.Data.Outputs) > 0 {
	//
	//			x.Url = result.Data.Outputs[0]
	//			x.Status = ExecuteStatusCompleted
	//
	//			_, err = t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId, mgz.Op().
	//				Set(fmt.Sprintf("jobs.%d.dataBus.keyFrames.frames.%d", jobState.Index, index), x))
	//
	//			if err != nil {
	//				return err
	//			}
	//
	//			return nil
	//		}
	//
	//		if result.Data.Status == "failed" {
	//			x.Url = ""
	//			x.Status = ExecuteStatusRunning
	//			x.TaskId = ""
	//
	//			_, err = t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId, mgz.Op().
	//				Set(fmt.Sprintf("jobs.%d.dataBus.keyFrames.frames.%d", jobState.Index, index), x))
	//
	//			if err != nil {
	//				return err
	//			}
	//
	//			return nil
	//		}
	//
	//		return nil
	//	}
	//	aspectRatio := helper.OrString(wfState.GetDataBus().GetSettings().GetAspectRatio(), "9:16")
	//	res, err := t.data.Wavespeed.Gemini3ProImage(ctx, wavespeed.Gemini3ProImageRequest{
	//		Prompt:      x.Prompt,
	//		Images:      x.Refs,
	//		AspectRatio: aspectRatio,
	//		Resolution:  "1k",
	//		//OutputFormat: "",
	//		//EnableSyncMode: true,
	//		//EnableBase64Output: false,
	//	})
	//	if err != nil {
	//		logger.Errorw("Gemini3ProImage err", err)
	//		return err
	//	}
	//
	//	logger.Debugw("keyFrames generation job task Gemini3ProImage", res)
	//
	//	x.TaskId = res.Data.Id
	//	x.Status = ExecuteStatusRunning
	//
	//	_, err = t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId,
	//		mgz.Op().
	//			Set(fmt.Sprintf("jobs.%d.dataBus.keyFrames.frames.%d", jobState.Index, index), x))
	//
	//	if err != nil {
	//		return err
	//	}
	//
	//	return nil
	//})

	wg.WaitGroupIndexed(ctx, frames, func(ctx context.Context, x *projpb.KeyFrames_Frame, index int) error {
		if x.Status != ExecuteStatusRunning {
			return nil
		}

		//refParts, err := gemini.NewImageParts(x.Refs)
		//if err != nil {
		//	return err
		//}

		aspectRatio := helper.OrString(wfState.GetDataBus().GetSettings().GetAspectRatio(), "9:16")
		blob, err := t.data.GenaiFactory.Get().GenerateImage(ctx, gemini.GenerateImageRequest{
			//Parts:       refParts,
			Images:      x.Refs,
			Prompt:      x.Prompt,
			AspectRatio: aspectRatio,
		})
		if err != nil {
			logger.Errorw("GenerateImage err", err)

			x.Error = err.Error()
			x.Status = ExecuteStatusFailed

			_, err = t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId,
				mgz.Op().
					Set(fmt.Sprintf("jobs.%d.dataBus.keyFrames.frames.%d", jobState.Index, index), x),
			)

			return err
		}

		tmpUrl, err := t.data.TOS.PutImageBytes(ctx, blob)
		if err != nil {
			return err
		}

		x.Url = tmpUrl
		x.Status = ExecuteStatusCompleted
		x.AspectRatio = aspectRatio

		_, err = t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId,
			mgz.Op().
				Set(fmt.Sprintf("jobs.%d.dataBus.keyFrames.frames.%d", jobState.Index, index), x),
		)

		if err != nil {
			return err
		}

		//go func() {
		//	defer helper.DeferFunc()
		//
		//	// 审计
		//	promptKey := "segment_first_frame_audit"
		//	if x.Category != "first_frame" {
		//		promptKey = "segment_last_frame_audit"
		//	}
		//
		//	prompt, err := t.data.Mongo.Settings.GetPrompt(ctx, promptKey)
		//	if err != nil {
		//		fmt.Println("GetPrompt err", err)
		//		return
		//	}
		//
		//	refParts = append(refParts,
		//		genai.NewPartFromBytes(blob, "image/jpeg"),
		//		gemini.NewTextPart(prompt.Content),
		//	)
		//
		//	auditResJson, err := t.data.GenaiFactory.Get().GenerateContent(ctx, gemini.GenerateContentRequest{
		//		Config: &genai.GenerateContentConfig{
		//			ResponseMIMEType: "application/json",
		//			ResponseSchema: &genai.Schema{
		//				Required: []string{"pass", "reason"},
		//				Type:     genai.TypeObject,
		//				Properties: map[string]*genai.Schema{
		//					"pass": {Type: genai.TypeBoolean, Description: "是否符合要求"},
		//					"desc": {Type: genai.TypeString, Description: "审计结果概述"},
		//				},
		//			},
		//		},
		//		Parts: refParts,
		//	})
		//
		//	fmt.Println("audit res", auditResJson)
		//
		//	var auditRes projpb.Review
		//	err = json.Unmarshal([]byte(auditResJson), &auditRes)
		//	//if auditRes["pass"] == false {
		//	//
		//	//}
		//	_, err = t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId,
		//		mgz.Op().
		//			Set(fmt.Sprintf("jobs.%d.dataBus.keyFrames.frames.%d.review", jobState.Index, index), &auditRes),
		//	)
		//}()
		return nil
	})

	return nil, nil
}
