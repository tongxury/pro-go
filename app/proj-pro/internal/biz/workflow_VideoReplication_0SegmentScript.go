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

type VideoReplication_SegmentScriptJob struct {
	data *data.Data
}

func (t VideoReplication_SegmentScriptJob) Initialize(ctx context.Context, options Options) error {
	return nil
}

func (t VideoReplication_SegmentScriptJob) GetName() string {
	return "segmentScriptJob"
}

func (t VideoReplication_SegmentScriptJob) Execute(ctx context.Context, jobState *projpb.Job, wfState *projpb.Workflow) (status *ExecuteResult, err error) {

	dataBus := GetDataBus(wfState)

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "VideoReplication_SegmentScriptJob.Execute",
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

	if len(dataBus.SegmentScript.GetScript()) > 0 {
		return &ExecuteResult{
			Status: ExecuteStatusCompleted,
		}, nil
	}

	logger.Debugw("start", "")

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
							//"style",
							"timeStart",
							"timeEnd",
							"subtitle",
							"intention",
							"coreAction",
							"elementTransformation",
							"visualChange",
							"contentStyle",
							"sceneStyle",
							"description",
							"shootingStyle",
							//"typedTags"
						},
						Properties: map[string]*genai.Schema{
							//"style": {
							//	Type:        genai.TypeString,
							//	Description: "当前分段的拍摄脚本描述",
							//},
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
								Description: "元素变化",
							},
							"visualChange": {
								Type:        genai.TypeString,
								Description: "视觉变化",
							},
							"contentStyle": {
								Type:        genai.TypeString,
								Description: "当前分段的内容风格描述",
							},
							"sceneStyle": {
								Type:        genai.TypeString,
								Description: "当前分段的场景设计描述",
							},
							"description": {
								Type:        genai.TypeString,
								Description: "画面描述",
							},
							"shootingStyle": {
								Type:        genai.TypeString,
								Description: "拍摄手法",
							},
							"keyFrames": {
								Type:     genai.TypeObject,
								Required: []string{"storyBeginning", "climax", "storyEnding"},
								Properties: map[string]*genai.Schema{
									"storyBeginning": {
										Type:     genai.TypeObject,
										Required: []string{"timestamp", "desc"},
										Properties: map[string]*genai.Schema{
											"timestamp": {
												Type:        genai.TypeNumber,
												Description: "关键帧的时间戳，精确到毫秒。例如 1.232",
											},
											"desc": {
												Type:        genai.TypeString,
												Description: "一句话描述",
											},
										},
									},
									"climax": {
										Type:     genai.TypeObject,
										Required: []string{"timestamp", "desc"},

										Properties: map[string]*genai.Schema{
											"timestamp": {
												Type:        genai.TypeNumber,
												Description: "关键帧的时间戳，精确到毫秒。例如 1.232",
											},
											"desc": {
												Type:        genai.TypeString,
												Description: "一句话描述",
											},
										},
									},
									"storyEnding": {
										Type:     genai.TypeObject,
										Required: []string{"timestamp", "desc"},

										Properties: map[string]*genai.Schema{
											"timestamp": {
												Type:        genai.TypeNumber,
												Description: "关键帧的时间戳，精确到毫秒。例如 1.232",
											},
											"desc": {
												Type:        genai.TypeString,
												Description: "一句话描述",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	//videoBytes, err := videoz.GetBytes(segment.Root.Url)
	//if err != nil {
	//	logger.Errorw("videoz.GetBytes", "err", err)
	//	return nil, err
	//}

	seg, err := videoz.GetSegmentByUrl(ctx, segment.Root.Url, segment.TimeStart, segment.TimeEnd)
	if err != nil {
		return nil, err
	}

	prompt, err := t.data.Mongo.Settings.GetPrompt(ctx, "segment_script_replication")
	if err != nil {
		return nil, err
	}

	p := fmt.Sprintf(`
%s
===
*新商品信息*: %s
`,
		prompt.Content,
		//dataBus.GetSegment().GetScript(),
		conv.S2J(dataBus.GetCommodity()),
	)

	logger.Debugw("p", p)

	text, err := t.data.GenaiFactory.Get().AnalyzeVideo(ctx, gemini.AnalyzeVideoRequest{
		VideoBytes: seg.Content,
		Prompt:     p,
		Config:     config,
	})

	//	text, err := t.data.GenaiFactory.Get().GenerateContentV2(ctx, gemini.GenerateContentRequestV2{
	//		//ImageUrls: []string{dataBus.GetCommodity().GetMedias()[0].Url},
	//		Prompt: fmt.Sprintf(`
	//%s
	//===
	//要复刻的爆款脚本: %s
	//新商品信息: %s
	//`,
	//			prompt.Content,
	//			dataBus.GetSegment().GetScript(),
	//			conv.S2J(dataBus.GetCommodity()),
	//		),
	//		Config: config,
	//	})
	if err != nil {
		logger.Errorw("gen content error", err)
		return nil, err
	}

	var segments projpb.Resource
	err = protojson.Unmarshal([]byte(text), &segments)
	if err != nil {
		return nil, err
	}

	// 补充图片
	for i := range segments.Segments {

		frame, err := videoz.GetFrame(seg.Content, segments.Segments[i].TimeStart+0.1)
		if err != nil {
			return nil, err
		}

		segments.Segments[i].StartFrame, err = t.data.TOS.PutImageBytes(ctx, frame)
		if err != nil {
			return nil, err
		}

		endFrameBytes, err := videoz.GetFrame(seg.Content, segments.Segments[i].TimeEnd-0.1)
		if err != nil {
			return nil, err
		}

		segments.Segments[i].EndFrame, err = t.data.TOS.PutImageBytes(ctx, endFrameBytes)
		if err != nil {
			return nil, err
		}
	}

	_, err = t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId, mgz.Op().
		Set(fmt.Sprintf("jobs.%d.dataBus.segmentScript", jobState.Index), &projpb.SegmentScript{
			Script:   text,
			Segments: segments.Segments,
		}))
	if err != nil {
		logger.Errorw("update segment script fail", "err", err)
		return nil, err
	}

	return &ExecuteResult{
		Status: ExecuteStatusCompleted,
	}, nil
}
