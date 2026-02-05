package service

import (
	"context"
	"fmt"
	projpb "store/api/proj"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/helper/videoz"
	"store/pkg/sdk/helper/wg"
	"store/pkg/sdk/third/bytedance/tos"
	"store/pkg/sdk/third/gemini"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/genai"
	"google.golang.org/protobuf/encoding/protojson"
)

func (t ProjService) AnalyzeSegment() {

	ctx := context.Background()

	hits, err := t.data.Mongo.TemplateSegment.List(ctx, bson.M{"status": "processing_segmented"})
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

	data := helper.Mapping(hits, func(x *projpb.ResourceSegment) *analyzeSegmentJobData {

		return &analyzeSegmentJobData{
			Item:     x,
			Settings: settings,
		}
	})

	wg.WaitGroup(ctx, data, t.analyzeSegment)
}

type analyzeSegmentJobData struct {
	Item     *projpb.ResourceSegment
	Settings *projpb.AppSettings
}

func (t ProjService) analyzeSegment(ctx context.Context, data *analyzeSegmentJobData) error {

	item := data.Item

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "analyzeSegment",
		"item", data.Item.XId,
		"prompt",
		helper.SubString(data.Settings.GetVideoHighlight().GetContent(), 0, 100),
	))

	logger.Debug("start to analyze segment")

	seg, err := t.analyzeSegmentByGemini(ctx, data)
	if err != nil {
		logger.Errorw("analyzeSegmentByGemini err", "err", err)

		return err
	}

	//seg, err := t.analyzeSegmentBySeed(ctx, data)
	//if err != nil {
	//	logger.Errorw("analyzeSegmentBySeed err", "err", err)
	//	return err
	//}

	item.Root.Commodity = seg.Root.GetCommodity()

	seg.TimeStart = item.TimeStart
	seg.TimeEnd = item.TimeEnd
	seg.Root = item.Root
	//seg.Status = "processing_analyzed"
	seg.Status = "completed"
	seg.CreatedAt = time.Now().Unix()
	seg.Category = item.Category

	seg.Tags = append(seg.Tags, seg.TypedTags.FocusOn...)
	seg.Tags = append(seg.Tags, seg.TypedTags.Text...)
	seg.Tags = append(seg.Tags, seg.TypedTags.Person...)
	seg.Tags = append(seg.Tags, seg.TypedTags.Picture...)
	seg.Tags = append(seg.Tags, seg.TypedTags.Scene...)
	seg.Tags = append(seg.Tags, seg.TypedTags.ShootingStyle...)
	seg.Tags = append(seg.Tags, seg.TypedTags.Emotion...)

	_, err = t.data.Mongo.TemplateSegment.ReplaceByID(ctx,
		item.XId,
		seg,
	)
	if err != nil {
		logger.Errorw("ReplaceByID err", "err", err)
		return err
	}

	return nil
}

//
//func (t ProjService) analyzeSegmentByGetGoAPI(ctx context.Context, data *analyzeSegmentJobData) (*projpb.ResourceSegment, error) {
//
//	logger := log.NewHelper(log.With(log.DefaultLogger,
//		"func", "analyzeSegmentByGetGoAPI",
//		"item", data.Item.XId,
//	))
//
//	return nil, nil
//}

func (t ProjService) analyzeSegmentByGemini(ctx context.Context, data *analyzeSegmentJobData) (*projpb.ResourceSegment, error) {

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "analyzeSegmentByGemini",
		"item", data.Item.XId,
	))

	logger.Debug("start to analyze segment by gemini")

	x := data.Item

	genaiClient := t.data.GenaiFactory.Get()

	segment, err := videoz.GetSegmentByUrl(ctx, x.Root.Url, x.TimeStart, x.TimeEnd)
	if err != nil {
		logger.Errorw("GetSegmentByUrl err", err)
		return nil, err
	}

	prompt, err := t.data.Mongo.Settings.GetPrompt(ctx, "highlight_video_analysis")
	if err != nil {
		return nil, err
	}

	parts := []*genai.Part{
		genai.NewPartFromBytes(segment.Content, "video/mp4"),
		{
			Text: prompt.Content,
		},
	}

	prompt, err = t.data.Mongo.Settings.GetPrompt(ctx, "product_analysis")
	if err != nil {
		return nil, err
	}

	parts = append(parts, gemini.NewTextPart(prompt.Content))

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
							Type: genai.TypeString,
						},
						"tags": {
							Type: genai.TypeArray,
							Items: &genai.Schema{
								Type: genai.TypeString,
							},
						},
					},
				},
				"segments": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type:        genai.TypeObject,
						Description: "只有1个，就是我提供给你的整个视频",
						Required: []string{
							//"highlightFrames",
							"style", "contentStyle", "sceneStyle", "description", "shootingStyle", "typedTags", "subtitle"},
						Properties: map[string]*genai.Schema{
							"formula": {
								Type:        genai.TypeString,
								Description: "爆款公式",
							},
							"productionMethod": {
								Type:        genai.TypeString,
								Description: "制作方式",
								Enum: []string{
									"creative",
									"realistic",
								},
							},
							"script": {
								Type:        genai.TypeString,
								Description: "详细分镜脚本",
							},
							"style": {
								Type:        genai.TypeString,
								Description: "拍摄脚本描述",
							},
							"contentStyle": {
								Type:        genai.TypeString,
								Description: "内容风格描述",
							},
							"sceneStyle": {
								Type:        genai.TypeString,
								Description: "场景设计描述",
							},

							"highlightFrames": {
								Type:        genai.TypeArray,
								Description: `关键画面帧`,
								Items: &genai.Schema{
									Type:     genai.TypeObject,
									Required: []string{"timestamp", "desc"},
									Properties: map[string]*genai.Schema{
										"timestamp": {
											Type:        genai.TypeNumber,
											Description: "当前帧对应的时间戳, 示例格式: 1.222",
										},
										"desc": {
											Type:        genai.TypeString,
											Description: "一句换描述这一帧",
										},
									},
								},
							},
							"description": {
								Type:        genai.TypeString,
								Description: "画面描述",
							},
							"shootingStyle": {
								Type:        genai.TypeString,
								Description: `拍摄手法`,
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
								Description: "口播文案",
							},
						},
					},
				},
			},
		},
	}
	text, err := genaiClient.GenerateContent(ctx, gemini.GenerateContentRequest{
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

	// 校验
	if len(tmpItem.Segments) != 1 {
		return nil, fmt.Errorf("invalid segments")
	}

	//if len(tmpItem.Segments[0].HighlightFrames) > 7 || len(tmpItem.Segments[0].HighlightFrames) == 0 {
	//	return nil, fmt.Errorf("invalid highlightFrames: %d", len(tmpItem.Segments[0].HighlightFrames))
	//}

	tmpItem.Segments[0].Root = &projpb.Resource{
		Commodity: tmpItem.Commodity,
	}

	for i := range tmpItem.Segments[0].HighlightFrames {

		xx := tmpItem.Segments[0].HighlightFrames[i]

		frame, err := videoz.GetFrame(segment.Content, xx.Timestamp)
		if err != nil {
			log.Debugw("GetFrame err", err, "segment", segment.From)
			t.data.Mongo.TemplateSegment.DeleteByID(ctx, data.Item.XId)
			return nil, err
		}

		if len(frame) == 0 {
			continue
		}

		tmpItem.Segments[0].HighlightFrames[i].Url, err = t.data.TOS.Put(ctx, tos.PutRequest{
			Bucket:  "yoozyres",
			Content: frame,
			Key:     helper.CreateUUID() + ".jpg",
		})

		if err != nil {
			return nil, err
		}
	}

	return tmpItem.Segments[0], nil
}

//func (t ProjService) analyzeSegmentBySeed(ctx context.Context, data *analyzeSegmentJobData) (*projpb.ResourceSegment, error) {
//
//	logger := log.NewHelper(log.With(log.DefaultLogger,
//		"func", "analyzeSegmentBySeed",
//		"item", data.Item.XId,
//	))
//
//	item := data.Item
//
//	logger.Debugw("start", "")
//
//	segment, err := videoz.GetSegmentByUrl(ctx, item.Root.Url, item.TimeStart, item.TimeEnd)
//	if err != nil {
//		logger.Errorw("GetSegmentByUrl err", err)
//		return nil, err
//	}
//	base64Video := base64.StdEncoding.EncodeToString(segment.Content)
//
//	var messages []*model.ChatCompletionMessage
//
//	messages = append(messages,
//		&model.ChatCompletionMessage{
//			Role: model.ChatMessageRoleSystem,
//			Content: &model.ChatCompletionMessageContent{
//				StringValue: volcengine.String(fmt.Sprintf(
//					`
//# 角色
//你是一位专业的抖音内容策略分析师和视频剪辑师。
//
//# 任务
//我将提供一个抖音带货视频的高光片段，请你完成以下任务：
//
//### 任务1: 根据【标签体系】为高光片段打标签：根据其内容和上述分类，打上相关的标签。
//
//### 任务2: 给出精细详尽的分镜脚本(*500字以上*)。分镜脚本要包含以下信息：
//	- 详细说明人物妆容、动作、神态，便于复刻原视频
//    - 详细说明画面构图，便于复刻原视频
//	- 详细说明转场及运镜方式，便于复刻原视频
//
//
//### 任务3: 提取片段的关键帧序列
//	为整个视频片段片段提取能够串联起该片段核心叙事逻辑的关键帧序列。提取时请遵循以下原则：
//	1. 分析叙事节奏: 把握每个片段内部的起承转合和节奏变化。
//	2. 识别关键节点: 识别出片段中的故事转折点、情绪高潮、核心信息展示等重要时刻。
//	3. 确保逻辑完整: 确保提取出的关键帧序列能够连贯、完整地复现该片段的叙事脉络。
//
//	输出要求
//	1. 关键帧个数：2-5个，不要超过这个范围
//	2. 帮我按要求确认以下字段：
//		- frame.timestamp: 关键帧的时间戳
//		- frame.desc: 一句话描述关键帧
//
//### 任务3：请详细给出拍摄手法。包括：
//	1. 镜头选择：具体说明使用什么焦距的镜头，如广角、标准、长焦等
//	2. 拍摄角度：详细描述拍摄角度，如正面、侧面、俯视、仰视等
//	3. 镜头运动：说明镜头的运动方式，如推、拉、摇、移、跟拍等
//	4. 构图方式：描述画面的构图原则和视觉重点
//	5. 演员指导：
//		- 眼神表达：具体说明眼神的方向、情感、专注度
//		- 表情管理：描述面部表情的细节，如微笑程度、情绪表达
//		- 语气语调：说明说话的语气、语调、语速、停顿等
//		- 肢体语言：描述手势、姿态、动作幅度等
//		- 状态调整：说明演员的整体状态和情绪准备
//	6. 表演细节：
//		- 情感层次：说明表演的情感层次和变化
//		- 互动方式：描述与镜头、道具、环境的互动
//		- 自然度：确保表演的自然和真实感
//	7. 技术细节：
//		- 对焦方式：说明手动对焦还是自动对焦
//		- 曝光控制：描述光圈、快门、ISO的设置
//		- 稳定方式：说明手持、三脚架、稳定器等使用方式",
//	8. 背景音乐:
//		- 背景音乐风格
//
//===
//
//【标签体系】
//
//- 主要阐述
//示例:
//	销量卖点: 描述产品受欢迎程度、复购情况。
//	赠品活动: 描述“买赠”、“送礼品”等非折扣优惠。
//	打折活动: 描述直接降价、折扣、限时优惠等。
//	价格信息: 描述具体价格或性价比。
//	服务卖点: 描述物流、售后、购物保障等。
//	人设背书: 描述品牌、产地、认证、推荐等建立信任的信息。
//	生产工艺: 描述制作方法、流程、配比等。
//	成分材质: 描述产品原料、营养成分等。
//	使用效果: 描述使用后的感受、变化、效果。
//	使用过程: 描述操作的便捷性、简易度。
//	外观样式: 描述产品的包装、形态、设计等。
//	适用人群: 描述产品的目标用户。
//	适用场景: 描述产品的使用场景、时机。
//	身体痛点: 描述用户身体上的不适或待解决的问题。
//	品质痛点: 描述市面上同类产品的缺点或用户的担忧。
//	价格痛点: 描述价格方面的顾虑或问题。
//	产品特点: 描述产品固有特性。
//	情感痛点: 描述与情感相关的痛点。
//
//- 文案类标签:
//
//输出示例: 正话反说、塑造冲突、提出疑问、威胁警告、直击痛点、引发好奇、情绪价值、价格优惠、场景需求、信用背书、种草推荐、效果展示、产品引入、身份介绍、产品材质、产品外观、产品工艺、产品背书、产品理念、产品情怀、产品价格、功能特点、痛点解决、使用场景、人群圈定、使用方法、种草体验、效果承诺
//
//- 人物类标签:
//
//格式: 形容词 + 人物角色
//
//示例: 穿毛衣的女演员, 穿旗袍的女主播, 精致妈妈（有妆容）
//
//- 动作类标签:
//
//示例: 展示洗衣液, 手指向商品
//
//- 画面类型类标签:
//
//示例:
//	视觉冲击类: 产品特写+动态细节画面, 对比冲击画面, 色彩暴击画面
//	痛点直击类: 场景还原画面（真人素材）, 痛点具像化画面, 情绪放大画面
//	利益前置型: 数字冲击画面, 结果呈现画面, 福利暗示画面
//	场景代入型: 场景全景+产品画面, 场景痛点+产品画面, 场景仪式感画面
//	悬念/冲突型: 矛盾行为画面, 悬念画面, 反转铺垫画面
//	信任背书型: 权威认证画面, 背书用户推荐画面, 品牌实力画面
//	产品细节特写: 原料 / 成分画面, 工艺 / 技术画面, 使用痕迹画面
//	场景化使用演示: 高频使用场景, 痛点解决场景, 仪式感场景
//	用户佐证 / 对比实验: 素人用户实测画面, 对比实验画面, 专家 / 达人背书画面
//	数据可视化: 成分数据画面, 效果数据画面, 效率数据画面
//	情感共鸣: 亲子 / 家庭场景, 自我变化场景, 怀旧 / 治愈场景
//	价格 / 福利暗示: 价格对比画面, 福利预告画面, 限时提示画面
//
//- 场景类标签 (类别、具体场景、典型元素):
//
//示例:
//	家庭环境场景、客厅/儿童房（活动场景）、婴儿床、爬行垫、玩具架、绘本、小桌椅
//	真实生活场景、办公室/工位、电脑、键盘、咖啡杯、外卖袋、工位隔板
//	功能需求场景、提神场景、揉眼睛、打哈欠、端起杯子猛喝
//	社交互动场景、朋友聚会、奶茶店、露营地、KTV包厢、拼盘零食
//	季节适配场景、夏季场景（冰爽/解暑）、空调、风扇、冰箱、冰块、汗水、西瓜
//	健康安全场景、权威认证场景、检测报告、认证证书、专家推荐
//	休闲零食场景、追剧/刷手机场景、沙发、靠垫、手机/平板、零食盘
//	原料溯源场景、产地直采场景（生鲜/水果）、果园、茶园、菜地、果农/菜农采摘
//	日常护理场景、浴室洗漱台、浴室镜（雾气水珠）、洗发水/沐浴露瓶身、毛巾
//	情感仪式场景、睡前护肤（治愈放松）、暖黄台灯、香薰机、面膜/晚霜
//	人群定制场景、职场精英、高端写字楼电梯、大牌护肤品、办公室（大班台、咖啡杯）
//	信任背书场景、用户真实测评（效果见证）、素人采访、前后对比图、聊天记录
//
//- 拍摄手法类标签 (运镜类型、具体动作):
//
//示例:
//	基础运镜: 推镜头, 拉镜头, 摇镜头, 移镜头
//	进阶运镜（动感&专业感）: 跟镜头, 升降镜头, 环绕镜头, 甩镜头, 变焦镜头
//	特殊视角运镜: 俯拍镜头, 仰拍镜头, 主观视角（第一人称视角）, 低角度平移/跟拍
//
//
//- 情绪类标签
//示例:
//	开心、兴奋、轻松、愉快、放松
//	温馨、感动、治愈、安心、满足
//	激励、自信、斗志、热血、振奋
//	紧张、悬疑、压迫、焦虑、惊悚
//	伤感、落寞、忧郁、失落、泪点
//	浪漫、甜蜜、心动、暧昧、期待
//	搞笑、俏皮、逗趣、恶搞、轻松诙谐
//	庄严、肃穆、敬畏、沉稳、正式
//
//**片段范围**:  %f-%f
//`, item.TimeStart, item.TimeEnd),
//				),
//			},
//		},
//		&model.ChatCompletionMessage{
//			Role: model.ChatMessageRoleSystem,
//			Content: &model.ChatCompletionMessageContent{
//				//StringValue: volcengine.String(
//				//	fmt.Sprintf(analysisPrompt2)),
//				ListValue: []*model.ChatCompletionMessageContentPart{
//					{
//						Type: "video_url",
//						VideoURL: &model.ChatMessageVideoURL{
//							URL: fmt.Sprintf("data:video/mp4;base64,%s", base64Video),
//						},
//					},
//				},
//			},
//		},
//	)
//
//	req := model.CreateChatCompletionRequest{
//		Model:    "doubao-seed-1-6-250615",
//		Messages: messages,
//		ResponseFormat: &model.ResponseFormat{
//			Type: model.ResponseFormatJSONSchema,
//			JSONSchema: &model.ResponseFormatJSONSchemaJSONSchemaParam{
//				Name:        "resource",
//				Description: "视频分析结果",
//				Schema:      arkr.GenerateSchema[AnalyzeResponse](),
//				Strict:      true,
//			},
//		},
//	}
//
//	resp, err := t.data.Arkr.C().CreateChatCompletion(ctx, req)
//	if err != nil {
//		log.Errorw("CreateChatCompletion err", err, "req", req)
//		return nil, err
//	}
//
//	var res *AnalyzeResponse
//	err = json.Unmarshal([]byte(*resp.Choices[0].Message.Content.StringValue), &res)
//	if err != nil {
//		return nil, err
//	}
//
//	if len(res.Segments) == 0 {
//		return nil, fmt.Errorf("no segments found")
//	}
//
//	seg := res.Segments[0]
//
//	tmpItem := &projpb.ResourceSegment{
//		Root: &projpb.Resource{
//			Commodity: &projpb.Commodity{
//				Name: res.Commodity.Name,
//				Tags: res.Commodity.Tags,
//			},
//		},
//		TypedTags: &projpb.TypedTags{
//			FocusOn:       seg.TypedTags.FocusOn,
//			Picture:       seg.TypedTags.Picture,
//			ShootingStyle: seg.TypedTags.ShootingStyle,
//			Scene:         seg.TypedTags.Scene,
//			Action:        seg.TypedTags.Action,
//			Person:        seg.TypedTags.Person,
//			Text:          seg.TypedTags.Text,
//			Emotion:       seg.TypedTags.Emotion,
//		},
//		Formula:       seg.Formula,
//		Script:        seg.Script,
//		Style:         seg.Style,
//		ContentStyle:  seg.ContentStyle,
//		SceneStyle:    seg.SceneStyle,
//		Description:   seg.Description,
//		ShootingStyle: seg.ShootingStyle,
//		Subtitle:      seg.Subtitle,
//	}
//
//	//segment, err := videoz.GetSegmentByUrl(ctx, item.Root.Url, item.TimeStart, item.TimeEnd)
//	////video, err := http.Get(item.Root.Url)
//	//if err != nil {
//	//	return nil, err
//	//}
//
//	for i := range seg.HighlightFrames {
//
//		xx := seg.HighlightFrames[i]
//
//		frame, err := videoz.GetFrame(segment.Content, xx.Timestamp)
//		if err != nil {
//			return nil, err
//		}
//
//		url, err := t.data.TOS.Put(ctx, tos.PutRequest{
//			Bucket:  "yoozyres",
//			Content: frame,
//			Key:     helper.CreateUUID() + ".jpg",
//		})
//
//		if err != nil {
//			return nil, err
//		}
//		tmpItem.HighlightFrames = append(tmpItem.HighlightFrames, &projpb.HighlightFrame{
//			Timestamp: xx.Timestamp,
//			Desc:      xx.Desc,
//			Url:       url,
//		})
//	}
//
//	return tmpItem, nil
//}
