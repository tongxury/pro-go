package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/helper/wg"
	"store/pkg/sdk/third/bytedance/arkr"
	"store/pkg/sdk/third/gemini"
	"time"
	"unicode/utf8"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/genai"
	"google.golang.org/protobuf/encoding/protojson"
)

func (t ProjAdminService) Segment() {
	//t.data.Elastics.Create()
	ctx := context.Background()

	hits, err := t.data.Mongo.Template.List(ctx, bson.M{"status": "created"})
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

	data := helper.Mapping(hits, func(x *projpb.Resource) *segmentJobData {

		return &segmentJobData{
			Item:     x,
			Settings: settings,
		}
	})

	wg.WaitGroup(ctx, data, t.segment)
}

type segmentJobData struct {
	Item     *projpb.Resource
	Settings *projpb.AppSettings
}

func (t ProjAdminService) segment(ctx context.Context, data *segmentJobData) error {
	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "segment",
		"item", data.Item.XId,
		"prompt",
		helper.SubString(data.Settings.GetVideoHighlight().GetContent(), 0, 100),
	))

	template, err := t.segmentByGemini(ctx, data)
	if err != nil {
		logger.Errorw("segmentByGemini err", err)
		//template, err = t.segmentBySeed(ctx, data)
		//if err != nil {
		//	logger.Errorw("segmentBySeed err", err)
		//	return err
		//}
		return err
	}

	//template, err := t.segmentBySeed(ctx, data)
	//if err != nil {
	//	logger.Errorw("segmentBySeed err", err)
	//	return err
	//}

	if !t.verify(template) {
		logger.Errorw("invalid resource")
		return fmt.Errorf("invalid resource")
	}

	//root := proto.Clone(template).(*projpb.Resource)
	root := data.Item
	root.Commodity = template.Commodity
	root.Segments = nil

	for _, x := range template.Segments {
		x.Tags = append(x.Tags, x.TypedTags.FocusOn...)
		x.Tags = append(x.Tags, x.TypedTags.Text...)
		x.Tags = append(x.Tags, x.TypedTags.Person...)
		x.Tags = append(x.Tags, x.TypedTags.Picture...)
		x.Tags = append(x.Tags, x.TypedTags.Scene...)
		x.Tags = append(x.Tags, x.TypedTags.ShootingStyle...)
		x.Tags = append(x.Tags, x.TypedTags.Emotion...)

		x.Status = "processing_segmented"
		x.Category = "templateSegment"
		x.Root = root
		x.CreatedAt = time.Now().Unix()
		//x.TimeStart = x.TimeStart
	}

	//Status: status,
	////Title:     x.Title,
	//	Root: &projpb.Resource{
	//		Url:      x.Url,
	//		CoverUrl: x.CoverUrl,
	//	},
	//
	//		TimeStart: x.TimeStart,
	//		TimeEnd:   x.TimeEnd,
	//
	//		CreatedAt: time.Now().Unix(),

	_, err = t.data.Mongo.TemplateSegment.InsertMany(ctx, template.Segments...)
	if err != nil {
		log.Errorw("InsertMany err", err, "data.Item.XId", data.Item.XId)
		return err
	}

	_, err = t.data.Mongo.Template.UpdateByIDIfExists(ctx,
		data.Item.XId,
		mgz.Op().
			Set("commodity", template.Commodity).
			Set("status", "segmented"),
		//Set("segments", template.Segments),
	)
	if err != nil {
		log.Errorw("Update err", err, "data.Item.XId", data.Item.XId)
		return err
	}

	return nil
}

func (t ProjAdminService) segmentByGemini(ctx context.Context, data *segmentJobData) (*projpb.Resource, error) {

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "segmentByGemini",
		"item", data.Item.XId,
		"prompt",
		helper.SubString(data.Settings.GetVideoTemplate().GetContent(), 0, 100),
	))

	logger.Debugw("start", "")

	genaiClient := t.data.GenaiFactory.Get()

	//genaiUrl, err := genaiClient.UploadFile(ctx, data.Item.Url, "video/mp4")
	//if err != nil {
	//	return nil, err
	//}

	part, err := gemini.NewVideoPart(data.Item.Url)
	if err != nil {
		return nil, err
	}

	parts := []*genai.Part{
		//{
		//	FileData: &genai.FileData{
		//		MIMEType: "video/mp4",
		//		FileURI:  genaiUrl,
		//	},
		//},
		part,
		{
			Text: data.Settings.GetVideoTemplate().GetContent(),
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
						Type:        genai.TypeObject,
						Description: "分镜结果",
						Required: []string{
							//"highlightFrames",
							"style", "contentStyle", "sceneStyle", "timeStart", "timeEnd", "description", "shootingStyle", "typedTags", "subtitle"},
						Properties: map[string]*genai.Schema{
							"style": {
								Type:        genai.TypeString,
								Description: "当前分段的拍摄脚本描述",
							},
							"contentStyle": {
								Type:        genai.TypeString,
								Description: "当前分段的内容风格描述",
							},
							"sceneStyle": {
								Type:        genai.TypeString,
								Description: "当前分段的场景设计描述",
							},
							"timeStart": {
								Type:        genai.TypeNumber,
								Description: "当前分段脚本的开始时间戳(秒, 示例：1.22)",
							},
							"timeEnd": {
								Type:        genai.TypeNumber,
								Description: "当前分段脚本的结束时间戳(秒, 示例：1.23)",
							},
							//"highlightFrames": {
							//	Type:        genai.TypeArray,
							//	Description: `关键画面帧`,
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
								Description: "画面描述",
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
										Description: "文案类标签",
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
								Description: "当前片段的口播文案",
							},
						},
					},
				},
			},
		},
	}
	text, err := genaiClient.GenerateContent(ctx, gemini.GenerateContentRequest{
		Model:  gemini.DefaultModel,
		Config: config,
		Parts:  parts,
	})

	if err != nil {
		logger.Errorw("GenerateContent err", err)
		return nil, err
	}

	var tmpItem projpb.Resource

	err = protojson.Unmarshal([]byte(text), &tmpItem)
	if err != nil {
		logger.Errorw("protojson.Unmarshal err", err, "text", text)
		return nil, err
	}

	return &tmpItem, nil
}

func (t ProjAdminService) segmentBySeed(ctx context.Context, data *segmentJobData) (*projpb.Resource, error) {

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "segmentBySeed",
		"item", data.Item.XId,
		"prompt",
		helper.SubString(data.Settings.GetVideoTemplate().GetContent(), 0, 100),
	))

	logger.Debugw("start", "")

	var messages []*model.ChatCompletionMessage

	v, err := http.Get(data.Item.Url)
	if err != nil {
		return nil, err
	}

	all, err := io.ReadAll(v.Body)
	if err != nil {
		return nil, err
	}

	//base64Video := base64.Encoding.EncodeToString(all)
	base64Video := base64.StdEncoding.EncodeToString(all)

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
				//StringValue: volcengine.String(
				//	fmt.Sprintf(analysisPrompt2)),
				ListValue: []*model.ChatCompletionMessageContentPart{
					{
						Type: "video_url",
						VideoURL: &model.ChatMessageVideoURL{
							URL: fmt.Sprintf("data:video/mp4;base64,%s", base64Video),
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
		logger.Errorw("CreateChatCompletion err", err, "req", conv.S2J(req))
		return nil, err
	}

	var res *Result
	err = json.Unmarshal([]byte(*resp.Choices[0].Message.Content.StringValue), &res)
	if err != nil {
		return nil, err
	}

	var segments []*projpb.ResourceSegment

	for i, x := range res.Segments {

		segments = append(segments, &projpb.ResourceSegment{
			Style:         x.Style,
			ContentStyle:  x.ContentStyle,
			SceneStyle:    x.SceneStyle,
			Description:   x.Description,
			ShootingStyle: x.ShootingStyle,
			TimeStart:     x.TimeStart,
			TimeEnd:       x.TimeEnd,
			Subtitle:      x.Subtitle,
			TypedTags: &projpb.TypedTags{
				FocusOn: x.TypedTags.FocusOn,
				Text:    x.TypedTags.Text,
				Person:  x.TypedTags.Person,
				Action:  x.TypedTags.Action,
				Picture: x.TypedTags.Picture,
				Scene:   x.TypedTags.Scene,
			},
			Index: int64(i),
		})

		//for i := range x.HighlightFrames {
		//
		//	xx := x.HighlightFrames[i]
		//
		//	fmt.Println(xx.Timestamp, xx.Desc)
		//
		//	frame, err := videoz.GetFrame(videoBytes, xx.Timestamp)
		//	if err != nil {
		//		return err
		//	}
		//
		//	x.HighlightFrames[i].Url, err = t.data.Alioss.UploadBytes(ctx, helper.CreateUUID()+".jpg", frame)
		//	if err != nil {
		//		return err
		//	}
		//
		//}
	}

	return &projpb.Resource{
		Commodity: &projpb.Commodity{
			Name: res.Commodity.Name,
			Tags: res.Commodity.Tags,
		},
		Segments: segments,
	}, nil

}

func (t ProjAdminService) verify(item *projpb.Resource) bool {

	if len(item.Segments) < 5 {
		return false
	}

	if len(item.Segments) > 10 {
		return false
	}

	for _, x := range item.Segments {

		if x.TimeEnd < 1 {
			return false
		}

		if x.TimeEnd-x.TimeStart < 1 {
			return false
		}

		for _, xx := range x.TypedTags.Text {
			if utf8.RuneCountInString(xx) > 8 {
				return false
			}
		}
		//
		//if t.englishRatioOver(x.Description, 0.5) {
		//	return false
		//}

	}

	//if t.englishRatioOver(item.Commodity.GetName(), 0.5) {
	//	return false
	//}

	return true
}
