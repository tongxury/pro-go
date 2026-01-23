package service

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/errors"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/genai"

	trackerpb "store/api/aiagent"
	"store/pkg/sdk/third/gemini"
	"strings"
)

func (t TrackerService) packSystemPromptV3(ctx context.Context, question *trackerpb.Question) ([]*genai.Part, error) {

	prompt := question.Prompt
	resources := trackerpb.Resources(question.Session.GetResources())

	configKey := fmt.Sprintf("%s", prompt.GetId())

	// todo
	if strings.HasPrefix(prompt.GetId(), "limitAnalysis") {
		personalProfile := resources.FindOneByCategory("personalProfile")
		if personalProfile != nil && personalProfile.Meta["platform"] == "douyin" {
			configKey += "_douyin"
		}
	}

	config, err := t.Data.Mongo.PromptConfig.Get(ctx, configKey)

	if err != nil {
		return nil, err
	}

	if config == nil {
		return nil, nil
	}

	var results []*genai.Part
	for _, x := range config.Parts {
		if x.Category == "" || x.Category == "text" {
			results = append(results, gemini.NewTextPart(x.Value))
		}

		if x.Category == "resource" {
			// 文件
			if x.Value == "" {

				for _, xx := range resources {
					if xx.Category != "" {
						continue
					}

					//toGenai, err := t.uploadToGenai(ctx, genaiClient, xx.Url, xx.MimeType)
					//if err != nil {
					//	log.Errorw("uploadToGenai err", err, "url", xx.Url)
					//	return nil, err
					//}

					part, err := gemini.NewMediaPart(xx.Url, xx.MimeType)
					if err != nil {
						return nil, err
					}

					results = append(results, part)
				}
			} else if x.Value == "history" {

				// 补充其他问题的答案 以备多轮
				questions, err := t.Data.Mongo.Question.List(ctx,
					bson.M{
						"session._id": question.GetSession().GetXId(),
						"status":      "completed",
					},
				)
				if err != nil {
					return nil, err
				}
				for _, xx := range questions {
					results = append(results, gemini.NewTextPart(xx.Text))
					results = append(results, gemini.NewTextPart(xx.Answer.GetText()))
				}
			} else if x.Value == "question" {
				if prompt.GetContent() != "" {
					results = append(results, gemini.NewTextPart(prompt.GetContent()))
				}

			} else {
				// 文本
				y := resources.FindOneByCategory(x.Value)
				if y != nil {

					//content := strings.ReplaceAll(y.Content, "")

					results = append(results, gemini.NewTextPart(y.Content))
				}
			}

		}
	}

	if len(results) == 0 {
		return nil, errors.BadRequest("invalidPromptId", "")
	}

	//c := genaiClient.GenerativeModel("").StartChat()
	//
	//c.SendMessage()

	// 前端传了 prompt的 content的话 要拼接上

	return results, nil
}
