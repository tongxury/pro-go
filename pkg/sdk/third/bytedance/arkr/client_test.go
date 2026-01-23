package arkr

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"store/pkg/sdk/conv"
	"testing"
)

type Result struct {
	Subtitles []string `json:"subtitles" jsonschema_description:"生成的新的文案"`
}

func TestName(t *testing.T) {

	c := NewClient()

	req := model.CreateChatCompletionRequest{
		Model: "doubao-seed-1-6-250615",
		Messages: []*model.ChatCompletionMessage{
			{
				Role: model.ChatMessageRoleUser,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String(`
 帮我生成两段口播文案
`,
					),
				},
			},
		},
		ResponseFormat: &model.ResponseFormat{
			Type: model.ResponseFormatJSONSchema,
			JSONSchema: &model.ResponseFormatJSONSchemaJSONSchemaParam{
				Name:        "subtitles",
				Description: "新的口播文案",
				Schema:      GenerateSchema[Result](),
				Strict:      true,
			},
		},
	}

	// 发送聊天完成请求，并将结果存储在 resp 中，将可能出现的错误存储在 err 中
	resp, err := c.C().CreateChatCompletion(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return
	}

	var res Result
	err = json.Unmarshal([]byte(*resp.Choices[0].Message.Content.StringValue), &res)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(conv.S2J(res))

}
