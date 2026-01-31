package biz

import (
	"context"
	"fmt"
	projpb "store/api/proj"
	"store/app/proj-pro/internal/data"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/helper/bytez"
	"store/pkg/sdk/helper/wg"
	"store/pkg/sdk/third/gemini"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/genai"
)

type KeyFramesGenerationJob struct {
	data *data.Data
}

func (t KeyFramesGenerationJob) Initialize(ctx context.Context, options Options) error {
	return nil
}

func (t KeyFramesGenerationJob) GetName() string {
	return "keyFramesGenerationJob"
}

func (t KeyFramesGenerationJob) Execute(ctx context.Context, jobState *projpb.Job, wfState *projpb.Workflow) (status *ExecuteResult, err error) {

	dataBus := GetDataBus(wfState)

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "KeyFramesGenerationJob.Execute",
		"workflowId ", wfState.XId,
		"jobState.Name ", jobState.Name,
		"jobState.Index ", jobState.Index,
		"jobState.Status ", jobState.Status,
	))

	logger.Debug("keyFrames generation job start")

	doings := helper.Filter(dataBus.KeyFrames.GetFrames(), func(x *projpb.KeyFrames_Frame) bool {
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
						helper.Mapping(dataBus.Commodity.GetMedias(),
							func(x *projpb.Media) string {
								return x.Url
							}),
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
					Status: ExecuteStatusRunning,
					Refs: []string{
						lastFrame.Url,
						helper.SliceElement[string](
							helper.Mapping(dataBus.Commodity.GetMedias(),
								func(x *projpb.Media) string {
									return x.Url
								}),
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

	//// Waved AI 版
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
	//
	//	res, err := t.data.Wavespeed.Gemini3ProImage(ctx, wavespeed.Gemini3ProImageRequest{
	//		Prompt:       x.Prompt,
	//		Images:       x.Refs,
	//		AspectRatio:  "9:16",
	//		Resolution:   "1k",
	//		OutputFormat: "",
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

		// 将整体实现改成 用下面的方法
		// t.data.GenaiFactory.Get().C().Models.EditImage(ctx)

		aspectRatio := helper.OrString(wfState.GetDataBus().GetSettings().GetAspectRatio(), "9:16")

		client := t.data.GenaiFactory.Get().C()
		model := gemini.ModelImagen3Capability001

		var refImages []genai.ReferenceImage
		// Multi-Raw Test: use all images but as NewRawReferenceImage to check slice handling
		for i, ref := range x.Refs {
			imgBytes, err := bytez.ReadFileBytes(ref)
			if err != nil {
				return err
			}

			if len(imgBytes) == 0 {
				return fmt.Errorf("reference image %d (%s) is empty", i, ref)
			}

			logger.Debugw("GenAI Multi-Raw Test", "index", i, "url", ref, "size", len(imgBytes), "mime", "image/png")

			img := &genai.Image{
				ImageBytes: imgBytes,
				MIMEType:   "image/png",
			}

			// Using RawReferenceImage for all images for this test phase
			refImages = append(refImages, genai.NewRawReferenceImage(img, int32(i+1)))
		}

		// Use an empty config instead of nil to be safe with SDK internal converters.
		resp, err := client.Models.EditImage(ctx, model, x.Prompt, refImages, &genai.EditImageConfig{})

		if err != nil {
			logger.Errorw("EditImage err", err)

			x.Error = err.Error()
			x.Status = ExecuteStatusFailed

			_, _ = t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId,
				mgz.Op().
					Set(fmt.Sprintf("jobs.%d.dataBus.keyFrames.frames.%d", jobState.Index, index), x),
			)

			return err
		}

		if len(resp.GeneratedImages) == 0 {
			logger.Errorw("no image generated")
			return fmt.Errorf("no image generated")
		}

		blob := resp.GeneratedImages[0].Image.ImageBytes

		tmpUrl, err := t.data.TOS.PutImageBytes(ctx, blob)
		if err != nil {
			logger.Errorw("PutImage err", err)
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

		return nil
	})

	return nil, nil
}
