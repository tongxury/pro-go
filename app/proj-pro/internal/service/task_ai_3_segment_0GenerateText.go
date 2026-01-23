package service

import (
	"context"
	"encoding/json"
	"fmt"
	projpb "store/api/proj"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper/wg"
	"store/pkg/sdk/third/bytedance/arkr"
	"store/pkg/sdk/third/gemini"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/genai"
)

func (t ProjService) GenerateText() {

	ctx := context.Background()

	list, err := t.data.Mongo.TaskSegment.List(ctx, bson.M{"status": "textGenerating"})
	if err != nil {
		log.Errorw("List err", err)
		return
	}

	if len(list) == 0 {
		return
	}

	wg.WaitGroup(ctx, list, t.generateText)
}

func (t ProjService) generateText(ctx context.Context, taskSegment *projpb.TaskSegment) error {

	res, err2 := t.generateBySeed(ctx, taskSegment)
	if err2 != nil {
		return err2
	}

	log.Debugw("doGenerateV2", conv.S2J(res))

	_, err := t.data.Mongo.TaskSegment.UpdateOneXXById(ctx,
		taskSegment.XId,
		bson.M{
			"subtitle": res.Subtitle,
			"status":   "textGenerated",
		},
	)
	if err != nil {
		return err
	}
	return nil
}

type Result struct {
	Subtitle string `json:"subtitle" jsonschema_description:"生成的新的文案"`
}

func (t ProjService) generateBySeed(ctx context.Context, taskSegment *projpb.TaskSegment) (*Result, error) {

	log.Debugw("generateBySeed generateText", "", "参考item", taskSegment.XId)

	prm := `
你是一位顶级的抖音电商短视频策划师，尤其擅长通过“黄金3秒”、“痛点共鸣”、“价值塑造”和“逼单”等技巧，为新产品快速打造爆款文案。

### 核心目标
	让用户看完能明明白白做决策
	内容需兼具干货价值与真实性。

### 输出要求:
	- 文案要说人话，口语化，接地气，需要符合抖音带货短视频的“快节奏”、“直切直给”的风格，完全复刻参考我提供的脚本文案
	- 文案中不能出现抖音广告审核卡审的词语，包括涉及医疗效果、因果关系的卡审词，例如:发育、提升保护力、解决排便不畅"
	- 口播文案的字数要和原文案相同
	- 只输出改写后的新的口播文案
`
	var messages []*model.ChatCompletionMessage

	messages = append(messages,
		&model.ChatCompletionMessage{
			Role: model.ChatMessageRoleSystem,
			Content: &model.ChatCompletionMessageContent{
				StringValue: volcengine.String(prm),
			},
		},
		&model.ChatCompletionMessage{
			Role: model.ChatMessageRoleSystem,
			Content: &model.ChatCompletionMessageContent{
				StringValue: volcengine.String(
					fmt.Sprintf(`
帮我生成新的口播文案，需要结合分镜头信息（%s）和新的商品信息（%s）, 
`,
						conv.S2J(taskSegment.Segment),
						conv.S2J(taskSegment.Task.Commodity))),
			},
		},
	)

	req := model.CreateChatCompletionRequest{
		Model:    "doubao-seed-1-6-250615",
		Messages: messages,
		ResponseFormat: &model.ResponseFormat{
			Type: model.ResponseFormatJSONSchema,
			JSONSchema: &model.ResponseFormatJSONSchemaJSONSchemaParam{
				Name:        "subtitle",
				Description: "新的口播文案",
				Schema:      arkr.GenerateSchema[Result](),
				Strict:      true,
			},
		},
	}

	resp, err := t.data.Arkr.C().CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, err
	}

	var res Result
	err = json.Unmarshal([]byte(*resp.Choices[0].Message.Content.StringValue), &res)
	if err != nil {
		return nil, err
	}

	fmt.Println("conv.S2J(res)", conv.S2J(res))

	return &res, nil
}
func (t ProjService) generateByGenai(ctx context.Context, taskSegment *projpb.TaskSegment) (*Result, error) {

	log.Debugw("generateByGenai generateText", "", "参考item", taskSegment.XId)

	prm := `
你是一位顶级的抖音电商短视频策划师，尤其擅长通过“黄金3秒”、“痛点共鸣”、“价值塑造”和“逼单”等技巧，为新产品快速打造爆款文案。

### 核心目标
	让用户看完能明明白白做决策
	内容需兼具干货价值与真实性。

### 输出要求:
	- 文案要说人话，口语化，接地气，需要符合抖音带货短视频的“快节奏”、“直切直给”的风格，完全复刻参考我提供的脚本文案
	- 文案中不能出现抖音广告审核卡审的词语，包括涉及医疗效果、因果关系的卡审词，例如:发育、提升保护力、解决排便不畅"
	- 口播文案的字数要和原文案相同
	- 只输出改写后的新的口播文案
`

	content, err := t.data.GenaiFactory.Get().GenerateContent(ctx, gemini.GenerateContentRequest{
		Model: "gemini-3-pro-preview",
		Config: &genai.GenerateContentConfig{
			ResponseMIMEType: "application/json",
			ResponseSchema: &genai.Schema{
				Type:     genai.TypeObject,
				Required: []string{"subtitle"},
				Properties: map[string]*genai.Schema{
					"subtitle": {
						Type:        genai.TypeString,
						Description: `口播文案`,
					},
				},
			},
		},
		Parts: []*genai.Part{
			{
				Text: prm,
			},
			{
				Text: fmt.Sprintf(`
帮我生成新的口播文案，需要结合分镜头信息（%s）和新的商品信息（%s）, 
`,
					conv.S2J(taskSegment.Segment),
					conv.S2J(taskSegment.Task.Commodity)),
			},
		},
	})
	if err != nil {
		log.Errorw("generateByGenai err", err)
		return nil, err
	}

	log.Debugw("GenerateContent ", content)

	var res Result
	err = json.Unmarshal([]byte(content), &res)
	if err != nil {
		log.Errorw("generateByGenai Unmarshal err", err)
		return nil, err
	}

	fmt.Println("conv.S2J(res)", conv.S2J(res))

	return &res, nil
}
