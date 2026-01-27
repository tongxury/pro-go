package biz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	projpb "store/api/proj"
	"store/app/proj-pro/internal/data"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/third/gemini"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/genai"
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

	// Audit
	promptAudit, err := t.data.Mongo.Settings.GetPrompt(ctx, "sora2_image_audit")
	if err != nil {
		logger.Warnw("GetPrompt sora2_image_audit fail, use default", "err", err)
		promptAudit = &projpb.Prompt{
			Content: `你是一个专业的视觉审核员。请判断生成的图片是否符合脚本描述和参考图风格。
1. 检查图片内容是否与脚本描述一致。
2. 检查图片风格是否与参考图保持连贯。
3. 检查图片中是否包含不符合要求的文字。
请返回 JSON 格式：{"pass": true/false, "reason": "通过或不通过的原因"}`,
		}
	}

	auditRes, err := t.data.GenaiFactory.Get().GenerateContent(ctx, gemini.GenerateContentRequest{
		Config: &genai.GenerateContentConfig{
			ResponseMIMEType: "application/json",
			ResponseSchema: &genai.Schema{
				Type:     genai.TypeObject,
				Required: []string{"pass", "reason"},
				Properties: map[string]*genai.Schema{
					"pass": {
						Type:        genai.TypeBoolean,
						Description: "是否符合要求",
					},
					"reason": {
						Type:        genai.TypeString,
						Description: "原因",
					},
				},
			},
		},
		Parts: []*genai.Part{
			{Text: promptAudit.Content},
			{Text: fmt.Sprintf("这是你要审核的参考信息：\n脚本：%s\n参考图列表：%v", script, keyFrame.Refs)},
			{FileData: &genai.FileData{MIMEType: "image/png", FileURI: imageUrl}},
		},
	})

	if err != nil {
		logger.Errorw("audit image fail", "err", err)
		return nil, err
	}

	var auditResult struct {
		Pass   bool   `json:"pass"`
		Reason string `json:"reason"`
	}
	err = json.Unmarshal([]byte(auditRes), &auditResult)
	if err != nil {
		logger.Errorw("unmarshal audit result fail", "err", err, "res", auditRes)
		return nil, err
	}

	if !auditResult.Pass {
		keyFrame.Status = ExecuteStatusFailed
		keyFrame.Error = auditResult.Reason
	} else {
		keyFrame.Status = ExecuteStatusCompleted
	}

	_, err = t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId, mgz.Op().
		Set(fmt.Sprintf("jobs.%d.dataBus.keyFrames", jobState.Index), &projpb.KeyFrames{
			Frames: []*projpb.KeyFrames_Frame{keyFrame},
		}))
	if err != nil {
		logger.Errorw("update segment keyframe fail", "err", err)
		return nil, err
	}

	if !auditResult.Pass {
		return &ExecuteResult{
			Status: ExecuteStatusFailed,
		}, nil
	}

	return &ExecuteResult{
		Status: ExecuteStatusCompleted,
	}, nil
}
