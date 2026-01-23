package service

import (
	"context"
	"encoding/json"
	"fmt"
	projpb "store/api/proj"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/helper/videoz"
	"store/pkg/sdk/helper/wg"
	"store/pkg/sdk/third/bytedance/arkr"
	"store/pkg/sdk/third/gemini"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/genai"
	"google.golang.org/protobuf/encoding/protojson"
)

func (t ProjService) ExtractSegment() {

	ctx := context.Background()

	hits, err := t.data.Mongo.TemplateSegment.List(ctx, bson.M{"status": "processing_created"})
	if err != nil {
		log.Errorw("Search err", err, "index", "items")
		return
	}

	if len(hits) == 0 {
		return
	}

	settings, err := t.data.Mongo.Settings.FindOne(ctx, bson.M{})
	if err != nil {
		log.Errorw("Settings.FindOne err", err)
		return
	}

	data := helper.Mapping(hits, func(x *projpb.ResourceSegment) *extractSegmentJobData {

		return &extractSegmentJobData{
			Item:     x,
			Settings: settings,
		}
	})

	wg.WaitGroup(ctx, data, t.extractSegment)
}

type extractSegmentJobData struct {
	Item     *projpb.ResourceSegment
	Settings *projpb.AppSettings
}

func (t ProjService) extractSegment(ctx context.Context, data *extractSegmentJobData) error {

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "extractSegment",
		"item", data.Item.XId,
		"prompt",
		helper.SubString(data.Settings.GetVideoHighlight().GetContent(), 0, 100),
	))

	resource, err := t.extractSegmentByGemini(ctx, data)
	if err != nil {
		logger.Errorw("ExtractSegmentByGemini err", err)

		return err
		//resource, err = t.extractSegmentBySeed(ctx, data)
		//if err != nil {
		//	logger.Errorw("ExtractSegmentBySeed err", err)
		//	return err
		//}
	}
	//
	//resource, err := t.extractSegmentBySeed(ctx, data)
	//if err != nil {
	//	logger.Errorw("ExtractSegmentBySeed err", err)
	//	return err
	//}

	metadata, err := videoz.GetMetadata(data.Item.Root.Url)
	if err != nil {
		logger.Errorw("GetMetadata err", err, "data", data)

		return err
	}

	// 校验
	if len(resource.Segments) != 3 {
		log.Errorw("invalid segment size", len(resource.Segments))
		return fmt.Errorf("invalid segment size")
	}

	for _, x := range resource.Segments {

		if !t.verifySegment(x, metadata) {
			log.Errorw("invalid segment err", err, "data", data)
			return fmt.Errorf("invalid segment")
		}
	}

	for _, x := range resource.Segments {
		x.Tags = append(x.Tags, x.TypedTags.FocusOn...)
		x.Tags = append(x.Tags, x.TypedTags.Text...)
		x.Tags = append(x.Tags, x.TypedTags.Person...)
		x.Tags = append(x.Tags, x.TypedTags.Picture...)
		x.Tags = append(x.Tags, x.TypedTags.Scene...)
		x.Tags = append(x.Tags, x.TypedTags.ShootingStyle...)
		x.Tags = append(x.Tags, x.TypedTags.Emotion...)
	}

	//item.Analysis2 = tmpItem.Analysis
	data.Item.Root.Commodity = resource.Commodity

	for i, x := range resource.Segments {
		x.Root = data.Item.Root
		x.Index = int64(i)
		x.CreatedAt = time.Now().Unix()
		x.Status = "processing_segmented"
	}

	_, err = t.data.Mongo.TemplateSegment.Delete(ctx, bson.M{
		"root.url": data.Item.Root.Url,
	})

	if err != nil {
		logger.Errorw("DeleteByRequest err", err)
		return err
	}

	_, err = t.data.Mongo.TemplateSegment.InsertMany(ctx, resource.Segments...)

	if err != nil {
		logger.Errorw("InsertMany err", err)
		return err
	}

	return nil

}

func (t ProjService) extractSegmentByGemini(ctx context.Context, data *extractSegmentJobData) (*projpb.Resource, error) {

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "analysisSegmentByGemini",
		"item", data.Item.XId,
		"prompt",
		helper.SubString(data.Settings.GetVideoHighlight().GetContent(), 0, 100),
	))

	logger.Debugw("start", "")

	genaiClient := t.data.GenaiFactory.Get()

	//genaiUrl, err := genaiClient.UploadFile(ctx, data.Item.Root.Url, "video/mp4")
	//if err != nil {
	//	return nil, err
	//}

	part, err := gemini.NewVideoPart(data.Item.Root.Url)
	if err != nil {
		return nil, err
	}

	parts := []*genai.Part{
		//{
		//	FileData: &genai.FileData{
		//		MIMEType: "video/mp4",
		//		FileURI:  genaiUrl,
		//	},
		//
		//},
		part,
		{
			Text: data.Settings.GetVideoHighlight().GetContent(),
		},
	}

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type:     genai.TypeObject,
			Required: []string{"segments", "commodity"},
			Properties: map[string]*genai.Schema{
				"commodity": {
					Type:     genai.TypeObject,
					Required: []string{"tags", "name"},
					Properties: map[string]*genai.Schema{
						"name": {
							Type:        genai.TypeString,
							Description: projpb.PromptCommodityName,
						},
						"tags": {
							Type: genai.TypeArray,
							Items: &genai.Schema{
								Type: genai.TypeString,
							},
							Description: projpb.PromptCommodityTags,
						},
					},
				},
				"segments": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeObject,
						Required: []string{
							//"highlightFrames",
							"style", "contentStyle", "sceneStyle", "timeStart", "timeEnd", "description", "shootingStyle", "typedTags", "subtitle"},
						Properties: map[string]*genai.Schema{
							"style": {
								Type:        genai.TypeString,
								Description: "高光片段的拍摄脚本描述",
							},
							"contentStyle": {
								Type:        genai.TypeString,
								Description: "高光片段的内容风格描述",
							},
							"sceneStyle": {
								Type:        genai.TypeString,
								Description: "高光片段的场景设计描述",
							},
							"timeStart": {
								Type:        genai.TypeNumber,
								Description: "高光片段的脚本的开始时间戳(秒, 示例：1.22)",
							},
							"timeEnd": {
								Type:        genai.TypeNumber,
								Description: "高光片段的脚本的结束时间戳(秒, 示例：1.23)",
							},
							//"highlightFrames": {
							//	Type:        genai.TypeArray,
							//	Description: `高光片段的关键画面帧`,
							//	Items: &genai.Schema{
							//		Type:     genai.TypeObject,
							//		Required: []string{"timestamp", "desc"},
							//		Properties: map[string]*genai.Schema{
							//			"timestamp": {
							//				Type:        genai.TypeNumber,
							//				Description: "当前帧对应的时间戳, 示例格式: 1.222",
							//			},
							//			"desc": {
							//				Type:        genai.TypeString,
							//				Description: "一句换描述这一帧",
							//			},
							//		},
							//	},
							//},
							"description": {
								Type:        genai.TypeString,
								Description: "高光片段画面描述",
							},
							//"focusOn": {
							//	Type: genai.TypeArray,
							//	Items: &genai.Schema{
							//		Type:        genai.TypeString,
							//		Description: "主要阐述",
							//	},
							//},
							"shootingStyle": {
								Type: genai.TypeString,
								Description: `
请详细给出拍摄手法，
包括：
1. 镜头选择：具体说明使用什么焦距的镜头，如广角、标准、长焦等
2. 拍摄角度：详细描述拍摄角度，如正面、侧面、俯视、仰视等
3. 镜头运动：说明镜头的运动方式，如推、拉、摇、移、跟拍等
4. 构图方式：描述画面的构图原则和视觉重点
5. 演员指导：
	- 眼神表达：具体说明眼神的方向、情感、专注度
	- 表情管理：描述面部表情的细节，如微笑程度、情绪表达
	- 语气语调：说明说话的语气、语调、语速、停顿等
	- 肢体语言：描述手势、姿态、动作幅度等
	- 状态调整：说明演员的整体状态和情绪准备
6. 表演细节：
	- 情感层次：说明表演的情感层次和变化
	- 互动方式：描述与镜头、道具、环境的互动
	- 自然度：确保表演的自然和真实感
7. 技术细节：
	- 对焦方式：说明手动对焦还是自动对焦
	- 曝光控制：描述光圈、快门、ISO的设置
	- 稳定方式：说明手持、三脚架、稳定器等使用方式",
8. 背景音乐:
	- 背景音乐风格
`,
							},
							"typedTags": {
								Type:        genai.TypeObject,
								Description: "标签",
								Required:    []string{"focusOn", "picture", "shootingStyle", "scene", "action", "person", "text", "emotion"},

								Properties: map[string]*genai.Schema{
									"focusOn": {
										Type:        genai.TypeArray,
										Description: "主要阐述类标签",
										Items: &genai.Schema{
											Type: genai.TypeString,
										},
									},
									"picture": {
										Type:        genai.TypeArray,
										Description: "画面类型类标签",
										Items: &genai.Schema{
											Type: genai.TypeString,
										},
									},
									"shootingStyle": {
										Type:        genai.TypeArray,
										Description: "拍摄手法类标签",
										Items: &genai.Schema{
											Type: genai.TypeString,
										},
									},
									"scene": {
										Type:        genai.TypeArray,
										Description: "场景类标签",
										Items: &genai.Schema{
											Type: genai.TypeString,
										},
									},
									"action": {
										Type:        genai.TypeArray,
										Description: "动作类标签",
										Items: &genai.Schema{
											Type: genai.TypeString,
										},
									},
									"person": {
										Type:        genai.TypeArray,
										Description: "人物类标签",
										Items: &genai.Schema{
											Type: genai.TypeString,
										},
									},
									"text": {
										Type:        genai.TypeArray,
										Description: "文案类标签（注意不是文案本身）",
										Items: &genai.Schema{
											Type: genai.TypeString,
										},
									},
									"emotion": {
										Type:        genai.TypeArray,
										Description: "情绪类标签",
										Items: &genai.Schema{
											Type: genai.TypeString,
										},
									},
								},
							},
							"subtitle": {
								Type:        genai.TypeString,
								Description: "高光片段的口播文案",
							},
						},
					},
				},
			},
		},
	}

	logger.Debugw("GenerateContent start", "")

	text, err := genaiClient.GenerateContent(ctx, gemini.GenerateContentRequest{
		//Model:  models[mathz.RandNumber(0, len(models)-1)],
		Config: config,
		Parts:  parts,
	})

	if err != nil {
		logger.Errorw("GenerateContent err", err)
		return nil, err
	}

	logger.Debugw("GenerateContent done", text)

	var tmpItem projpb.Resource

	err = protojson.Unmarshal([]byte(text), &tmpItem)
	if err != nil {
		logger.Errorw("protojson.Unmarshal err", err, "text", text)
		return nil, err
	}

	return &tmpItem, nil
}

func (t ProjService) extractSegmentBySeed(ctx context.Context, data *extractSegmentJobData) (*projpb.Resource, error) {

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "extractSegmentBySeed",
		"item", data.Item.XId,
		"prompt",
		helper.SubString(data.Settings.GetVideoHighlight().GetContent(), 0, 100),
	))

	logger.Debugw("start", "")

	var messages []*model.ChatCompletionMessage

	messages = append(messages,
		&model.ChatCompletionMessage{
			Role: model.ChatMessageRoleSystem,
			Content: &model.ChatCompletionMessageContent{
				StringValue: volcengine.String(data.Settings.GetVideoTemplate().GetContent()),
			},
		},
		&model.ChatCompletionMessage{
			Role: model.ChatMessageRoleSystem,
			Content: &model.ChatCompletionMessageContent{
				ListValue: []*model.ChatCompletionMessageContentPart{
					{
						Type: "video_url",
						VideoURL: &model.ChatMessageVideoURL{
							URL: data.Item.Root.Url,
						},
					},
				},
			},
		},
	)

	type Result struct {
		Commodity struct {
			Name string   `json:"name"  jsonschema_description:"commodity.name" required:"true"`
			Tags []string `json:"tags"  jsonschema_description:"commodity.tags"`
		} `json:"commodity" jsonschema_description:"commodity"  required:"true"`
		Segments []struct {
			Style         string  `json:"style"  jsonschema_description:"segment.style"`
			ContentStyle  string  `json:"contentStyle"  jsonschema_description:"segment.contentStyle"`
			SceneStyle    string  `json:"sceneStyle"  jsonschema_description:"segment.sceneStyle"`
			Description   string  `json:"description"  jsonschema_description:"segment.description"`
			ShootingStyle string  `json:"shootingStyle"  jsonschema_description:"segment.shootingStyle"`
			TimeStart     float64 `json:"timeStart"  jsonschema_description:"segment.timeStart"`
			TimeEnd       float64 `json:"timeEnd"  jsonschema_description:"segment.timeEnd"`
			Subtitle      string  `json:"subtitle"  jsonschema_description:"segment.subtitle"`
			TypedTags     struct {
				FocusOn       []string `json:"focusOn"  jsonschema_description:"segment.typedTags.focusOn"`
				Text          []string `json:"text"  jsonschema_description:"segment.typedTags.text"`
				Person        []string `json:"person"  jsonschema_description:"segment.typedTags.person"`
				Action        []string `json:"action"  jsonschema_description:"segment.typedTags.action"`
				Picture       []string `json:"picture"  jsonschema_description:"segment.typedTags.picture"`
				Scene         []string `json:"scene"  jsonschema_description:"segment.typedTags.scene"`
				ShootingStyle []string `json:"shootingStyle"  jsonschema_description:"segment.typedTags.shootingStyle"`
				Emotion       []string `json:"emotion"  jsonschema_description:"segment.typedTags.emotion"`
			} `json:"typedTags"  jsonschema_description:"segment.typedTags"`
		} `json:"segments" jsonschema_description:"segments"  required:"true"`
	}

	req := model.CreateChatCompletionRequest{
		Model:    "doubao-seed-1-6-250615",
		Messages: messages,
		ResponseFormat: &model.ResponseFormat{
			Type: model.ResponseFormatJSONSchema,
			JSONSchema: &model.ResponseFormatJSONSchemaJSONSchemaParam{
				Name:        "resource",
				Description: "视频分析结果",
				Schema:      arkr.GenerateSchema[Result](),
				Strict:      true,
			},
		},
	}

	resp, err := t.data.Arkr.C().CreateChatCompletion(ctx, req)
	if err != nil {
		logger.Errorw("CreateChatCompletion err", err, "req", req)
		return nil, err
	}

	var res *Result
	err = json.Unmarshal([]byte(*resp.Choices[0].Message.Content.StringValue), &res)
	if err != nil {
		return nil, err
	}

	var segments []*projpb.ResourceSegment
	for _, x := range res.Segments {

		segments = append(segments, &projpb.ResourceSegment{
			XId:           "",
			Root:          nil,
			Style:         x.Style,
			ContentStyle:  x.ContentStyle,
			SceneStyle:    x.SceneStyle,
			Description:   x.Description,
			ShootingStyle: x.ShootingStyle,
			TimeStart:     x.TimeStart,
			TimeEnd:       x.TimeEnd,
			Subtitle:      x.Subtitle,
			TypedTags: &projpb.TypedTags{
				FocusOn:       x.TypedTags.FocusOn,
				Text:          x.TypedTags.Text,
				Person:        x.TypedTags.Person,
				Action:        x.TypedTags.Action,
				Picture:       x.TypedTags.Picture,
				Scene:         x.TypedTags.Scene,
				ShootingStyle: x.TypedTags.ShootingStyle,
				Emotion:       x.TypedTags.ShootingStyle,
			},
		})
	}

	return &projpb.Resource{
		Commodity: &projpb.Commodity{
			Name: res.Commodity.Name,
			Tags: res.Commodity.Tags,
		},
		Segments: segments,
	}, nil
}

func (t ProjService) verifySegment(x *projpb.ResourceSegment, metadata *videoz.Metadata) bool {

	if x.TimeEnd > metadata.Duration {
		return false
	}

	if x.TimeStart > metadata.Duration {
		return false
	}

	if x.TimeEnd-x.TimeStart < 5 {
		return false
	}

	if x.TimeEnd-x.TimeStart > 8 {
		return false
	}

	if x.TimeEnd < 1 {
		return false
	}

	for _, xx := range x.TypedTags.Text {
		if utf8.RuneCountInString(xx) > 8 {
			return false
		}
	}

	if t.englishRatioOver(x.Description, 0.5) {
		return false
	}

	if t.englishRatioOver(x.Root.GetCommodity().GetName(), 0.5) {
		return false
	}

	return true
}

func (t ProjService) englishRatioOver(s string, threshold float64) bool {
	if threshold < 0 {
		threshold = 0
	} else if threshold > 1 {
		threshold = 1
	}

	var englishCount, totalCount float64

	for _, r := range s {
		if unicode.IsSpace(r) || unicode.IsPunct(r) {
			continue
		}
		totalCount++
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
			englishCount++
		}
	}

	if totalCount == 0 {
		return false
	}
	return englishCount/totalCount > threshold
}
