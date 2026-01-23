package service

import (
	"context"
	"errors"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/helper/wg"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
)

func (t ProjService) ConfirmResource() {

	ctx := context.Background()

	list, err := t.data.Mongo.Task.List(ctx, bson.M{"status": "templateSelecting"})
	if err != nil {
		log.Errorw("List err", err)
		return
	}

	if len(list) == 0 {
		return
	}

	wg.WaitGroup(ctx, list, t.confirmResource)
}

func (t ProjService) confirmResource(ctx context.Context, task *projpb.Task) error {

	log.Debugw("confirmResource", task.Commodity.GetTags())

	// 拉取参考视频
	items, err := t.data.GrpcClients.ProjAdminClient.SearchMatchedItems(ctx, &projpb.SearchMatchedItemsParams{
		Keyword: strings.Join(task.Commodity.GetTags(), " "),
	})
	if err != nil {
		log.Error("SearchMatchedItems err", err, "keyword", task.Commodity.GetTags())
		return err
	}

	if len(items.List) == 0 {
		log.Errorw("docs is empty: search by", task.Commodity.GetTitle())
		return errors.New("docs is empty: search by " + task.Commodity.GetTitle())
	}

	item := items.List[0]

	var taskSegments []interface{}
	for _, x := range item.GetSegments() {

		x.Root = &projpb.Resource{
			XId:       item.XId,
			Url:       item.Url,
			CoverUrl:  item.CoverUrl,
			Commodity: item.Commodity,
		}

		taskSegments = append(taskSegments, &projpb.TaskSegment{
			Segment: x,
			TaskId:  task.XId,
			Task:    task,
			Status:  "textGenerating",
		})
	}

	_, err = t.data.Mongo.TaskSegment.GetCoreCollection().InsertMany(ctx, taskSegments)
	if err != nil {
		return err
	}
	//	script, err := t.doGenerateV2(ctx, task, item)
	//if err != nil {
	//	log.Errorw("generateBySeed err", err, "item", item.Id)
	//	return err
	//}
	//
	//log.Debugw("generateBySeed result", "", "generate script", script)
	//

	_, err = t.data.Mongo.Task.UpdateByIDIfExists(ctx,
		task.XId,
		mgz.Op().
			Set("status", "generating"),
	)

	if err != nil {
		log.Errorw("UpdateOneXX err", err, "task", task)
		return err
	}

	return nil
}

//type Result struct {
//	Subtitles []string `json:"subtitles" jsonschema_description:"生成的新的文案"`
//}
//
//func (t ProjService) doGenerateV2(ctx context.Context, task *projpb.Task, item *projpb.Item) (*projpb.Script, error) {
//
//	log.Debugw("generateBySeed", "", "参考item", item.Commodity.GetTitle())
//
//	prm := `
//你是一位顶级的抖音电商短视频策划师，尤其擅长通过“黄金3秒”、“痛点共鸣”、“价值塑造”和“逼单”等技巧，为新产品快速打造爆款文案。
//
//### 核心目标
//	让用户看完能明明白白做决策
//	内容需兼具干货价值与真实性。
//
//### 输出要求:
//	- 文案要说人话，口语化，接地气，需要符合抖音带货短视频的“快节奏”、“直切直给”的风格，完全复刻参考我提供的脚本文案
//	- 文案中不能出现抖音广告审核卡审的词语，包括涉及医疗效果、因果关系的卡审词，例如:发育、提升保护力、解决排便不畅"
//	- 口播文案的字数要和原文案相同
//	- 只输出改写后的新的口播文案
//	- 如果原文案是空的或者无意义的，输出空白字符串即可
//
//`
//	var messages []*model.ChatCompletionMessage
//
//	messages = append(messages,
//		&model.ChatCompletionMessage{
//			Role: model.ChatMessageRoleSystem,
//			Content: &model.ChatCompletionMessageContent{
//				StringValue: volcengine.String(prm),
//			},
//		},
//		&model.ChatCompletionMessage{
//			Role: model.ChatMessageRoleSystem,
//			Content: &model.ChatCompletionMessageContent{
//				StringValue: volcengine.String(
//					fmt.Sprintf(`
//帮我为每一个分镜头生成新的口播文案，需要结合分镜头信息和新的商品信息：%s,
//`,
//						conv.S2J(task.Commodity))),
//			},
//		},
//		&model.ChatCompletionMessage{
//			Role: model.ChatMessageRoleSystem,
//			Content: &model.ChatCompletionMessageContent{
//				ListValue: helper.Mapping(item.Analysis.GetSegments(), func(x *projpb.SegmentAnalysis) *model.ChatCompletionMessageContentPart {
//					return &model.ChatCompletionMessageContentPart{
//						Type: "text",
//						Text: conv.S2J(x),
//					}
//				}),
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
//				Name:        "subtitles",
//				Description: "新的口播文案",
//				Schema:      arkr.GenerateSchema[Result](),
//				Strict:      true,
//			},
//		},
//	}
//
//	resp, err := t.data.Arkr.C().CreateChatCompletion(ctx, req)
//	if err != nil {
//		return nil, err
//	}
//
//	var res Result
//	err = json.Unmarshal([]byte(*resp.Choices[0].Message.Content.StringValue), &res)
//	if err != nil {
//		return nil, err
//	}
//
//	fmt.Println("conv.S2J(res)", conv.S2J(res))
//
//	if len(res.Subtitles) != len(item.Analysis.GetSegments()) {
//		return nil, errors.New("generate failed")
//	}
//
//	var segments []*projpb.Script_Segment
//
//	for _, x := range res.Subtitles {
//		segments = append(segments, &projpb.Script_Segment{
//			Subtitle: x,
//		})
//	}
//
//	return &projpb.Script{
//		Segments: segments,
//	}, nil
//}
//
//func (t ProjService) generateBySeed(ctx context.Context, task *projpb.Task, item *projpb.Item) (*projpb.Script, error) {
//
//	log.Debugw("generateBySeed", "", "参考item", item.Commodity.GetTitle())
//
//	//prompts, err := t.data.Mongo.Settings.GetPrompts(ctx)
//	//if err != nil {
//	//	return nil, err
//	//}
//
//	//log.Debugw("GetPrompts", "success", "prompts", len(prompts))
//
//	task.Commodity.Url = ""
//
//	genaiClient := t.data.GenaiFactory.Get()
//
//	parts := []genai.Part{
//		genai.Text(`
//角色\n你是一位顶级的抖音电商短视频策划师，尤其擅长通过“黄金3秒”、“痛点共鸣”、“价值塑造”和“逼单”等技巧，为新产品快速打造爆款带货脚本。
//严格参考提供给你的参数视频的分镜头脚本（segments），为每一个镜头都按照新的产品和新的卖点复刻一段和参考视频相同结构爆款脚本`,
//		),
//		genai.Text("这提供给你的商品信息:" + conv.S2J(task.Commodity)),
//		genai.Text("这是选定的目标受众和目标卖点:" + conv.S2J(task.TargetChance)),
//		genai.Text(fmt.Sprintf("这是给你的参考的分镜脚本:%s, 镜头数为: %d, 最后生成的新的分镜脚本的镜头数要和这个保持一致", conv.S2J(item), len(item.Analysis.GetSegments()))),
//	}
//
//	//	"getCategory": "请仔细分析商品介绍中的图片，准确分析商品类型、使用场景、目标用户群体，选择最精确的类目。",
//	//		"sellingPoints": "请仔细分析商品介绍中的图片，提炼出该产品的**核心卖点**（至少3-5个）。",
//	//		"targetAudience": "请仔细分析商品介绍中的图片，**推导出**该产品的**核心目标用户群体**。",
//	//		"usageScenarios": "请仔细分析商品介绍中的图片，**推导出**该产品的**需求场景**。",
//	//		"shootingStyle": "\n请详细给出拍摄手法，\n包括：\n\n1. 镜头选择：具体说明使用什么焦距的镜头，如广角、标准、长焦等\n2. 拍摄角度：详细描述拍摄角度，如正面、侧面、俯视、仰视等\n3. 镜头运动：说明镜头的运动方式，如推、拉、摇、移、跟拍等\n4. 构图方式：描述画面的构图原则和视觉重点\n5. 演员指导：\n   - 眼神表达：具体说明眼神的方向、情感、专注度\n   - 表情管理：描述面部表情的细节，如微笑程度、情绪表达\n   - 语气语调：说明说话的语气、语调、语速、停顿等\n   - 肢体语言：描述手势、姿态、动作幅度等\n   - 状态调整：说明演员的整体状态和情绪准备\n6. 表演细节：\n   - 情感层次：说明表演的情感层次和变化\n   - 互动方式：描述与镜头、道具、环境的互动\n   - 自然度：确保表演的自然和真实感\n7. 技术细节：\n   - 对焦方式：说明手动对焦还是自动对焦\n   - 曝光控制：描述光圈、快门、ISO的设置\n   - 稳定方式：说明手持、三脚架、稳定器等使用方式\n\n多个镜头要用换行符号隔开",
//	//		"global": "# 角色\n你是一位顶级的抖音电商短视频策划师，尤其擅长通过“黄金3秒”、“痛点共鸣”、“价值塑造”和“逼单”等技巧，为新产品快速打造爆款带货脚本。\n\n# 任务\n请严格遵循我提供的 [参考爆款脚本] 的结构、节奏和每一句文案的“底层作用”，结合我给出的 [新产品信息]，为新产品创作一个全新的、能够快速起量的抖音带货短视频分镜头脚本。\n\n## [输出要求]\n1.  严格复刻卖货结构：新脚本必须每个镜头的“文案作用”必须与[参考爆款脚本]中的一一对应。\n2.  文案： 新脚本的所有文案必须是完全围绕[新产品信息]进行创作，必须用新产品的卖点和场景重写。\n3.  保持风格：保持原作 “带货话术”、“带货语气”、“直切直给”、“快节奏”、“强冲击力”的风格和语气。\n4.  语言： 使用中文。\n\n请基于以上信息，开始创作。",
//	//		"style": "整体拍摄脚本要求：风格需简洁、直接、节奏快，注重视觉冲击力和信息传达效率，符合抖音爆款带货视频的特点。",
//	//		"contentStyle": "整体内容风格要求：突出产品核心优势，解决用户痛点，引发购买欲望，语言口语化，具有亲和力。",
//	//		"sceneStyle": "整体场景设计要求：根据产品特性和目标人群，选择真实、贴近生活或能突出产品功能的场景，避免冗余和分散注意力的元素。",
//	//		"segments": "\n分镜头脚本：\n要求：\n1. 每个分段镜头都要有完整的拍摄指导，包括具体的画面、文案、拍摄手法（可参考shootingStyle但需结合本镜头具体化）。\n2. 注意段落之间的连贯性和过渡，但保持直给风格。\n3. 确保每个分段的重点突出，避免信息堆砌。\n4. 考虑整体节奏的平衡，快节奏为主。\n5. 在参考提供的视频卖货结构和标签的同时，**务必保证文案和画面内容的原创度，完全基于上方已分析的新产品信息，严禁照搬参考脚本的产品内容。**\n6. 结合商品的卖点对文案进行合理地改写和优化，使其更具吸引力。\n"
//	//}
//	//
//
//	generationConfig := &genai.GenerationConfig{
//		ResponseMIMEType: "application/json",
//		ResponseSchema: &genai.Schema{
//			Type:     genai.TypeObject,
//			Required: []string{"segments"},
//			Properties: map[string]*genai.Schema{
//				"segments": {
//					Type:        genai.TypeArray,
//					Description: "\n分镜头脚本：\n要求：\n1. 每个分段镜头都要有完整的拍摄指导，包括具体的画面、文案、拍摄手法（可参考shootingStyle但需结合本镜头具体化）。\n2. 注意段落之间的连贯性和过渡，但保持直给风格。\n3. 确保每个分段的重点突出，避免信息堆砌。\n4. 考虑整体节奏的平衡，快节奏为主。\n5. 在参考提供的视频卖货结构和标签的同时，**务必保证文案和画面内容的原创度，完全基于上方已分析的新产品信息，严禁照搬参考脚本的产品内容。**\n6. 结合商品的卖点对文案进行合理地改写和优化，使其更具吸引力。\n",
//					Items: &genai.Schema{
//						Type:     genai.TypeObject,
//						Required: []string{"content"},
//						Properties: map[string]*genai.Schema{
//							"content": {
//								Type:     genai.TypeObject,
//								Required: []string{"subtitle"},
//								Properties: map[string]*genai.Schema{
//									"subtitle": {
//										Type:        genai.TypeString,
//										Description: `口播文案`,
//									},
//								},
//							},
//						},
//					},
//				},
//			},
//		},
//	}
//
//	response, err := genaiClient.GenerateContent(ctx, geminiai.GenerateContentRequest{
//		Model:  "gemini-2.5-pro-preview-05-06",
//		Config: generationConfig,
//		Parts:  parts,
//	})
//
//	if err != nil {
//		log.Errorw("GenerateContent err", err)
//		return nil, err
//	}
//
//	var script projpb.Script
//
//	err = json.Unmarshal([]byte(response), &script)
//	if err != nil {
//		return nil, err
//	}
//
//	if script.Style == "" {
//		return nil, fmt.Errorf("empty script")
//	}
//
//	return &script, nil
//}
