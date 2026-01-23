package service

import (
	"context"
	"encoding/json"
	"fmt"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"
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

func (t *SessionService) GenerateSubtitle() {

	ctx := context.Background()

	list, err := t.data.Mongo.SessionSegment.List(ctx, bson.M{"status": "created"})
	if err != nil {
		log.Errorw("List err", err)
		return
	}

	if len(list) == 0 {
		return
	}

	wg.WaitGroup(ctx, list, t.generateSubtitle)
}

func (t *SessionService) generateSubtitle(ctx context.Context, sessionSegment *projpb.SessionSegment) error {

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "generateSubtitle",
		"item", sessionSegment.XId,
	))

	logger.Debugw("start generate subtitle")
	res, err2 := t.generateBySeed(ctx, sessionSegment)
	if err2 != nil {
		return err2
	}

	_, err := t.data.Mongo.SessionSegment.UpdateByIDIfExists(ctx,
		sessionSegment.XId,
		mgz.Op().Sets(
			bson.M{
				"subtitle": res.Subtitle,
				"status":   "subtitleGenerated",
			}),
	)
	if err != nil {
		return err
	}
	return nil
}

func (t *SessionService) generateBySeed(ctx context.Context, sessionSegment *projpb.SessionSegment) (*Result, error) {

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "generateBySeed",
		"item", sessionSegment.XId,
	))

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
						conv.S2J(sessionSegment.Segment),
						conv.S2J(sessionSegment.Session.Commodity))),
			},
		},
	)

	req := model.CreateChatCompletionRequest{
		Model:    arkr.Seed16,
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

	logger.Debugw("res", res)

	return &res, nil
}
func (t *SessionService) generateByGenai(ctx context.Context, sessionSegment *projpb.SessionSegment) (*Result, error) {
	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "generateByGenai",
		"item", sessionSegment.XId,
	))
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
					conv.S2J(sessionSegment.Segment),
					conv.S2J(sessionSegment.Session.Commodity)),
			},
		},
	})
	if err != nil {
		logger.Errorw("generateByGenai err", err)
		return nil, err
	}

	logger.Debugw("GenerateContent ", content)

	var res Result
	err = json.Unmarshal([]byte(content), &res)
	if err != nil {
		logger.Errorw("generateByGenai Unmarshal err", err)
		return nil, err
	}

	logger.Debugw("res", res)

	return &res, nil
}
