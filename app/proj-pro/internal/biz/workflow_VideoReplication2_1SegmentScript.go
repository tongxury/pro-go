package biz

import (
	"context"
	"errors"
	"fmt"
	projpb "store/api/proj"
	"store/app/proj-pro/internal/data"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper/videoz"
	"store/pkg/sdk/third/gemini"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/genai"
	"google.golang.org/protobuf/encoding/protojson"
)

type VideoReplication2_SegmentScriptJob struct {
	data *data.Data
}

func (t VideoReplication2_SegmentScriptJob) Initialize(ctx context.Context, options Options) error {
	return nil
}

func (t VideoReplication2_SegmentScriptJob) GetName() string {
	return "segmentScriptJob"
}

func (t VideoReplication2_SegmentScriptJob) Execute(ctx context.Context, jobState *projpb.Job, wfState *projpb.Workflow) (status *ExecuteResult, err error) {

	dataBus := GetDataBus(wfState)

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "VideoReplication2_SegmentScriptJob.Execute",
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

	if dataBus.SegmentScript.GetScript() != "" {
		return &ExecuteResult{
			Status: ExecuteStatusCompleted,
		}, nil
	}

	promptObj, err := t.data.Mongo.Settings.GetPrompt(ctx, "sora2_video_prompt_generation")
	if err != nil {
		return nil, err
	}
	prompt := promptObj.GetContent()

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type:     genai.TypeObject,
			Required: []string{"segments"},
			Properties: map[string]*genai.Schema{
				"segments": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type:        genai.TypeObject,
						Description: "分镜结果",
						Required: []string{
							"timeStart",
							"timeEnd",
							"subtitle",
							"intention",
							"coreAction",
							"elementTransformation",
							"visualChange",
							"contentStyle",
							"sceneStyle",
						},
						Properties: map[string]*genai.Schema{
							"timeStart": {
								Type:        genai.TypeNumber,
								Description: "当前分镜脚本的开始时间戳(秒, 示例：1.22)",
							},
							"timeEnd": {
								Type:        genai.TypeNumber,
								Description: "当前分镜脚本的结束时间戳(秒, 示例：1.23)",
							},
							"subtitle": {
								Type:        genai.TypeString,
								Description: "文案",
							},
							"intention": {
								Type:        genai.TypeString,
								Description: "意图",
							},
							"coreAction": {
								Type:        genai.TypeString,
								Description: "核心动作",
							},
							"elementTransformation": {
								Type:        genai.TypeString,
								Description: "元素変化",
							},
							"visualChange": {
								Type:        genai.TypeString,
								Description: "视觉変化",
							},
							"contentStyle": {
								Type:        genai.TypeString,
								Description: "当前分段的内容风格描述",
							},
							"sceneStyle": {
								Type:        genai.TypeString,
								Description: "当前分段的场景设计描述",
							},
						},
					},
				},
			},
		},
	}

	seg, err := videoz.GetSegmentByUrl(ctx, segment.Root.GetUrl(), segment.TimeStart, segment.TimeEnd)
	if err != nil {
		return nil, err
	}

	text, err := t.data.GenaiFactory.Get().AnalyzeVideo(ctx, gemini.AnalyzeVideoRequest{
		VideoBytes: seg.Content,
		Prompt: fmt.Sprintf(`
%s
===
要复刻的爆款脚本: %s
新商品信息: %s
`,
			prompt,
			dataBus.GetSegment().GetScript(),
			conv.S2J(dataBus.GetCommodity()),
		),
		Config: config,
	})
	if err != nil {
		logger.Errorw("gen content error", err)
		return nil, err
	}

	var scriptData projpb.Resource
	err = protojson.Unmarshal([]byte(text), &scriptData)
	if err != nil {
		logger.Errorw("unmarshal error", err, "text", text)
		return nil, err
	}

	// 补充图片
	for i := range scriptData.Segments {
		frame, err := videoz.GetFrame(seg.Content, scriptData.Segments[i].TimeStart+0.1)
		if err == nil {
			scriptData.Segments[i].StartFrame, _ = t.data.TOS.PutImageBytes(ctx, frame)
		}

		endFrameBytes, err := videoz.GetFrame(seg.Content, scriptData.Segments[i].TimeEnd-0.1)
		if err == nil {
			scriptData.Segments[i].EndFrame, _ = t.data.TOS.PutImageBytes(ctx, endFrameBytes)
		}
	}

	_, err = t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId, mgz.Op().
		Set(fmt.Sprintf("jobs.%d.dataBus.segmentScript", jobState.Index), &projpb.SegmentScript{
			//Script:   text,
			Segments: scriptData.Segments,
		}))
	if err != nil {
		logger.Errorw("update segment script fail", "err", err)
		return nil, err
	}

	return &ExecuteResult{
		Status: ExecuteStatusCompleted,
	}, nil
}
